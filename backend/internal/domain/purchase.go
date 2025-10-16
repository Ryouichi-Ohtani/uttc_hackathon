package domain

import (
	"time"

	"github.com/google/uuid"
)

type PurchaseStatus string

const (
	PurchaseStatusPending   PurchaseStatus = "pending"
	PurchaseStatusCompleted PurchaseStatus = "completed"
	PurchaseStatusCancelled PurchaseStatus = "cancelled"
)

type Purchase struct {
	ID              uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ProductID       uuid.UUID      `json:"product_id" gorm:"type:uuid;not null;index"`
	Product         *Product       `json:"product,omitempty" gorm:"foreignKey:ProductID"`
	BuyerID         uuid.UUID      `json:"buyer_id" gorm:"type:uuid;not null;index"`
	Buyer           *User          `json:"buyer,omitempty" gorm:"foreignKey:BuyerID"`
	SellerID        uuid.UUID      `json:"seller_id" gorm:"type:uuid;not null;index"`
	Seller          *User          `json:"seller,omitempty" gorm:"foreignKey:SellerID"`
	Price           int            `json:"price" gorm:"not null"`
	CO2SavedKg      float64        `json:"co2_saved_kg" gorm:"type:decimal(10,2);not null"`
	Status          PurchaseStatus `json:"status" gorm:"default:pending"`
	PaymentMethod   string         `json:"payment_method"`
	ShippingAddress string         `json:"shipping_address"`
	CompletedAt     *time.Time     `json:"completed_at"`
	CreatedAt       time.Time      `json:"created_at"`
}

type PurchaseRepository interface {
	Create(purchase *Purchase) error
	FindByID(id uuid.UUID) (*Purchase, error)
	FindByUser(userID uuid.UUID, role string, page, limit int) ([]*Purchase, *PaginationResponse, error)
	UpdateStatus(id uuid.UUID, status PurchaseStatus) error
}

type CreatePurchaseRequest struct {
	ProductID       uuid.UUID `json:"product_id" binding:"required"`
	ShippingAddress string    `json:"shipping_address" binding:"required"`
	PaymentMethod   string    `json:"payment_method" binding:"required"`
}
