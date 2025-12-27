package domain

import (
	"time"

	"github.com/google/uuid"
)

type NotificationType string

const (
	NotificationTypeMessage  NotificationType = "message"
	NotificationTypePurchase NotificationType = "purchase"
	NotificationTypeFavorite NotificationType = "favorite"
	NotificationTypeReview   NotificationType = "review"
)

type Notification struct {
	ID        uuid.UUID        `json:"id" gorm:"type:char(36);primary_key"`
	UserID    uuid.UUID        `json:"user_id" gorm:"type:char(36);not null;index"`
	Type      NotificationType `json:"type" gorm:"not null"`
	Title     string           `json:"title" gorm:"not null"`
	Message   string           `json:"message"`
	Link      string           `json:"link"`
	IsRead    bool             `json:"is_read" gorm:"default:false"`
	CreatedAt time.Time        `json:"created_at"`
}

type Review struct {
	ID         uuid.UUID  `json:"id" gorm:"type:char(36);primary_key"`
	ProductID  uuid.UUID  `json:"product_id" gorm:"type:char(36);not null;index"`
	Product    *Product   `json:"product,omitempty" gorm:"foreignKey:ProductID"`
	PurchaseID uuid.UUID  `json:"purchase_id" gorm:"type:char(36);not null;index"`
	ReviewerID uuid.UUID  `json:"reviewer_id" gorm:"type:char(36);not null;index"`
	Reviewer   *User      `json:"reviewer,omitempty" gorm:"foreignKey:ReviewerID"`
	Rating     int        `json:"rating" gorm:"not null"` // 1-5
	Comment    string     `json:"comment"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

type NotificationRepository interface {
	Create(notification *Notification) error
	FindByUserID(userID uuid.UUID, limit int) ([]*Notification, error)
	MarkAsRead(id uuid.UUID) error
	MarkAllAsRead(userID uuid.UUID) error
	GetUnreadCount(userID uuid.UUID) (int64, error)
}

type ReviewRepository interface {
	Create(review *Review) error
	FindByProductID(productID uuid.UUID) ([]*Review, error)
	FindByID(id uuid.UUID) (*Review, error)
	GetAverageRating(productID uuid.UUID) (float64, error)
}
