package infrastructure

import (
	"github.com/google/uuid"
	"github.com/yourusername/ecomate/backend/internal/domain"
	"gorm.io/gorm"
)

type AutoPurchaseWatchRepositoryImpl struct {
	db *gorm.DB
}

func NewAutoPurchaseWatchRepository(db *gorm.DB) domain.AutoPurchaseWatchRepository {
	return &AutoPurchaseWatchRepositoryImpl{db: db}
}

func (r *AutoPurchaseWatchRepositoryImpl) Create(watch *domain.AutoPurchaseWatch) error {
	return r.db.Create(watch).Error
}

func (r *AutoPurchaseWatchRepositoryImpl) FindByID(id uuid.UUID) (*domain.AutoPurchaseWatch, error) {
	var watch domain.AutoPurchaseWatch
	err := r.db.Preload("User").
		Preload("Product").
		Preload("Product.Images").
		Preload("Purchase").
		Where("id = ?", id).
		First(&watch).Error
	if err != nil {
		return nil, err
	}
	return &watch, nil
}

func (r *AutoPurchaseWatchRepositoryImpl) FindByUserID(userID uuid.UUID) ([]*domain.AutoPurchaseWatch, error) {
	var watches []*domain.AutoPurchaseWatch
	err := r.db.Preload("Product").
		Preload("Product.Images").
		Preload("Purchase").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&watches).Error
	return watches, err
}

func (r *AutoPurchaseWatchRepositoryImpl) FindByProductID(productID uuid.UUID) ([]*domain.AutoPurchaseWatch, error) {
	var watches []*domain.AutoPurchaseWatch
	err := r.db.Preload("User").
		Where("product_id = ? AND status = ?", productID, domain.AutoPurchaseWatchStatusActive).
		Order("max_price DESC, created_at ASC"). // Higher max price and earlier creation gets priority
		Find(&watches).Error
	return watches, err
}

func (r *AutoPurchaseWatchRepositoryImpl) FindActiveWatches() ([]*domain.AutoPurchaseWatch, error) {
	var watches []*domain.AutoPurchaseWatch
	err := r.db.Preload("User").
		Preload("Product").
		Where("status = ? AND payment_authorized = ? AND expires_at > NOW()",
			domain.AutoPurchaseWatchStatusActive, true).
		Find(&watches).Error
	return watches, err
}

func (r *AutoPurchaseWatchRepositoryImpl) Update(watch *domain.AutoPurchaseWatch) error {
	return r.db.Save(watch).Error
}

func (r *AutoPurchaseWatchRepositoryImpl) Delete(id uuid.UUID) error {
	return r.db.Delete(&domain.AutoPurchaseWatch{}, "id = ?", id).Error
}

// AutoPurchaseLogRepository implementation
type AutoPurchaseLogRepositoryImpl struct {
	db *gorm.DB
}

func NewAutoPurchaseLogRepository(db *gorm.DB) domain.AutoPurchaseLogRepository {
	return &AutoPurchaseLogRepositoryImpl{db: db}
}

func (r *AutoPurchaseLogRepositoryImpl) Create(log *domain.AutoPurchaseLog) error {
	return r.db.Create(log).Error
}

func (r *AutoPurchaseLogRepositoryImpl) FindByWatchID(watchID uuid.UUID) ([]*domain.AutoPurchaseLog, error) {
	var logs []*domain.AutoPurchaseLog
	err := r.db.Where("watch_id = ?", watchID).
		Order("created_at DESC").
		Find(&logs).Error
	return logs, err
}
