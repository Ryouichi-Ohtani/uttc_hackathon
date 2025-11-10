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

	// Step 1: AI画像分析（複数画像対応）
	imageAnalysis := ""
	if len(req.ImageURLs) > 0 {
		fmt.Printf("[AI_AGENT] Starting multi-image analysis for %d images\n", len(req.ImageURLs))

		// 複数画像を分析（最大3枚まで）
		maxImages := 3
		if len(req.ImageURLs) < maxImages {
			maxImages = len(req.ImageURLs)
		}

		var imageAnalysisResults []string
		for i := 0; i < maxImages; i++ {
			imageURL := req.ImageURLs[i]
			fmt.Printf("[AI_AGENT] Analyzing image %d/%d: %s\n", i+1, maxImages, imageURL)

			// 画像URLを実際のファイルパスに変換してbase64エンコード
			imageData, mimeType, err := uc.downloadAndEncodeImage(imageURL)
			if err != nil {
				fmt.Printf("[AI_AGENT ERROR] Failed to download/encode image %d: %v\n", i+1, err)
				continue
			}
			fmt.Printf("[AI_AGENT] Image %d encoded successfully, mimeType: %s, dataLength: %d\n", i+1, mimeType, len(imageData))

			// Gemini APIに画像分析をリクエスト
			fmt.Printf("[AI_AGENT] Calling Gemini API for image %d analysis...\n", i+1)
			analysis, err := uc.geminiClient.AnalyzeProductImage(ctx, imageData, mimeType)
			if err != nil {
				fmt.Printf("[AI_AGENT ERROR] Gemini API error for image %d: %v\n", i+1, err)
				continue
			}
			fmt.Printf("[AI_AGENT] Gemini analysis result for image %d: %s\n", i+1, analysis)
			imageAnalysisResults = append(imageAnalysisResults, analysis)
		}

		if len(imageAnalysisResults) == 0 {
			fmt.Println("[AI_AGENT ERROR] All image analyses failed")
			return uc.generateFallbackListing(userID, req)
		}

		// 複数画像の分析結果を統合
		if len(imageAnalysisResults) > 1 {
			imageAnalysis = uc.mergeMultipleAnalyses(imageAnalysisResults)
		} else {
			imageAnalysis = imageAnalysisResults[0]
		}
		fmt.Printf("[AI_AGENT] Final merged analysis: %s\n", imageAnalysis)
	} else {
		return nil, fmt.Errorf("at least one image is required")
	}

	// Step 2: 分析結果をパース
	suggestedProduct, confidence, err := uc.parseAIAnalysis(imageAnalysis)
	if err != nil {
		return nil, fmt.Errorf("failed to parse AI analysis: %w", err)
	}

	// Step 2.5: 市場価格データで価格を補正
	if suggestedProduct.Title != "" && suggestedProduct.Category != "" {
		marketPrice := uc.enhancePriceWithMarketData(ctx, suggestedProduct)
		if marketPrice > 0 && marketPrice != suggestedProduct.Price {
			fmt.Printf("[AI_AGENT] Market price correction: %d -> %d\n", suggestedProduct.Price, marketPrice)
			suggestedProduct.Price = marketPrice
			suggestedProduct.PricingRationale += fmt.Sprintf(" 市場価格データに基づき¥%d に調整しました。", marketPrice)
			confidence += 10 // 市場データ参照で確信度アップ
		}
	}

	// Step 2.6: ML価格予測で最終調整
	pricePrediction, err := uc.PredictOptimalPrice(ctx, suggestedProduct)
	if err == nil && pricePrediction != nil {
		fmt.Printf("[AI_AGENT] ML Price prediction: ¥%d (confidence: %d%%, expected days: %d)\n",
			pricePrediction.PredictedPrice, pricePrediction.Confidence, pricePrediction.ExpectedDaysToSell)

		// 高確信度の予測価格を採用
		if pricePrediction.Confidence >= 75 {
			suggestedProduct.Price = pricePrediction.PredictedPrice
			suggestedProduct.PricingRationale = fmt.Sprintf("ML予測により最適価格¥%d を提案（%d日以内に売却確率%.0f%%）。%s",
				pricePrediction.PredictedPrice,
				pricePrediction.ExpectedDaysToSell,
				pricePrediction.SellProbability.AtPredictedPrice*100,
				pricePrediction.Reasoning)
			confidence += 5 // ML予測で確信度アップ
		}
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
	feedbackData := make(map[string]interface{})

	// 修正内容を記録してフィードバック学習
	if title, ok := modifications["title"].(string); ok && title != product.Title {
		feedbackData["title_original"] = product.Title
		feedbackData["title_modified"] = title
		product.Title = title
		modifiedFields = append(modifiedFields, "title")
	}

	if desc, ok := modifications["description"].(string); ok && desc != product.Description {
		feedbackData["description_original"] = product.Description
		feedbackData["description_modified"] = desc
		product.Description = desc
		modifiedFields = append(modifiedFields, "description")
	}

	if price, ok := modifications["price"].(float64); ok && int(price) != product.Price {
		feedbackData["price_original"] = product.Price
		feedbackData["price_modified"] = int(price)
		product.Price = int(price)
		modifiedFields = append(modifiedFields, "price")
	}

	if category, ok := modifications["category"].(string); ok && category != product.Category {
		feedbackData["category_original"] = product.Category
		feedbackData["category_modified"] = category
		product.Category = category
		modifiedFields = append(modifiedFields, "category")
	}

	if condition, ok := modifications["condition"].(string); ok {
		feedbackData["condition_original"] = string(product.Condition)
		feedbackData["condition_modified"] = condition
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

	// === 新機能: AIフィードバック学習 ===
	// ユーザーの修正内容からAIが学習
	if len(modifiedFields) > 0 {
		go func() {
			uc.learnFromUserFeedback(context.Background(), product.Category, feedbackData, modifiedFields)
		}()
	}

	// ログ記録
	uc.aiAgentRepo.(*infrastructure.AIAgentRepositoryImpl).TrackAction(
		userID,
		domain.AgentTypeListing,
		"approved",
		productID,
		map[string]interface{}{
			"modified_fields": modifiedFields,
			"feedback_data":   feedbackData,
		},
		true,
		0,
	)

	return nil
}

// ========== AI Feedback Loop ==========

// learnFromUserFeedback - ユーザーの修正内容からAIが学習
func (uc *AIAgentUseCase) learnFromUserFeedback(ctx context.Context, category string, feedbackData map[string]interface{}, modifiedFields []string) {
	if uc.geminiClient == nil {
		return
	}

	feedbackJSON, _ := json.MarshalIndent(feedbackData, "", "  ")
	prompt := fmt.Sprintf(`あなたはAI出品エージェントの学習システムです。ユーザーがAI生成結果を修正しました。
この修正内容から学習し、今後の予測精度を向上させるためのインサイトを抽出してください:

【カテゴリー】%s
【修正されたフィールド】%v
【修正内容】
%s

【分析タスク】
1. ユーザーがなぜこの修正を行ったかを推測
2. AI予測の問題点を特定
3. カテゴリー別の傾向を分析
4. 今後の改善ポイントを提案

以下のJSON形式で学習結果を返してください:
{
  "learned_patterns": ["学習したパターン1", "パターン2"],
  "improvement_suggestions": ["改善提案1", "提案2"],
  "category_insights": "カテゴリー固有のインサイト",
  "confidence_adjustment": "予測確信度の調整方針"
}`,
		category,
		modifiedFields,
		string(feedbackJSON),
	)

	response, err := uc.geminiClient.GenerateContent(ctx, prompt)
	if err != nil {
		fmt.Printf("[FEEDBACK_LEARNING] Learning failed: %v\n", err)
		return
	}

	fmt.Printf("[FEEDBACK_LEARNING] Learned from user modifications:\n%s\n", response)

	// 実際の実装では、学習結果をデータベースに保存して
	// 次回のAI予測時に活用する
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

// ========== Shipping Tracking ==========

// UpdateShippingTracking - 配送追跡情報を更新
func (uc *AIAgentUseCase) UpdateShippingTracking(ctx context.Context, purchaseID uuid.UUID, trackingNumber string, carrier string) error {
	prep, err := uc.aiAgentRepo.GetShippingPreparation(purchaseID)
	if err != nil {
		return fmt.Errorf("shipping preparation not found: %w", err)
	}

	// 配送追跡番号を記録
	// 実際の実装では、配送業者APIと連携して自動追跡
	fmt.Printf("[SHIPPING_TRACKING] Purchase %s: Tracking number %s registered with %s\n",
		purchaseID, trackingNumber, carrier)

	now := time.Now()
	prep.ShippedAt = &now

	if err := uc.aiAgentRepo.UpdateShippingPreparation(prep); err != nil {
		return fmt.Errorf("failed to update shipping tracking: %w", err)
	}

	// 配送追跡の自動更新を開始（非同期）
	go uc.trackShipmentStatus(ctx, purchaseID, trackingNumber, carrier)

	return nil
}

// trackShipmentStatus - 配送状況を自動追跡
func (uc *AIAgentUseCase) trackShipmentStatus(ctx context.Context, purchaseID uuid.UUID, trackingNumber string, carrier string) {
	// 実際の実装では、配送業者APIを定期的にポーリング
	// ここでは概念的な実装

	fmt.Printf("[SHIPPING_TRACKING] Started tracking shipment %s\n", trackingNumber)

	// 配送状況の例
	stages := []string{
		"荷物を受け付けました",
		"配送センターに到着しました",
		"配送中です",
		"配達店に到着しました",
		"配達完了",
	}

	for i, stage := range stages {
		// 実際はAPIポーリング間隔（例: 1時間ごと）
		// ここでは概念的にログ出力のみ
		fmt.Printf("[SHIPPING_TRACKING] %s: %s (stage %d/%d)\n",
			trackingNumber, stage, i+1, len(stages))

		// 最終ステージで通知
		if i == len(stages)-1 {
			fmt.Printf("[SHIPPING_TRACKING] Shipment %s delivered successfully\n", trackingNumber)
		}
	}
}

// GetShippingStatus - 配送状況を取得（AI分析付き）
func (uc *AIAgentUseCase) GetShippingStatus(ctx context.Context, purchaseID uuid.UUID) (*ShippingStatus, error) {
	prep, err := uc.aiAgentRepo.GetShippingPreparation(purchaseID)
	if err != nil {
		return nil, fmt.Errorf("shipping preparation not found: %w", err)
	}

	// 配送状況を構築
	status := &ShippingStatus{
		PurchaseID: purchaseID,
		Carrier:    prep.SuggestedCarrier,
		Status:     "準備中",
		UpdatedAt:  time.Now(),
	}

	if prep.ShippedAt != nil {
		status.ShippedAt = prep.ShippedAt
		status.Status = "配送中"

		// AI予測: 配達予定日を計算
		estimatedDays := uc.estimateDeliveryDays(prep.SuggestedCarrier, prep.EstimatedWeight)
		estimatedDelivery := prep.ShippedAt.Add(time.Duration(estimatedDays*24) * time.Hour)
		status.EstimatedDelivery = &estimatedDelivery
		status.DeliveryDaysEstimate = estimatedDays
	}

	return status, nil
}

// estimateDeliveryDays - AI配送日数予測（Gemini API使用）
func (uc *AIAgentUseCase) estimateDeliveryDays(carrier string, weight float64) int {
	ctx := context.Background()

	if uc.geminiClient == nil {
		// フォールバック
		return uc.estimateDeliveryDaysFallback(carrier, weight)
	}

	prompt := fmt.Sprintf(`日本の配送業者「%s」で重量%.2fkgの荷物を送る場合、配送に何日かかりますか？

以下のJSON形式で返してください（前置きなしで、JSONのみを返してください）:
{
  "estimated_days": 配送日数（整数）,
  "reasoning": "推定理由"
}`,
		carrier,
		weight,
	)

	response, err := uc.geminiClient.GenerateContent(ctx, prompt)
	if err != nil {
		fmt.Printf("[AI_SHIPPING] Delivery days estimation failed, using fallback: %v\n", err)
		return uc.estimateDeliveryDaysFallback(carrier, weight)
	}

	var result struct {
		EstimatedDays int    `json:"estimated_days"`
		Reasoning     string `json:"reasoning"`
	}

	jsonStr := response
	if strings.Contains(response, "```json") {
		jsonStr = strings.Split(strings.Split(response, "```json")[1], "```")[0]
	} else if strings.Contains(response, "{") {
		start := strings.Index(response, "{")
		end := strings.LastIndex(response, "}")
		if start != -1 && end != -1 {
			jsonStr = response[start : end+1]
		}
	}

	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		fmt.Printf("[AI_SHIPPING] Failed to parse delivery days response, using fallback: %v\n", err)
		return uc.estimateDeliveryDaysFallback(carrier, weight)
	}

	fmt.Printf("[AI_SHIPPING] Estimated delivery: %d days (%s)\n", result.EstimatedDays, result.Reasoning)
	return result.EstimatedDays
}

// estimateDeliveryDaysFallback - フォールバック配送日数予測
func (uc *AIAgentUseCase) estimateDeliveryDaysFallback(carrier string, weight float64) int {
	// 配送業者・重量別の配送日数予測
	baseDays := 2 // デフォルト2日

	if strings.Contains(carrier, "ネコポス") {
		baseDays = 1
	} else if strings.Contains(carrier, "宅急便") {
		if weight > 5.0 {
			baseDays = 3
		} else {
			baseDays = 2
		}
	}

	return baseDays
}

// ShippingStatus represents shipping tracking status
type ShippingStatus struct {
	PurchaseID           uuid.UUID  `json:"purchase_id"`
	Carrier              string     `json:"carrier"`
	TrackingNumber       string     `json:"tracking_number"`
	Status               string     `json:"status"` // 準備中/配送中/配達完了
	ShippedAt            *time.Time `json:"shipped_at"`
	EstimatedDelivery    *time.Time `json:"estimated_delivery"`
	DeliveryDaysEstimate int        `json:"delivery_days_estimate"`
	UpdatedAt            time.Time  `json:"updated_at"`
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

	// === 新機能: AI市場価格分析による動的判定 ===
	ctx := context.Background()
	aiDecision := uc.makeSmartNegotiationDecision(ctx, offer, product, settings)
	if aiDecision != nil {
		return aiDecision
	}

	// === フォールバック: 従来の戦略ベース判定 ===
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

// makeSmartNegotiationDecision - AI市場価格分析 + 過去データ学習による高度な交渉判定
func (uc *AIAgentUseCase) makeSmartNegotiationDecision(ctx context.Context, offer *domain.Offer, product *domain.Product, settings *domain.AINegotiationSettings) *NegotiationDecision {
	if uc.geminiClient == nil {
		return nil
	}

	// 過去の交渉履歴を取得
	historicalOffers, _ := uc.offerRepo.FindByProductID(product.ID)
	acceptedCount := 0
	rejectedCount := 0
	avgAcceptedPrice := 0
	acceptedPrices := []int{}

	for _, histOffer := range historicalOffers {
		if histOffer.Status == domain.OfferStatusAccepted {
			acceptedCount++
			acceptedPrices = append(acceptedPrices, histOffer.OfferPrice)
		} else if histOffer.Status == domain.OfferStatusRejected {
			rejectedCount++
		}
	}

	if len(acceptedPrices) > 0 {
		sum := 0
		for _, p := range acceptedPrices {
			sum += p
		}
		avgAcceptedPrice = sum / len(acceptedPrices)
	}

	// AIに判定依頼
	prompt := fmt.Sprintf(`あなたはフリマアプリの価格交渉AIアシスタントです。以下の情報をもとに、オファーを承認すべきか判定してください:

【商品情報】
- 商品名: %s
- カテゴリー: %s
- 状態: %s
- 出品価格: ¥%d

【オファー情報】
- 提示価格: ¥%d
- 出品価格の%.1f%%
- 購入者のメッセージ: "%s"

【交渉設定】
- 戦略: %s (aggressive=積極的, moderate=中立, conservative=保守的)
- 最低許容価格: ¥%d
- 自動承認閾値: ¥%d以上
- 自動拒否閾値: ¥%d未満

【過去の交渉履歴】
- 承認した交渉: %d件
- 拒否した交渉: %d件
- 平均承認価格: ¥%d

【重要指示】
1. メルカリ、ヤフオク、Amazon、楽天で類似商品の市場価格を調査
2. 商品の状態と市場価格を考慮
3. 購入者のメッセージから真剣度・緊急度を分析
4. 過去の承認パターンから学習
5. 出品者の戦略(%s)を最優先で考慮

以下のJSON形式で返してください（前置きなしで、JSONのみを返してください）:
{
  "decision": "accept" または "reject",
  "confidence": 判定の確信度(0-100),
  "reasoning": "判定理由（市場価格データと過去データを具体的に引用）",
  "suggested_message": "購入者への返信メッセージ案"
}`,
		product.Title,
		product.Category,
		product.Condition,
		product.Price,
		offer.OfferPrice,
		float64(offer.OfferPrice)/float64(product.Price)*100,
		offer.Message,
		settings.NegotiationStrategy,
		settings.MinAcceptablePrice,
		settings.AutoAcceptThreshold,
		settings.AutoRejectThreshold,
		acceptedCount,
		rejectedCount,
		avgAcceptedPrice,
		settings.NegotiationStrategy,
	)

	response, err := uc.geminiClient.GenerateContent(ctx, prompt)
	if err != nil {
		fmt.Printf("[AI_NEGOTIATION] Smart decision failed: %v\n", err)
		return nil
	}

	// JSONパース
	var result struct {
		Decision         string `json:"decision"`
		Confidence       int    `json:"confidence"`
		Reasoning        string `json:"reasoning"`
		SuggestedMessage string `json:"suggested_message"`
	}

	jsonStr := response
	if strings.Contains(response, "```json") {
		jsonStr = strings.Split(strings.Split(response, "```json")[1], "```")[0]
	} else if strings.Contains(response, "{") {
		start := strings.Index(response, "{")
		end := strings.LastIndex(response, "}")
		if start != -1 && end != -1 {
			jsonStr = response[start : end+1]
		}
	}

	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		fmt.Printf("[AI_NEGOTIATION] Failed to parse smart decision: %v\n", err)
		return nil
	}

	fmt.Printf("[AI_NEGOTIATION] Smart decision: %s (confidence: %d%%) - %s\n",
		result.Decision, result.Confidence, result.Reasoning)

	// 確信度が低い場合はフォールバック
	if result.Confidence < 70 {
		fmt.Println("[AI_NEGOTIATION] Low confidence, using fallback strategy")
		return nil
	}

	status := domain.OfferStatusRejected
	if result.Decision == "accept" {
		status = domain.OfferStatusAccepted
	}

	return &NegotiationDecision{
		Status:  status,
		Message: result.SuggestedMessage,
		Reason:  fmt.Sprintf("ai_smart_decision_%d%%", result.Confidence),
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
	// Gemini APIを使って配送情報をAI推定
	ctx := context.Background()

	if uc.geminiClient == nil {
		// フォールバック: Gemini利用不可時は簡易アルゴリズム
		return uc.estimateShippingInfoFallback(product, purchase)
	}

	prompt := fmt.Sprintf(`あなたは日本の配送物流エキスパートAIです。以下の商品と配送先情報から、最適な配送方法を提案してください:

【商品情報】
- カテゴリー: %s
- 重量: %.2f kg
- 商品名: %s
- 状態: %s

【配送先】
- 住所: %s

【分析タスク】
1. 配送先の住所から地域（北海道/東北/関東/中部/関西/中国/四国/九州/沖縄）を判定
2. 商品の重量とサイズから最適な配送方法を選択（ヤマト運輸、佐川急便、日本郵便など）
3. パッケージサイズを決定（60/80/100/120サイズ以上）
4. 地域別の配送料金を計算（遠隔地は追加料金）
5. カテゴリ別の梱包指示を生成
6. 配送日数を予測

以下のJSON形式で返してください（前置きなしで、JSONのみを返してください）:
{
  "carrier": "推奨配送業者（例: ヤマト運輸（宅急便））",
  "package_size": "パッケージサイズ（例: 80サイズ）",
  "estimated_weight": 推定重量（kg）,
  "estimated_cost": 配送料（円、整数）,
  "instructions": "梱包指示（カテゴリと配送先を考慮した具体的な指示）",
  "estimated_delivery_days": 配送日数（整数）,
  "region": "配送先地域",
  "reasoning": "推奨理由"
}`,
		product.Category,
		product.WeightKg,
		product.Title,
		product.Condition,
		purchase.ShippingAddress,
	)

	response, err := uc.geminiClient.GenerateContent(ctx, prompt)
	if err != nil {
		fmt.Printf("[AI_SHIPPING ERROR] Gemini API error: %v, using fallback\n", err)
		return uc.estimateShippingInfoFallback(product, purchase)
	}

	// JSONパース
	var result struct {
		Carrier              string  `json:"carrier"`
		PackageSize          string  `json:"package_size"`
		EstimatedWeight      float64 `json:"estimated_weight"`
		EstimatedCost        int     `json:"estimated_cost"`
		Instructions         string  `json:"instructions"`
		EstimatedDeliveryDays int    `json:"estimated_delivery_days"`
		Region               string  `json:"region"`
		Reasoning            string  `json:"reasoning"`
	}

	jsonStr := response
	if strings.Contains(response, "```json") {
		jsonStr = strings.Split(strings.Split(response, "```json")[1], "```")[0]
	} else if strings.Contains(response, "{") {
		start := strings.Index(response, "{")
		end := strings.LastIndex(response, "}")
		if start != -1 && end != -1 {
			jsonStr = response[start : end+1]
		}
	}

	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		fmt.Printf("[AI_SHIPPING ERROR] Failed to parse Gemini response: %v, using fallback\n", err)
		return uc.estimateShippingInfoFallback(product, purchase)
	}

	fmt.Printf("[AI_SHIPPING] Gemini estimate: %s, %s, ¥%d (region: %s)\n",
		result.Carrier, result.PackageSize, result.EstimatedCost, result.Region)
	fmt.Printf("[AI_SHIPPING] Reasoning: %s\n", result.Reasoning)

	return &ShippingInfo{
		Carrier:      result.Carrier,
		PackageSize:  result.PackageSize,
		Weight:       result.EstimatedWeight,
		Cost:         result.EstimatedCost,
		Instructions: result.Instructions,
	}
}

// estimateShippingInfoFallback - Gemini利用不可時のフォールバック
func (uc *AIAgentUseCase) estimateShippingInfoFallback(product *domain.Product, purchase *domain.Purchase) *ShippingInfo {
	// カテゴリと重量から配送情報を推定（地域別対応版）
	info := &ShippingInfo{
		Weight: product.WeightKg,
	}

	// 配送先住所から地域を判定
	region := uc.detectRegion(purchase.ShippingAddress)
	fmt.Printf("[AI_SHIPPING] Detected region: %s from address: %s\n", region, purchase.ShippingAddress)

	// 重量とサイズから基本料金を算出
	baseSize := ""
	baseCost := 0

	if info.Weight <= 1.0 {
		baseSize = "60サイズ"
		baseCost = 800
		info.Carrier = "ヤマト運輸（ネコポス）"
	} else if info.Weight <= 3.0 {
		baseSize = "80サイズ"
		baseCost = 1000
		info.Carrier = "ヤマト運輸（宅急便）"
	} else if info.Weight <= 5.0 {
		baseSize = "100サイズ"
		baseCost = 1300
		info.Carrier = "ヤマト運輸（宅急便）"
	} else {
		baseSize = "120サイズ以上"
		baseCost = 1600
		info.Carrier = "ヤマト運輸（宅急便）"
	}

	// 地域別料金調整
	regionalCost := uc.adjustCostByRegion(baseCost, region)
	info.PackageSize = baseSize
	info.Cost = regionalCost

	fmt.Printf("[AI_SHIPPING] Base cost: ¥%d -> Regional cost: ¥%d (region: %s)\n", baseCost, regionalCost, region)

	// カテゴリ別の梱包指示
	switch product.Category {
	case "electronics":
		info.Instructions = "電子機器のため、緩衝材で厳重に包装してください。「精密機器」「取扱注意」のシールを貼付推奨。"
	case "clothing":
		info.Instructions = "衣類用ビニール袋に入れ、防水対策をしてください。"
	case "books":
		info.Instructions = "本の角を保護し、防水対策をしてください。"
	case "furniture":
		info.Instructions = "大型家具のため、複数箇所に緩衝材を使用し、エッジ保護を徹底してください。"
	default:
		info.Instructions = "商品が破損しないよう、適切な緩衝材を使用してください。"
	}

	// 配送先が遠隔地の場合は追加指示
	if region == "沖縄" || region == "北海道" {
		info.Instructions += " 遠隔地配送のため、配送日数に余裕を持ってください（通常+1〜2日）。"
	}

	return info
}

// detectRegion - 住所から地域を判定
func (uc *AIAgentUseCase) detectRegion(address string) string {
	regions := map[string][]string{
		"北海道": {"北海道"},
		"東北":   {"青森", "岩手", "宮城", "秋田", "山形", "福島"},
		"関東":   {"東京", "神奈川", "千葉", "埼玉", "茨城", "栃木", "群馬"},
		"中部":   {"新潟", "富山", "石川", "福井", "山梨", "長野", "岐阜", "静岡", "愛知"},
		"関西":   {"三重", "滋賀", "京都", "大阪", "兵庫", "奈良", "和歌山"},
		"中国":   {"鳥取", "島根", "岡山", "広島", "山口"},
		"四国":   {"徳島", "香川", "愛媛", "高知"},
		"九州":   {"福岡", "佐賀", "長崎", "熊本", "大分", "宮崎", "鹿児島"},
		"沖縄":   {"沖縄"},
	}

	for region, prefectures := range regions {
		for _, pref := range prefectures {
			if strings.Contains(address, pref) {
				return region
			}
		}
	}

	return "関東" // デフォルト
}

// adjustCostByRegion - 地域別に配送料を調整
func (uc *AIAgentUseCase) adjustCostByRegion(baseCost int, region string) int {
	// 地域別の料金調整率
	adjustments := map[string]float64{
		"北海道": 1.3,  // +30%
		"東北":   1.1,  // +10%
		"関東":   1.0,  // 基準
		"中部":   1.0,  // 基準
		"関西":   1.05, // +5%
		"中国":   1.1,  // +10%
		"四国":   1.15, // +15%
		"九州":   1.2,  // +20%
		"沖縄":   1.5,  // +50%
	}

	adjustment, exists := adjustments[region]
	if !exists {
		adjustment = 1.0
	}

	return int(float64(baseCost) * adjustment)
}

// GetAgentStats - ユーザーのAIエージェント利用統計
func (uc *AIAgentUseCase) GetAgentStats(ctx context.Context, userID uuid.UUID) (*domain.AIAgentStats, error) {
	return uc.aiAgentRepo.GetAgentStats(userID)
}

// ========== Price Prediction ML Model ==========

// PredictOptimalPrice - 過去取引データから最適価格を予測
func (uc *AIAgentUseCase) PredictOptimalPrice(ctx context.Context, product *domain.SuggestedProductData) (*PricePrediction, error) {
	if uc.geminiClient == nil {
		return nil, fmt.Errorf("Gemini client not available")
	}

	// 類似商品の過去取引データを取得
	historicalData := uc.getHistoricalPriceData(product.Category, product.Condition)

	prompt := fmt.Sprintf(`あなたは機械学習ベースの価格予測エンジンです。以下の商品の最適販売価格を予測してください:

【商品情報】
- 商品名: %s
- カテゴリー: %s
- 状態: %s
- ブランド: %s
- モデル: %s
- 重量: %.2f kg

【過去取引データ】
%s

【分析指示】
1. カテゴリー別の平均販売価格を分析
2. 商品の状態による価格補正率を計算
3. ブランド価値を考慮
4. 季節性・トレンドを分析
5. 売れるまでの平均日数を予測
6. 価格帯別の売却確率を計算

以下のJSON形式で返してください:
{
  "predicted_price": 予測価格（整数）,
  "confidence": 予測の確信度(0-100),
  "price_range": {
    "min": 最低推奨価格,
    "max": 最高推奨価格
  },
  "expected_days_to_sell": 売れるまでの予測日数,
  "sell_probability": {
    "at_predicted_price": 予測価格での売却確率(0-1),
    "at_10_percent_lower": 10%%安い場合の売却確率,
    "at_10_percent_higher": 10%%高い場合の売却確率
  },
  "reasoning": "予測根拠",
  "pricing_strategy": "recommended（推奨）/aggressive（高め）/conservative（安め）のいずれか"
}`,
		product.Title,
		product.Category,
		product.Condition,
		product.DetectedBrand,
		product.DetectedModel,
		product.WeightKg,
		historicalData,
	)

	response, err := uc.geminiClient.GenerateContent(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("price prediction failed: %w", err)
	}

	// JSONパース
	var result PricePrediction
	jsonStr := response
	if strings.Contains(response, "```json") {
		jsonStr = strings.Split(strings.Split(response, "```json")[1], "```")[0]
	} else if strings.Contains(response, "{") {
		start := strings.Index(response, "{")
		end := strings.LastIndex(response, "}")
		if start != -1 && end != -1 {
			jsonStr = response[start : end+1]
		}
	}

	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		return nil, fmt.Errorf("failed to parse price prediction: %w", err)
	}

	fmt.Printf("[PRICE_PREDICTION] Predicted: ¥%d (confidence: %d%%, days: %d)\n",
		result.PredictedPrice, result.Confidence, result.ExpectedDaysToSell)

	return &result, nil
}

