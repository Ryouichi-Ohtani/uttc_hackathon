package infrastructure

import (
	"github.com/google/uuid"
	"github.com/yourusername/ecomate/backend/internal/domain"
	"gorm.io/gorm"
)

type ShippingLabelRepositoryImpl struct {
	db *gorm.DB
}

func NewShippingLabelRepository(db *gorm.DB) domain.ShippingLabelRepository {
	return &ShippingLabelRepositoryImpl{db: db}
}

func (r *ShippingLabelRepositoryImpl) Create(label *domain.ShippingLabel) error {
	return r.db.Create(label).Error
}

func (r *ShippingLabelRepositoryImpl) FindByPurchaseID(purchaseID uuid.UUID) (*domain.ShippingLabel, error) {
	var label domain.ShippingLabel
	err := r.db.Preload("Purchase").
		Preload("Purchase.Product").
		Preload("Purchase.Buyer").
		Preload("Purchase.Seller").
		Where("purchase_id = ?", purchaseID).
		First(&label).Error
	if err != nil {
		return nil, err
	}
	return &label, nil
}

func (r *ShippingLabelRepositoryImpl) Update(label *domain.ShippingLabel) error {
	return r.db.Save(label).Error
}
