package usecase

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/yourusername/ecomate/backend/internal/domain"
)

type AutoPurchaseUseCase interface {
	CreateWatch(userID uuid.UUID, req *domain.CreateAutoPurchaseWatchRequest) (*domain.AutoPurchaseWatch, error)
	GetUserWatches(userID uuid.UUID) ([]*domain.AutoPurchaseWatch, error)
	GetWatchByID(id uuid.UUID, userID uuid.UUID) (*domain.AutoPurchaseWatch, error)
	CancelWatch(id uuid.UUID, userID uuid.UUID) error
	AuthorizePayment(req *domain.PaymentAuthorizationRequest) (*domain.PaymentAuthorizationResponse, error)
	CheckAndExecuteAutoPurchases() (int, error)
	ExecuteAutoPurchase(watchID uuid.UUID) error
}

type autoPurchaseUseCase struct {
	autoPurchaseRepo domain.AutoPurchaseWatchRepository
	autoPurchaseLogRepo domain.AutoPurchaseLogRepository
	productRepo      domain.ProductRepository
	userRepo         domain.UserRepository
	purchaseRepo     domain.PurchaseRepository
	shippingLabelRepo domain.ShippingLabelRepository
	notificationRepo domain.NotificationRepository
}

func NewAutoPurchaseUseCase(
	autoPurchaseRepo domain.AutoPurchaseWatchRepository,
	autoPurchaseLogRepo domain.AutoPurchaseLogRepository,
	productRepo domain.ProductRepository,
	userRepo domain.UserRepository,
	purchaseRepo domain.PurchaseRepository,
	shippingLabelRepo domain.ShippingLabelRepository,
	notificationRepo domain.NotificationRepository,
) AutoPurchaseUseCase {
	return &autoPurchaseUseCase{
		autoPurchaseRepo:    autoPurchaseRepo,
		autoPurchaseLogRepo: autoPurchaseLogRepo,
		productRepo:         productRepo,
		userRepo:            userRepo,
		purchaseRepo:        purchaseRepo,
		shippingLabelRepo:   shippingLabelRepo,
		notificationRepo:    notificationRepo,
	}
}

func (u *autoPurchaseUseCase) CreateWatch(userID uuid.UUID, req *domain.CreateAutoPurchaseWatchRequest) (*domain.AutoPurchaseWatch, error) {
	// Validate product exists and is available
	product, err := u.productRepo.FindByID(req.ProductID)
	if err != nil {
		return nil, errors.New("product not found")
	}

	if product.Status != domain.StatusActive {
		return nil, errors.New("product is not available for auto-purchase")
	}

	// Can't auto-purchase own product
	if product.SellerID == userID {
		return nil, errors.New("cannot auto-purchase your own product")
	}

	// Validate max price is reasonable
	if req.MaxPrice <= 0 {
		return nil, errors.New("max price must be greater than 0")
	}

	// Get user information for delivery
	user, err := u.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Calculate expiration date
	expiresInDays := req.ExpiresInDays
	if expiresInDays == 0 {
		expiresInDays = 30 // Default 30 days
	}
	expiresAt := time.Now().AddDate(0, 0, expiresInDays)

	// Create auto-purchase watch
	watch := &domain.AutoPurchaseWatch{
		UserID:               userID,
		ProductID:            req.ProductID,
		MaxPrice:             req.MaxPrice,
		Status:               domain.AutoPurchaseWatchStatusActive,
		PaymentAuthorized:    true, // Already authorized via payment auth token
		PaymentMethodID:      req.PaymentMethodID,
		PaymentAuthToken:     req.PaymentAuthToken,
		UseRegisteredAddress: req.UseRegisteredAddress,
		ShippingAddress:      req.ShippingAddress,
		DeliveryTimeSlot:     req.DeliveryTimeSlot,
		ExpiresAt:            expiresAt,
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}

	// Set delivery information
	if req.UseRegisteredAddress {
		watch.RecipientName = user.DisplayName
		watch.RecipientPhoneNumber = user.PhoneNumber
		watch.RecipientPostalCode = user.PostalCode
		watch.RecipientPrefecture = user.Prefecture
		watch.RecipientCity = user.City
		watch.RecipientAddressLine1 = user.AddressLine1
		watch.RecipientAddressLine2 = user.AddressLine2
	} else {
		watch.RecipientName = req.RecipientName
		watch.RecipientPhoneNumber = req.RecipientPhoneNumber
		watch.RecipientPostalCode = req.RecipientPostalCode
		watch.RecipientPrefecture = req.RecipientPrefecture
		watch.RecipientCity = req.RecipientCity
		watch.RecipientAddressLine1 = req.RecipientAddressLine1
		watch.RecipientAddressLine2 = req.RecipientAddressLine2
	}

	if err := u.autoPurchaseRepo.Create(watch); err != nil {
		return nil, err
	}

	// Create notification
	u.notificationRepo.Create(&domain.Notification{
		UserID:  userID,
		Type:    "auto_purchase_watch_created",
		Title:   "自動購入監視を開始しました",
		Message: fmt.Sprintf("%s が ¥%d 以下になったら自動購入します", product.Title, req.MaxPrice),
	})

	return u.autoPurchaseRepo.FindByID(watch.ID)
}

