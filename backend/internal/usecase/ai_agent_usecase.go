package usecase

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/yourusername/ecomate/backend/internal/domain"
	"github.com/yourusername/ecomate/backend/internal/infrastructure"
)

type AIAgentUseCase struct {
	aiAgentRepo   domain.AIAgentRepository
	productRepo   domain.ProductRepository
	offerRepo     domain.OfferRepository
	purchaseRepo  domain.PurchaseRepository
	geminiClient  *infrastructure.GeminiClient
}

func NewAIAgentUseCase(
	aiAgentRepo domain.AIAgentRepository,
	productRepo domain.ProductRepository,
	offerRepo domain.OfferRepository,
	purchaseRepo domain.PurchaseRepository,
	geminiClient *infrastructure.GeminiClient,
) *AIAgentUseCase {
	return &AIAgentUseCase{
		aiAgentRepo:  aiAgentRepo,
		productRepo:  productRepo,
		offerRepo:    offerRepo,
		purchaseRepo: purchaseRepo,
		geminiClient: geminiClient,
	}
}

// ========== AI Listing Agent ==========

// GenerateListingFromImages - 画像からAIが商品情報を全自動生成
func (uc *AIAgentUseCase) GenerateListingFromImages(ctx context.Context, userID uuid.UUID, req *domain.AIListingGenerationRequest) (*domain.AIListingGenerationResponse, error) {
	startTime := time.Now()

	// Step 1: AI画像分析
	imageAnalysis := ""
	if len(req.ImageURLs) > 0 {
		fmt.Printf("[AI_AGENT] Starting image analysis for URL: %s\n", req.ImageURLs[0])

		// 画像URLを実際のファイルパスに変換してbase64エンコード
		imageData, mimeType, err := uc.downloadAndEncodeImage(req.ImageURLs[0])
		if err != nil {
			fmt.Printf("[AI_AGENT ERROR] Failed to download/encode image: %v\n", err)
			return uc.generateFallbackListing(userID, req)
		}
		fmt.Printf("[AI_AGENT] Image encoded successfully, mimeType: %s, dataLength: %d\n", mimeType, len(imageData))

		// Gemini APIに画像分析をリクエスト
		fmt.Println("[AI_AGENT] Calling Gemini API for image analysis...")
		analysis, err := uc.geminiClient.AnalyzeProductImage(ctx, imageData, mimeType)
		if err != nil {
			// AIエラー時はデフォルト値を返す
			fmt.Printf("[AI_AGENT ERROR] Gemini API error: %v\n", err)
			return uc.generateFallbackListing(userID, req)
		}
		fmt.Printf("[AI_AGENT] Gemini analysis result: %s\n", analysis)
		imageAnalysis = analysis
	} else {
		return nil, fmt.Errorf("at least one image is required")
	}

	// Step 2: 分析結果をパース
	suggestedProduct, confidence, err := uc.parseAIAnalysis(imageAnalysis)
	if err != nil {
		return nil, fmt.Errorf("failed to parse AI analysis: %w", err)
	}

	// Step 3: 商品を仮作成（ステータスはdraft）
	product := &domain.Product{
		SellerID:               userID,
		Title:                  suggestedProduct.Title,
		Description:            suggestedProduct.Description,
		Price:                  suggestedProduct.Price,
		Category:               suggestedProduct.Category,
		Condition:              suggestedProduct.Condition,
		WeightKg:               suggestedProduct.WeightKg,
		Status:                 domain.StatusDraft, // ユーザー承認待ち
		AIGeneratedDescription: suggestedProduct.Description,
		AISuggestedPrice:       suggestedProduct.Price,
	}

	if err := uc.productRepo.Create(product); err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	// 画像をProductImageとして保存
	for i, imageURL := range req.ImageURLs {
		productImage := &domain.ProductImage{
			ProductID:    product.ID,
			ImageURL:     imageURL,
			DisplayOrder: i,
			IsPrimary:    i == 0, // 最初の画像をプライマリに設定
		}
		product.Images = append(product.Images, *productImage)
	}

	// 画像を保存
	if len(product.Images) > 0 {
		if err := uc.productRepo.Update(product); err != nil {
			fmt.Printf("[AI_AGENT ERROR] Failed to save product images: %v\n", err)
		}
	}

	// Step 4: AI生成データを保存
	listingData := &domain.AIListingData{
		ProductID:            product.ID,
		IsAIGenerated:        true,
		AIConfidenceScore:    confidence,
		GeneratedTitle:       suggestedProduct.Title,
		GeneratedDescription: suggestedProduct.Description,
		GeneratedCategory:    suggestedProduct.Category,
		GeneratedCondition:   string(suggestedProduct.Condition),
		GeneratedPrice:       suggestedProduct.Price,
		UserModifiedFields:   "[]",
		ImageAnalysisResult:  imageAnalysis,
	}

	if err := uc.aiAgentRepo.CreateListingData(listingData); err != nil {
		return nil, fmt.Errorf("failed to save listing data: %w", err)
	}

	// Step 5: ログ記録
	processTime := int(time.Since(startTime).Milliseconds())
	uc.aiAgentRepo.(*infrastructure.AIAgentRepositoryImpl).TrackAction(
		userID,
		domain.AgentTypeListing,
		"generated",
		product.ID,
		map[string]interface{}{
			"confidence": confidence,
			"images":     len(req.ImageURLs),
		},
		true,
		processTime,
	)

	// Step 6: レスポンス構築
	response := &domain.AIListingGenerationResponse{
		ProductID:        product.ID,
		ListingData:      listingData,
		SuggestedProduct: suggestedProduct,
		ConfidenceBreakdown: map[string]float64{
			"title":       confidence,
			"description": confidence - 5,
			"category":    confidence + 5,
			"price":       confidence - 10,
		},
		RequiresApproval: !req.AutoPublish || confidence < 80,
	}

	return response, nil
}

