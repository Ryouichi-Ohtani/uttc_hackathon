package infrastructure

import (
	"time"

	"github.com/google/uuid"
	"github.com/yourusername/ecomate/backend/internal/domain"
	"gorm.io/gorm"
)

type sustainabilityRepository struct {
	db *gorm.DB
}

func NewSustainabilityRepository(db *gorm.DB) domain.SustainabilityRepository {
	return &sustainabilityRepository{db: db}
}

// Achievements
func (r *sustainabilityRepository) GetUserAchievements(userID uuid.UUID) ([]*domain.UserAchievement, error) {
	var achievements []*domain.UserAchievement
	if err := r.db.
		Preload("Achievement").
		Where("user_id = ?", userID).
		Order("earned_at DESC").
		Find(&achievements).Error; err != nil {
		return nil, err
	}
	return achievements, nil
}

func (r *sustainabilityRepository) CheckAndAwardAchievements(userID uuid.UUID) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Get user stats
		var user domain.User
		if err := tx.Where("id = ?", userID).First(&user).Error; err != nil {
			return err
		}

		// Get transaction count
		var transactionCount int64
		if err := tx.Model(&domain.Purchase{}).
			Where("buyer_id = ? OR seller_id = ?", userID, userID).
			Where("status = ?", domain.PurchaseStatusCompleted).
			Count(&transactionCount).Error; err != nil {
			return err
		}

		// Get all achievements
		var allAchievements []domain.Achievement
		if err := tx.Find(&allAchievements).Error; err != nil {
			return err
		}

		// Get already earned achievements
		var earnedAchievements []domain.UserAchievement
		if err := tx.Where("user_id = ?", userID).Find(&earnedAchievements).Error; err != nil {
			return err
		}

		earnedMap := make(map[uuid.UUID]bool)
		for _, ea := range earnedAchievements {
			earnedMap[ea.AchievementID] = true
		}

		// Check each achievement
		for _, achievement := range allAchievements {
			if earnedMap[achievement.ID] {
				continue // Already earned
			}

			shouldAward := false
			switch achievement.RequirementType {
			case "transaction_count":
				if int(transactionCount) >= achievement.RequirementValue {
					shouldAward = true
				}
				// Removed CO2 and level based achievements
			}

			if shouldAward {
				userAchievement := domain.UserAchievement{
					UserID:        userID,
					AchievementID: achievement.ID,
					EarnedAt:      time.Now(),
				}
				if err := tx.Create(&userAchievement).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})
}

// Logs
func (r *sustainabilityRepository) CreateLog(log *domain.SustainabilityLog) error {
	return r.db.Create(log).Error
}

func (r *sustainabilityRepository) GetUserLogs(userID uuid.UUID, limit int) ([]*domain.SustainabilityLog, error) {
	var logs []*domain.SustainabilityLog
	if err := r.db.
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Find(&logs).Error; err != nil {
		return nil, err
	}
	return logs, nil
}

func (r *sustainabilityRepository) GetMonthlyStats(userID uuid.UUID) (*domain.MonthlyStats, error) {
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	var totalCO2 float64
	var transactionCount int64

	// Get CO2 saved this month
	if err := r.db.Model(&domain.SustainabilityLog{}).
		Where("user_id = ? AND created_at >= ?", userID, startOfMonth).
		Select("COALESCE(SUM(co2_saved_kg), 0)").
		Scan(&totalCO2).Error; err != nil {
		return nil, err
	}

	// Get transaction count this month
	if err := r.db.Model(&domain.Purchase{}).
		Where("(buyer_id = ? OR seller_id = ?) AND created_at >= ?", userID, userID, startOfMonth).
		Where("status = ?", domain.PurchaseStatusCompleted).
		Count(&transactionCount).Error; err != nil {
		return nil, err
	}

	return &domain.MonthlyStats{
		CurrentMonthCO2Saved: totalCO2,
		Transactions:         int(transactionCount),
	}, nil
}

// Favorites
func (r *sustainabilityRepository) AddFavorite(userID, productID uuid.UUID) error {
	favorite := domain.Favorite{
		UserID:    userID,
		ProductID: productID,
	}

	// Also increment favorite_count on product
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&favorite).Error; err != nil {
			return err
		}
		return tx.Model(&domain.Product{}).
			Where("id = ?", productID).
			Update("favorite_count", gorm.Expr("favorite_count + 1")).
			Error
	})
}

func (r *sustainabilityRepository) RemoveFavorite(userID, productID uuid.UUID) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_id = ? AND product_id = ?", userID, productID).
			Delete(&domain.Favorite{}).Error; err != nil {
			return err
		}
		return tx.Model(&domain.Product{}).
			Where("id = ?", productID).
			Update("favorite_count", gorm.Expr("GREATEST(favorite_count - 1, 0)")).
			Error
	})
}

func (r *sustainabilityRepository) IsFavorited(userID, productID uuid.UUID) (bool, error) {
	var count int64
	if err := r.db.Model(&domain.Favorite{}).
		Where("user_id = ? AND product_id = ?", userID, productID).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *sustainabilityRepository) GetUserFavorites(userID uuid.UUID) ([]*domain.Product, error) {
	var products []*domain.Product
	if err := r.db.
		Joins("JOIN favorites ON favorites.product_id = products.id").
		Where("favorites.user_id = ?", userID).
		Preload("Seller").
		Preload("Images", "is_primary = true").
		Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}