func (u *autoPurchaseUseCase) GetUserWatches(userID uuid.UUID) ([]*domain.AutoPurchaseWatch, error) {
	return u.autoPurchaseRepo.FindByUserID(userID)
}

func (u *autoPurchaseUseCase) GetWatchByID(id uuid.UUID, userID uuid.UUID) (*domain.AutoPurchaseWatch, error) {
	watch, err := u.autoPurchaseRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Verify ownership
	if watch.UserID != userID {
		return nil, errors.New("unauthorized")
	}

	return watch, nil
}

func (u *autoPurchaseUseCase) CancelWatch(id uuid.UUID, userID uuid.UUID) error {
	watch, err := u.autoPurchaseRepo.FindByID(id)
	if err != nil {
		return errors.New("watch not found")
	}

	// Verify ownership
	if watch.UserID != userID {
		return errors.New("unauthorized")
	}

	// Update status to cancelled
	watch.Status = domain.AutoPurchaseWatchStatusCancelled
	watch.UpdatedAt = time.Now()

	if err := u.autoPurchaseRepo.Update(watch); err != nil {
		return err
	}

	// Create notification
	u.notificationRepo.Create(&domain.Notification{
		UserID:  userID,
		Type:    "auto_purchase_watch_cancelled",
		Title:   "自動購入監視をキャンセルしました",
		Message: fmt.Sprintf("%s の自動購入監視をキャンセルしました", watch.Product.Title),
	})

	return nil
}

// AuthorizePayment simulates payment authorization (like Google's pre-authorization)
// In production, this would integrate with Stripe, PayPal, or Japanese payment gateways
func (u *autoPurchaseUseCase) AuthorizePayment(req *domain.PaymentAuthorizationRequest) (*domain.PaymentAuthorizationResponse, error) {
	// Simulate payment validation
	if len(req.CardNumber) < 13 || len(req.CardNumber) > 19 {
		return &domain.PaymentAuthorizationResponse{
			Authorized: false,
			Message:    "Invalid card number",
		}, nil
	}

	if req.ExpiryYear < time.Now().Year() {
		return &domain.PaymentAuthorizationResponse{
			Authorized: false,
			Message:    "Card has expired",
		}, nil
	}

	if len(req.CVV) != 3 {
		return &domain.PaymentAuthorizationResponse{
			Authorized: false,
			Message:    "Invalid CVV",
		}, nil
	}

	// Generate mock payment method ID and auth token
	paymentMethodID := generateRandomToken(16)
	authToken := generateRandomToken(32)

	// In production, this would:
	// 1. Create a payment method with Stripe/PayPal
	// 2. Place a hold (authorization) on the card for the amount
	// 3. Return the payment method ID and authorization token

	return &domain.PaymentAuthorizationResponse{
		Authorized:      true,
		PaymentMethodID: paymentMethodID,
		AuthToken:       authToken,
		ExpiresAt:       time.Now().AddDate(0, 1, 0), // Expires in 1 month
		Message:         "Payment authorized successfully",
	}, nil
}