// ApproveAndModifyListing - ユーザーが承認画面で修正した内容を反映
func (uc *AIAgentUseCase) ApproveAndModifyListing(ctx context.Context, userID uuid.UUID, productID uuid.UUID, modifications map[string]interface{}) error {
	// 商品取得
	product, err := uc.productRepo.FindByID(productID)
	if err != nil {
		return fmt.Errorf("product not found: %w", err)
	}

	if product.SellerID != userID {
		return fmt.Errorf("unauthorized")
	}

	// AI生成データ取得
	listingData, err := uc.aiAgentRepo.GetListingData(productID)
	if err != nil {
		return fmt.Errorf("listing data not found: %w", err)
	}

	// 修正内容を反映
	modifiedFields := []string{}

	if title, ok := modifications["title"].(string); ok && title != product.Title {
		product.Title = title
		modifiedFields = append(modifiedFields, "title")
	}

	if desc, ok := modifications["description"].(string); ok && desc != product.Description {
		product.Description = desc
		modifiedFields = append(modifiedFields, "description")
	}

	if price, ok := modifications["price"].(float64); ok && int(price) != product.Price {
		product.Price = int(price)
		modifiedFields = append(modifiedFields, "price")
	}

	if category, ok := modifications["category"].(string); ok && category != product.Category {
		product.Category = category
		modifiedFields = append(modifiedFields, "category")
	}

	if condition, ok := modifications["condition"].(string); ok {
		product.Condition = domain.ProductCondition(condition)
		modifiedFields = append(modifiedFields, "condition")
	}

	// 承認時にステータスをactiveに変更
	product.Status = domain.StatusActive

	// 商品更新
	if err := uc.productRepo.Update(product); err != nil {
		return fmt.Errorf("failed to update product: %w", err)
	}

	// 修正フィールドを記録
	modifiedFieldsJSON, _ := json.Marshal(modifiedFields)
	listingData.UserModifiedFields = string(modifiedFieldsJSON)

	if err := uc.aiAgentRepo.UpdateListingData(listingData); err != nil {
		return fmt.Errorf("failed to update listing data: %w", err)
	}

	// ログ記録
	uc.aiAgentRepo.(*infrastructure.AIAgentRepositoryImpl).TrackAction(
		userID,
		domain.AgentTypeListing,
		"approved",
		productID,
		map[string]interface{}{
			"modified_fields": modifiedFields,
		},
		true,
		0,
	)

	return nil
}

// ========== AI Negotiation Agent ==========

