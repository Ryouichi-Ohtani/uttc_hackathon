package infrastructure

import (
	"github.com/google/uuid"
	"github.com/yourusername/ecomate/backend/internal/domain"
	"gorm.io/gorm"
)

type offerRepository struct {
	db *gorm.DB
}

func NewOfferRepository(db *gorm.DB) domain.OfferRepository {
	return &offerRepository{db: db}
}

func (r *offerRepository) Create(offer *domain.Offer) error {
	return r.db.Create(offer).Error
}

func (r *offerRepository) FindByID(id uuid.UUID) (*domain.Offer, error) {
	var offer domain.Offer
	err := r.db.Preload("Product").Preload("Product.Seller").Preload("Buyer").
		First(&offer, "id = ?", id).Error
	return &offer, err
}

func (r *offerRepository) FindByProductID(productID uuid.UUID) ([]*domain.Offer, error) {
	var offers []*domain.Offer
	err := r.db.Preload("Buyer").
		Where("product_id = ?", productID).
		Order("created_at DESC").
		Find(&offers).Error
	return offers, err
}

func (r *offerRepository) FindByBuyerID(buyerID uuid.UUID) ([]*domain.Offer, error) {
	var offers []*domain.Offer
	err := r.db.Preload("Product").Preload("Product.Seller").
		Where("buyer_id = ?", buyerID).
		Order("created_at DESC").
		Find(&offers).Error
	return offers, err
}

func (r *offerRepository) FindBySellerID(sellerID uuid.UUID) ([]*domain.Offer, error) {
	var offers []*domain.Offer
	err := r.db.Preload("Product").Preload("Buyer").
		Joins("JOIN products ON products.id = offers.product_id").
		Where("products.seller_id = ?", sellerID).
		Order("offers.created_at DESC").
		Find(&offers).Error
	return offers, err
}

func (r *offerRepository) Update(offer *domain.Offer) error {
	return r.db.Save(offer).Error
}