// CheckAndExecuteAutoPurchases is called by a background job to check prices and execute purchases
func (u *autoPurchaseUseCase) CheckAndExecuteAutoPurchases() (int, error) {
	// Get all active watches
	watches, err := u.autoPurchaseRepo.FindActiveWatches()
	if err != nil {
		return 0, err
	}

	executed := 0
	for _, watch := range watches {
		// Check if product is still available and price meets criteria
		product, err := u.productRepo.FindByID(watch.ProductID)
		if err != nil {
			continue
		}

		// Log price check
		now := time.Now()
		watch.LastCheckedAt = &now
		u.autoPurchaseRepo.Update(watch)

		u.autoPurchaseLogRepo.Create(&domain.AutoPurchaseLog{
			WatchID:      watch.ID,
			Action:       "price_check",
			CurrentPrice: product.Price,
			Success:      true,
			CreatedAt:    time.Now(),
		})

		// Check if product is available and price is acceptable
		if product.Status == domain.StatusActive && product.Price <= watch.MaxPrice {
			// Execute auto-purchase
			startTime := time.Now()
			err := u.ExecuteAutoPurchase(watch.ID)
			executionTime := int(time.Since(startTime).Milliseconds())

			if err == nil {
				executed++
				u.autoPurchaseLogRepo.Create(&domain.AutoPurchaseLog{
					WatchID:         watch.ID,
					Action:          "purchase_success",
					CurrentPrice:    product.Price,
					Success:         true,
					ExecutionTimeMs: executionTime,
					CreatedAt:       time.Now(),
				})
			} else {
				u.autoPurchaseLogRepo.Create(&domain.AutoPurchaseLog{
					WatchID:         watch.ID,
					Action:          "purchase_failed",
					CurrentPrice:    product.Price,
					Success:         false,
					ErrorMessage:    err.Error(),
					ExecutionTimeMs: executionTime,
					CreatedAt:       time.Now(),
				})
			}
		}
	}

	return executed, nil
}