// EnableAINegotiation - AI交渉モードを有効化
func (uc *AIAgentUseCase) EnableAINegotiation(ctx context.Context, userID uuid.UUID, req *domain.ToggleAINegotiationRequest) error {
	// 商品の所有権確認
	product, err := uc.productRepo.FindByID(req.ProductID)
	if err != nil {
		return fmt.Errorf("product not found: %w", err)
	}

	if product.SellerID != userID {
		return fmt.Errorf("unauthorized")
	}

	// AI交渉設定を作成/更新
	settings, err := uc.aiAgentRepo.GetNegotiationSettings(req.ProductID)
	if err != nil {
		// 新規作成
		settings = &domain.AINegotiationSettings{
			ProductID:           req.ProductID,
			Mode:                req.Mode,
			IsEnabled:           true,
			MinAcceptablePrice:  req.MinAcceptablePrice,
			AutoAcceptThreshold: req.AutoAcceptThreshold,
			AutoRejectThreshold: req.AutoRejectThreshold,
			NegotiationStrategy: req.Strategy,
		}

		if err := uc.aiAgentRepo.CreateNegotiationSettings(settings); err != nil {
			return fmt.Errorf("failed to create negotiation settings: %w", err)
		}
	} else {
		// 既存設定を更新
		settings.Mode = req.Mode
		settings.IsEnabled = true
		settings.MinAcceptablePrice = req.MinAcceptablePrice
		settings.AutoAcceptThreshold = req.AutoAcceptThreshold
		settings.AutoRejectThreshold = req.AutoRejectThreshold
		settings.NegotiationStrategy = req.Strategy

		if err := uc.aiAgentRepo.UpdateNegotiationSettings(settings); err != nil {
			return fmt.Errorf("failed to update negotiation settings: %w", err)
		}
	}

	return nil
}

// ProcessOfferWithAI - オファーをAIが自動処理
func (uc *AIAgentUseCase) ProcessOfferWithAI(ctx context.Context, offer *domain.Offer) error {
	// AI交渉設定を取得
	settings, err := uc.aiAgentRepo.GetNegotiationSettings(offer.ProductID)
	if err != nil || !settings.IsEnabled || settings.Mode != domain.NegotiationModeAI {
		// AI処理が無効の場合はスキップ
		return nil
	}

	// 商品情報取得
	product, err := uc.productRepo.FindByID(offer.ProductID)
	if err != nil {
		return err
	}

	// AI判断ロジック
	decision := uc.makeAINegotiationDecision(offer, product, settings)

	// オファーステータス更新
	offer.Status = decision.Status
	offer.ResponseMessage = decision.Message

	if decision.Status == domain.OfferStatusAccepted {
		settings.AIAcceptedCount++
	} else if decision.Status == domain.OfferStatusRejected {
		settings.AIRejectedCount++
	}

	settings.TotalOffersProcessed++

	// 更新
	if err := uc.offerRepo.Update(offer); err != nil {
		return err
	}

	if err := uc.aiAgentRepo.UpdateNegotiationSettings(settings); err != nil {
		return err
	}

	// ログ記録
	uc.aiAgentRepo.(*infrastructure.AIAgentRepositoryImpl).TrackAction(
		product.SellerID,
		domain.AgentTypeNegotiation,
		string(decision.Status),
		offer.ID,
		map[string]interface{}{
			"offer_price": offer.OfferPrice,
			"decision":    decision.Reason,
		},
		true,
		0,
	)

	return nil
}

// ========== AI Shipping Preparation Agent ==========

// PrepareShipping - 購入後に配送情報をAIが自動準備
func (uc *AIAgentUseCase) PrepareShipping(ctx context.Context, purchaseID uuid.UUID) (*domain.AIShippingPreparation, error) {
	// 購入情報取得
	purchase, err := uc.purchaseRepo.FindByID(purchaseID)
	if err != nil {
		return nil, fmt.Errorf("purchase not found: %w", err)
	}

	// 商品情報取得
	product, err := uc.productRepo.FindByID(purchase.ProductID)
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}

	// AIで配送情報を推定
	shippingInfo := uc.estimateShippingInfo(product, purchase)

	// 配送準備データを作成
	prep := &domain.AIShippingPreparation{
		PurchaseID:           purchaseID,
		IsAIPrepared:         true,
		SuggestedCarrier:     shippingInfo.Carrier,
		SuggestedPackageSize: shippingInfo.PackageSize,
		EstimatedWeight:      shippingInfo.Weight,
		EstimatedCost:        shippingInfo.Cost,
		ShippingInstructions: shippingInfo.Instructions,
		UserApproved:         false,
	}

	if err := uc.aiAgentRepo.CreateShippingPreparation(prep); err != nil {
		return nil, fmt.Errorf("failed to create shipping preparation: %w", err)
	}

	// ログ記録
	uc.aiAgentRepo.(*infrastructure.AIAgentRepositoryImpl).TrackAction(
		purchase.SellerID,
		domain.AgentTypeShipping,
		"prepared",
		purchaseID,
		shippingInfo,
		true,
		0,
	)

	return prep, nil
}

