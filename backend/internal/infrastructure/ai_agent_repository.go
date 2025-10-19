package infrastructure

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/yourusername/ecomate/backend/internal/domain"
	"gorm.io/gorm"
)

type AIAgentRepositoryImpl struct {
	db *gorm.DB
}

func NewAIAgentRepository(db *gorm.DB) domain.AIAgentRepository {
	return &AIAgentRepositoryImpl{db: db}
}

// ========== Listing Agent ==========

func (r *AIAgentRepositoryImpl) CreateListingData(data *domain.AIListingData) error {
	return r.db.Create(data).Error
}

func (r *AIAgentRepositoryImpl) GetListingData(productID uuid.UUID) (*domain.AIListingData, error) {
	var data domain.AIListingData
	err := r.db.Where("product_id = ?", productID).First(&data).Error
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *AIAgentRepositoryImpl) UpdateListingData(data *domain.AIListingData) error {
	return r.db.Save(data).Error
}

// ========== Negotiation Agent ==========

func (r *AIAgentRepositoryImpl) CreateNegotiationSettings(settings *domain.AINegotiationSettings) error {
	return r.db.Create(settings).Error
}

func (r *AIAgentRepositoryImpl) GetNegotiationSettings(productID uuid.UUID) (*domain.AINegotiationSettings, error) {
	var settings domain.AINegotiationSettings
	err := r.db.Where("product_id = ?", productID).First(&settings).Error
	if err != nil {
		return nil, err
	}
	return &settings, nil
}

func (r *AIAgentRepositoryImpl) UpdateNegotiationSettings(settings *domain.AINegotiationSettings) error {
	return r.db.Save(settings).Error
}

func (r *AIAgentRepositoryImpl) DeleteNegotiationSettings(productID uuid.UUID) error {
	return r.db.Where("product_id = ?", productID).Delete(&domain.AINegotiationSettings{}).Error
}

// ========== Shipping Agent ==========

func (r *AIAgentRepositoryImpl) CreateShippingPreparation(prep *domain.AIShippingPreparation) error {
	return r.db.Create(prep).Error
}

func (r *AIAgentRepositoryImpl) GetShippingPreparation(purchaseID uuid.UUID) (*domain.AIShippingPreparation, error) {
	var prep domain.AIShippingPreparation
	err := r.db.Where("purchase_id = ?", purchaseID).
		Preload("Purchase").
		Preload("Purchase.Product").
		First(&prep).Error
	if err != nil {
		return nil, err
	}
	return &prep, nil
}

func (r *AIAgentRepositoryImpl) UpdateShippingPreparation(prep *domain.AIShippingPreparation) error {
	return r.db.Save(prep).Error
}

func (r *AIAgentRepositoryImpl) ApproveShipping(purchaseID uuid.UUID, modifications string) error {
	now := time.Now()
	return r.db.Model(&domain.AIShippingPreparation{}).
		Where("purchase_id = ?", purchaseID).
		Updates(map[string]interface{}{
			"user_approved":      true,
			"user_modifications": modifications,
			"approved_at":        now,
		}).Error
}

// ========== Logs ==========

func (r *AIAgentRepositoryImpl) CreateLog(log *domain.AIAgentLog) error {
	return r.db.Create(log).Error
}

func (r *AIAgentRepositoryImpl) GetUserLogs(userID uuid.UUID, agentType domain.AgentType, limit int) ([]*domain.AIAgentLog, error) {
	var logs []*domain.AIAgentLog
	query := r.db.Where("user_id = ?", userID)

	if agentType != "" {
		query = query.Where("agent_type = ?", agentType)
	}

	err := query.Order("created_at DESC").Limit(limit).Find(&logs).Error
	return logs, err
}

func (r *AIAgentRepositoryImpl) GetAgentStats(userID uuid.UUID) (*domain.AIAgentStats, error) {
	stats := &domain.AIAgentStats{}

	// Total AI generations
	r.db.Model(&domain.AIAgentLog{}).
		Where("user_id = ? AND success = ?", userID, true).
		Count(&stats.TotalAIGenerations)

	// Listings created
	r.db.Model(&domain.AIAgentLog{}).
		Where("user_id = ? AND agent_type = ? AND success = ?", userID, domain.AgentTypeListing, true).
		Count(&stats.ListingsCreated)

	// Negotiations handled
	var negSettings []domain.AINegotiationSettings
	r.db.Joins("JOIN products ON products.id = ai_negotiation_settings.product_id").
		Where("products.seller_id = ?", userID).
		Find(&negSettings)

	totalNegotiations := 0
	totalAccepted := 0
	for _, setting := range negSettings {
		totalNegotiations += setting.TotalOffersProcessed
		totalAccepted += setting.AIAcceptedCount
	}
	stats.NegotiationsHandled = int64(totalNegotiations)

	if totalNegotiations > 0 {
		stats.AcceptanceRate = float64(totalAccepted) / float64(totalNegotiations) * 100
	}

	// Shipments prepared
	r.db.Model(&domain.AIShippingPreparation{}).
		Joins("JOIN purchases ON purchases.id = ai_shipping_preparations.purchase_id").
		Where("purchases.seller_id = ? AND ai_shipping_preparations.user_approved = ?", userID, true).
		Count(&stats.ShipmentsPrepared)

	// Average confidence
	var avgConfidence struct {
		Avg float64
	}
	r.db.Model(&domain.AIListingData{}).
		Joins("JOIN products ON products.id = ai_listing_data.product_id").
		Where("products.seller_id = ?", userID).
		Select("AVG(ai_confidence_score) as avg").
		Scan(&avgConfidence)
	stats.AverageConfidence = avgConfidence.Avg

	// Estimate time saved (1 listing = 15 min, 1 negotiation = 5 min, 1 shipment = 10 min)
	stats.TimeSavedMinutes = stats.ListingsCreated*15 + stats.NegotiationsHandled*5 + stats.ShipmentsPrepared*10

	return stats, nil
}

// Helper: Track AI agent action
func (r *AIAgentRepositoryImpl) TrackAction(userID uuid.UUID, agentType domain.AgentType, action string, targetID uuid.UUID, details interface{}, success bool, processTimeMs int) error {
	detailsJSON, _ := json.Marshal(details)

	log := &domain.AIAgentLog{
		UserID:      userID,
		AgentType:   agentType,
		Action:      action,
		TargetID:    targetID,
		Details:     string(detailsJSON),
		Success:     success,
		ProcessTime: processTimeMs,
	}

	return r.CreateLog(log)
}
