package infrastructure

import (
	"context"
	"fmt"
	"log"
	"reflect"

	"github.com/google/uuid"
	"github.com/yourusername/ecomate/backend/internal/config"
	"github.com/yourusername/ecomate/backend/internal/domain"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func NewDatabase(cfg *config.DatabaseConfig) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(cfg.DSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	registerUUIDCallbacks(db)

	log.Println("Database connected successfully")
	return db, nil
}

func registerUUIDCallbacks(db *gorm.DB) {
	db.Callback().Create().Before("gorm:create").Register("ecomate:set_uuid", func(tx *gorm.DB) {
		if tx.Statement == nil || tx.Statement.Schema == nil {
			return
		}

		field := tx.Statement.Schema.LookUpField("ID")
		if field == nil || field.FieldType != reflect.TypeOf(uuid.UUID{}) {
			return
		}

		setUUIDIfZero(tx.Statement.Context, field, tx.Statement.ReflectValue)
	})
}

func setUUIDIfZero(ctx context.Context, field *schema.Field, value reflect.Value) {
	if !value.IsValid() {
		return
	}

	if value.Kind() == reflect.Ptr {
		if value.IsNil() {
			return
		}
		value = value.Elem()
	}

	switch value.Kind() {
	case reflect.Struct:
		if _, isZero := field.ValueOf(ctx, value); isZero {
			_ = field.Set(ctx, value, uuid.New())
		}
	case reflect.Slice, reflect.Array:
		for i := 0; i < value.Len(); i++ {
			setUUIDIfZero(ctx, field, value.Index(i))
		}
	}
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
