package infrastructure

import (
	"github.com/google/uuid"
	"github.com/yourusername/ecomate/backend/internal/domain"
	"gorm.io/gorm"
)

type reviewRepository struct {
	db *gorm.DB
}

func NewReviewRepository(db *gorm.DB) domain.ReviewRepository {
	return &reviewRepository{db: db}
}

func (r *reviewRepository) Create(review *domain.Review) error {
	return r.db.Create(review).Error
}

func (r *reviewRepository) FindByProductID(productID uuid.UUID) ([]*domain.Review, error) {
	var reviews []*domain.Review
	if err := r.db.
		Where("product_id = ?", productID).
		Preload("Reviewer").
		Order("created_at DESC").
		Find(&reviews).Error; err != nil {
		return nil, err
	}
	return reviews, nil
}

func (r *reviewRepository) FindByID(id uuid.UUID) (*domain.Review, error) {
	var review domain.Review
	if err := r.db.
		Preload("Reviewer").
		Where("id = ?", id).
		First(&review).Error; err != nil {
		return nil, err
	}
	return &review, nil
}

func (r *reviewRepository) GetAverageRating(productID uuid.UUID) (float64, error) {
	var avg float64
	if err := r.db.Model(&domain.Review{}).
		Where("product_id = ?", productID).
		Select("COALESCE(AVG(rating), 0)").
		Scan(&avg).Error; err != nil {
		return 0, err
	}
	return avg, nil
}
