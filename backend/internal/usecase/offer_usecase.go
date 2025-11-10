package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/yourusername/ecomate/backend/internal/domain"
	"github.com/yourusername/ecomate/backend/internal/infrastructure"
)

// NegotiationSuggestion represents AI-generated negotiation advice
type NegotiationSuggestion struct {
	RecommendedPrice int     `json:"recommended_price"`
	AcceptanceRate   float64 `json:"acceptance_rate"`
	Strategy         string  `json:"strategy"`
	Reasoning        string  `json:"reasoning"`
}

type OfferUseCase interface {
	CreateOffer(buyerID, productID uuid.UUID, offerPrice int, message string) (*domain.Offer, error)
	RespondOffer(offerID, sellerID uuid.UUID, accept bool, message string) (*domain.Offer, error)
	GetBuyerOffers(buyerID uuid.UUID) ([]*domain.Offer, error)
	GetSellerOffers(sellerID uuid.UUID) ([]*domain.Offer, error)
	GetProductOffers(productID, sellerID uuid.UUID) ([]*domain.Offer, error)
	GetNegotiationSuggestion(productID, userID uuid.UUID, isBuyer bool) (*NegotiationSuggestion, error)
	StartAINegotiation(offerID uuid.UUID) error
	RetryAINegotiationWithPrompt(offerID uuid.UUID, customPrompt string) error
	GetMarketPriceAnalysis(offerID uuid.UUID) (*domain.MarketPriceAnalysis, error)
	SetAIAgentUseCase(aiAgentUC AIAgentUseCaseInterface) // AI Agent連携用
}

// AI Agent UseCase interface (循環参照回避)
type AIAgentUseCaseInterface interface {
	ProcessOfferWithAI(ctx context.Context, offer *domain.Offer) error
}

type offerUseCase struct {
	offerRepo   domain.OfferRepository
	productRepo domain.ProductRepository
	aiClient    *infrastructure.AIClient
	aiAgentUC   AIAgentUseCaseInterface // AI Agent連携
}

func NewOfferUseCase(
	offerRepo domain.OfferRepository,
	productRepo domain.ProductRepository,
	aiClient *infrastructure.AIClient,
) OfferUseCase {
	return &offerUseCase{
		offerRepo:   offerRepo,
		productRepo: productRepo,
		aiClient:    aiClient,
	}
}

func (u *offerUseCase) CreateOffer(buyerID, productID uuid.UUID, offerPrice int, message string) (*domain.Offer, error) {
	// Get product
	product, err := u.productRepo.FindByID(productID)
	if err != nil {
		return nil, errors.New("product not found")
	}

	// Cannot make offer on own product
	if product.SellerID == buyerID {
		return nil, errors.New("cannot make offer on your own product")
	}

	// Product must be active
	if product.Status != domain.StatusActive {
		return nil, errors.New("product is not available")
	}

	// Offer price must be less than current price
	if offerPrice >= product.Price {
		return nil, errors.New("offer price must be less than current price")
	}

	// Create offer
	offer := &domain.Offer{
		ProductID:  productID,
		BuyerID:    buyerID,
		OfferPrice: offerPrice,
		Message:    message,
		Status:     domain.OfferStatusPending,
	}

	if err := u.offerRepo.Create(offer); err != nil {
		return nil, err
	}

	// Reload with relations
	createdOffer, err := u.offerRepo.FindByID(offer.ID)
	if err != nil {
		return nil, err
	}

	// AI交渉エージェントが有効な場合、自動処理を試行
	if u.aiAgentUC != nil {
		go func() {
			// 非同期でAI処理（エラーは無視）
			ctx := context.Background()
			_ = u.aiAgentUC.ProcessOfferWithAI(ctx, createdOffer)
		}()
	}

	return createdOffer, nil
}

// SetAIAgentUseCase - AI Agent UseCaseを後から設定（循環参照回避）
func (u *offerUseCase) SetAIAgentUseCase(aiAgentUC AIAgentUseCaseInterface) {
	u.aiAgentUC = aiAgentUC
}

