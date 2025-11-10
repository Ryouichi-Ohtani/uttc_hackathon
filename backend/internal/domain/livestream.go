package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// LiveStream represents a live streaming session
type LiveStream struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	SellerID    uuid.UUID      `gorm:"type:uuid;not null;index" json:"seller_id"`
	Title       string         `gorm:"not null" json:"title"`
	Description string         `gorm:"type:text" json:"description"`
	ProductIDs  []uuid.UUID    `gorm:"type:uuid[];not null" json:"product_ids"`    // Products featured in stream
	Status      string         `gorm:"not null;default:'scheduled'" json:"status"` // scheduled, live, ended
	StreamURL   string         `json:"stream_url"`
	ViewerCount int            `gorm:"default:0" json:"viewer_count"`
	StartedAt   *time.Time     `json:"started_at"`
	EndedAt     *time.Time     `json:"ended_at"`
	ScheduledAt time.Time      `json:"scheduled_at"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Seller   *User            `gorm:"foreignKey:SellerID" json:"seller,omitempty"`
	Comments []*StreamComment `gorm:"foreignKey:StreamID" json:"comments,omitempty"`
}

// StreamComment represents a comment in a live stream
type StreamComment struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	StreamID  uuid.UUID      `gorm:"type:uuid;not null;index" json:"stream_id"`
	UserID    uuid.UUID      `gorm:"type:uuid;not null" json:"user_id"`
	Comment   string         `gorm:"type:text;not null" json:"comment"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

type LiveStreamRepository interface {
	Create(stream *LiveStream) error
	GetByID(id uuid.UUID) (*LiveStream, error)
	Update(stream *LiveStream) error
	GetLiveStreams() ([]*LiveStream, error)
	GetUpcomingStreams(limit int) ([]*LiveStream, error)
	GetBySellerID(sellerID uuid.UUID) ([]*LiveStream, error)
	IncrementViewerCount(id uuid.UUID) error
	DecrementViewerCount(id uuid.UUID) error
}

type StreamCommentRepository interface {
	Create(comment *StreamComment) error
	GetByStream(streamID uuid.UUID, limit int) ([]*StreamComment, error)
}
