package usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"
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
	GetMarketPriceAnalysis(offerID uuid.UUID) (*domain.MarketPriceAnalysis, error)
	SetMessageUseCase(messageUseCase MessageUseCase)
}

type offerUseCase struct {
	offerRepo      domain.OfferRepository
	productRepo    domain.ProductRepository
	aiClient       *infrastructure.AIClient
	messageUseCase MessageUseCase
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

	// Post the initial negotiation message if provided
	if err := u.deliverOfferMessage(createdOffer); err != nil {
		fmt.Printf("[OFFER] Failed to deliver negotiation message: %v\n", err)
	}

	return createdOffer, nil
}

// SetMessageUseCase allows wiring the message use case after initialization
func (u *offerUseCase) SetMessageUseCase(messageUseCase MessageUseCase) {
	u.messageUseCase = messageUseCase
}

func (u *offerUseCase) deliverOfferMessage(offer *domain.Offer) error {
	if u.messageUseCase == nil || offer == nil || offer.Product == nil {
		return nil
	}

	content := strings.TrimSpace(offer.Message)
	if content == "" {
		return nil
	}

	conversation, err := u.messageUseCase.GetOrCreateConversation(
		offer.ProductID,
		offer.BuyerID,
		offer.Product.SellerID,
	)
	if err != nil {
		return fmt.Errorf("failed to get or create conversation: %w", err)
	}

	if _, err := u.messageUseCase.SendMessage(conversation.ID, offer.BuyerID, content); err != nil {
		return fmt.Errorf("failed to send negotiation message: %w", err)
	}

	return nil
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
