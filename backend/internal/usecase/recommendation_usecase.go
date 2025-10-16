package usecase

import (
	"github.com/google/uuid"
	"github.com/yourusername/ecomate/backend/internal/domain"
)

type RecommendationUseCase interface {
	GetRecommendations(userID uuid.UUID, limit int) ([]*domain.Product, error)
}

type recommendationUseCase struct {
	productRepo  domain.ProductRepository
	purchaseRepo domain.PurchaseRepository
}

func NewRecommendationUseCase(
	productRepo domain.ProductRepository,
	purchaseRepo domain.PurchaseRepository,
) RecommendationUseCase {
	return &recommendationUseCase{
		productRepo:  productRepo,
		purchaseRepo: purchaseRepo,
	}
}

// GetRecommendations provides simple recommendations based on user's purchase history
// In production, this would use collaborative filtering or ML models
func (u *recommendationUseCase) GetRecommendations(userID uuid.UUID, limit int) ([]*domain.Product, error) {
	// Simple recommendation: get products in categories the user has purchased from
	// excluding products they've already purchased or listed

	// For now, return latest active products with high CO2 impact (eco-friendly)
	filters := &domain.ProductFilters{
		Page:  1,
		Limit: limit,
		Sort:  "eco_impact_desc", // Sort by CO2 impact
	}

	products, _, err := u.productRepo.List(filters)
	if err != nil {
		return nil, err
	}

	// Filter out user's own products
	recommendations := make([]*domain.Product, 0)
	for _, product := range products {
		if product.SellerID != userID && product.Status == domain.StatusActive {
			recommendations = append(recommendations, product)
			if len(recommendations) >= limit {
				break
			}
		}
	}

	return recommendations, nil
}
