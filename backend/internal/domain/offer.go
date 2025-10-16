package domain

import (
	"time"

	"github.com/google/uuid"
)

type OfferStatus string

const (
	OfferStatusPending  OfferStatus = "pending"
	OfferStatusAccepted OfferStatus = "accepted"
	OfferStatusRejected OfferStatus = "rejected"
	OfferStatusCancelled OfferStatus = "cancelled"
)

type Offer struct {
	ID            uuid.UUID   `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ProductID     uuid.UUID   `json:"product_id" gorm:"type:uuid;not null;index"`
	Product       *Product    `json:"product,omitempty" gorm:"foreignKey:ProductID"`
	BuyerID       uuid.UUID   `json:"buyer_id" gorm:"type:uuid;not null;index"`
	Buyer         *User       `json:"buyer,omitempty" gorm:"foreignKey:BuyerID"`
	OfferPrice    int         `json:"offer_price" gorm:"not null"`
	Message       string      `json:"message"`
	Status        OfferStatus `json:"status" gorm:"default:pending"`
	ResponseMessage string    `json:"response_message"`
	CreatedAt     time.Time   `json:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at"`
	RespondedAt   *time.Time  `json:"responded_at"`
}

type OfferRepository interface {
	Create(offer *Offer) error
	FindByID(id uuid.UUID) (*Offer, error)
	FindByProductID(productID uuid.UUID) ([]*Offer, error)
	FindByBuyerID(buyerID uuid.UUID) ([]*Offer, error)
	FindBySellerID(sellerID uuid.UUID) ([]*Offer, error)
	Update(offer *Offer) error
}

type CreateOfferRequest struct {
	ProductID  string `json:"product_id" binding:"required"`
	OfferPrice int    `json:"offer_price" binding:"required,min=1"`
	Message    string `json:"message"`
}

type RespondOfferRequest struct {
	Accept  bool   `json:"accept" binding:"required"`
	Message string `json:"message"`
}