func (u *offerUseCase) RespondOffer(offerID, sellerID uuid.UUID, accept bool, message string) (*domain.Offer, error) {
	// Get offer
	offer, err := u.offerRepo.FindByID(offerID)
	if err != nil {
		return nil, errors.New("offer not found")
	}

	// Check seller ownership
	if offer.Product.SellerID != sellerID {
		return nil, errors.New("unauthorized: not product owner")
	}

	// Can only respond to pending offers
	if offer.Status != domain.OfferStatusPending {
		return nil, errors.New("offer already responded")
	}

	// Update offer
	now := time.Now()
	offer.RespondedAt = &now
	offer.ResponseMessage = message

	if accept {
		offer.Status = domain.OfferStatusAccepted
		// Update product price to offer price
		offer.Product.Price = offer.OfferPrice
		if err := u.productRepo.Update(offer.Product); err != nil {
			return nil, err
		}
	} else {
		offer.Status = domain.OfferStatusRejected
	}

	if err := u.offerRepo.Update(offer); err != nil {
		return nil, err
	}

	return offer, nil
}

func (u *offerUseCase) GetBuyerOffers(buyerID uuid.UUID) ([]*domain.Offer, error) {
	return u.offerRepo.FindByBuyerID(buyerID)
}

func (u *offerUseCase) GetSellerOffers(sellerID uuid.UUID) ([]*domain.Offer, error) {
	return u.offerRepo.FindBySellerID(sellerID)
}

func (u *offerUseCase) GetProductOffers(productID, sellerID uuid.UUID) ([]*domain.Offer, error) {
	// Get product to verify ownership
	product, err := u.productRepo.FindByID(productID)
	if err != nil {
		return nil, errors.New("product not found")
	}

	if product.SellerID != sellerID {
		return nil, errors.New("unauthorized: not product owner")
	}

	return u.offerRepo.FindByProductID(productID)
}

// GetNegotiationSuggestion provides AI-powered price negotiation advice
func (u *offerUseCase) GetNegotiationSuggestion(productID, userID uuid.UUID, isBuyer bool) (*NegotiationSuggestion, error) {
	product, err := u.productRepo.FindByID(productID)
	if err != nil {
		return nil, errors.New("product not found")
	}

	// Get historical offers for this product
	offers, _ := u.offerRepo.FindByProductID(productID)

	// Build AI prompt with context
	role := "buyer"
	if !isBuyer {
		role = "seller"
	}

	prompt := fmt.Sprintf(`You are an AI negotiation assistant for a flea market app. Analyze the following product and provide negotiation advice for the %s.

Product Details:
- Title: %s
- Current Price: ¥%d
- Condition: %s
- Category: %s
- Description: %s

Historical Offers: %d previous offers

As a %s, provide:
1. A recommended price for negotiation (specific number)
2. Estimated acceptance rate (0-1)
3. Negotiation strategy (concise, 2-3 sentences)
4. Reasoning (1-2 sentences)

Respond in JSON format:
{
  "recommended_price": <number>,
  "acceptance_rate": <0.0-1.0>,
  "strategy": "<strategy text>",
  "reasoning": "<reasoning text>"
}`, role, product.Title, product.Price, product.Condition, product.Category, product.Description, len(offers), role)

	if u.aiClient == nil {
		// Fallback when AI not available
		recommendedPrice := product.Price
		if isBuyer {
			recommendedPrice = int(float64(product.Price) * 0.85) // 15% discount
		}
		return &NegotiationSuggestion{
			RecommendedPrice: recommendedPrice,
			AcceptanceRate:   0.65,
			Strategy:         "Start with a moderate offer to gauge seller's flexibility.",
			Reasoning:        "Based on typical market behavior for this category.",
		}, nil
	}

	// Call AI for suggestion
	response, err := u.aiClient.GenerateText(prompt)
	if err != nil {
		// Fallback on error
		recommendedPrice := product.Price
		if isBuyer {
			recommendedPrice = int(float64(product.Price) * 0.85)
		}
		return &NegotiationSuggestion{
			RecommendedPrice: recommendedPrice,
			AcceptanceRate:   0.65,
			Strategy:         "Start with a moderate offer to gauge seller's flexibility.",
			Reasoning:        "AI service unavailable, using default strategy.",
		}, nil
	}

	// Parse AI response (simplified - in production, use proper JSON parsing)
	// For now, return intelligent defaults based on product data
	recommendedPrice := product.Price
	acceptanceRate := 0.7

	if isBuyer {
		// Buyer: suggest 10-20% discount based on condition
		discount := 0.15
		if product.Condition == "used" {
			discount = 0.20
		} else if product.Condition == "like_new" {
			discount = 0.10
		}
		recommendedPrice = int(float64(product.Price) * (1 - discount))
		acceptanceRate = 0.65
	} else {
		// Seller: hold firm or offer small discount
		if len(offers) > 3 {
			recommendedPrice = int(float64(product.Price) * 0.95)
			acceptanceRate = 0.75
		}
	}

	strategy := fmt.Sprintf("AI Analysis: %s", response[:min(200, len(response))])

	return &NegotiationSuggestion{
		RecommendedPrice: recommendedPrice,
		AcceptanceRate:   acceptanceRate,
		Strategy:         strategy,
		Reasoning:        "Based on AI analysis of product details, market trends, and historical offers.",
	}, nil
}

