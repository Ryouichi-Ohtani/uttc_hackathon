package domain

import (
	"time"

	"github.com/google/uuid"
)

type EventType string

const (
	EventProductView   EventType = "product_view"
	EventProductSearch EventType = "product_search"
	EventProductLike   EventType = "product_like"
	EventProductOffer  EventType = "product_offer"
	EventPurchase      EventType = "purchase"
	EventMessageSent   EventType = "message_sent"
)

type UserEvent struct {
	ID         uuid.UUID `json:"id" gorm:"type:char(36);primary_key"`
	UserID     uuid.UUID `json:"user_id" gorm:"type:char(36);not null;index"`
	EventType  EventType `json:"event_type" gorm:"not null;index"`
	ProductID  *uuid.UUID `json:"product_id" gorm:"type:char(36);index"`
	Metadata   string    `json:"metadata" gorm:"type:json"`
	IPAddress  string    `json:"ip_address"`
	UserAgent  string    `json:"user_agent"`
	CreatedAt  time.Time `json:"created_at" gorm:"index"`
}

type AnalyticsRepository interface {
	TrackEvent(event *UserEvent) error
	GetUserEvents(userID uuid.UUID, eventType EventType, limit int) ([]*UserEvent, error)
	GetPopularProducts(limit int, since time.Time) ([]*ProductAnalytics, error)
	GetUserBehaviorSummary(userID uuid.UUID) (*UserBehaviorSummary, error)
	GetSearchKeywords(limit int, since time.Time) ([]*SearchKeyword, error)
}

type ProductAnalytics struct {
	ProductID  uuid.UUID `json:"product_id"`
	Product    *Product  `json:"product"`
	ViewCount  int       `json:"view_count"`
	LikeCount  int       `json:"like_count"`
	OfferCount int       `json:"offer_count"`
}

type UserBehaviorSummary struct {
	TotalViews     int      `json:"total_views"`
	TotalSearches  int      `json:"total_searches"`
	TotalLikes     int      `json:"total_likes"`
	TotalOffers    int      `json:"total_offers"`
	FavoriteCategories []CategoryCount `json:"favorite_categories"`
	RecentlyViewed []*Product `json:"recently_viewed"`
}

type CategoryCount struct {
	Category string `json:"category"`
	Count    int    `json:"count"`
}

type SearchKeyword struct {
	Keyword string `json:"keyword"`
	Count   int    `json:"count"`
}
