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
	List(page, limit int) ([]*domain.Purchase, *domain.PaginationResponse, error)
	CompletePurchase(id uuid.UUID, userID uuid.UUID) error
	UpdatePurchaseStatus(id uuid.UUID, status domain.PurchaseStatus) error
	GetShippingLabel(purchaseID uuid.UUID, userID uuid.UUID) (*domain.ShippingLabel, error)
	GenerateShippingLabel(purchaseID uuid.UUID, userID uuid.UUID) (*domain.ShippingLabel, error)
}

type purchaseUseCase struct {
	purchaseRepo      domain.PurchaseRepository
	productRepo       domain.ProductRepository
	userRepo          domain.UserRepository
	shippingLabelRepo domain.ShippingLabelRepository
}

func NewPurchaseUseCase(
	purchaseRepo domain.PurchaseRepository,
	productRepo domain.ProductRepository,
	userRepo domain.UserRepository,
	shippingLabelRepo domain.ShippingLabelRepository,
) PurchaseUseCase {
	return &purchaseUseCase{
		purchaseRepo:      purchaseRepo,
		productRepo:       productRepo,
		userRepo:          userRepo,
		shippingLabelRepo: shippingLabelRepo,
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

	// Get buyer information
	buyer, err := u.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("buyer not found")
	}

	// Create purchase with delivery information
	purchase := &domain.Purchase{
		ProductID:            product.ID,
		BuyerID:              userID,
		SellerID:             product.SellerID,
		Price:                product.Price,
		Status:               domain.PurchaseStatusPending,
		ShippingAddress:      req.ShippingAddress,
		PaymentMethod:        req.PaymentMethod,
		DeliveryDate:         req.DeliveryDate,
		DeliveryTimeSlot:     req.DeliveryTimeSlot,
		UseRegisteredAddress: req.UseRegisteredAddress,
		CreatedAt:            time.Now(),
	}

	// Set recipient information
	if req.UseRegisteredAddress {
		// Use buyer's registered address
		purchase.RecipientName = buyer.DisplayName
		purchase.RecipientPhoneNumber = buyer.PhoneNumber
		purchase.RecipientPostalCode = buyer.PostalCode
		purchase.RecipientPrefecture = buyer.Prefecture
		purchase.RecipientCity = buyer.City
		purchase.RecipientAddressLine1 = buyer.AddressLine1
		purchase.RecipientAddressLine2 = buyer.AddressLine2
	} else {
		// Use provided address
		purchase.RecipientName = req.RecipientName
		purchase.RecipientPhoneNumber = req.RecipientPhoneNumber
		purchase.RecipientPostalCode = req.RecipientPostalCode
		purchase.RecipientPrefecture = req.RecipientPrefecture
		purchase.RecipientCity = req.RecipientCity
		purchase.RecipientAddressLine1 = req.RecipientAddressLine1
		purchase.RecipientAddressLine2 = req.RecipientAddressLine2
	}

	if err := u.purchaseRepo.Create(purchase); err != nil {
		return nil, err
	}

	// Update product status to sold
	product.Status = domain.StatusSold
	if err := u.productRepo.Update(product); err != nil {
		return nil, err
	}

	// Automatically generate shipping label for the seller
	_, err = u.GenerateShippingLabel(purchase.ID, product.SellerID)
	if err != nil {
		// Log error but don't fail the purchase
		// The label can be generated later
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

func (u *purchaseUseCase) List(page, limit int) ([]*domain.Purchase, *domain.PaginationResponse, error) {
	return u.purchaseRepo.List(page, limit)
}

func (u *purchaseUseCase) UpdatePurchaseStatus(id uuid.UUID, status domain.PurchaseStatus) error {
	return u.purchaseRepo.UpdateStatus(id, status)
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

	now := time.Now()
	purchase.CompletedAt = &now

	return nil
}

func (u *purchaseUseCase) GetShippingLabel(purchaseID uuid.UUID, userID uuid.UUID) (*domain.ShippingLabel, error) {
	// Get purchase to verify user is seller
	purchase, err := u.purchaseRepo.FindByID(purchaseID)
	if err != nil {
		return nil, errors.New("purchase not found")
	}

	// Only seller can view shipping label
	if purchase.SellerID != userID {
		return nil, errors.New("only seller can view shipping label")
	}

	return u.shippingLabelRepo.FindByPurchaseID(purchaseID)
}

func (u *purchaseUseCase) GenerateShippingLabel(purchaseID uuid.UUID, userID uuid.UUID) (*domain.ShippingLabel, error) {
	// Get purchase with all relations
	purchase, err := u.purchaseRepo.FindByID(purchaseID)
	if err != nil {
		return nil, errors.New("purchase not found")
	}

	// Only seller can generate shipping label
	if purchase.SellerID != userID {
		return nil, errors.New("only seller can generate shipping label")
	}

	// Get seller information
	seller, err := u.userRepo.FindByID(purchase.SellerID)
	if err != nil {
		return nil, errors.New("seller not found")
	}

	// Check if label already exists
	existingLabel, _ := u.shippingLabelRepo.FindByPurchaseID(purchaseID)
	if existingLabel != nil {
		return existingLabel, nil
	}

	// Determine package size based on product category and weight
	packageSize := u.determinePackageSize(purchase.Product)

	// Determine carrier (default to yamato)
	carrier := "yamato"

	// Create shipping label
	label := &domain.ShippingLabel{
		PurchaseID: purchaseID,
		// Sender (Seller) information
		SenderName:         seller.DisplayName,
		SenderPostalCode:   seller.PostalCode,
		SenderPrefecture:   seller.Prefecture,
		SenderCity:         seller.City,
		SenderAddressLine1: seller.AddressLine1,
		SenderAddressLine2: seller.AddressLine2,
		SenderPhoneNumber:  seller.PhoneNumber,
		// Recipient (Buyer) information
		RecipientName:         purchase.RecipientName,
		RecipientPostalCode:   purchase.RecipientPostalCode,
		RecipientPrefecture:   purchase.RecipientPrefecture,
		RecipientCity:         purchase.RecipientCity,
		RecipientAddressLine1: purchase.RecipientAddressLine1,
		RecipientAddressLine2: purchase.RecipientAddressLine2,
		RecipientPhoneNumber:  purchase.RecipientPhoneNumber,
		// Delivery details
		DeliveryDate:     purchase.DeliveryDate,
		DeliveryTimeSlot: purchase.DeliveryTimeSlot,
		ProductName:      purchase.Product.Title,
		PackageSize:      packageSize,
		Weight:           purchase.Product.WeightKg,
		Carrier:          carrier,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	if err := u.shippingLabelRepo.Create(label); err != nil {
		return nil, err
	}

	return u.shippingLabelRepo.FindByPurchaseID(purchaseID)
}

func (u *purchaseUseCase) determinePackageSize(product *domain.Product) string {
	if product == nil {
		return "80"
	}

	// Determine size based on weight
	weight := product.WeightKg
	if weight == 0 {
		weight = 1.0 // Default weight
	}

	switch {
	case weight <= 2:
		return "60"
	case weight <= 5:
		return "80"
	case weight <= 10:
		return "100"
	case weight <= 15:
		return "120"
	case weight <= 20:
		return "140"
	default:
		return "160"
	}
}