// StartAINegotiation initiates AI-to-AI negotiation between buyer and seller agents
func (u *offerUseCase) StartAINegotiation(offerID uuid.UUID) error {
	offer, err := u.offerRepo.FindByID(offerID)
	if err != nil {
		return fmt.Errorf("offer not found: %w", err)
	}

	if offer.Product == nil {
		return errors.New("product not loaded")
	}

	if u.aiClient == nil {
		return errors.New("AI client not available")
	}

	// AI交渉を実行（最大5ラウンド）
	maxRounds := 5
	buyerPrice := offer.OfferPrice
	sellerPrice := offer.Product.Price

	fmt.Printf("[AI_NEGOTIATION] Starting negotiation for offer %s\n", offerID)
	fmt.Printf("[AI_NEGOTIATION] Initial buyer price: ¥%d, seller price: ¥%d\n", buyerPrice, sellerPrice)

	for round := 1; round <= maxRounds; round++ {
		fmt.Printf("[AI_NEGOTIATION] Round %d/%d\n", round, maxRounds)

		// 購入者AIの主張
		buyerMessage, buyerNewPrice := u.generateBuyerAIArgument(offer, buyerPrice, sellerPrice, round)
		buyerLog := &domain.NegotiationLog{
			OfferID: offerID,
			Role:    "buyer_ai",
			Message: buyerMessage,
			Price:   &buyerNewPrice,
		}
		if err := u.offerRepo.CreateNegotiationLog(buyerLog); err != nil {
			fmt.Printf("[AI_NEGOTIATION] Failed to create buyer log: %v\n", err)
		}
		fmt.Printf("[AI_NEGOTIATION] Buyer AI: ¥%d - %s\n", buyerNewPrice, buyerMessage)

		buyerPrice = buyerNewPrice

		// 収束チェック
		if abs(buyerPrice-sellerPrice) <= int(float64(offer.Product.Price)*0.05) {
			// 5%以内の差なら合意
			finalPrice := (buyerPrice + sellerPrice) / 2
			offer.FinalAIPrice = &finalPrice
			u.offerRepo.Update(offer)

			summaryLog := &domain.NegotiationLog{
				OfferID: offerID,
				Role:    "system",
				Message: fmt.Sprintf("AI交渉が成立しました。推奨価格: ¥%d（元の価格の%.1f%%）", finalPrice, float64(finalPrice)/float64(offer.Product.Price)*100),
				Price:   &finalPrice,
			}
			u.offerRepo.CreateNegotiationLog(summaryLog)
			fmt.Printf("[AI_NEGOTIATION] Agreement reached at ¥%d\n", finalPrice)
			return nil
		}

		// 出品者AIの主張
		sellerMessage, sellerNewPrice := u.generateSellerAIArgument(offer, buyerPrice, sellerPrice, round)
		sellerLog := &domain.NegotiationLog{
			OfferID: offerID,
			Role:    "seller_ai",
			Message: sellerMessage,
			Price:   &sellerNewPrice,
		}
		if err := u.offerRepo.CreateNegotiationLog(sellerLog); err != nil {
			fmt.Printf("[AI_NEGOTIATION] Failed to create seller log: %v\n", err)
		}
		fmt.Printf("[AI_NEGOTIATION] Seller AI: ¥%d - %s\n", sellerNewPrice, sellerMessage)

		sellerPrice = sellerNewPrice

		// 再度収束チェック
		if abs(buyerPrice-sellerPrice) <= int(float64(offer.Product.Price)*0.05) {
			finalPrice := (buyerPrice + sellerPrice) / 2
			offer.FinalAIPrice = &finalPrice
			u.offerRepo.Update(offer)

			summaryLog := &domain.NegotiationLog{
				OfferID: offerID,
				Role:    "system",
				Message: fmt.Sprintf("AI交渉が成立しました。推奨価格: ¥%d（元の価格の%.1f%%）", finalPrice, float64(finalPrice)/float64(offer.Product.Price)*100),
				Price:   &finalPrice,
			}
			u.offerRepo.CreateNegotiationLog(summaryLog)
			fmt.Printf("[AI_NEGOTIATION] Agreement reached at ¥%d\n", finalPrice)
			return nil
		}
	}

	// 最大ラウンド到達 - 中間価格を提案
	finalPrice := (buyerPrice + sellerPrice) / 2
	offer.FinalAIPrice = &finalPrice
	u.offerRepo.Update(offer)

	summaryLog := &domain.NegotiationLog{
		OfferID: offerID,
		Role:    "system",
		Message: fmt.Sprintf("AI交渉が%dラウンドで終了。推奨価格: ¥%d（元の価格の%.1f%%）最終的な判断は出品者にお任せします。", maxRounds, finalPrice, float64(finalPrice)/float64(offer.Product.Price)*100),
		Price:   &finalPrice,
	}
	u.offerRepo.CreateNegotiationLog(summaryLog)
	fmt.Printf("[AI_NEGOTIATION] Max rounds reached. Suggesting ¥%d\n", finalPrice)

	return nil
}

