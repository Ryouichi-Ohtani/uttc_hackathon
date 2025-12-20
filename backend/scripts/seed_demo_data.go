//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Demo data entities (simplified)
type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Email        string    `gorm:"uniqueIndex;not null"`
	Username     string    `gorm:"uniqueIndex;not null"`
	PasswordHash string    `gorm:"not null"`
	Role         string    `gorm:"default:'user'"`
	CreatedAt    time.Time
}

type Product struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	SellerID    uuid.UUID `gorm:"type:uuid;not null"`
	Title       string    `gorm:"not null"`
	Description string
	Price       int    `gorm:"not null"`
	Category    string `gorm:"not null"`
	Condition   string `gorm:"not null"`
	Status      string `gorm:"default:'active'"`
	ViewCount   int    `gorm:"default:0"`
	CreatedAt   time.Time
}

func main() {
	// Connect to database
	dsn := "host=localhost user=ecomate password=ecomate_password dbname=ecomate_db port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	fmt.Println("🌱 Seeding demo data for Automate...")

	// Create demo users
	users := []User{
		{
			Email:        "demo.seller@automate.com",
			Username:     "EcoSeller",
			PasswordHash: "$2a$10$lxgFaUxK7nvh0TyoOiEFJ.UBvorIGTDZFxiRQXUBB1TJMJ8dyFida", // bcrypt hash for "password123"
			Role:         "user",
		},
		{
			Email:        "demo.buyer@automate.com",
			Username:     "GreenBuyer",
			PasswordHash: "$2a$10$lxgFaUxK7nvh0TyoOiEFJ.UBvorIGTDZFxiRQXUBB1TJMJ8dyFida",
			Role:         "user",
		},
		{
			Email:        "admin@automate.com",
			Username:     "AdminUser",
			PasswordHash: "$2a$10$lxgFaUxK7nvh0TyoOiEFJ.UBvorIGTDZFxiRQXUBB1TJMJ8dyFida",
			Role:         "admin",
		},
	}

	for _, user := range users {
		if err := db.FirstOrCreate(&user, User{Email: user.Email}).Error; err != nil {
			log.Printf("Failed to create user %s: %v", user.Email, err)
		} else {
			fmt.Printf("✓ Created user: %s\n", user.Username)
		}
	}

	// Get seller ID
	var seller User
	db.Where("email = ?", "demo.seller@automate.com").First(&seller)

	// Create demo products
	products := []Product{
		{
			SellerID:    seller.ID,
			Title:       "ヴィンテージ カメラ Canon AE-1",
			Description: "1970年代の名機。完動品、レンズ付き。フィルムカメラ愛好家に最適。",
			Price:       25000,
			Category:    "electronics",
			Condition:   "used",
			Status:      "active",
			ViewCount:   156,
		},
		{
			SellerID:    seller.ID,
			Title:       "北欧デザイン 木製チェア",
			Description: "ハンス・ウェグナー風のミッドセンチュリーチェア。美品。",
			Price:       35000,
			Category:    "furniture",
			Condition:   "like_new",
			Status:      "active",
			ViewCount:   243,
		},
		{
			SellerID:    seller.ID,
			Title:       "無印良品 ダウンジャケット メンズL",
			Description: "昨シーズン購入、数回着用のみ。クリーニング済み。",
			Price:       8000,
			Category:    "fashion",
			Condition:   "like_new",
			Status:      "active",
			ViewCount:   89,
		},
		{
			SellerID:    seller.ID,
			Title:       "Kindle Paperwhite 第10世代 8GB",
			Description: "2020年モデル、広告なし。保護ケース付き。",
			Price:       12000,
			Category:    "electronics",
			Condition:   "like_new",
			Status:      "active",
			ViewCount:   412,
		},
		{
			SellerID:    seller.ID,
			Title:       "STAUB ココット 20cm チェリーレッド",
			Description: "ストウブの鋳物ホーロー鍋。数回使用、美品。",
			Price:       18000,
			Category:    "home",
			Condition:   "like_new",
			Status:      "active",
			ViewCount:   178,
		},
	}

	for _, product := range products {
		if err := db.Create(&product).Error; err != nil {
			log.Printf("Failed to create product %s: %v", product.Title, err)
		} else {
			fmt.Printf("✓ Created product: %s (¥%d)\n", product.Title, product.Price)
		}
	}

	fmt.Println("\n✨ Demo data seeding completed!")
	fmt.Println("\nDemo Accounts:")
	fmt.Println("  Seller: demo.seller@automate.com / password123")
	fmt.Println("  Buyer:  demo.buyer@automate.com / password123")
	fmt.Println("  Admin:  admin@automate.com / password123")
	fmt.Println("\nUse these accounts for your demo presentation!")
}
