package infrastructure

import (
	"github.com/google/uuid"
	"github.com/yourusername/ecomate/backend/internal/domain"
	"gorm.io/gorm"
)

type FollowRepository struct {
	db *gorm.DB
}

func NewFollowRepository(db *gorm.DB) *FollowRepository {
	return &FollowRepository{db: db}
}

func (r *FollowRepository) Create(follow *domain.Follow) error {
	return r.db.Create(follow).Error
}

func (r *FollowRepository) Delete(followerID, followingID uuid.UUID) error {
	return r.db.Where("follower_id = ? AND following_id = ?", followerID, followingID).
		Delete(&domain.Follow{}).Error
}

func (r *FollowRepository) GetFollowers(userID uuid.UUID, limit int) ([]*domain.Follow, error) {
	var follows []*domain.Follow
	err := r.db.Preload("Follower").
		Where("following_id = ?", userID).
		Limit(limit).
		Order("created_at DESC").
		Find(&follows).Error
	return follows, err
}

func (r *FollowRepository) GetFollowing(userID uuid.UUID, limit int) ([]*domain.Follow, error) {
	var follows []*domain.Follow
	err := r.db.Preload("Following").
		Where("follower_id = ?", userID).
		Limit(limit).
		Order("created_at DESC").
		Find(&follows).Error
	return follows, err
}

func (r *FollowRepository) IsFollowing(followerID, followingID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.Model(&domain.Follow{}).
		Where("follower_id = ? AND following_id = ?", followerID, followingID).
		Count(&count).Error
	return count > 0, err
}

func (r *FollowRepository) GetFollowerCount(userID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.Model(&domain.Follow{}).
		Where("following_id = ?", userID).
		Count(&count).Error
	return count, err
}

func (r *FollowRepository) GetFollowingCount(userID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.Model(&domain.Follow{}).
		Where("follower_id = ?", userID).
		Count(&count).Error
	return count, err
}

type ProductShareRepository struct {
	db *gorm.DB
}

func NewProductShareRepository(db *gorm.DB) *ProductShareRepository {
	return &ProductShareRepository{db: db}
}

func (r *ProductShareRepository) Create(share *domain.ProductShare) error {
	return r.db.Create(share).Error
}

func (r *ProductShareRepository) GetByProduct(productID uuid.UUID) ([]*domain.ProductShare, error) {
	var shares []*domain.ProductShare
	err := r.db.Preload("User").
		Where("product_id = ?", productID).
		Order("created_at DESC").
		Find(&shares).Error
	return shares, err
}

func (r *ProductShareRepository) GetShareCount(productID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.Model(&domain.ProductShare{}).
		Where("product_id = ?", productID).
		Count(&count).Error
	return count, err
}

type UserFeedRepository struct {
	db *gorm.DB
}

func NewUserFeedRepository(db *gorm.DB) *UserFeedRepository {
	return &UserFeedRepository{db: db}
}

func (r *UserFeedRepository) Create(feed *domain.UserFeed) error {
	return r.db.Create(feed).Error
}

func (r *UserFeedRepository) GetFeed(userID uuid.UUID, limit int) ([]*domain.UserFeed, error) {
	var feeds []*domain.UserFeed
	err := r.db.Preload("Actor").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Find(&feeds).Error
	return feeds, err
}

func (r *UserFeedRepository) CreateFeedForFollowers(actorID uuid.UUID, actionType string, targetID uuid.UUID, targetType, description string) error {
	// Get all followers
	var follows []*domain.Follow
	if err := r.db.Where("following_id = ?", actorID).Find(&follows).Error; err != nil {
		return err
	}

	// Create feed for each follower
	for _, follow := range follows {
		feed := &domain.UserFeed{
			UserID:      follow.FollowerID,
			ActorID:     actorID,
			ActionType:  actionType,
			TargetID:    targetID,
			TargetType:  targetType,
			Description: description,
		}
		if err := r.Create(feed); err != nil {
			return err
		}
	}

	return nil
}