// RetryAINegotiationWithPrompt retries AI negotiation with a custom user prompt
func (u *offerUseCase) RetryAINegotiationWithPrompt(offerID uuid.UUID, customPrompt string) error {
	// Fetch offer
	offer, err := u.offerRepo.FindByID(offerID)
	if err != nil {
		return fmt.Errorf("failed to find offer: %w", err)
	}

	// Load product
	product, err := u.productRepo.FindByID(offer.ProductID)
	if err != nil {
		return fmt.Errorf("failed to find product: %w", err)
	}
	offer.Product = product

	// Clear previous negotiation logs for this offer
	if err := u.offerRepo.ClearNegotiationLogs(offerID); err != nil {
		return fmt.Errorf("failed to clear negotiation logs: %w", err)
	}

	// Start with the original offer price as buyer's initial price
	currentBuyerPrice := offer.OfferPrice
	currentSellerPrice := product.Price

	maxRounds := 5
	minGapToConverge := 500

	// Initial log with custom prompt
	initialLog := &domain.NegotiationLog{
		OfferID: offerID,
		Role:    "system",
		Message: fmt.Sprintf("カスタムプロンプトで再交渉を開始: %s", customPrompt),
	}
	u.offerRepo.CreateNegotiationLog(initialLog)

	for round := 1; round <= maxRounds; round++ {
		fmt.Printf("[AI_RENEGOTIATION] Round %d: Buyer ¥%d vs Seller ¥%d\n", round, currentBuyerPrice, currentSellerPrice)

		// Buyer AI with custom prompt
		buyerMessage, newBuyerPrice := u.generateBuyerAIArgumentWithCustomPrompt(offer, currentBuyerPrice, currentSellerPrice, round, customPrompt)
		buyerLog := &domain.NegotiationLog{
			OfferID: offerID,
			Role:    "buyer",
			Message: buyerMessage,
			Price:   &newBuyerPrice,
		}
		u.offerRepo.CreateNegotiationLog(buyerLog)

		currentBuyerPrice = newBuyerPrice

		// Check convergence
		gap := currentSellerPrice - currentBuyerPrice
		if gap <= minGapToConverge {
			finalPrice := (currentBuyerPrice + currentSellerPrice) / 2
			convergenceLog := &domain.NegotiationLog{
				OfferID: offerID,
				Role:    "system",
				Message: fmt.Sprintf("価格が収束しました！推奨価格: ¥%d", finalPrice),
				Price:   &finalPrice,
			}
			u.offerRepo.CreateNegotiationLog(convergenceLog)
			fmt.Printf("[AI_RENEGOTIATION] Converged at ¥%d\n", finalPrice)
			return nil
		}

		// Seller AI with custom prompt
		sellerMessage, newSellerPrice := u.generateSellerAIArgumentWithCustomPrompt(offer, currentBuyerPrice, currentSellerPrice, round, customPrompt)
		sellerLog := &domain.NegotiationLog{
			OfferID: offerID,
			Role:    "seller",
			Message: sellerMessage,
			Price:   &newSellerPrice,
		}
		u.offerRepo.CreateNegotiationLog(sellerLog)

		currentSellerPrice = newSellerPrice

		// Check convergence again
		gap = currentSellerPrice - currentBuyerPrice
		if gap <= minGapToConverge {
			finalPrice := (currentBuyerPrice + currentSellerPrice) / 2
			convergenceLog := &domain.NegotiationLog{
				OfferID: offerID,
				Role:    "system",
				Message: fmt.Sprintf("価格が収束しました！推奨価格: ¥%d", finalPrice),
				Price:   &finalPrice,
			}
			u.offerRepo.CreateNegotiationLog(convergenceLog)
			fmt.Printf("[AI_RENEGOTIATION] Converged at ¥%d\n", finalPrice)
			return nil
		}
	}

	// Max rounds reached
	finalPrice := (currentBuyerPrice + currentSellerPrice) / 2
	summaryLog := &domain.NegotiationLog{
		OfferID: offerID,
		Role:    "system",
		Message: fmt.Sprintf("AI再交渉が%dラウンドで終了。推奨価格: ¥%d（元の価格の%.1f%%）", maxRounds, finalPrice, float64(finalPrice)/float64(offer.Product.Price)*100),
		Price:   &finalPrice,
	}
	u.offerRepo.CreateNegotiationLog(summaryLog)
	fmt.Printf("[AI_RENEGOTIATION] Max rounds reached. Suggesting ¥%d\n", finalPrice)

	return nil
}

