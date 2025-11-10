package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Follow represents a user following relationship
type Follow struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	FollowerID  uuid.UUID      `gorm:"type:uuid;not null;index:idx_follow_follower" json:"follower_id"`
	FollowingID uuid.UUID      `gorm:"type:uuid;not null;index:idx_follow_following" json:"following_id"`
	CreatedAt   time.Time      `json:"created_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Follower  *User `gorm:"foreignKey:FollowerID" json:"follower,omitempty"`
	Following *User `gorm:"foreignKey:FollowingID" json:"following,omitempty"`
}

// ProductShare represents a product share on social media
type ProductShare struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	ProductID uuid.UUID      `gorm:"type:uuid;not null;index" json:"product_id"`
	Platform  string         `gorm:"not null" json:"platform"` // twitter, facebook, line, etc.
	ShareURL  string         `json:"share_url"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	User    *User    `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Product *Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
}

// UserFeed represents a feed item in user's timeline
type UserFeed struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID      uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	ActorID     uuid.UUID      `gorm:"type:uuid;not null" json:"actor_id"` // Who performed the action
	ActionType  string         `gorm:"not null" json:"action_type"`        // listed, purchased, reviewed, followed
	TargetID    uuid.UUID      `gorm:"type:uuid" json:"target_id"`         // Product ID or User ID
	TargetType  string         `json:"target_type"`                        // product, user
	Description string         `gorm:"type:text" json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	User  *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Actor *User `gorm:"foreignKey:ActorID" json:"actor,omitempty"`
}

type FollowRepository interface {
	Create(follow *Follow) error
	Delete(followerID, followingID uuid.UUID) error
	GetFollowers(userID uuid.UUID, limit int) ([]*Follow, error)
	GetFollowing(userID uuid.UUID, limit int) ([]*Follow, error)
	IsFollowing(followerID, followingID uuid.UUID) (bool, error)
	GetFollowerCount(userID uuid.UUID) (int64, error)
	GetFollowingCount(userID uuid.UUID) (int64, error)
}

type ProductShareRepository interface {
	Create(share *ProductShare) error
	GetByProduct(productID uuid.UUID) ([]*ProductShare, error)
	GetShareCount(productID uuid.UUID) (int64, error)
}

type UserFeedRepository interface {
	Create(feed *UserFeed) error
	GetFeed(userID uuid.UUID, limit int) ([]*UserFeed, error)
	CreateFeedForFollowers(actorID uuid.UUID, actionType string, targetID uuid.UUID, targetType, description string) error
}
