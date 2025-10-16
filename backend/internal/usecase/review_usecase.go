package usecase

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/yourusername/ecomate/backend/internal/domain"
)

type ReviewUseCase interface {
	Create(reviewerID, productID, purchaseID uuid.UUID, rating int, comment string) (*domain.Review, error)
	GetProductReviews(productID uuid.UUID) ([]*domain.Review, float64, error)
}

type reviewUseCase struct {
	reviewRepo   domain.ReviewRepository
	purchaseRepo domain.PurchaseRepository
}

func NewReviewUseCase(
	reviewRepo domain.ReviewRepository,
	purchaseRepo domain.PurchaseRepository,
) ReviewUseCase {
	return &reviewUseCase{
		reviewRepo:   reviewRepo,
		purchaseRepo: purchaseRepo,
	}
}

func (u *reviewUseCase) Create(reviewerID, productID, purchaseID uuid.UUID, rating int, comment string) (*domain.Review, error) {
	// Validate purchase exists and belongs to reviewer
	purchase, err := u.purchaseRepo.FindByID(purchaseID)
	if err != nil {
		return nil, errors.New("purchase not found")
	}

	if purchase.BuyerID != reviewerID {
		return nil, errors.New("only buyer can review")
	}

	if purchase.ProductID != productID {
		return nil, errors.New("product mismatch")
	}

	if purchase.Status != domain.PurchaseStatusCompleted {
		return nil, errors.New("can only review completed purchases")
	}

	// Validate rating
	if rating < 1 || rating > 5 {
		return nil, errors.New("rating must be between 1 and 5")
	}

	review := &domain.Review{
		ProductID:  productID,
		PurchaseID: purchaseID,
		ReviewerID: reviewerID,
		Rating:     rating,
		Comment:    comment,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := u.reviewRepo.Create(review); err != nil {
		return nil, err
	}

	return u.reviewRepo.FindByID(review.ID)
}

func (u *reviewUseCase) GetProductReviews(productID uuid.UUID) ([]*domain.Review, float64, error) {
	reviews, err := u.reviewRepo.FindByProductID(productID)
	if err != nil {
		return nil, 0, err
	}

	avgRating, err := u.reviewRepo.GetAverageRating(productID)
	if err != nil {
		return nil, 0, err
	}

	return reviews, avgRating, nil
}