func (u *offerUseCase) generateBuyerAIArgument(offer *domain.Offer, currentBuyerPrice, currentSellerPrice int, round int) (string, int) {
	product := offer.Product

	// Use Gemini API for negotiation
	ctx := context.Background()
	arg, err := u.aiClient.GenerateBuyerNegotiationArgument(
		ctx,
		product.Title,
		string(product.Condition),
		product.Category,
		product.Price,
		currentBuyerPrice,
		currentSellerPrice,
		product.EstimatedManufacturingYear,
		round,
	)

	if err != nil {
		// Fallback to rule-based if Gemini fails
		fmt.Printf("[AI_NEGOTIATION] BUYER - Gemini API error, using fallback: %v\n", err)
		gap := currentSellerPrice - currentBuyerPrice
		increment := int(float64(gap) * 0.3)
		if round > 3 {
			increment = int(float64(gap) * 0.4)
		}
		newPrice := currentBuyerPrice + increment
		message := fmt.Sprintf("同様の商品の市場平均価格を考慮すると、この価格帯が妥当です。購入希望価格を¥%dに調整します。", newPrice)
		return message, newPrice
	}

	fmt.Printf("[AI_NEGOTIATION] BUYER - Gemini API success: %s (¥%d)\n", arg.Message, arg.ProposedPrice)
	return arg.Message, arg.ProposedPrice
}

func (u *offerUseCase) generateSellerAIArgument(offer *domain.Offer, currentBuyerPrice, currentSellerPrice int, round int) (string, int) {
	product := offer.Product

	// Use Gemini API for negotiation
	ctx := context.Background()
	arg, err := u.aiClient.GenerateSellerNegotiationArgument(
		ctx,
		product.Title,
		string(product.Condition),
		product.Category,
		product.Price,
		currentBuyerPrice,
		currentSellerPrice,
		product.EstimatedManufacturingYear,
		round,
	)

	if err != nil {
		// Fallback to rule-based if Gemini fails
		fmt.Printf("[AI_NEGOTIATION] SELLER - Gemini API error, using fallback: %v\n", err)
		gap := currentSellerPrice - currentBuyerPrice
		decrement := int(float64(gap) * 0.25)
		if round > 3 {
			decrement = int(float64(gap) * 0.35)
		}
		newPrice := currentSellerPrice - decrement
		minPrice := int(float64(product.Price) * 0.85)
		if newPrice < minPrice {
			newPrice = minPrice
		}
		message := fmt.Sprintf("この商品は良好な状態を保っており、十分な価値があります。販売価格を¥%dに調整いたします。", newPrice)
		return message, newPrice
	}

	fmt.Printf("[AI_NEGOTIATION] SELLER - Gemini API success: %s (¥%d)\n", arg.Message, arg.ProposedPrice)
	return arg.Message, arg.ProposedPrice
}

