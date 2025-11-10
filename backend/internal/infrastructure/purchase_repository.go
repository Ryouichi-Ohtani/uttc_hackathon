package infrastructure

import (
	"github.com/google/uuid"
	"github.com/yourusername/ecomate/backend/internal/domain"
	"gorm.io/gorm"
)

type purchaseRepository struct {
	db *gorm.DB
}

func NewPurchaseRepository(db *gorm.DB) domain.PurchaseRepository {
	return &purchaseRepository{db: db}
}

func (r *purchaseRepository) Create(purchase *domain.Purchase) error {
	return r.db.Create(purchase).Error
}

func (r *purchaseRepository) FindByID(id uuid.UUID) (*domain.Purchase, error) {
	var purchase domain.Purchase
	if err := r.db.
		Preload("Product").
		Preload("Product.Images", "is_primary = true").
		Preload("Buyer").
		Preload("Seller").
		Where("id = ?", id).
		First(&purchase).Error; err != nil {
		return nil, err
	}
	return &purchase, nil
}

func (r *purchaseRepository) FindByUser(userID uuid.UUID, role string, page, limit int) ([]*domain.Purchase, *domain.PaginationResponse, error) {
	var purchases []*domain.Purchase
	var total int64

	query := r.db.Model(&domain.Purchase{}).
		Preload("Product").
		Preload("Product.Images", "is_primary = true").
		Preload("Buyer").
		Preload("Seller")

	if role == "buyer" {
		query = query.Where("buyer_id = ?", userID)
	} else if role == "seller" {
		query = query.Where("seller_id = ?", userID)
	} else {
		query = query.Where("buyer_id = ? OR seller_id = ?", userID, userID)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, nil, err
	}

	// Pagination
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}

	offset := (page - 1) * limit
	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&purchases).Error; err != nil {
		return nil, nil, err
	}

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	pagination := &domain.PaginationResponse{
		Page:       page,
		Limit:      limit,
		Total:      int(total),
		TotalPages: totalPages,
	}

	return purchases, pagination, nil
}

func (r *purchaseRepository) List(page, limit int) ([]*domain.Purchase, *domain.PaginationResponse, error) {
	var purchases []*domain.Purchase
	var total int64

	query := r.db.Model(&domain.Purchase{}).
		Preload("Product").
		Preload("Product.Images", "is_primary = true").
		Preload("Buyer").
		Preload("Seller")

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, nil, err
	}

	// Pagination
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}

	offset := (page - 1) * limit
	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&purchases).Error; err != nil {
		return nil, nil, err
	}

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	pagination := &domain.PaginationResponse{
		Page:       page,
		Limit:      limit,
		Total:      int(total),
		TotalPages: totalPages,
	}

	return purchases, pagination, nil
}

func (r *purchaseRepository) UpdateStatus(id uuid.UUID, status domain.PurchaseStatus) error {
	return r.db.Model(&domain.Purchase{}).
		Where("id = ?", id).
		Update("status", status).
		Error
}
