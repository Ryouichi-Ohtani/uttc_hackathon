package infrastructure

import (
	"time"

	"github.com/google/uuid"
	"github.com/yourusername/ecomate/backend/internal/domain"
	"gorm.io/gorm"
)

type LiveStreamRepository struct {
	db *gorm.DB
}

func NewLiveStreamRepository(db *gorm.DB) *LiveStreamRepository {
	return &LiveStreamRepository{db: db}
}

func (r *LiveStreamRepository) Create(stream *domain.LiveStream) error {
	return r.db.Create(stream).Error
}

func (r *LiveStreamRepository) GetByID(id uuid.UUID) (*domain.LiveStream, error) {
	var stream domain.LiveStream
	err := r.db.Preload("Seller").Preload("Comments.User").First(&stream, "id = ?", id).Error
	return &stream, err
}

func (r *LiveStreamRepository) Update(stream *domain.LiveStream) error {
	return r.db.Save(stream).Error
}

func (r *LiveStreamRepository) GetLiveStreams() ([]*domain.LiveStream, error) {
	var streams []*domain.LiveStream
	err := r.db.Preload("Seller").
		Where("status = ?", "live").
		Order("viewer_count DESC").
		Find(&streams).Error
	return streams, err
}

func (r *LiveStreamRepository) GetUpcomingStreams(limit int) ([]*domain.LiveStream, error) {
	var streams []*domain.LiveStream
	err := r.db.Preload("Seller").
		Where("status = ? AND scheduled_at > ?", "scheduled", time.Now()).
		Order("scheduled_at ASC").
		Limit(limit).
		Find(&streams).Error
	return streams, err
}

func (r *LiveStreamRepository) GetBySellerID(sellerID uuid.UUID) ([]*domain.LiveStream, error) {
	var streams []*domain.LiveStream
	err := r.db.Where("seller_id = ?", sellerID).
		Order("created_at DESC").
		Find(&streams).Error
	return streams, err
}

func (r *LiveStreamRepository) IncrementViewerCount(id uuid.UUID) error {
	return r.db.Model(&domain.LiveStream{}).
		Where("id = ?", id).
		UpdateColumn("viewer_count", gorm.Expr("viewer_count + ?", 1)).Error
}

func (r *LiveStreamRepository) DecrementViewerCount(id uuid.UUID) error {
	return r.db.Model(&domain.LiveStream{}).
		Where("id = ?", id).
		UpdateColumn("viewer_count", gorm.Expr("GREATEST(viewer_count - 1, 0)")).Error
}

type StreamCommentRepository struct {
	db *gorm.DB
}

func NewStreamCommentRepository(db *gorm.DB) *StreamCommentRepository {
	return &StreamCommentRepository{db: db}
}

func (r *StreamCommentRepository) Create(comment *domain.StreamComment) error {
	return r.db.Create(comment).Error
}

func (r *StreamCommentRepository) GetByStream(streamID uuid.UUID, limit int) ([]*domain.StreamComment, error) {
	var comments []*domain.StreamComment
	err := r.db.Preload("User").
		Where("stream_id = ?", streamID).
		Order("created_at DESC").
		Limit(limit).
		Find(&comments).Error
	return comments, err
}
