package usecase_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/yourusername/ecomate/backend/internal/domain"
	"github.com/yourusername/ecomate/backend/internal/usecase"
)

// Mock OfferRepository
type MockOfferRepository struct {
	mock.Mock
}

func (m *MockOfferRepository) Create(offer *domain.Offer) error {
	args := m.Called(offer)
	return args.Error(0)
}

func (m *MockOfferRepository) FindByID(id uuid.UUID) (*domain.Offer, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Offer), args.Error(1)
}

func (m *MockOfferRepository) FindByProductID(productID uuid.UUID) ([]*domain.Offer, error) {
	args := m.Called(productID)
	return args.Get(0).([]*domain.Offer), args.Error(1)
}

func (m *MockOfferRepository) FindByBuyerID(buyerID uuid.UUID) ([]*domain.Offer, error) {
	args := m.Called(buyerID)
	return args.Get(0).([]*domain.Offer), args.Error(1)
}

func (m *MockOfferRepository) FindBySellerID(sellerID uuid.UUID) ([]*domain.Offer, error) {
	args := m.Called(sellerID)
	return args.Get(0).([]*domain.Offer), args.Error(1)
}

func (m *MockOfferRepository) Update(offer *domain.Offer) error {
	args := m.Called(offer)
	return args.Error(0)
}

func TestOfferUseCase_CreateOffer_Success(t *testing.T) {
	// Arrange
	mockOfferRepo := new(MockOfferRepository)
	mockProductRepo := new(MockProductRepository)
	useCase := usecase.NewOfferUseCase(mockOfferRepo, mockProductRepo)

	buyerID := uuid.New()
	sellerID := uuid.New()
	productID := uuid.New()
	product := &domain.Product{
		ID:       productID,
		SellerID: sellerID,
		Price:    1000,
		Status:   domain.StatusActive,
	}

	offerPrice := 800

	mockProductRepo.On("FindByID", productID).Return(product, nil)
	mockOfferRepo.On("Create", mock.AnythingOfType("*domain.Offer")).Return(nil)
	mockOfferRepo.On("FindByID", mock.AnythingOfType("uuid.UUID")).Return(&domain.Offer{
		ProductID:  productID,
		BuyerID:    buyerID,
		OfferPrice: offerPrice,
		Status:     domain.OfferStatusPending,
	}, nil)

	// Act
	result, err := useCase.CreateOffer(buyerID, productID, offerPrice, "Please consider my offer")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, offerPrice, result.OfferPrice)
	assert.Equal(t, domain.OfferStatusPending, result.Status)
	mockProductRepo.AssertExpectations(t)
	mockOfferRepo.AssertExpectations(t)
}

func TestOfferUseCase_CreateOffer_CannotOfferOnOwnProduct(t *testing.T) {
	// Arrange
	mockOfferRepo := new(MockOfferRepository)
	mockProductRepo := new(MockProductRepository)
	useCase := usecase.NewOfferUseCase(mockOfferRepo, mockProductRepo)

	userID := uuid.New()
	productID := uuid.New()
	product := &domain.Product{
		ID:       productID,
		SellerID: userID, // Same as buyer
		Price:    1000,
		Status:   domain.StatusActive,
	}

	mockProductRepo.On("FindByID", productID).Return(product, nil)

	// Act
	result, err := useCase.CreateOffer(userID, productID, 800, "")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "cannot make offer on your own product")
	mockProductRepo.AssertExpectations(t)
}

func TestOfferUseCase_CreateOffer_PriceTooHigh(t *testing.T) {
	// Arrange
	mockOfferRepo := new(MockOfferRepository)
	mockProductRepo := new(MockProductRepository)
	useCase := usecase.NewOfferUseCase(mockOfferRepo, mockProductRepo)

	buyerID := uuid.New()
	sellerID := uuid.New()
	productID := uuid.New()
	product := &domain.Product{
		ID:       productID,
		SellerID: sellerID,
		Price:    1000,
		Status:   domain.StatusActive,
	}

	offerPrice := 1100 // Higher than current price

	mockProductRepo.On("FindByID", productID).Return(product, nil)

	// Act
	result, err := useCase.CreateOffer(buyerID, productID, offerPrice, "")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "offer price must be less than current price")
	mockProductRepo.AssertExpectations(t)
}
