package infrastructure

import (
	"github.com/google/uuid"
	"github.com/yourusername/ecomate/backend/internal/domain"
	"gorm.io/gorm"
)

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) domain.ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(product *domain.Product) error {
	return r.db.Create(product).Error
}

func (r *productRepository) FindByID(id uuid.UUID) (*domain.Product, error) {
	var product domain.Product
	if err := r.db.Preload("Seller").Preload("Images").Where("id = ?", id).First(&product).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) List(filters *domain.ProductFilters) ([]*domain.Product, *domain.PaginationResponse, error) {
	var products []*domain.Product
	var total int64

	query := r.db.Model(&domain.Product{}).
		Preload("Seller").
		Preload("Images").
		Where("status = ?", domain.StatusActive)

	// Apply filters
	if filters.Category != "" {
		query = query.Where("category = ?", filters.Category)
	}
	if filters.MinPrice > 0 {
		query = query.Where("price >= ?", filters.MinPrice)
	}
	if filters.MaxPrice > 0 {
		query = query.Where("price <= ?", filters.MaxPrice)
	}
	if filters.Condition != "" {
		query = query.Where("condition = ?", filters.Condition)
	}
	if filters.Search != "" {
		query = query.Where(
			"to_tsvector('english', title || ' ' || COALESCE(description, '')) @@ plainto_tsquery(?)",
			filters.Search,
		)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, nil, err
	}

	// Apply sorting
	switch filters.Sort {
	case "price_asc":
		query = query.Order("price ASC")
	case "price_desc":
		query = query.Order("price DESC")
	case "eco_impact_desc":
		query = query.Order("co2_impact_kg DESC")
	default:
		query = query.Order("created_at DESC")
	}

	// Apply pagination
	page := filters.Page
	if page < 1 {
		page = 1
	}
	limit := filters.Limit
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Find(&products).Error; err != nil {
		return nil, nil, err
	}

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	pagination := &domain.PaginationResponse{
		Page:       page,
		Limit:      limit,
		Total:      int(total),
		TotalPages: totalPages,
	}

	return products, pagination, nil
}

func (r *productRepository) Update(product *domain.Product) error {
	return r.db.Save(product).Error
}

func (r *productRepository) Delete(id uuid.UUID) error {
	return r.db.Model(&domain.Product{}).
		Where("id = ?", id).
		Update("status", domain.StatusDeleted).
		Error
}

func (r *productRepository) IncrementViewCount(id uuid.UUID) error {
	return r.db.Model(&domain.Product{}).
		Where("id = ?", id).
		Update("view_count", gorm.Expr("view_count + 1")).
		Error
}
