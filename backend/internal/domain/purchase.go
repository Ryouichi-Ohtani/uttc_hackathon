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
	Status          PurchaseStatus `json:"status" gorm:"default:pending"`
	PaymentMethod   string         `json:"payment_method"`
	ShippingAddress string         `json:"shipping_address"`
	// Delivery information
	DeliveryDate           *time.Time `json:"delivery_date"`
	DeliveryTimeSlot       string     `json:"delivery_time_slot"` // "morning", "afternoon", "evening", "anytime"
	UseRegisteredAddress   bool       `json:"use_registered_address" gorm:"default:false"`
	RecipientName          string     `json:"recipient_name"`
	RecipientPhoneNumber   string     `json:"recipient_phone_number"`
	RecipientPostalCode    string     `json:"recipient_postal_code"`
	RecipientPrefecture    string     `json:"recipient_prefecture"`
	RecipientCity          string     `json:"recipient_city"`
	RecipientAddressLine1  string     `json:"recipient_address_line1"`
	RecipientAddressLine2  string     `json:"recipient_address_line2"`
	CompletedAt            *time.Time `json:"completed_at"`
	CreatedAt              time.Time  `json:"created_at"`
}

type PurchaseRepository interface {
	Create(purchase *Purchase) error
	FindByID(id uuid.UUID) (*Purchase, error)
	FindByUser(userID uuid.UUID, role string, page, limit int) ([]*Purchase, *PaginationResponse, error)
	UpdateStatus(id uuid.UUID, status PurchaseStatus) error
}

type CreatePurchaseRequest struct {
	ProductID              uuid.UUID  `json:"product_id" binding:"required"`
	ShippingAddress        string     `json:"shipping_address" binding:"required"`
	PaymentMethod          string     `json:"payment_method" binding:"required"`
	// Delivery information
	DeliveryDate           *time.Time `json:"delivery_date"`
	DeliveryTimeSlot       string     `json:"delivery_time_slot"`
	UseRegisteredAddress   bool       `json:"use_registered_address"`
	RecipientName          string     `json:"recipient_name"`
	RecipientPhoneNumber   string     `json:"recipient_phone_number"`
	RecipientPostalCode    string     `json:"recipient_postal_code"`
	RecipientPrefecture    string     `json:"recipient_prefecture"`
	RecipientCity          string     `json:"recipient_city"`
	RecipientAddressLine1  string     `json:"recipient_address_line1"`
	RecipientAddressLine2  string     `json:"recipient_address_line2"`
}

// ShippingLabel represents the shipping label information for sellers
type ShippingLabel struct {
	ID                    uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	PurchaseID            uuid.UUID  `json:"purchase_id" gorm:"type:uuid;not null;index;unique"`
	Purchase              *Purchase  `json:"purchase,omitempty" gorm:"foreignKey:PurchaseID"`
	// Sender (Seller) information
	SenderName            string     `json:"sender_name"`
	SenderPostalCode      string     `json:"sender_postal_code"`
	SenderPrefecture      string     `json:"sender_prefecture"`
	SenderCity            string     `json:"sender_city"`
	SenderAddressLine1    string     `json:"sender_address_line1"`
	SenderAddressLine2    string     `json:"sender_address_line2"`
	SenderPhoneNumber     string     `json:"sender_phone_number"`
	// Recipient (Buyer) information
	RecipientName         string     `json:"recipient_name"`
	RecipientPostalCode   string     `json:"recipient_postal_code"`
	RecipientPrefecture   string     `json:"recipient_prefecture"`
	RecipientCity         string     `json:"recipient_city"`
	RecipientAddressLine1 string     `json:"recipient_address_line1"`
	RecipientAddressLine2 string     `json:"recipient_address_line2"`
	RecipientPhoneNumber  string     `json:"recipient_phone_number"`
	// Delivery details
	DeliveryDate          *time.Time `json:"delivery_date"`
	DeliveryTimeSlot      string     `json:"delivery_time_slot"`
	ProductName           string     `json:"product_name"`
	PackageSize           string     `json:"package_size"` // "60", "80", "100", "120", "140", "160"
	Weight                float64    `json:"weight"` // in kg
	TrackingNumber        string     `json:"tracking_number"`
	Carrier               string     `json:"carrier"` // "yamato", "sagawa", "japan_post"
	CreatedAt             time.Time  `json:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at"`
}

type ShippingLabelRepository interface {
	Create(label *ShippingLabel) error
	FindByPurchaseID(purchaseID uuid.UUID) (*ShippingLabel, error)
	Update(label *ShippingLabel) error
}