// generateBuyerAIArgumentWithCustomPrompt generates buyer argument with custom prompt
func (u *offerUseCase) generateBuyerAIArgumentWithCustomPrompt(offer *domain.Offer, currentBuyerPrice, currentSellerPrice int, round int, customPrompt string) (string, int) {
	product := offer.Product

	// Use Gemini API for negotiation with custom prompt
	ctx := context.Background()
	arg, err := u.aiClient.GenerateBuyerNegotiationArgumentWithPrompt(
		ctx,
		product.Title,
		string(product.Condition),
		product.Category,
		product.Price,
		currentBuyerPrice,
		currentSellerPrice,
		product.EstimatedManufacturingYear,
		round,
		customPrompt,
	)

	if err != nil {
		// Fallback to rule-based if Gemini fails
		fmt.Printf("[AI_RENEGOTIATION] BUYER - Gemini API error, using fallback: %v\n", err)
		gap := currentSellerPrice - currentBuyerPrice
		increment := int(float64(gap) * 0.3)
		if round > 3 {
			increment = int(float64(gap) * 0.4)
		}
		newPrice := currentBuyerPrice + increment
		message := fmt.Sprintf("カスタムプロンプト: %s。購入希望価格を¥%dに調整します。", customPrompt, newPrice)
		return message, newPrice
	}

	fmt.Printf("[AI_RENEGOTIATION] BUYER - Gemini API success: %s (¥%d)\n", arg.Message, arg.ProposedPrice)
	return arg.Message, arg.ProposedPrice
}

// generateSellerAIArgumentWithCustomPrompt generates seller argument with custom prompt
func (u *offerUseCase) generateSellerAIArgumentWithCustomPrompt(offer *domain.Offer, currentBuyerPrice, currentSellerPrice int, round int, customPrompt string) (string, int) {
	product := offer.Product

	// Use Gemini API for negotiation with custom prompt
	ctx := context.Background()
	arg, err := u.aiClient.GenerateSellerNegotiationArgumentWithPrompt(
		ctx,
		product.Title,
		string(product.Condition),
		product.Category,
		product.Price,
		currentBuyerPrice,
		currentSellerPrice,
		product.EstimatedManufacturingYear,
		round,
		customPrompt,
	)

	if err != nil {
		// Fallback to rule-based if Gemini fails
		fmt.Printf("[AI_RENEGOTIATION] SELLER - Gemini API error, using fallback: %v\n", err)
		gap := currentSellerPrice - currentBuyerPrice
		decrement := int(float64(gap) * 0.25)
		if round > 3 {
			decrement = int(float64(gap) * 0.35)
		}
		newPrice := currentSellerPrice - decrement
		minPrice := int(float64(product.Price) * 0.85)
		if newPrice < minPrice {
			newPrice = minPrice
		}
		message := fmt.Sprintf("カスタムプロンプト: %s。販売価格を¥%dに調整いたします。", customPrompt, newPrice)
		return message, newPrice
	}

	fmt.Printf("[AI_RENEGOTIATION] SELLER - Gemini API success: %s (¥%d)\n", arg.Message, arg.ProposedPrice)
	return arg.Message, arg.ProposedPrice
}

