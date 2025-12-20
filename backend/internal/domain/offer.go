package domain

import (
	"time"

	"github.com/google/uuid"
)

type OfferStatus string

const (
	OfferStatusPending   OfferStatus = "pending"
	OfferStatusAccepted  OfferStatus = "accepted"
	OfferStatusRejected  OfferStatus = "rejected"
	OfferStatusCancelled OfferStatus = "cancelled"
)

type Offer struct {
	ID              uuid.UUID   `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ProductID       uuid.UUID   `json:"product_id" gorm:"type:uuid;not null;index"`
	Product         *Product    `json:"product,omitempty" gorm:"foreignKey:ProductID"`
	BuyerID         uuid.UUID   `json:"buyer_id" gorm:"type:uuid;not null;index"`
	Buyer           *User       `json:"buyer,omitempty" gorm:"foreignKey:BuyerID"`
	OfferPrice      int         `json:"offer_price" gorm:"not null"`
	Message         string      `json:"message"`
	Status          OfferStatus `json:"status" gorm:"default:pending"`
	ResponseMessage string      `json:"response_message"`
	CreatedAt       time.Time   `json:"created_at"`
	UpdatedAt       time.Time   `json:"updated_at"`
	RespondedAt     *time.Time  `json:"responded_at"`
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
	Accept  *bool  `json:"accept" binding:"required"`
	Message string `json:"message"`
}

// MarketPriceAnalysis represents AI-powered market price analysis
type MarketPriceAnalysis struct {
	ProductTitle      string             `json:"product_title"`
	Category          string             `json:"category"`
	Condition         string             `json:"condition"`
	ListingPrice      int                `json:"listing_price"`
	RecommendedPrice  int                `json:"recommended_price"`
	MinPrice          int                `json:"min_price"`
	MaxPrice          int                `json:"max_price"`
	MarketDataSources []MarketDataSource `json:"market_data_sources"`
	Analysis          string             `json:"analysis"`
	ConfidenceLevel   string             `json:"confidence_level"` // "high", "medium", "low"
	AnalyzedAt        time.Time          `json:"analyzed_at"`
}

// MarketDataSource represents a single market data point from external sources
type MarketDataSource struct {
	Platform  string `json:"platform"` // "メルカリ", "ヤフオク", "Amazon", "楽天"
	Price     int    `json:"price"`
	Condition string `json:"condition"`
	URL       string `json:"url,omitempty"`
}