// ApproveShipping - ユーザーが配送情報を承認
func (uc *AIAgentUseCase) ApproveShipping(ctx context.Context, userID uuid.UUID, purchaseID uuid.UUID, req *domain.ApproveShippingRequest) error {
	prep, err := uc.aiAgentRepo.GetShippingPreparation(purchaseID)
	if err != nil {
		return fmt.Errorf("shipping preparation not found: %w", err)
	}

	// 権限確認
	if prep.Purchase.SellerID != userID {
		return fmt.Errorf("unauthorized")
	}

	if req.Approved {
		// 修正内容を反映
		if req.Carrier != "" {
			prep.SuggestedCarrier = req.Carrier
		}
		if req.PackageSize != "" {
			prep.SuggestedPackageSize = req.PackageSize
		}

		// 承認
		if err := uc.aiAgentRepo.ApproveShipping(purchaseID, req.Modifications); err != nil {
			return fmt.Errorf("failed to approve shipping: %w", err)
		}

		// ログ記録
		uc.aiAgentRepo.(*infrastructure.AIAgentRepositoryImpl).TrackAction(
			userID,
			domain.AgentTypeShipping,
			"approved",
			purchaseID,
			req,
			true,
			0,
		)
	}

	return nil
}

// ========== Helper Methods ==========

func (uc *AIAgentUseCase) buildImageAnalysisPrompt(userHints string) string {
	prompt := `この商品画像を詳細に分析し、以下のJSON形式で返してください:

{
  "title": "商品名（簡潔に20文字以内）",
  "description": "詳細説明（200-300文字）",
  "category": "clothing/electronics/furniture/books/toys/sports のいずれか",
  "condition": "new/like_new/good/fair のいずれか",
  "price": 推定価格（円）,
  "weight_kg": 推定重量（kg）,
  "brand": "ブランド名（あれば）",
  "model": "モデル名（あれば）",
  "features": ["特徴1", "特徴2", "特徴3"],
  "pricing_rationale": "価格設定の根拠",
  "category_rationale": "カテゴリ選定の根拠"
}`

	if userHints != "" {
		prompt += "\n\nユーザーからのヒント: " + userHints
	}

	return prompt
}

func (uc *AIAgentUseCase) parseAIAnalysis(analysis string) (*domain.SuggestedProductData, float64, error) {
	// JSONをパース
	var result struct {
		Title             string   `json:"title"`
		Description       string   `json:"description"`
		Category          string   `json:"category"`
		Condition         string   `json:"condition"`
		Price             int      `json:"price"`
		WeightKg          float64  `json:"weight_kg"`
		Brand             string   `json:"brand"`
		Model             string   `json:"model"`
		Features          []string `json:"features"`
		PricingRationale  string   `json:"pricing_rationale"`
		CategoryRationale string   `json:"category_rationale"`
	}

	// ```json ... ``` の部分を抽出
	jsonStr := analysis
	if strings.Contains(analysis, "```json") {
		jsonStr = strings.Split(strings.Split(analysis, "```json")[1], "```")[0]
	}

	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		// パースエラーの場合はデフォルト値
		return &domain.SuggestedProductData{
			Title:       "商品名（要修正）",
			Description: "説明文（要修正）",
			Category:    "electronics",
			Condition:   domain.ConditionGood,
			Price:       1000,
		}, 50.0, nil
	}

	suggested := &domain.SuggestedProductData{
		Title:             result.Title,
		Description:       result.Description,
		Category:          result.Category,
		Condition:         domain.ProductCondition(result.Condition),
		Price:             result.Price,
		WeightKg:          result.WeightKg,
		DetectedBrand:     result.Brand,
		DetectedModel:     result.Model,
		KeyFeatures:       result.Features,
		PricingRationale:  result.PricingRationale,
		CategoryRationale: result.CategoryRationale,
	}

	// 信頼度を計算（簡易版）
	confidence := 85.0
	if result.Brand != "" {
		confidence += 5
	}
	if len(result.Features) >= 3 {
		confidence += 5
	}

	return suggested, confidence, nil
}

type NegotiationDecision struct {
	Status  domain.OfferStatus
	Message string
	Reason  string
}

