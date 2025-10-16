package usecase

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/yourusername/ecomate/backend/internal/domain"
)

type PurchaseUseCase interface {
	Create(userID uuid.UUID, req *domain.CreatePurchaseRequest) (*domain.Purchase, error)
	GetByID(id uuid.UUID) (*domain.Purchase, error)
	ListByUser(userID uuid.UUID, role string, page, limit int) ([]*domain.Purchase, *domain.PaginationResponse, error)
	CompletePurchase(id uuid.UUID, userID uuid.UUID) error
}

type purchaseUseCase struct {
	purchaseRepo domain.PurchaseRepository
	productRepo  domain.ProductRepository
	userRepo     domain.UserRepository
}

func NewPurchaseUseCase(
	purchaseRepo domain.PurchaseRepository,
	productRepo domain.ProductRepository,
	userRepo domain.UserRepository,
) PurchaseUseCase {
	return &purchaseUseCase{
		purchaseRepo: purchaseRepo,
		productRepo:  productRepo,
		userRepo:     userRepo,
	}
}

func (u *purchaseUseCase) Create(userID uuid.UUID, req *domain.CreatePurchaseRequest) (*domain.Purchase, error) {
	// Get product
	product, err := u.productRepo.FindByID(req.ProductID)
	if err != nil {
		return nil, errors.New("product not found")
	}

	// Check if product is available
	if product.Status != domain.StatusActive {
		return nil, errors.New("product is not available for purchase")
	}

	// Can't buy own product
	if product.SellerID == userID {
		return nil, errors.New("cannot purchase your own product")
	}

	// Create purchase
	purchase := &domain.Purchase{
		ProductID:       product.ID,
		BuyerID:         userID,
		SellerID:        product.SellerID,
		Price:           product.Price,
		CO2SavedKg:      product.CO2ImpactKg,
		Status:          domain.PurchaseStatusPending,
		ShippingAddress: req.ShippingAddress,
		PaymentMethod:   req.PaymentMethod,
		CreatedAt:       time.Now(),
	}

	if err := u.purchaseRepo.Create(purchase); err != nil {
		return nil, err
	}

	// Update product status to sold
	product.Status = domain.StatusSold
	if err := u.productRepo.Update(product); err != nil {
		return nil, err
	}

	// Reload with relations
	return u.purchaseRepo.FindByID(purchase.ID)
}

func (u *purchaseUseCase) GetByID(id uuid.UUID) (*domain.Purchase, error) {
	return u.purchaseRepo.FindByID(id)
}

func (u *purchaseUseCase) ListByUser(userID uuid.UUID, role string, page, limit int) ([]*domain.Purchase, *domain.PaginationResponse, error) {
	return u.purchaseRepo.FindByUser(userID, role, page, limit)
}

func (u *purchaseUseCase) CompletePurchase(id uuid.UUID, userID uuid.UUID) error {
	purchase, err := u.purchaseRepo.FindByID(id)
	if err != nil {
		return errors.New("purchase not found")
	}

	// Only seller can complete
	if purchase.SellerID != userID {
		return errors.New("only seller can complete the purchase")
	}

	if purchase.Status != domain.PurchaseStatusPending {
		return errors.New("purchase is not in pending status")
	}

	// Update status
	if err := u.purchaseRepo.UpdateStatus(id, domain.PurchaseStatusCompleted); err != nil {
		return err
	}

	// Update buyer's CO2 saved
	buyer, err := u.userRepo.FindByID(purchase.BuyerID)
	if err == nil {
		buyer.TotalCO2SavedKg += purchase.CO2SavedKg
		_ = u.userRepo.Update(buyer)
	}

	now := time.Now()
	purchase.CompletedAt = &now

	return nil
}
