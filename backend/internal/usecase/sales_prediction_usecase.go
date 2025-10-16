package usecase

import (
	"math"
	

	"github.com/google/uuid"
	"github.com/yourusername/ecomate/backend/internal/domain"
)

type SalesPredictionUseCase interface {
	PredictProductSalePrice(productID uuid.UUID) (*SalesPrediction, error)
	PredictSellerRevenue(sellerID uuid.UUID, days int) (*RevenuePrediction, error)
	GetMarketTrends(category string, days int) (*MarketTrend, error)
}

type SalesPrediction struct {
	ProductID           uuid.UUID `json:"product_id"`
	CurrentPrice        int       `json:"current_price"`
	PredictedPrice      int       `json:"predicted_price"`
	ConfidenceLevel     float64   `json:"confidence_level"`
	SaleProbability     float64   `json:"sale_probability"`
	RecommendedPrice    int       `json:"recommended_price"`
	EstimatedDaysToSell int       `json:"estimated_days_to_sell"`
}

type RevenuePrediction struct {
	SellerID             uuid.UUID `json:"seller_id"`
	PredictionPeriodDays int       `json:"prediction_period_days"`
	PredictedRevenue     int       `json:"predicted_revenue"`
	ExpectedSales        int       `json:"expected_sales"`
	ConfidenceLevel      float64   `json:"confidence_level"`
}

type MarketTrend struct {
	Category          string    `json:"category"`
	AveragePrice      int       `json:"average_price"`
	MedianPrice       int       `json:"median_price"`
	PriceGrowthRate   float64   `json:"price_growth_rate"`
	DemandScore       float64   `json:"demand_score"`
	CompetitionLevel  string    `json:"competition_level"`
	TrendDirection    string    `json:"trend_direction"` // "up", "down", "stable"
}

type salesPredictionUseCase struct {
	productRepo  domain.ProductRepository
	purchaseRepo domain.PurchaseRepository
	userRepo     domain.UserRepository
}

func NewSalesPredictionUseCase(
	productRepo domain.ProductRepository,
	purchaseRepo domain.PurchaseRepository,
	userRepo domain.UserRepository,
) SalesPredictionUseCase {
	return &salesPredictionUseCase{
		productRepo:  productRepo,
		purchaseRepo: purchaseRepo,
		userRepo:     userRepo,
	}
}

func (u *salesPredictionUseCase) PredictProductSalePrice(productID uuid.UUID) (*SalesPrediction, error) {
	product, err := u.productRepo.FindByID(productID)
	if err != nil {
		return nil, err
	}

	// Simple prediction based on category and condition
	// In production, this would use ML models
	filters := &domain.ProductFilters{
		Category: product.Category,
		Limit:    100,
	}

	similarProducts, _, err := u.productRepo.List(filters)
	if err != nil {
		return nil, err
	}

	// Calculate average price of similar products
	var totalPrice int
	var soldCount int
	for _, p := range similarProducts {
		if p.Status == domain.StatusSold {
			totalPrice += p.Price
			soldCount++
		}
	}

	var predictedPrice int
	var confidenceLevel float64

	if soldCount > 0 {
		avgPrice := totalPrice / soldCount
		predictedPrice = avgPrice
		confidenceLevel = math.Min(float64(soldCount)/20.0, 0.95)
	} else {
		predictedPrice = product.Price
		confidenceLevel = 0.3
	}

	// Adjust for condition
	conditionMultipliers := map[domain.ProductCondition]float64{
		domain.ConditionNew:     1.0,
		domain.ConditionLikeNew: 0.85,
		domain.ConditionGood:    0.7,
		domain.ConditionFair:    0.5,
		domain.ConditionPoor:    0.3,
	}

	multiplier := conditionMultipliers[product.Condition]
	predictedPrice = int(float64(predictedPrice) * multiplier)

	// Calculate sale probability
	priceDiff := float64(product.Price - predictedPrice)
	saleProbability := 0.5 + (priceDiff / float64(predictedPrice) * 0.3)
	saleProbability = math.Max(0.1, math.Min(0.95, saleProbability))

	// Recommended price (slightly below predicted)
	recommendedPrice := int(float64(predictedPrice) * 0.95)

	// Estimate days to sell
	daysToSell := int(30 / saleProbability)
	if daysToSell > 90 {
		daysToSell = 90
	}

	return &SalesPrediction{
		ProductID:           productID,
		CurrentPrice:        product.Price,
		PredictedPrice:      predictedPrice,
		ConfidenceLevel:     confidenceLevel,
		SaleProbability:     saleProbability,
		RecommendedPrice:    recommendedPrice,
		EstimatedDaysToSell: daysToSell,
	}, nil
}

