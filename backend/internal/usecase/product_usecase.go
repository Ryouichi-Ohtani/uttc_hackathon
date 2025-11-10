package usecase

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"

	"github.com/google/uuid"
	"github.com/yourusername/ecomate/backend/internal/domain"
	"github.com/yourusername/ecomate/backend/internal/infrastructure"
	"gorm.io/gorm"
)

type ProductUseCase struct {
	productRepo domain.ProductRepository
	aiClient    *infrastructure.AIClient
}

func NewProductUseCase(productRepo domain.ProductRepository, aiClient *infrastructure.AIClient) *ProductUseCase {
	return &ProductUseCase{
		productRepo: productRepo,
		aiClient:    aiClient,
	}
}

func (uc *ProductUseCase) CreateProduct(
	ctx context.Context,
	sellerID uuid.UUID,
	req *domain.CreateProductRequest,
	imageFiles []*multipart.FileHeader,
) (*domain.Product, error) {
	product := &domain.Product{
		SellerID:    sellerID,
		Title:       req.Title,
		Description: req.Description,
		Price:       req.Price,
		Category:    req.Category,
		Condition:   req.Condition,
		WeightKg:    req.WeightKg,
		Status:      domain.StatusActive,
	}

	// Read image files
	var imageBytes [][]byte
	for _, fileHeader := range imageFiles {
		file, err := fileHeader.Open()
		if err != nil {
			return nil, fmt.Errorf("failed to open image file: %w", err)
		}
		defer file.Close()

		data, err := io.ReadAll(file)
		if err != nil {
			return nil, fmt.Errorf("failed to read image file: %w", err)
		}
		imageBytes = append(imageBytes, data)
	}

	// Use AI assistance if requested
	if req.UseAIAssistance && uc.aiClient != nil {
		analysisReq := &infrastructure.ProductAnalysisRequest{
			Images:                  imageBytes,
			Title:                   req.Title,
			UserProvidedDescription: req.Description,
			Category:                req.Category,
		}

		analysisResp, err := uc.aiClient.AnalyzeProduct(ctx, analysisReq)
		if err != nil {
			return nil, fmt.Errorf("AI analysis failed: %w", err)
		}

		// Check for inappropriate content
		if analysisResp.IsInappropriate {
			return nil, errors.New("inappropriate content detected: " + analysisResp.InappropriateReason)
		}

		// Use AI-generated data
		if req.Description == "" {
			product.Description = analysisResp.GeneratedDescription
			product.AIGeneratedDescription = analysisResp.GeneratedDescription
		}
		product.AISuggestedPrice = analysisResp.SuggestedPrice
		product.WeightKg = analysisResp.EstimatedWeightKg
		product.ManufacturerCountry = analysisResp.ManufacturerCountry
		product.EstimatedManufacturingYear = analysisResp.EstimatedManufacturingYear
	}

	// Create product images
	// In production, upload to GCS and get CDN URLs
	for i := range imageBytes {
		// Mock: save to local storage or upload to cloud
		imageURL := fmt.Sprintf("/uploads/products/%s_%d.jpg", uuid.New().String(), i)
		cdnURL := imageURL // In production, this would be CDN URL

		image := domain.ProductImage{
			ImageURL:     imageURL,
			CDNURL:       cdnURL,
			DisplayOrder: i,
			IsPrimary:    i == 0,
		}
		product.Images = append(product.Images, image)
	}

	// Save to database
	if err := uc.productRepo.Create(product); err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	return product, nil
}

func (uc *ProductUseCase) GetProduct(id uuid.UUID) (*domain.Product, error) {
	product, err := uc.productRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}

	// Increment view count asynchronously
	go uc.productRepo.IncrementViewCount(id)

	return product, nil
}

func (uc *ProductUseCase) ListProducts(filters *domain.ProductFilters) ([]*domain.Product, *domain.PaginationResponse, error) {
	return uc.productRepo.List(filters)
}

func (uc *ProductUseCase) GetCO2Comparison(product *domain.Product) interface{} {
	// CO2 comparison feature has been removed
	return nil
}

