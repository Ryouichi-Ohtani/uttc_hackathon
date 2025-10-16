package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ChatHistory represents a chatbot conversation history
type ChatHistory struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	Message   string         `gorm:"type:text;not null" json:"message"`
	Response  string         `gorm:"type:text;not null" json:"response"`
	Context   string         `gorm:"type:text" json:"context"` // Product ID or other context
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type ChatHistoryRepository interface {
	Create(history *ChatHistory) error
	GetByUserID(userID uuid.UUID, limit int) ([]*ChatHistory, error)
	Delete(id uuid.UUID) error
}

// CO2Goal represents a user's CO2 reduction goal
type CO2Goal struct {
	ID             uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID         uuid.UUID      `gorm:"type:uuid;not null;unique;index" json:"user_id"`
	TargetKG       float64        `gorm:"not null" json:"target_kg"`
	CurrentKG      float64        `gorm:"default:0" json:"current_kg"`
	TargetDate     time.Time      `gorm:"not null" json:"target_date"`
	StartDate      time.Time      `gorm:"not null" json:"start_date"`
	Status         string         `gorm:"not null;default:'active'" json:"status"` // active, completed, expired
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

type CO2GoalRepository interface {
	Create(goal *CO2Goal) error
	GetByUserID(userID uuid.UUID) (*CO2Goal, error)
	Update(goal *CO2Goal) error
	UpdateProgress(userID uuid.UUID, additionalKG float64) error
}

// ShippingTracking represents shipping information for a purchase
type ShippingTracking struct {
	ID               uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	PurchaseID       uuid.UUID      `gorm:"type:uuid;not null;unique;index" json:"purchase_id"`
	TrackingNumber   string         `gorm:"not null" json:"tracking_number"`
	Carrier          string         `gorm:"not null" json:"carrier"` // yamato, sagawa, yupack, etc.
	Status           string         `gorm:"not null;default:'pending'" json:"status"` // pending, shipped, in_transit, delivered
	ShippedAt        *time.Time     `json:"shipped_at"`
	DeliveredAt      *time.Time     `json:"delivered_at"`
	EstimatedArrival *time.Time     `json:"estimated_arrival"`
	ShippingMethod   string         `json:"shipping_method"` // standard, eco, express
	CO2Saved         float64        `gorm:"default:0" json:"co2_saved"` // CO2 saved by eco shipping
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Purchase         *Purchase      `gorm:"foreignKey:PurchaseID" json:"purchase,omitempty"`
}

type ShippingTrackingRepository interface {
	Create(tracking *ShippingTracking) error
	GetByPurchaseID(purchaseID uuid.UUID) (*ShippingTracking, error)
	Update(tracking *ShippingTracking) error
	UpdateStatus(id uuid.UUID, status string) error
}
