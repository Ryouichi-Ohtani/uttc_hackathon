package domain

import (
	"time"

	"github.com/google/uuid"
)

type Achievement struct {
	ID              uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name            string    `json:"name" gorm:"uniqueIndex;not null"`
	Description     string    `json:"description"`
	BadgeIconURL    string    `json:"badge_icon_url"`
	RequirementType string    `json:"requirement_type" gorm:"not null"` // co2_saved, transaction_count, level
	RequirementValue int      `json:"requirement_value" gorm:"not null"`
	CreatedAt       time.Time `json:"created_at"`
}

type UserAchievement struct {
	ID            uuid.UUID    `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID        uuid.UUID    `json:"user_id" gorm:"type:uuid;not null;index"`
	AchievementID uuid.UUID    `json:"achievement_id" gorm:"type:uuid;not null"`
	Achievement   *Achievement `json:"achievement,omitempty" gorm:"foreignKey:AchievementID"`
	EarnedAt      time.Time    `json:"earned_at"`
}

type SustainabilityLog struct {
	ID          uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID      uuid.UUID  `json:"user_id" gorm:"type:uuid;not null;index"`
	PurchaseID  *uuid.UUID `json:"purchase_id" gorm:"type:uuid"`
	ActionType  string     `json:"action_type" gorm:"not null"` // purchase, sale
	CO2SavedKg  float64    `json:"co2_saved_kg" gorm:"type:decimal(10,2);not null"`
	Description string     `json:"description"`
	CreatedAt   time.Time  `json:"created_at"`
}

type Favorite struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID    uuid.UUID `json:"user_id" gorm:"type:uuid;not null;index"`
	ProductID uuid.UUID `json:"product_id" gorm:"type:uuid;not null;index"`
	CreatedAt time.Time `json:"created_at"`
}

type SustainabilityRepository interface {
	// Achievements
	GetUserAchievements(userID uuid.UUID) ([]*UserAchievement, error)
	CheckAndAwardAchievements(userID uuid.UUID) error

	// Logs
	CreateLog(log *SustainabilityLog) error
	GetUserLogs(userID uuid.UUID, limit int) ([]*SustainabilityLog, error)
	GetMonthlyStats(userID uuid.UUID) (*MonthlyStats, error)

	// Favorites
	AddFavorite(userID, productID uuid.UUID) error
	RemoveFavorite(userID, productID uuid.UUID) error
	IsFavorited(userID, productID uuid.UUID) (bool, error)
	GetUserFavorites(userID uuid.UUID) ([]*Product, error)
}

type DashboardResponse struct {
	TotalCO2SavedKg     float64              `json:"total_co2_saved_kg"`
	Level               int                  `json:"level"`
	SustainabilityScore int                  `json:"sustainability_score"`
	NextLevelThreshold  int                  `json:"next_level_threshold"`
	Achievements        []*UserAchievement   `json:"achievements"`
	RecentLogs          []*SustainabilityLog `json:"recent_logs"`
	MonthlyStats        *MonthlyStats        `json:"monthly_stats"`
	Comparisons         *EnvironmentComparison `json:"comparisons"`
}

type MonthlyStats struct {
	CurrentMonthCO2Saved float64 `json:"current_month_co2_saved"`
	Transactions         int     `json:"transactions"`
}

type EnvironmentComparison struct {
	EquivalentTrees float64 `json:"equivalent_trees"`
	CarKmAvoided    float64 `json:"car_km_avoided"`
}
