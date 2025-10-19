package usecase

import (
	"github.com/google/uuid"
	"github.com/yourusername/ecomate/backend/internal/domain"
)

type SustainabilityUseCase interface {
	GetDashboard(userID uuid.UUID) (*domain.DashboardResponse, error)
	GetLeaderboard(limit int, period string) ([]*domain.LeaderboardEntry, error)
	GetUserFavorites(userID uuid.UUID) ([]*domain.Product, error)
}

type sustainabilityUseCase struct {
	sustainabilityRepo domain.SustainabilityRepository
	userRepo           domain.UserRepository
}

func NewSustainabilityUseCase(
	sustainabilityRepo domain.SustainabilityRepository,
	userRepo domain.UserRepository,
) SustainabilityUseCase {
	return &sustainabilityUseCase{
		sustainabilityRepo: sustainabilityRepo,
		userRepo:           userRepo,
	}
}

func (u *sustainabilityUseCase) GetDashboard(userID uuid.UUID) (*domain.DashboardResponse, error) {
	_, err := u.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	achievements, err := u.sustainabilityRepo.GetUserAchievements(userID)
	if err != nil {
		return nil, err
	}

	logs, err := u.sustainabilityRepo.GetUserLogs(userID, 10)
	if err != nil {
		return nil, err
	}

	monthlyStats, err := u.sustainabilityRepo.GetMonthlyStats(userID)
	if err != nil {
		return nil, err
	}

	// CO2 tracking removed - return minimal dashboard
	dashboard := &domain.DashboardResponse{
		TotalCO2SavedKg:     0,
		Level:               0,
		SustainabilityScore: 0,
		NextLevelThreshold:  0,
		Achievements:        achievements,
		RecentLogs:          logs,
		MonthlyStats:        monthlyStats,
		Comparisons: &domain.EnvironmentComparison{
			EquivalentTrees: 0,
			CarKmAvoided:    0,
		},
	}

	return dashboard, nil
}

func (u *sustainabilityUseCase) GetLeaderboard(limit int, period string) ([]*domain.LeaderboardEntry, error) {
	return u.userRepo.GetLeaderboard(limit, period)
}

func (u *sustainabilityUseCase) GetUserFavorites(userID uuid.UUID) ([]*domain.Product, error) {
	return u.sustainabilityRepo.GetUserFavorites(userID)
}