func (u *salesPredictionUseCase) PredictSellerRevenue(sellerID uuid.UUID, days int) (*RevenuePrediction, error) {
	// Get seller's active products
	filters := &domain.ProductFilters{
		Limit: 1000,
	}

	allProducts, _, err := u.productRepo.List(filters)
	if err != nil {
		return nil, err
	}

	var sellerProducts []*domain.Product
	for _, p := range allProducts {
		if p.SellerID == sellerID && p.Status == domain.StatusActive {
			sellerProducts = append(sellerProducts, p)
		}
	}

	// Calculate expected sales
	var totalPredictedRevenue int
	var expectedSales int

	for _, product := range sellerProducts {
		prediction, err := u.PredictProductSalePrice(product.ID)
		if err != nil {
			continue
		}

		// If estimated days to sell is within prediction period
		if prediction.EstimatedDaysToSell <= days {
			totalPredictedRevenue += prediction.PredictedPrice
			expectedSales++
		}
	}

	confidenceLevel := 0.6
	if len(sellerProducts) > 10 {
		confidenceLevel = 0.75
	}

	return &RevenuePrediction{
		SellerID:             sellerID,
		PredictionPeriodDays: days,
		PredictedRevenue:     totalPredictedRevenue,
		ExpectedSales:        expectedSales,
		ConfidenceLevel:      confidenceLevel,
	}, nil
}

func (u *salesPredictionUseCase) GetMarketTrends(category string, days int) (*MarketTrend, error) {
	filters := &domain.ProductFilters{
		Category: category,
		Limit:    500,
	}

	products, _, err := u.productRepo.List(filters)
	if err != nil {
		return nil, err
	}

	// Calculate statistics
	var prices []int
	var soldCount int
	var activeCount int

	for _, p := range products {
		if p.Status == domain.StatusSold {
			prices = append(prices, p.Price)
			soldCount++
		} else if p.Status == domain.StatusActive {
			activeCount++
		}
	}

	if len(prices) == 0 {
		return &MarketTrend{
			Category:         category,
			AveragePrice:     0,
			MedianPrice:      0,
			PriceGrowthRate:  0,
			DemandScore:      0.5,
			CompetitionLevel: "low",
			TrendDirection:   "stable",
		}, nil
	}

	// Calculate average
	var sum int
	for _, price := range prices {
		sum += price
	}
	avgPrice := sum / len(prices)

	// Calculate median (simplified)
	medianPrice := avgPrice

	// Demand score based on sales rate
	demandScore := float64(soldCount) / float64(soldCount+activeCount)

	// Competition level
	competitionLevel := "low"
	if activeCount > 50 {
		competitionLevel = "high"
	} else if activeCount > 20 {
		competitionLevel = "medium"
	}

	// Trend direction (simplified)
	trendDirection := "stable"
	if demandScore > 0.6 {
		trendDirection = "up"
	} else if demandScore < 0.3 {
		trendDirection = "down"
	}

	return &MarketTrend{
		Category:         category,
		AveragePrice:     avgPrice,
		MedianPrice:      medianPrice,
		PriceGrowthRate:  0.05, // Simplified
		DemandScore:      demandScore,
		CompetitionLevel: competitionLevel,
		TrendDirection:   trendDirection,
	}, nil
}
