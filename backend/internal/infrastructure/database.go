package infrastructure

import (
	"fmt"
	"log"

	"github.com/yourusername/ecomate/backend/internal/config"
	"github.com/yourusername/ecomate/backend/internal/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDatabase(cfg *config.DatabaseConfig) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Println("Database connected successfully")
	return db, nil
}

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&domain.User{},
		&domain.Product{},
		&domain.ProductImage{},
		&domain.Purchase{},
		&domain.Conversation{},
		&domain.ConversationParticipant{},
		&domain.Message{},
		&domain.Achievement{},
		&domain.UserAchievement{},
		&domain.SustainabilityLog{},
		&domain.Favorite{},
		&domain.Notification{},
		&domain.Review{},
		&domain.Offer{},
		&domain.UserEvent{},
		&domain.Auction{},
		&domain.Bid{},
		&domain.BlockchainTransaction{},
		&domain.NFTOwnership{},
		&domain.ChatHistory{},
		&domain.CO2Goal{},
		&domain.ShippingTracking{},
		&domain.Follow{},
		&domain.ProductShare{},
		&domain.UserFeed{},
		&domain.LiveStream{},
		&domain.StreamComment{},
	)
}

func SeedAchievements(db *gorm.DB) error {
	achievements := []domain.Achievement{
		{
			Name:             "First Step",
			Description:      "Complete your first transaction",
			RequirementType:  "transaction_count",
			RequirementValue: 1,
		},
		{
			Name:             "Eco Warrior",
			Description:      "Save 10kg of CO2",
			RequirementType:  "co2_saved",
			RequirementValue: 10,
		},
		{
			Name:             "Planet Hero",
			Description:      "Save 50kg of CO2",
			RequirementType:  "co2_saved",
			RequirementValue: 50,
		},
		{
			Name:             "Climate Champion",
			Description:      "Save 100kg of CO2",
			RequirementType:  "co2_saved",
			RequirementValue: 100,
		},
		{
			Name:             "Master Trader",
			Description:      "Complete 50 transactions",
			RequirementType:  "transaction_count",
			RequirementValue: 50,
		},
	}

	for _, achievement := range achievements {
		var existing domain.Achievement
		if err := db.Where("name = ?", achievement.Name).First(&existing).Error; err == gorm.ErrRecordNotFound {
			if err := db.Create(&achievement).Error; err != nil {
				return err
			}
		}
	}

	log.Println("Achievements seeded successfully")
	return nil
}
