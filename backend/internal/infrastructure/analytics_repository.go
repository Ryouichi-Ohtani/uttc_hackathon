package infrastructure

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/yourusername/ecomate/backend/internal/domain"
	"gorm.io/gorm"
)

type analyticsRepository struct {
	db *gorm.DB
}

func NewAnalyticsRepository(db *gorm.DB) domain.AnalyticsRepository {
	return &analyticsRepository{db: db}
}

func (r *analyticsRepository) TrackEvent(event *domain.UserEvent) error {
	return r.db.Create(event).Error
}

func (r *analyticsRepository) GetUserEvents(userID uuid.UUID, eventType domain.EventType, limit int) ([]*domain.UserEvent, error) {
	var events []*domain.UserEvent
	query := r.db.Where("user_id = ?", userID)

	if eventType != "" {
		query = query.Where("event_type = ?", eventType)
	}

	err := query.Order("created_at DESC").Limit(limit).Find(&events).Error
	return events, err
}

func (r *analyticsRepository) GetPopularProducts(limit int, since time.Time) ([]*domain.ProductAnalytics, error) {
	var results []struct {
		ProductID  uuid.UUID
		ViewCount  int
		LikeCount  int
		OfferCount int
	}

	err := r.db.Raw(`
		SELECT
			product_id,
			COUNT(CASE WHEN event_type = 'product_view' THEN 1 END) as view_count,
			COUNT(CASE WHEN event_type = 'product_like' THEN 1 END) as like_count,
			COUNT(CASE WHEN event_type = 'product_offer' THEN 1 END) as offer_count
		FROM user_events
		WHERE product_id IS NOT NULL AND created_at >= ?
		GROUP BY product_id
		ORDER BY view_count DESC, like_count DESC
		LIMIT ?
	`, since, limit).Scan(&results).Error

	if err != nil {
		return nil, err
	}

	analytics := make([]*domain.ProductAnalytics, 0, len(results))
	for _, result := range results {
		var product domain.Product
		if err := r.db.Preload("Images").First(&product, "id = ?", result.ProductID).Error; err != nil {
			continue
		}

		analytics = append(analytics, &domain.ProductAnalytics{
			ProductID:  result.ProductID,
			Product:    &product,
			ViewCount:  result.ViewCount,
			LikeCount:  result.LikeCount,
			OfferCount: result.OfferCount,
		})
	}

	return analytics, nil
}

func (r *analyticsRepository) GetUserBehaviorSummary(userID uuid.UUID) (*domain.UserBehaviorSummary, error) {
	summary := &domain.UserBehaviorSummary{}

	// Count events by type
	var counts []struct {
		EventType domain.EventType
		Count     int
	}

	err := r.db.Model(&domain.UserEvent{}).
		Select("event_type, COUNT(*) as count").
		Where("user_id = ?", userID).
		Group("event_type").
		Find(&counts).Error

	if err != nil {
		return nil, err
	}

	for _, count := range counts {
		switch count.EventType {
		case domain.EventProductView:
			summary.TotalViews = count.Count
		case domain.EventProductSearch:
			summary.TotalSearches = count.Count
		case domain.EventProductLike:
			summary.TotalLikes = count.Count
		case domain.EventProductOffer:
			summary.TotalOffers = count.Count
		}
	}

	// Get favorite categories
	var categoryResults []struct {
		Category string
		Count    int
	}

	err = r.db.Raw(`
		SELECT p.category, COUNT(*) as count
		FROM user_events ue
		JOIN products p ON p.id = ue.product_id
		WHERE ue.user_id = ? AND ue.event_type IN ('product_view', 'product_like')
		GROUP BY p.category
		ORDER BY count DESC
		LIMIT 5
	`, userID).Scan(&categoryResults).Error

	if err == nil {
		for _, result := range categoryResults {
			summary.FavoriteCategories = append(summary.FavoriteCategories, domain.CategoryCount{
				Category: result.Category,
				Count:    result.Count,
			})
		}
	}

	// Get recently viewed products
	var productIDs []uuid.UUID
	err = r.db.Model(&domain.UserEvent{}).
		Select("DISTINCT product_id").
		Where("user_id = ? AND event_type = ? AND product_id IS NOT NULL", userID, domain.EventProductView).
		Order("created_at DESC").
		Limit(10).
		Pluck("product_id", &productIDs).Error

	if err == nil && len(productIDs) > 0 {
		var products []*domain.Product
		r.db.Preload("Images").Where("id IN ?", productIDs).Find(&products)
		summary.RecentlyViewed = products
	}

	return summary, nil
}

func (r *analyticsRepository) GetSearchKeywords(limit int, since time.Time) ([]*domain.SearchKeyword, error) {
	var results []*domain.SearchKeyword

	rows, err := r.db.Raw(`
		SELECT
			metadata->>'keyword' as keyword,
			COUNT(*) as count
		FROM user_events
		WHERE event_type = 'product_search'
		  AND created_at >= ?
		  AND metadata->>'keyword' IS NOT NULL
		GROUP BY metadata->>'keyword'
		ORDER BY count DESC
		LIMIT ?
	`, since, limit).Rows()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var keyword domain.SearchKeyword
		if err := rows.Scan(&keyword.Keyword, &keyword.Count); err != nil {
			continue
		}
		results = append(results, &keyword)
	}

	return results, nil
}

// Helper to track product view
func TrackProductView(db *gorm.DB, userID, productID uuid.UUID, ipAddress, userAgent string) {
	event := &domain.UserEvent{
		UserID:    userID,
		EventType: domain.EventProductView,
		ProductID: &productID,
		IPAddress: ipAddress,
		UserAgent: userAgent,
	}
	db.Create(event)
}

// Helper to track search
func TrackSearch(db *gorm.DB, userID uuid.UUID, keyword, ipAddress, userAgent string) {
	metadata, _ := json.Marshal(map[string]string{"keyword": keyword})
	event := &domain.UserEvent{
		UserID:    userID,
		EventType: domain.EventProductSearch,
		Metadata:  string(metadata),
		IPAddress: ipAddress,
		UserAgent: userAgent,
	}
	db.Create(event)
}