func (u *autoPurchaseUseCase) ExecuteAutoPurchase(watchID uuid.UUID) error {
	watch, err := u.autoPurchaseRepo.FindByID(watchID)
	if err != nil {
		return errors.New("watch not found")
	}

	// Validate watch is still active
	if watch.Status != domain.AutoPurchaseWatchStatusActive {
		return errors.New("watch is not active")
	}

	// Get product
	product, err := u.productRepo.FindByID(watch.ProductID)
	if err != nil {
		return errors.New("product not found")
	}

	// Final validation
	if product.Status != domain.StatusActive {
		return errors.New("product is no longer available")
	}

	if product.Price > watch.MaxPrice {
		return errors.New("price exceeds max price")
	}

	// Create purchase
	purchase := &domain.Purchase{
		ProductID:              product.ID,
		BuyerID:                watch.UserID,
		SellerID:               product.SellerID,
		Price:                  product.Price,
		Status:                 domain.PurchaseStatusPending,
		PaymentMethod:          watch.PaymentMethodID,
		ShippingAddress:        watch.ShippingAddress,
		UseRegisteredAddress:   watch.UseRegisteredAddress,
		RecipientName:          watch.RecipientName,
		RecipientPhoneNumber:   watch.RecipientPhoneNumber,
		RecipientPostalCode:    watch.RecipientPostalCode,
		RecipientPrefecture:    watch.RecipientPrefecture,
		RecipientCity:          watch.RecipientCity,
		RecipientAddressLine1:  watch.RecipientAddressLine1,
		RecipientAddressLine2:  watch.RecipientAddressLine2,
		DeliveryTimeSlot:       watch.DeliveryTimeSlot,
		CreatedAt:              time.Now(),
	}

	if err := u.purchaseRepo.Create(purchase); err != nil {
		return err
	}

	// Update product status
	product.Status = domain.StatusSold
	if err := u.productRepo.Update(product); err != nil {
		return err
	}

	// Update watch status
	executedAt := time.Now()
	watch.Status = domain.AutoPurchaseWatchStatusExecuted
	watch.ExecutedAt = &executedAt
	watch.PurchaseID = &purchase.ID
	watch.UpdatedAt = time.Now()
	if err := u.autoPurchaseRepo.Update(watch); err != nil {
		return err
	}

	// Generate shipping label
	_, err = u.generateShippingLabelForAutoPurchase(purchase.ID, product.SellerID)
	if err != nil {
		// Log but don't fail
	}

	// Send notification to buyer
	u.notificationRepo.Create(&domain.Notification{
		UserID:  watch.UserID,
		Type:    "auto_purchase_executed",
		Title:   "🎉 自動購入が完了しました",
		Message: fmt.Sprintf("%s を ¥%d で自動購入しました！", product.Title, product.Price),
	})

	// Send notification to seller
	u.notificationRepo.Create(&domain.Notification{
		UserID:  product.SellerID,
		Type:    "product_sold",
		Title:   "商品が売れました",
		Message: fmt.Sprintf("%s が売れました！", product.Title),
	})

	return nil
}

func (u *autoPurchaseUseCase) generateShippingLabelForAutoPurchase(purchaseID uuid.UUID, sellerID uuid.UUID) (*domain.ShippingLabel, error) {
	purchase, err := u.purchaseRepo.FindByID(purchaseID)
	if err != nil {
		return nil, err
	}

	seller, err := u.userRepo.FindByID(sellerID)
	if err != nil {
		return nil, err
	}

	packageSize := determinePackageSize(purchase.Product)
	carrier := "yamato"

	label := &domain.ShippingLabel{
		PurchaseID:            purchaseID,
		SenderName:            seller.DisplayName,
		SenderPostalCode:      seller.PostalCode,
		SenderPrefecture:      seller.Prefecture,
		SenderCity:            seller.City,
		SenderAddressLine1:    seller.AddressLine1,
		SenderAddressLine2:    seller.AddressLine2,
		SenderPhoneNumber:     seller.PhoneNumber,
		RecipientName:         purchase.RecipientName,
		RecipientPostalCode:   purchase.RecipientPostalCode,
		RecipientPrefecture:   purchase.RecipientPrefecture,
		RecipientCity:         purchase.RecipientCity,
		RecipientAddressLine1: purchase.RecipientAddressLine1,
		RecipientAddressLine2: purchase.RecipientAddressLine2,
		RecipientPhoneNumber:  purchase.RecipientPhoneNumber,
		DeliveryDate:          purchase.DeliveryDate,
		DeliveryTimeSlot:      purchase.DeliveryTimeSlot,
		ProductName:           purchase.Product.Title,
		PackageSize:           packageSize,
		Weight:                purchase.Product.WeightKg,
		Carrier:               carrier,
		CreatedAt:             time.Now(),
		UpdatedAt:             time.Now(),
	}

	if err := u.shippingLabelRepo.Create(label); err != nil {
		return nil, err
	}

	return label, nil
}

func determinePackageSize(product *domain.Product) string {
	if product == nil {
		return "80"
	}

	weight := product.WeightKg
	if weight == 0 {
		weight = 1.0
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

func generateRandomToken(length int) string {
	bytes := make([]byte, length)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
