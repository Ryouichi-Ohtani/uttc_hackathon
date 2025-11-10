package infrastructure

import (
	"github.com/google/uuid"
	"github.com/yourusername/ecomate/backend/internal/domain"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) domain.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *domain.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) FindByID(id uuid.UUID) (*domain.User, error) {
	var user domain.User
	if err := r.db.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByEmail(email string) (*domain.User, error) {
	var user domain.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByUsername(username string) (*domain.User, error) {
	var user domain.User
	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(user *domain.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) UpdateSustainabilityStats(userID uuid.UUID, co2SavedKg float64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Update total CO2 saved
		if err := tx.Model(&domain.User{}).
			Where("id = ?", userID).
			Update("total_co2_saved_kg", gorm.Expr("total_co2_saved_kg + ?", co2SavedKg)).
			Error; err != nil {
			return err
		}

		// CO2 tracking removed - AI Agent focused app
		return nil
	})
}

func (r *userRepository) GetLeaderboard(limit int, period string) ([]*domain.LeaderboardEntry, error) {
	// Leaderboard feature removed - AI Agent focused app
	return []*domain.LeaderboardEntry{}, nil
}
