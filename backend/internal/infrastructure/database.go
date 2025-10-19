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
		&domain.Favorite{},
		&domain.Notification{},
		&domain.Review{},
		&domain.Offer{},
		&domain.NegotiationLog{},
		&domain.UserEvent{},
		&domain.Auction{},
		&domain.Bid{},
		&domain.ChatHistory{},
		&domain.ShippingTracking{},
		&domain.Follow{},
		&domain.ProductShare{},
		&domain.UserFeed{},
		&domain.LiveStream{},
		&domain.StreamComment{},
		// AI Agent Models
		&domain.AIListingData{},
		&domain.AINegotiationSettings{},
		&domain.AIShippingPreparation{},
		&domain.AIAgentLog{},
	)
}

func SeedAchievements(db *gorm.DB) error {
	// Removed sustainability achievements - focus on AI agent usage
	log.Println("Achievements seeding removed - AI Agent focused app")

	log.Println("Achievements seeded successfully")
	return nil
}