func (uc *AIAgentUseCase) makeAINegotiationDecision(offer *domain.Offer, product *domain.Product, settings *domain.AINegotiationSettings) *NegotiationDecision {
	// 自動承認判定
	if offer.OfferPrice >= settings.AutoAcceptThreshold {
		return &NegotiationDecision{
			Status:  domain.OfferStatusAccepted,
			Message: "ご提案ありがとうございます！ご希望の価格で承認いたしました。",
			Reason:  "auto_accept_threshold_met",
		}
	}

	// 自動拒否判定
	if offer.OfferPrice < settings.AutoRejectThreshold {
		return &NegotiationDecision{
			Status:  domain.OfferStatusRejected,
			Message: "申し訳ございませんが、ご提示の価格では承認できません。",
			Reason:  "below_auto_reject_threshold",
		}
	}

	// 中間価格の場合は戦略に応じて判定
	priceRatio := float64(offer.OfferPrice) / float64(product.Price)

	switch settings.NegotiationStrategy {
	case "aggressive":
		// 攻撃的: 90%以上で承認
		if priceRatio >= 0.90 {
			return &NegotiationDecision{
				Status:  domain.OfferStatusAccepted,
				Message: "ご提案を承認いたしました！",
				Reason:  "aggressive_strategy",
			}
		}
	case "moderate":
		// 中立的: 80%以上で承認
		if priceRatio >= 0.80 {
			return &NegotiationDecision{
				Status:  domain.OfferStatusAccepted,
				Message: "お待ちいただきありがとうございます。ご提案を承認いたしました。",
				Reason:  "moderate_strategy",
			}
		}
	case "conservative":
		// 保守的: 95%以上で承認
		if priceRatio >= 0.95 {
			return &NegotiationDecision{
				Status:  domain.OfferStatusAccepted,
				Message: "ご提案を承認いたしました。",
				Reason:  "conservative_strategy",
			}
		}
	}

	// デフォルトは拒否
	return &NegotiationDecision{
		Status:  domain.OfferStatusRejected,
		Message: fmt.Sprintf("ご提案ありがとうございます。恐れ入りますが、¥%d 以上でご検討いただけますと幸いです。", settings.MinAcceptablePrice),
		Reason:  "strategy_threshold_not_met",
	}
}

type ShippingInfo struct {
	Carrier      string
	PackageSize  string
	Weight       float64
	Cost         int
	Instructions string
}

func (uc *AIAgentUseCase) estimateShippingInfo(product *domain.Product, purchase *domain.Purchase) *ShippingInfo {
	// カテゴリと重量から配送情報を推定（簡易版）
	info := &ShippingInfo{
		Weight: product.WeightKg,
	}

	// 重量でサイズ判定
	if info.Weight <= 1.0 {
		info.PackageSize = "60サイズ"
		info.Cost = 800
		info.Carrier = "ヤマト運輸（ネコポス）"
	} else if info.Weight <= 3.0 {
		info.PackageSize = "80サイズ"
		info.Cost = 1000
		info.Carrier = "ヤマト運輸（宅急便）"
	} else if info.Weight <= 5.0 {
		info.PackageSize = "100サイズ"
		info.Cost = 1300
		info.Carrier = "ヤマト運輸（宅急便）"
	} else {
		info.PackageSize = "120サイズ以上"
		info.Cost = 1600
		info.Carrier = "ヤマト運輸（宅急便）"
	}

	// カテゴリ別の梱包指示
	switch product.Category {
	case "electronics":
		info.Instructions = "電子機器のため、緩衝材で厳重に包装してください。「精密機器」「取扱注意」のシールを貼付推奨。"
	case "clothing":
		info.Instructions = "衣類用ビニール袋に入れ、防水対策をしてください。"
	case "books":
		info.Instructions = "本の角を保護し、防水対策をしてください。"
	default:
		info.Instructions = "商品が破損しないよう、適切な緩衝材を使用してください。"
	}

	return info
}

// GetAgentStats - ユーザーのAIエージェント利用統計
func (uc *AIAgentUseCase) GetAgentStats(ctx context.Context, userID uuid.UUID) (*domain.AIAgentStats, error) {
	return uc.aiAgentRepo.GetAgentStats(userID)
}