// List is a wrapper for ListProducts for compatibility
func (uc *ProductUseCase) List(filters domain.ProductFilters) ([]*domain.Product, int64, error) {
	products, pagination, err := uc.ListProducts(&filters)
	if err != nil {
		return nil, 0, err
	}
	return products, int64(pagination.Total), nil
}

// GetByID is a wrapper for GetProduct for compatibility
func (uc *ProductUseCase) GetByID(id uuid.UUID) (*domain.Product, error) {
	return uc.GetProduct(id)
}

// Create is a simplified wrapper for CreateProduct
func (uc *ProductUseCase) Create(sellerID uuid.UUID, title, description, category string, price int, condition string, imageBytes [][]byte) (*domain.Product, error) {
	product := &domain.Product{
		SellerID:    sellerID,
		Title:       title,
		Description: description,
		Price:       int(price),
		Category:    category,
		Condition:   domain.ProductCondition(condition),
		Status:      domain.StatusActive,
	}

	// Create product images from bytes
	for i := range imageBytes {
		imageURL := fmt.Sprintf("/uploads/products/%s_%d.jpg", uuid.New().String(), i)
		image := domain.ProductImage{
			ImageURL:     imageURL,
			CDNURL:       imageURL,
			DisplayOrder: i,
			IsPrimary:    i == 0,
		}
		product.Images = append(product.Images, image)
	}

	if err := uc.productRepo.Create(product); err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	return product, nil
}

// AnswerQuestion uses AI to answer questions about a product
func (uc *ProductUseCase) AnswerQuestion(productID uuid.UUID, question string) (string, error) {
	product, err := uc.productRepo.FindByID(productID)
	if err != nil {
		return "", errors.New("product not found")
	}

	if uc.aiClient == nil {
		return "AI service is currently unavailable. Please contact the seller directly for more information.", nil
	}

	answer, err := uc.aiClient.AnswerProductQuestion(
		context.Background(),
		product.Title,
		product.Description,
		question,
	)
	if err != nil {
		return "", fmt.Errorf("failed to get AI answer: %w", err)
	}

	return answer, nil
}

// Update updates a product
func (uc *ProductUseCase) Update(productID, sellerID uuid.UUID, updates map[string]interface{}) (*domain.Product, error) {
	// Get existing product
	product, err := uc.productRepo.FindByID(productID)
	if err != nil {
		return nil, errors.New("product not found")
	}

	// Check ownership
	if product.SellerID != sellerID {
		return nil, errors.New("unauthorized: not product owner")
	}

	// Cannot update sold or deleted products
	if product.Status == domain.StatusSold || product.Status == domain.StatusDeleted {
		return nil, errors.New("cannot update sold or deleted products")
	}

	// Apply updates
	if title, ok := updates["title"].(string); ok && title != "" {
		product.Title = title
	}
	if description, ok := updates["description"].(string); ok {
		product.Description = description
	}
	if price, ok := updates["price"].(int); ok && price > 0 {
		product.Price = price
	}
	if category, ok := updates["category"].(string); ok && category != "" {
		product.Category = category
	}
	if condition, ok := updates["condition"].(string); ok && condition != "" {
		product.Condition = domain.ProductCondition(condition)
	}

	// Save to database
	if err := uc.productRepo.Update(product); err != nil {
		return nil, fmt.Errorf("failed to update product: %w", err)
	}

	return product, nil
}

// Delete deletes a product (soft delete)
func (uc *ProductUseCase) Delete(productID, sellerID uuid.UUID) error {
	// Get existing product
	product, err := uc.productRepo.FindByID(productID)
	if err != nil {
		return errors.New("product not found")
	}

	// Check ownership
	if product.SellerID != sellerID {
		return errors.New("unauthorized: not product owner")
	}

	// Cannot delete sold products
	if product.Status == domain.StatusSold {
		return errors.New("cannot delete sold products")
	}

	// Delete from database (soft delete)
	if err := uc.productRepo.Delete(productID); err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	return nil
}
