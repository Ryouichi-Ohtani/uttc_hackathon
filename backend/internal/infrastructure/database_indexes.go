package infrastructure

import (
	"gorm.io/gorm"
)

// AddIndexes adds performance indexes to the database
func AddIndexes(db *gorm.DB) error {
	// Products indexes
	db.Exec("CREATE INDEX IF NOT EXISTS idx_products_status_created ON products(status, created_at DESC)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_products_category_price ON products(category, price)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_products_seller_status ON products(seller_id, status)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_products_co2_impact ON products(co2_impact_kg DESC)")

	// Purchases indexes
	db.Exec("CREATE INDEX IF NOT EXISTS idx_purchases_buyer_created ON purchases(buyer_id, created_at DESC)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_purchases_seller_created ON purchases(seller_id, created_at DESC)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_purchases_status ON purchases(status)")

	// Messages indexes
	db.Exec("CREATE INDEX IF NOT EXISTS idx_messages_conversation_created ON messages(conversation_id, created_at DESC)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_messages_sender ON messages(sender_id)")

	// Favorites indexes
	db.Exec("CREATE INDEX IF NOT EXISTS idx_favorites_user_created ON favorites(user_id, created_at DESC)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_favorites_product_user ON favorites(product_id, user_id)")

	// User events indexes for analytics
	db.Exec("CREATE INDEX IF NOT EXISTS idx_user_events_user_type ON user_events(user_id, event_type, created_at DESC)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_user_events_product_type ON user_events(product_id, event_type, created_at DESC)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_user_events_created ON user_events(created_at DESC)")

	// Offers indexes
	db.Exec("CREATE INDEX IF NOT EXISTS idx_offers_product_status ON offers(product_id, status)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_offers_buyer_status ON offers(buyer_id, status, created_at DESC)")

	// Notifications indexes
	db.Exec("CREATE INDEX IF NOT EXISTS idx_notifications_user_read ON notifications(user_id, is_read, created_at DESC)")

	// Reviews indexes
	db.Exec("CREATE INDEX IF NOT EXISTS idx_reviews_product_created ON reviews(product_id, created_at DESC)")

	// Full-text search index for products
	db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_products_search
		ON products USING gin(to_tsvector('english', title || ' ' || COALESCE(description, '')))
	`)

	return nil
}
