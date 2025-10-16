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
	user, err := u.userRepo.FindByID(userID)
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

	// Calculate next level threshold (exponential)
	nextLevelThreshold := user.Level * 100

	// Calculate comparisons
	equivalentTrees := user.TotalCO2SavedKg / 20  // 1 tree absorbs ~20kg CO2/year
	carKmAvoided := user.TotalCO2SavedKg / 0.12   // Car emits ~0.12kg CO2/km

	dashboard := &domain.DashboardResponse{
		TotalCO2SavedKg:     user.TotalCO2SavedKg,
		Level:               user.Level,
		SustainabilityScore: user.SustainabilityScore,
		NextLevelThreshold:  nextLevelThreshold,
		Achievements:        achievements,
		RecentLogs:          logs,
		MonthlyStats:        monthlyStats,
		Comparisons: &domain.EnvironmentComparison{
			EquivalentTrees: equivalentTrees,
			CarKmAvoided:    carKmAvoided,
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
