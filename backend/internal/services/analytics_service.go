package services

import (
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/yourusername/ecomate/backend/internal/domain"
)

type AnalyticsService struct {
	productRepo domain.ProductRepository
	userRepo    domain.UserRepository
}

func NewAnalyticsService(productRepo domain.ProductRepository, userRepo domain.UserRepository) *AnalyticsService {
	return &AnalyticsService{
		productRepo: productRepo,
		userRepo:    userRepo,
	}
}

// UserBehaviorAnalytics represents user behavior metrics
type UserBehaviorAnalytics struct {
	UserID              uuid.UUID            `json:"user_id"`
	ViewedProducts      int                  `json:"viewed_products"`
	FavoriteProducts    int                  `json:"favorite_products"`
	PurchasedProducts   int                  `json:"purchased_products"`
	AverageSessionTime  float64              `json:"average_session_time_minutes"`
	CategoryPreferences map[string]int       `json:"category_preferences"`
	PriceRange          PriceRangePreference `json:"price_range"`
	ActivityPattern     []HourlyActivity     `json:"activity_pattern"`
	CO2Impact           float64              `json:"co2_impact"`
	RecommendedProducts []uuid.UUID          `json:"recommended_products"`
}

type PriceRangePreference struct {
	Min     int     `json:"min"`
	Max     int     `json:"max"`
	Average float64 `json:"average"`
}

type HourlyActivity struct {
	Hour     int `json:"hour"`
	Activity int `json:"activity"`
}

// SalesPrediction represents sales forecast
type SalesPrediction struct {
	ProductID        uuid.UUID              `json:"product_id"`
	PredictedPrice   int                    `json:"predicted_price"`
	ConfidenceScore  float64                `json:"confidence_score"`
	OptimalListPrice int                    `json:"optimal_list_price"`
	DemandScore      float64                `json:"demand_score"`
	TimeToSell       int                    `json:"estimated_days_to_sell"`
	SimilarProducts  []SimilarProductMetric `json:"similar_products"`
	TrendAnalysis    TrendData              `json:"trend_analysis"`
}

type SimilarProductMetric struct {
	ProductID  uuid.UUID `json:"product_id"`
	SoldPrice  int       `json:"sold_price"`
	DaysToSell int       `json:"days_to_sell"`
	Similarity float64   `json:"similarity_score"`
}

type TrendData struct {
	CategoryTrend string  `json:"category_trend"` // "rising", "stable", "declining"
	SeasonalBoost float64 `json:"seasonal_boost"`
	MarketDemand  float64 `json:"market_demand"`
}

// AnalyzeUserBehavior analyzes user behavior patterns
func (s *AnalyticsService) AnalyzeUserBehavior(userID uuid.UUID) (*UserBehaviorAnalytics, error) {
	// Simplified implementation - in production, this would query actual user activity logs

	categoryPreferences := map[string]int{
		"electronics": 15,
		"clothing":    10,
		"sports":      5,
	}

	activityPattern := make([]HourlyActivity, 24)
	for i := 0; i < 24; i++ {
		// Simulate activity pattern (peak hours: 12-14, 19-22)
		activity := 0
		if i >= 12 && i <= 14 {
			activity = 8 + (i-12)*2
		} else if i >= 19 && i <= 22 {
			activity = 10 + (i-19)*3
		} else {
			activity = 2
		}
		activityPattern[i] = HourlyActivity{Hour: i, Activity: activity}
	}

	analytics := &UserBehaviorAnalytics{
		UserID:              userID,
		ViewedProducts:      45,
		FavoriteProducts:    12,
		PurchasedProducts:   3,
		AverageSessionTime:  15.5,
		CategoryPreferences: categoryPreferences,
		PriceRange: PriceRangePreference{
			Min:     5000,
			Max:     80000,
			Average: 25000,
		},
		ActivityPattern:     activityPattern,
		CO2Impact:           125.5,
		RecommendedProducts: []uuid.UUID{},
	}

	return analytics, nil
}

