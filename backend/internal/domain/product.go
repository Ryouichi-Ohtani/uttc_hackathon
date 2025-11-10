package domain

import (
	"time"

	"github.com/google/uuid"
)

type ProductStatus string
type ProductCondition string

const (
	StatusDraft    ProductStatus = "draft"
	StatusActive   ProductStatus = "active"
	StatusSold     ProductStatus = "sold"
	StatusReserved ProductStatus = "reserved"
	StatusDeleted  ProductStatus = "deleted"

	ConditionNew     ProductCondition = "new"
	ConditionLikeNew ProductCondition = "like_new"
	ConditionGood    ProductCondition = "good"
	ConditionFair    ProductCondition = "fair"
	ConditionPoor    ProductCondition = "poor"
)

type Product struct {
	ID                         uuid.UUID        `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	SellerID                   uuid.UUID        `json:"seller_id" gorm:"type:uuid;not null;index"`
	Seller                     *User            `json:"seller,omitempty" gorm:"foreignKey:SellerID"`
	Title                      string           `json:"title" gorm:"not null"`
	Description                string           `json:"description"`
	Price                      int              `json:"price" gorm:"not null"` // in cents/yen
	Category                   string           `json:"category" gorm:"not null;index"`
	Condition                  ProductCondition `json:"condition" gorm:"not null"`
	Status                     ProductStatus    `json:"status" gorm:"default:active;index"`
	WeightKg                   float64          `json:"weight_kg" gorm:"type:decimal(8,2)"`
	ManufacturerCountry        string           `json:"manufacturer_country"`
	EstimatedManufacturingYear int              `json:"estimated_manufacturing_year"`
	AIGeneratedDescription     string           `json:"ai_generated_description"`
	AISuggestedPrice           int              `json:"ai_suggested_price"`
	ViewCount                  int              `json:"view_count" gorm:"default:0"`
	FavoriteCount              int              `json:"favorite_count" gorm:"default:0"`
	Has3DModel                 bool             `json:"has_3d_model" gorm:"default:false"`
	ModelURL                   string           `json:"model_url"`
	Images                     []ProductImage   `json:"images,omitempty" gorm:"foreignKey:ProductID"`
	CreatedAt                  time.Time        `json:"created_at"`
	UpdatedAt                  time.Time        `json:"updated_at"`
	SoldAt                     *time.Time       `json:"sold_at"`
	DeletedAt                  *time.Time       `json:"-" gorm:"index"`
}

type ProductImage struct {
	ID           uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ProductID    uuid.UUID `json:"product_id" gorm:"type:uuid;not null;index"`
	ImageURL     string    `json:"image_url" gorm:"not null"`
	CDNURL       string    `json:"cdn_url"`
	DisplayOrder int       `json:"display_order" gorm:"default:0"`
	IsPrimary    bool      `json:"is_primary" gorm:"default:false"`
	Width        int       `json:"width"`
	Height       int       `json:"height"`
	CreatedAt    time.Time `json:"created_at"`
}

type ProductRepository interface {
	Create(product *Product) error
	FindByID(id uuid.UUID) (*Product, error)
	List(filters *ProductFilters) ([]*Product, *PaginationResponse, error)
	Update(product *Product) error
	Delete(id uuid.UUID) error
	IncrementViewCount(id uuid.UUID) error
}

type ProductFilters struct {
	Category  string
	MinPrice  int
	MaxPrice  int
	Condition ProductCondition
	Search    string
	Sort      string // price_asc, price_desc, created_desc
	Page      int
	Limit     int
}

type PaginationResponse struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

type CreateProductRequest struct {
	Title           string           `form:"title" binding:"required,min=3,max=255"`
	Description     string           `form:"description"`
	Price           int              `form:"price" binding:"required,min=0"`
	Category        string           `form:"category" binding:"required"`
	Condition       ProductCondition `form:"condition" binding:"required"`
	WeightKg        float64          `form:"weight_kg"`
	UseAIAssistance bool             `form:"use_ai_assistance"`
}
