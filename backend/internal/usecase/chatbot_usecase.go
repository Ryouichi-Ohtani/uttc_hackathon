package usecase

import (
	"time"

	"github.com/google/uuid"
	"github.com/yourusername/ecomate/backend/internal/domain"
)

type ChatHistoryUseCase struct {
	chatHistoryRepo domain.ChatHistoryRepository
}

func NewChatHistoryUseCase(repo domain.ChatHistoryRepository) *ChatHistoryUseCase {
	return &ChatHistoryUseCase{chatHistoryRepo: repo}
}

func (uc *ChatHistoryUseCase) SaveHistory(userID uuid.UUID, message, response, context string) error {
	history := &domain.ChatHistory{
		UserID:   userID,
		Message:  message,
		Response: response,
		Context:  context,
	}
	return uc.chatHistoryRepo.Create(history)
}

func (uc *ChatHistoryUseCase) GetHistory(userID uuid.UUID, limit int) ([]*domain.ChatHistory, error) {
	if limit <= 0 {
		limit = 50
	}
	return uc.chatHistoryRepo.GetByUserID(userID, limit)
}

func (uc *ChatHistoryUseCase) DeleteHistory(id uuid.UUID) error {
	return uc.chatHistoryRepo.Delete(id)
}

type CO2GoalUseCase struct {
	goalRepo domain.CO2GoalRepository
}

func NewCO2GoalUseCase(goalRepo domain.CO2GoalRepository) *CO2GoalUseCase {
	return &CO2GoalUseCase{goalRepo: goalRepo}
}

func (uc *CO2GoalUseCase) CreateGoal(userID uuid.UUID, targetKG float64, targetDate time.Time) (*domain.CO2Goal, error) {
	goal := &domain.CO2Goal{
		UserID:     userID,
		TargetKG:   targetKG,
		CurrentKG:  0,
		TargetDate: targetDate,
		StartDate:  time.Now(),
		Status:     "active",
	}
	err := uc.goalRepo.Create(goal)
	return goal, err
}

func (uc *CO2GoalUseCase) GetGoal(userID uuid.UUID) (*domain.CO2Goal, error) {
	goal, err := uc.goalRepo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	// Check if goal expired
	if goal != nil && goal.Status == "active" && time.Now().After(goal.TargetDate) {
		if goal.CurrentKG >= goal.TargetKG {
			goal.Status = "completed"
		} else {
			goal.Status = "expired"
		}
		uc.goalRepo.Update(goal)
	}

	return goal, nil
}

func (uc *CO2GoalUseCase) UpdateProgress(userID uuid.UUID, additionalKG float64) error {
	return uc.goalRepo.UpdateProgress(userID, additionalKG)
}

type ShippingTrackingUseCase struct {
	shippingRepo domain.ShippingTrackingRepository
}

func NewShippingTrackingUseCase(repo domain.ShippingTrackingRepository) *ShippingTrackingUseCase {
	return &ShippingTrackingUseCase{shippingRepo: repo}
}

func (uc *ShippingTrackingUseCase) CreateTracking(
	purchaseID uuid.UUID,
	carrier string,
	shippingMethod string,
) (*domain.ShippingTracking, error) {
	// Generate tracking number (in real app, this would come from carrier API)
	trackingNumber := generateTrackingNumber(carrier)

	// Calculate CO2 saved for eco shipping
	co2Saved := 0.0
	if shippingMethod == "eco" {
		co2Saved = 0.5 // 500g CO2 saved
	}

	tracking := &domain.ShippingTracking{
		PurchaseID:     purchaseID,
		TrackingNumber: trackingNumber,
		Carrier:        carrier,
		Status:         "pending",
		ShippingMethod: shippingMethod,
		CO2Saved:       co2Saved,
	}

	err := uc.shippingRepo.Create(tracking)
	return tracking, err
}

func (uc *ShippingTrackingUseCase) GetTracking(purchaseID uuid.UUID) (*domain.ShippingTracking, error) {
	return uc.shippingRepo.GetByPurchaseID(purchaseID)
}

func (uc *ShippingTrackingUseCase) UpdateStatus(id uuid.UUID, status string) error {
	tracking, err := uc.shippingRepo.GetByPurchaseID(id)
	if err != nil {
		return err
	}

	tracking.Status = status
	now := time.Now()

	switch status {
	case "shipped":
		tracking.ShippedAt = &now
		estimatedArrival := now.Add(72 * time.Hour) // 3 days
		tracking.EstimatedArrival = &estimatedArrival
	case "delivered":
		tracking.DeliveredAt = &now
	}

	return uc.shippingRepo.Update(tracking)
}

func generateTrackingNumber(carrier string) string {
	// Simple tracking number generation
	// In production, this would integrate with carrier APIs
	prefix := "TRK"
	switch carrier {
	case "yamato":
		prefix = "YMT"
	case "sagawa":
		prefix = "SGW"
	case "yupack":
		prefix = "YUP"
	}

	return prefix + uuid.New().String()[:8]
}