// PricePrediction represents ML-based price prediction results
type PricePrediction struct {
	PredictedPrice      int                    `json:"predicted_price"`
	Confidence          int                    `json:"confidence"`
	PriceRange          PriceRange             `json:"price_range"`
	ExpectedDaysToSell  int                    `json:"expected_days_to_sell"`
	SellProbability     SellProbabilityMatrix  `json:"sell_probability"`
	Reasoning           string                 `json:"reasoning"`
	PricingStrategy     string                 `json:"pricing_strategy"`
}

type PriceRange struct {
	Min int `json:"min"`
	Max int `json:"max"`
}

type SellProbabilityMatrix struct {
	AtPredictedPrice  float64 `json:"at_predicted_price"`
	At10PercentLower  float64 `json:"at_10_percent_lower"`
	At10PercentHigher float64 `json:"at_10_percent_higher"`
}

// ========== Demand Prediction ==========

// PredictDemand - 商品の需要と売れ行きを予測
func (uc *AIAgentUseCase) PredictDemand(ctx context.Context, product *domain.SuggestedProductData, price int) (*DemandPrediction, error) {
	if uc.geminiClient == nil {
		return nil, fmt.Errorf("Gemini client not available")
	}

	// 現在の日時を取得（季節性分析用）
	now := time.Now()
	month := now.Month()
	season := uc.getSeason(month)

	prompt := fmt.Sprintf(`あなたは需要予測AIアナリストです。以下の商品の需要を予測してください:

【商品情報】
- 商品名: %s
- カテゴリー: %s
- 状態: %s
- 販売価格: ¥%d
- ブランド: %s
- 現在の季節: %s (%d月)

【分析指示】
1. カテゴリー別の季節需要を分析
2. 価格帯別の購入者層を特定
3. トレンド分析（人気度の推移）
4. 競合商品の数を推定
5. ターゲット顧客の購買意欲を評価
6. 最適な出品タイミングを提案

以下のJSON形式で返してください:
{
  "demand_level": "high/medium/low",
  "demand_score": 需要スコア(0-100),
  "expected_views": 予想閲覧数,
  "expected_inquiries": 予想問い合わせ数,
  "competition_level": "high/medium/low",
  "best_listing_time": "朝/昼/夕方/夜",
  "seasonal_factor": 季節要因の影響度(0-1),
  "trend_direction": "increasing/stable/decreasing",
  "target_audience": "ターゲット顧客層の説明",
  "optimization_tips": ["最適化のヒント1", "ヒント2", "ヒント3"],
  "reasoning": "需要予測の根拠"
}`,
		product.Title,
		product.Category,
		product.Condition,
		price,
		product.DetectedBrand,
		season,
		int(month),
	)

	response, err := uc.geminiClient.GenerateContent(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("demand prediction failed: %w", err)
	}

	// JSONパース
	var result DemandPrediction
	jsonStr := response
	if strings.Contains(response, "```json") {
		jsonStr = strings.Split(strings.Split(response, "```json")[1], "```")[0]
	} else if strings.Contains(response, "{") {
		start := strings.Index(response, "{")
		end := strings.LastIndex(response, "}")
		if start != -1 && end != -1 {
			jsonStr = response[start : end+1]
		}
	}

	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		return nil, fmt.Errorf("failed to parse demand prediction: %w", err)
	}

	fmt.Printf("[DEMAND_PREDICTION] Demand: %s (score: %d, views: %d)\n",
		result.DemandLevel, result.DemandScore, result.ExpectedViews)

	return &result, nil
}

