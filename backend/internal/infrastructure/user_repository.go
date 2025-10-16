package infrastructure

import (
	"time"

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

		// Fetch updated user
		var user domain.User
		if err := tx.Where("id = ?", userID).First(&user).Error; err != nil {
			return err
		}

		// Calculate new level (every 20kg CO2 = 1 level)
		newLevel := int(user.TotalCO2SavedKg/20) + 1
		newScore := int(user.TotalCO2SavedKg * 6) // 1kg = 6 points

		// Update level and score
		if err := tx.Model(&domain.User{}).
			Where("id = ?", userID).
			Updates(map[string]interface{}{
				"level":               newLevel,
				"sustainability_score": newScore,
			}).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *userRepository) GetLeaderboard(limit int, period string) ([]*domain.LeaderboardEntry, error) {
	var users []domain.User
	query := r.db.Order("sustainability_score DESC").Limit(limit)

	// Filter by period if specified
	if period == "week" {
		oneWeekAgo := time.Now().AddDate(0, 0, -7)
		query = query.Where("updated_at >= ?", oneWeekAgo)
	} else if period == "month" {
		oneMonthAgo := time.Now().AddDate(0, -1, 0)
		query = query.Where("updated_at >= ?", oneMonthAgo)
	}

	if err := query.Find(&users).Error; err != nil {
		return nil, err
	}

	entries := make([]*domain.LeaderboardEntry, len(users))
	for i, user := range users {
		entries[i] = &domain.LeaderboardEntry{
			Rank:                i + 1,
			User:                &user,
			TotalCO2SavedKg:     user.TotalCO2SavedKg,
			SustainabilityScore: user.SustainabilityScore,
			Level:               user.Level,
		}
	}

	return entries, nil
}