// downloadAndEncodeImage - 画像URLからファイルを読み込みbase64エンコード
func (uc *AIAgentUseCase) downloadAndEncodeImage(imageURL string) (string, string, error) {
	// URLからファイルパスを抽出 (例: http://localhost:8080/uploads/xxx.png -> ./uploads/xxx.png)
	filePath := strings.Replace(imageURL, "http://localhost:8080/uploads/", "./uploads/", 1)

	// ファイルを読み込み
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", "", fmt.Errorf("failed to read image file: %w", err)
	}

	// base64エンコード
	encoded := base64.StdEncoding.EncodeToString(data)

	// MIMEタイプを判定
	mimeType := "image/jpeg"
	if strings.HasSuffix(filePath, ".png") {
		mimeType = "image/png"
	} else if strings.HasSuffix(filePath, ".jpg") || strings.HasSuffix(filePath, ".jpeg") {
		mimeType = "image/jpeg"
	} else if strings.HasSuffix(filePath, ".webp") {
		mimeType = "image/webp"
	}

	return encoded, mimeType, nil
}

// generateFallbackListing - AI失敗時のフォールバック
func (uc *AIAgentUseCase) generateFallbackListing(userID uuid.UUID, req *domain.AIListingGenerationRequest) (*domain.AIListingGenerationResponse, error) {
	// デフォルト商品情報
	defaultProduct := &domain.SuggestedProductData{
		Title:             "商品名を入力してください",
		Description:       "商品の詳細を入力してください。状態、サイズ、色などの情報を記載すると売れやすくなります。",
		Category:          "electronics",
		Condition:         domain.ConditionGood,
		Price:             1000,
		WeightKg:          0.5,
		KeyFeatures:       []string{},
		PricingRationale:  "AI分析に失敗したため、デフォルト価格を設定しています",
		CategoryRationale: "画像から判定できませんでした",
	}

	// ユーザーヒントがあれば反映
	if req.UserHints != "" {
		defaultProduct.Title = req.UserHints
		defaultProduct.Description = "ヒント: " + req.UserHints + "\n\n商品の詳細を追加で入力してください。"
	}

	// 仮の商品を作成
	product := &domain.Product{
		SellerID:               userID,
		Title:                  defaultProduct.Title,
		Description:            defaultProduct.Description,
		Price:                  defaultProduct.Price,
		Category:               defaultProduct.Category,
		Condition:              defaultProduct.Condition,
		WeightKg:               defaultProduct.WeightKg,
		Status:                 domain.StatusActive,
		AIGeneratedDescription: defaultProduct.Description,
		AISuggestedPrice:       defaultProduct.Price,
	}

	if err := uc.productRepo.Create(product); err != nil {
		return nil, fmt.Errorf("failed to create fallback product: %w", err)
	}

	listingData := &domain.AIListingData{
		ProductID:            product.ID,
		IsAIGenerated:        false, // フォールバックなのでfalse
		AIConfidenceScore:    30.0,  // 低い確信度
		GeneratedTitle:       defaultProduct.Title,
		GeneratedDescription: defaultProduct.Description,
		GeneratedCategory:    defaultProduct.Category,
		GeneratedCondition:   string(defaultProduct.Condition),
		GeneratedPrice:       defaultProduct.Price,
		UserModifiedFields:   "[]",
		ImageAnalysisResult:  "AI analysis failed - using fallback values",
	}

	if err := uc.aiAgentRepo.CreateListingData(listingData); err != nil {
		return nil, fmt.Errorf("failed to save fallback listing data: %w", err)
	}

	return &domain.AIListingGenerationResponse{
		ProductID:        product.ID,
		ListingData:      listingData,
		SuggestedProduct: defaultProduct,
		ConfidenceBreakdown: map[string]float64{
			"title":       30.0,
			"description": 30.0,
			"category":    30.0,
			"price":       30.0,
		},
		RequiresApproval: true, // 必ず承認が必要
	}, nil
}

// ========== Public Methods for Handlers ==========

// GetListingData - AIリスティングデータを取得
func (uc *AIAgentUseCase) GetListingData(productID uuid.UUID) (*domain.AIListingData, error) {
	return uc.aiAgentRepo.GetListingData(productID)
}

// GetNegotiationSettings - AI交渉設定を取得
func (uc *AIAgentUseCase) GetNegotiationSettings(productID uuid.UUID) (*domain.AINegotiationSettings, error) {
	return uc.aiAgentRepo.GetNegotiationSettings(productID)
}

// DisableNegotiation - AI交渉を無効化
func (uc *AIAgentUseCase) DisableNegotiation(productID uuid.UUID) error {
	return uc.aiAgentRepo.DeleteNegotiationSettings(productID)
}

// GetShippingPreparation - AI配送準備データを取得
func (uc *AIAgentUseCase) GetShippingPreparation(purchaseID uuid.UUID) (*domain.AIShippingPreparation, error) {
	return uc.aiAgentRepo.GetShippingPreparation(purchaseID)
}