// DemandPrediction represents demand prediction results
type DemandPrediction struct {
	DemandLevel        string   `json:"demand_level"`
	DemandScore        int      `json:"demand_score"`
	ExpectedViews      int      `json:"expected_views"`
	ExpectedInquiries  int      `json:"expected_inquiries"`
	CompetitionLevel   string   `json:"competition_level"`
	BestListingTime    string   `json:"best_listing_time"`
	SeasonalFactor     float64  `json:"seasonal_factor"`
	TrendDirection     string   `json:"trend_direction"`
	TargetAudience     string   `json:"target_audience"`
	OptimizationTips   []string `json:"optimization_tips"`
	Reasoning          string   `json:"reasoning"`
}

// getSeason - 月から季節を判定
func (uc *AIAgentUseCase) getSeason(month time.Month) string {
	switch month {
	case time.December, time.January, time.February:
		return "冬"
	case time.March, time.April, time.May:
		return "春"
	case time.June, time.July, time.August:
		return "夏"
	case time.September, time.October, time.November:
		return "秋"
	default:
		return "不明"
	}
}

// getHistoricalPriceData - カテゴリー・状態別の過去取引データを取得
func (uc *AIAgentUseCase) getHistoricalPriceData(category string, condition domain.ProductCondition) string {
	// 実際の実装では、データベースから過去30日の取引データを取得
	// ここでは簡易版として統計データを生成

	// カテゴリー別の平均価格（サンプルデータ）
	categoryAvgPrices := map[string]map[string]int{
		"electronics": {
			"new":      50000,
			"like_new": 35000,
			"good":     25000,
			"fair":     15000,
		},
		"clothing": {
			"new":      8000,
			"like_new": 5000,
			"good":     3000,
			"fair":     1500,
		},
		"books": {
			"new":      1500,
			"like_new": 1000,
			"good":     600,
			"fair":     300,
		},
		"furniture": {
			"new":      80000,
			"like_new": 50000,
			"good":     30000,
			"fair":     15000,
		},
		"toys": {
			"new":      5000,
			"like_new": 3000,
			"good":     1500,
			"fair":     800,
		},
		"sports": {
			"new":      15000,
			"like_new": 10000,
			"good":     6000,
			"fair":     3000,
		},
	}

	avgPrice := 10000 // デフォルト
	if prices, exists := categoryAvgPrices[category]; exists {
		if price, ok := prices[string(condition)]; ok {
			avgPrice = price
		}
	}

	return fmt.Sprintf(`過去30日間の取引データ（%s カテゴリー, %s 状態）:
- 平均販売価格: ¥%d
- 取引件数: 150件
- 平均売却日数: 7日
- 価格レンジ: ¥%d - ¥%d
- 最頻価格帯: ¥%d - ¥%d`,
		category, condition, avgPrice,
		int(float64(avgPrice)*0.7), int(float64(avgPrice)*1.3),
		int(float64(avgPrice)*0.9), int(float64(avgPrice)*1.1))
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

// enhancePriceWithMarketData - 市場価格データで価格を補正
func (uc *AIAgentUseCase) enhancePriceWithMarketData(ctx context.Context, product *domain.SuggestedProductData) int {
	if uc.geminiClient == nil {
		return product.Price
	}

	prompt := fmt.Sprintf(`以下の商品の適正な市場価格を分析してください:

商品名: %s
カテゴリー: %s
状態: %s
AI推定価格: ¥%d
ブランド: %s
モデル: %s

【重要指示】
1. メルカリ、ヤフオク、Amazon中古、楽天市場で実際に販売されている類似商品の価格を調査
2. 商品の状態(%s)を考慮した価格補正
3. 以下のJSON形式で返してください（前置きなしで、JSONのみを返してください）:

{
  "recommended_price": 推奨価格（整数）,
  "price_range_min": 最低価格（整数）,
  "price_range_max": 最高価格（整数）,
  "rationale": "価格根拠（市場データを具体的に引用）"
}`, product.Title, product.Category, product.Condition, product.Price, product.DetectedBrand, product.DetectedModel, product.Condition)

	response, err := uc.geminiClient.GenerateContent(ctx, prompt)
	if err != nil {
		fmt.Printf("[AI_AGENT] Market price analysis failed: %v\n", err)
		return product.Price
	}

	// JSONパース
	var result struct {
		RecommendedPrice int    `json:"recommended_price"`
		PriceRangeMin    int    `json:"price_range_min"`
		PriceRangeMax    int    `json:"price_range_max"`
		Rationale        string `json:"rationale"`
	}

	// JSONを抽出
	jsonStr := response
	if strings.Contains(response, "```json") {
		jsonStr = strings.Split(strings.Split(response, "```json")[1], "```")[0]
	} else if strings.Contains(response, "{") {
		start := strings.Index(response, "{")
		end := strings.LastIndex(response, "}")
		if start != -1 && end != -1 {
			jsonStr = response[start : end+1]
		}
	}

	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		fmt.Printf("[AI_AGENT] Failed to parse market price response: %v\n", err)
		return product.Price
	}

	fmt.Printf("[AI_AGENT] Market price analysis: ¥%d (range: ¥%d-¥%d) - %s\n",
		result.RecommendedPrice, result.PriceRangeMin, result.PriceRangeMax, result.Rationale)

	// 推奨価格が妥当な範囲内かチェック
	if result.RecommendedPrice > 100 && result.RecommendedPrice < 10000000 {
		return result.RecommendedPrice
	}

	return product.Price
}

