package usecase

import (
	"time"

	"github.com/google/uuid"
	"github.com/yourusername/ecomate/backend/internal/domain"
)

type AnalyticsUseCase interface {
	TrackEvent(event *domain.UserEvent) error
	GetUserBehaviorSummary(userID uuid.UUID) (*domain.UserBehaviorSummary, error)
	GetPopularProducts(days int) ([]*domain.ProductAnalytics, error)
	GetSearchTrends(days int) ([]*domain.SearchKeyword, error)
}

type analyticsUseCase struct {
	analyticsRepo domain.AnalyticsRepository
}

func NewAnalyticsUseCase(analyticsRepo domain.AnalyticsRepository) AnalyticsUseCase {
	return &analyticsUseCase{
		analyticsRepo: analyticsRepo,
	}
}

func (u *analyticsUseCase) TrackEvent(event *domain.UserEvent) error {
	return u.analyticsRepo.TrackEvent(event)
}

func (u *analyticsUseCase) GetUserBehaviorSummary(userID uuid.UUID) (*domain.UserBehaviorSummary, error) {
	return u.analyticsRepo.GetUserBehaviorSummary(userID)
}

func (u *analyticsUseCase) GetPopularProducts(days int) ([]*domain.ProductAnalytics, error) {
	since := time.Now().AddDate(0, 0, -days)
	return u.analyticsRepo.GetPopularProducts(20, since)
}

func (u *analyticsUseCase) GetSearchTrends(days int) ([]*domain.SearchKeyword, error) {
	since := time.Now().AddDate(0, 0, -days)
	return u.analyticsRepo.GetSearchKeywords(20, since)
}
