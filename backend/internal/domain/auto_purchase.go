package domain

import (
	"time"

	"github.com/google/uuid"
)

type AutoPurchaseWatchStatus string

const (
	AutoPurchaseWatchStatusActive    AutoPurchaseWatchStatus = "active"
	AutoPurchaseWatchStatusExecuted  AutoPurchaseWatchStatus = "executed"
	AutoPurchaseWatchStatusCancelled AutoPurchaseWatchStatus = "cancelled"
	AutoPurchaseWatchStatusExpired   AutoPurchaseWatchStatus = "expired"
)

// AutoPurchaseWatch represents a user's automatic purchase watch on a product
type AutoPurchaseWatch struct {
	ID                uuid.UUID               `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID            uuid.UUID               `json:"user_id" gorm:"type:uuid;not null;index"`
	User              *User                   `json:"user,omitempty" gorm:"foreignKey:UserID"`
	ProductID         uuid.UUID               `json:"product_id" gorm:"type:uuid;not null;index"`
	Product           *Product                `json:"product,omitempty" gorm:"foreignKey:ProductID"`
	MaxPrice          int                     `json:"max_price" gorm:"not null"` // Maximum price willing to pay
	Status            AutoPurchaseWatchStatus `json:"status" gorm:"default:'active'"`
	PaymentAuthorized bool                    `json:"payment_authorized" gorm:"default:false"`
	PaymentMethodID   string                  `json:"payment_method_id"`  // ID of the authorized payment method
	PaymentAuthToken  string                  `json:"payment_auth_token"` // Token for payment authorization
	// Delivery preferences
	UseRegisteredAddress  bool   `json:"use_registered_address" gorm:"default:true"`
	RecipientName         string `json:"recipient_name"`
	RecipientPhoneNumber  string `json:"recipient_phone_number"`
	RecipientPostalCode   string `json:"recipient_postal_code"`
	RecipientPrefecture   string `json:"recipient_prefecture"`
	RecipientCity         string `json:"recipient_city"`
	RecipientAddressLine1 string `json:"recipient_address_line1"`
	RecipientAddressLine2 string `json:"recipient_address_line2"`
	ShippingAddress       string `json:"shipping_address"`
	DeliveryTimeSlot      string `json:"delivery_time_slot"`
	// Tracking
	LastCheckedAt    *time.Time `json:"last_checked_at"`
	ExecutedAt       *time.Time `json:"executed_at"`
	PurchaseID       *uuid.UUID `json:"purchase_id" gorm:"type:uuid"` // ID of the executed purchase
	Purchase         *Purchase  `json:"purchase,omitempty" gorm:"foreignKey:PurchaseID"`
	ExpiresAt        time.Time  `json:"expires_at"` // Auto-cancel after this date
	NotificationSent bool       `json:"notification_sent" gorm:"default:false"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

// AutoPurchaseLog represents a log entry for auto-purchase attempts
type AutoPurchaseLog struct {
	ID              uuid.UUID          `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	WatchID         uuid.UUID          `json:"watch_id" gorm:"type:uuid;not null;index"`
	Watch           *AutoPurchaseWatch `json:"watch,omitempty" gorm:"foreignKey:WatchID"`
	Action          string             `json:"action"` // "price_check", "purchase_attempt", "purchase_success", "purchase_failed"
	CurrentPrice    int                `json:"current_price"`
	Success         bool               `json:"success"`
	ErrorMessage    string             `json:"error_message"`
	ExecutionTimeMs int                `json:"execution_time_ms"`
	CreatedAt       time.Time          `json:"created_at"`
}

type AutoPurchaseWatchRepository interface {
	Create(watch *AutoPurchaseWatch) error
	FindByID(id uuid.UUID) (*AutoPurchaseWatch, error)
	FindByUserID(userID uuid.UUID) ([]*AutoPurchaseWatch, error)
	FindByProductID(productID uuid.UUID) ([]*AutoPurchaseWatch, error)
	FindActiveWatches() ([]*AutoPurchaseWatch, error)
	Update(watch *AutoPurchaseWatch) error
	Delete(id uuid.UUID) error
}

type AutoPurchaseLogRepository interface {
	Create(log *AutoPurchaseLog) error
	FindByWatchID(watchID uuid.UUID) ([]*AutoPurchaseLog, error)
}

type CreateAutoPurchaseWatchRequest struct {
	ProductID             uuid.UUID `json:"product_id" binding:"required"`
	MaxPrice              int       `json:"max_price" binding:"required,gt=0"`
	PaymentMethodID       string    `json:"payment_method_id" binding:"required"`
	PaymentAuthToken      string    `json:"payment_auth_token" binding:"required"`
	UseRegisteredAddress  bool      `json:"use_registered_address"`
	RecipientName         string    `json:"recipient_name"`
	RecipientPhoneNumber  string    `json:"recipient_phone_number"`
	RecipientPostalCode   string    `json:"recipient_postal_code"`
	RecipientPrefecture   string    `json:"recipient_prefecture"`
	RecipientCity         string    `json:"recipient_city"`
	RecipientAddressLine1 string    `json:"recipient_address_line1"`
	RecipientAddressLine2 string    `json:"recipient_address_line2"`
	ShippingAddress       string    `json:"shipping_address"`
	DeliveryTimeSlot      string    `json:"delivery_time_slot"`
	ExpiresInDays         int       `json:"expires_in_days"` // Default 30 days
}

type PaymentAuthorizationRequest struct {
	CardNumber     string `json:"card_number" binding:"required"`
	ExpiryMonth    int    `json:"expiry_month" binding:"required,min=1,max=12"`
	ExpiryYear     int    `json:"expiry_year" binding:"required"`
	CVV            string `json:"cvv" binding:"required,len=3"`
	CardholderName string `json:"cardholder_name" binding:"required"`
	Amount         int    `json:"amount" binding:"required,gt=0"` // Pre-authorization amount
}

type PaymentAuthorizationResponse struct {
	Authorized      bool      `json:"authorized"`
	PaymentMethodID string    `json:"payment_method_id"`
	AuthToken       string    `json:"auth_token"`
	ExpiresAt       time.Time `json:"expires_at"`
	Message         string    `json:"message"`
}
