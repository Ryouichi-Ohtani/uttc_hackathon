package usecase_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/yourusername/ecomate/backend/internal/domain"
	"github.com/yourusername/ecomate/backend/internal/usecase"
)

// Mock ProductRepository
type MockProductRepository struct {
	mock.Mock
}

func (m *MockProductRepository) Create(product *domain.Product) error {
	args := m.Called(product)
	return args.Error(0)
}

func (m *MockProductRepository) FindByID(id uuid.UUID) (*domain.Product, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Product), args.Error(1)
}

func (m *MockProductRepository) List(filters *domain.ProductFilters) ([]*domain.Product, *domain.PaginationResponse, error) {
	args := m.Called(filters)
	return args.Get(0).([]*domain.Product), args.Get(1).(*domain.PaginationResponse), args.Error(2)
}

func (m *MockProductRepository) Update(product *domain.Product) error {
	args := m.Called(product)
	return args.Error(0)
}

func (m *MockProductRepository) Delete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockProductRepository) IncrementViewCount(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestProductUseCase_GetByID(t *testing.T) {
	// Arrange
	mockRepo := new(MockProductRepository)
	useCase := usecase.NewProductUseCase(mockRepo, nil)

	productID := uuid.New()
	expectedProduct := &domain.Product{
		ID:          productID,
		Title:       "Test Product",
		Description: "Test Description",
		Price:       1000,
		Status:      domain.StatusActive,
	}

	mockRepo.On("FindByID", productID).Return(expectedProduct, nil)
	mockRepo.On("IncrementViewCount", productID).Return(nil)

	// Act
	result, err := useCase.GetByID(productID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedProduct.ID, result.ID)
	assert.Equal(t, expectedProduct.Title, result.Title)
	mockRepo.AssertExpectations(t)
}

func TestProductUseCase_Delete_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockProductRepository)
	useCase := usecase.NewProductUseCase(mockRepo, nil)

	productID := uuid.New()
	sellerID := uuid.New()
	product := &domain.Product{
		ID:       productID,
		SellerID: sellerID,
		Status:   domain.StatusActive,
	}

	mockRepo.On("FindByID", productID).Return(product, nil)
	mockRepo.On("Delete", productID).Return(nil)

	// Act
	err := useCase.Delete(productID, sellerID)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestProductUseCase_Delete_Unauthorized(t *testing.T) {
	// Arrange
	mockRepo := new(MockProductRepository)
	useCase := usecase.NewProductUseCase(mockRepo, nil)

	productID := uuid.New()
	sellerID := uuid.New()
	differentUserID := uuid.New()
	product := &domain.Product{
		ID:       productID,
		SellerID: sellerID,
		Status:   domain.StatusActive,
	}

	mockRepo.On("FindByID", productID).Return(product, nil)

	// Act
	err := useCase.Delete(productID, differentUserID)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unauthorized")
	mockRepo.AssertExpectations(t)
}