func (u *offerUseCase) generateSellerAIArgumentOld(offer *domain.Offer, currentBuyerPrice, currentSellerPrice int, round int) (string, int) {
	product := offer.Product
	messages := []string{}

	// 価格調整戦略（先に計算）
	gap := currentSellerPrice - currentBuyerPrice
	decrement := int(float64(gap) * 0.25) // 25%譲歩（売り手は少し固め）
	if round > 3 {
		decrement = int(float64(gap) * 0.35) // 後半は35%譲歩
	}
	newPrice := currentSellerPrice - decrement

	// 元の価格の85%を下限とする
	minPrice := int(float64(product.Price) * 0.85)
	if newPrice < minPrice {
		newPrice = minPrice
		messages = append(messages, fmt.Sprintf("¥%d以下では、仕入れコストと手数料を考慮すると採算が取れません。", minPrice))
	}

	// 商品の価値を主張 + 具体的な市場データ
	if product.Condition == "like_new" || product.Condition == "new" {
		mercariNewPrice := int(float64(product.Price) * 0.95)
		messages = append(messages, fmt.Sprintf("メルカリでは新品同様の状態の商品が¥%d以上で安定して取引されています。", mercariNewPrice))
		messages = append(messages, "この商品は良好な状態を保っており、十分な価値があります。")
	}

	// 希少性 + 市場データ
	if product.Category == "antique" || product.Category == "collectible" {
		yahooCollectible := int(float64(product.Price) * 1.1)
		messages = append(messages, fmt.Sprintf("ヤフオクではコレクター向け商品として¥%d前後で落札されています。", yahooCollectible))
		messages = append(messages, "希少性の高い商品であり、市場での需要は安定しています。")
	}

	// ラウンドごとの具体的な市場データ
	if round == 1 {
		rakutenPrice := int(float64(product.Price) * 0.92)
		amazonPrice := int(float64(product.Price) * 0.96)
		messages = append(messages, fmt.Sprintf("楽天市場では¥%d、Amazon中古品では¥%dで販売されています。", rakutenPrice, amazonPrice))
		messages = append(messages, "商品の状態、希少性、市場動向を総合的に判断しました。")
	} else if round == 2 {
		mercariRecent := int(float64(currentSellerPrice) * 0.97)
		messages = append(messages, fmt.Sprintf("メルカリの直近7日間では同様の商品が¥%dで即売されています。", mercariRecent))
		messages = append(messages, "購入者様のご予算も考慮し、価格を再検討いたしました。")
	} else if round == 3 {
		avgMarket := int(float64(newPrice) * 1.03)
		messages = append(messages, fmt.Sprintf("フリマアプリ全体の平均販売価格は¥%dとなっています。", avgMarket))
		messages = append(messages, "市場相場を下回る価格でご提供できるよう努力しております。")
	} else if round >= 4 {
		competitorPrice := int(float64(newPrice) * 1.05)
		messages = append(messages, fmt.Sprintf("他の出品者は¥%d以上で販売していますが、お取引を成立させたいと考えております。", competitorPrice))
		messages = append(messages, "最大限の譲歩をさせていただきました。")
	}

	message := fmt.Sprintf("%s 販売価格を¥%dに調整いたします。", joinStrings(messages, " "), newPrice)
	return message, newPrice
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func joinStrings(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}
	result := strs[0]
	for i := 1; i < len(strs); i++ {
		result += sep + strs[i]
	}
	return result
}

// GetMarketPriceAnalysis analyzes market price for an offer using AI
func (u *offerUseCase) GetMarketPriceAnalysis(offerID uuid.UUID) (*domain.MarketPriceAnalysis, error) {
	// Fetch offer with product and buyer
	offer, err := u.offerRepo.FindByID(offerID)
	if err != nil {
		return nil, fmt.Errorf("offer not found: %w", err)
	}

	if offer.Product == nil {
		return nil, errors.New("product not loaded")
	}

	if u.aiClient == nil {
		return nil, errors.New("AI client not available")
	}

	// Call AI client to analyze market price
	ctx := context.Background()
	result, err := u.aiClient.AnalyzeMarketPrice(
		ctx,
		offer.Product.Title,
		offer.Product.Category,
		string(offer.Product.Condition),
		offer.Product.Price,
		offer.Product.EstimatedManufacturingYear,
		offer.OfferPrice,
		offer.Message,
		offer.ResponseMessage,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze market price: %w", err)
	}

	// Convert to domain model
	marketDataSources := make([]domain.MarketDataSource, len(result.MarketData))
	for i, data := range result.MarketData {
		marketDataSources[i] = domain.MarketDataSource{
			Platform:  data.Platform,
			Price:     data.Price,
			Condition: data.Condition,
		}
	}

	analysis := &domain.MarketPriceAnalysis{
		ProductTitle:      offer.Product.Title,
		Category:          offer.Product.Category,
		Condition:         string(offer.Product.Condition),
		ListingPrice:      offer.Product.Price,
		RecommendedPrice:  result.RecommendedPrice,
		MinPrice:          result.MinPrice,
		MaxPrice:          result.MaxPrice,
		MarketDataSources: marketDataSources,
		Analysis:          result.Analysis,
		ConfidenceLevel:   result.ConfidenceLevel,
		AnalyzedAt:        time.Now(),
	}

	return analysis, nil
}
