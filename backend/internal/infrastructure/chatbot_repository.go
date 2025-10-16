package infrastructure

import (
	"github.com/google/uuid"
	"github.com/yourusername/ecomate/backend/internal/domain"
	"gorm.io/gorm"
)

type ChatHistoryRepository struct {
	db *gorm.DB
}

func NewChatHistoryRepository(db *gorm.DB) *ChatHistoryRepository {
	return &ChatHistoryRepository{db: db}
}

func (r *ChatHistoryRepository) Create(history *domain.ChatHistory) error {
	return r.db.Create(history).Error
}

func (r *ChatHistoryRepository) GetByUserID(userID uuid.UUID, limit int) ([]*domain.ChatHistory, error) {
	var histories []*domain.ChatHistory
	err := r.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Find(&histories).Error
	return histories, err
}

func (r *ChatHistoryRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&domain.ChatHistory{}, "id = ?", id).Error
}

type CO2GoalRepository struct {
	db *gorm.DB
}

func NewCO2GoalRepository(db *gorm.DB) *CO2GoalRepository {
	return &CO2GoalRepository{db: db}
}

func (r *CO2GoalRepository) Create(goal *domain.CO2Goal) error {
	return r.db.Create(goal).Error
}

func (r *CO2GoalRepository) GetByUserID(userID uuid.UUID) (*domain.CO2Goal, error) {
	var goal domain.CO2Goal
	err := r.db.Where("user_id = ? AND status = ?", userID, "active").
		First(&goal).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &goal, err
}

func (r *CO2GoalRepository) Update(goal *domain.CO2Goal) error {
	return r.db.Save(goal).Error
}

func (r *CO2GoalRepository) UpdateProgress(userID uuid.UUID, additionalKG float64) error {
	return r.db.Model(&domain.CO2Goal{}).
		Where("user_id = ? AND status = ?", userID, "active").
		UpdateColumn("current_kg", gorm.Expr("current_kg + ?", additionalKG)).Error
}

type ShippingTrackingRepository struct {
	db *gorm.DB
}

func NewShippingTrackingRepository(db *gorm.DB) *ShippingTrackingRepository {
	return &ShippingTrackingRepository{db: db}
}

func (r *ShippingTrackingRepository) Create(tracking *domain.ShippingTracking) error {
	return r.db.Create(tracking).Error
}

func (r *ShippingTrackingRepository) GetByPurchaseID(purchaseID uuid.UUID) (*domain.ShippingTracking, error) {
	var tracking domain.ShippingTracking
	err := r.db.Preload("Purchase").
		Where("purchase_id = ?", purchaseID).
		First(&tracking).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &tracking, err
}

func (r *ShippingTrackingRepository) Update(tracking *domain.ShippingTracking) error {
	return r.db.Save(tracking).Error
}

func (r *ShippingTrackingRepository) UpdateStatus(id uuid.UUID, status string) error {
	return r.db.Model(&domain.ShippingTracking{}).
		Where("id = ?", id).
		Update("status", status).Error
}