// mergeMultipleAnalyses - 複数画像の分析結果を統合
func (uc *AIAgentUseCase) mergeMultipleAnalyses(analyses []string) string {
	if len(analyses) == 0 {
		return ""
	}
	if len(analyses) == 1 {
		return analyses[0]
	}

	// すべての分析結果をGemini APIで統合
	prompt := `以下は同じ商品の複数の画像から得られた分析結果です。これらを統合して、最も正確で詳細な商品情報を1つのJSONとして出力してください:

`
	for i, analysis := range analyses {
		prompt += fmt.Sprintf("\n画像%d の分析結果:\n%s\n", i+1, analysis)
	}

	prompt += `

統合時のルール:
1. 最も信頼性の高い情報を選択
2. 複数の画像で共通する情報は確実性が高い
3. 詳細な説明は複数の画像の情報を組み合わせる
4. 価格は複数の推定値の中央値を採用
5. 特徴(features)は重複を除いてマージ

以下のJSON形式で返してください:
{
  "title": "商品名",
  "description": "詳細説明",
  "category": "カテゴリ",
  "condition": "状態",
  "price": 価格,
  "weight_kg": 重量,
  "brand": "ブランド",
  "model": "モデル",
  "features": ["特徴1", "特徴2"],
  "pricing_rationale": "価格根拠",
  "category_rationale": "カテゴリ根拠"
}`

	ctx := context.Background()
	merged, err := uc.geminiClient.GenerateContent(ctx, prompt)
	if err != nil {
		fmt.Printf("[AI_AGENT ERROR] Failed to merge analyses: %v\n", err)
		// エラー時は最初の分析結果を返す
		return analyses[0]
	}

	return merged
}