// PredictSales predicts sales metrics for a product
func (s *AnalyticsService) PredictSales(productID uuid.UUID) (*SalesPrediction, error) {
	product, err := s.productRepo.FindByID(productID)
	if err != nil {
		return nil, err
	}

	// Advanced price prediction using multiple factors
	basePrice := float64(product.Price)

	// Factor 1: Category trend
	categoryTrend := s.getCategoryTrend(product.Category)
	trendMultiplier := 1.0
	if categoryTrend == "rising" {
		trendMultiplier = 1.15
	} else if categoryTrend == "declining" {
		trendMultiplier = 0.85
	}

	// Factor 2: Seasonal adjustment
	seasonalBoost := s.getSeasonalBoost(product.Category, time.Now())

	// Factor 3: Condition factor
	conditionMultiplier := 1.0
	switch product.Condition {
	case "excellent":
		conditionMultiplier = 1.0
	case "good":
		conditionMultiplier = 0.85
	case "fair":
		conditionMultiplier = 0.70
	}

	// Factor 4: Demand score (based on view count, favorite count)
	demandScore := s.calculateDemandScore(product.ViewCount, product.FavoriteCount)

	// Calculate predicted price
	predictedPrice := int(basePrice * trendMultiplier * (1 + seasonalBoost) * conditionMultiplier * (1 + demandScore*0.1))

	// Optimal list price (slightly higher than predicted)
	optimalListPrice := int(float64(predictedPrice) * 1.1)

	// Estimate time to sell
	timeToSell := s.estimateTimeToSell(demandScore, float64(product.Price), float64(predictedPrice))

	// Find similar products (mock data)
	similarProducts := []SimilarProductMetric{
		{ProductID: uuid.New(), SoldPrice: predictedPrice - 5000, DaysToSell: 7, Similarity: 0.92},
		{ProductID: uuid.New(), SoldPrice: predictedPrice + 3000, DaysToSell: 12, Similarity: 0.85},
		{ProductID: uuid.New(), SoldPrice: predictedPrice - 2000, DaysToSell: 5, Similarity: 0.88},
	}

	prediction := &SalesPrediction{
		ProductID:        productID,
		PredictedPrice:   predictedPrice,
		ConfidenceScore:  0.78 + demandScore*0.15,
		OptimalListPrice: optimalListPrice,
		DemandScore:      demandScore,
		TimeToSell:       timeToSell,
		SimilarProducts:  similarProducts,
		TrendAnalysis: TrendData{
			CategoryTrend: categoryTrend,
			SeasonalBoost: seasonalBoost,
			MarketDemand:  demandScore,
		},
	}

	return prediction, nil
}

func (s *AnalyticsService) getCategoryTrend(category string) string {
	// Simplified - in production, this would analyze historical data
	trends := map[string]string{
		"electronics": "stable",
		"clothing":    "rising",
		"furniture":   "declining",
		"sports":      "rising",
		"books":       "stable",
	}

	if trend, ok := trends[category]; ok {
		return trend
	}
	return "stable"
}

func (s *AnalyticsService) getSeasonalBoost(category string, date time.Time) float64 {
	month := date.Month()

	seasonalFactors := map[string]map[time.Month]float64{
		"electronics": {
			time.November: 0.15,
			time.December: 0.25,
			time.January:  -0.10,
		},
		"clothing": {
			time.March:     0.10,
			time.April:     0.15,
			time.September: 0.12,
			time.October:   0.18,
		},
		"sports": {
			time.April: 0.20,
			time.May:   0.25,
			time.June:  0.15,
		},
	}

	if factors, ok := seasonalFactors[category]; ok {
		if boost, ok := factors[month]; ok {
			return boost
		}
	}
	return 0.0
}

func (s *AnalyticsService) calculateDemandScore(viewCount, favoriteCount int) float64 {
	// Normalize scores
	viewScore := math.Min(float64(viewCount)/100.0, 1.0)
	favoriteScore := math.Min(float64(favoriteCount)/20.0, 1.0)

	// Weighted average
	return viewScore*0.6 + favoriteScore*0.4
}

func (s *AnalyticsService) estimateTimeToSell(demandScore, currentPrice, predictedPrice float64) int {
	baseDays := 14.0

	// Adjust based on demand
	demandFactor := 1.0 - demandScore*0.5

	// Adjust based on pricing
	priceFactor := 1.0
	if currentPrice > predictedPrice {
		priceFactor = 1.0 + ((currentPrice - predictedPrice) / predictedPrice)
	} else {
		priceFactor = 1.0 - ((predictedPrice - currentPrice) / predictedPrice * 0.5)
	}

	estimatedDays := int(baseDays * demandFactor * priceFactor)

	if estimatedDays < 1 {
		return 1
	}
	if estimatedDays > 60 {
		return 60
	}

	return estimatedDays
}
