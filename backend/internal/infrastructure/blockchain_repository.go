package infrastructure

import (
	"github.com/google/uuid"
	"github.com/yourusername/ecomate/backend/internal/domain"
	"gorm.io/gorm"
)

type blockchainRepository struct {
	db *gorm.DB
}

func NewBlockchainRepository(db *gorm.DB) domain.BlockchainRepository {
	return &blockchainRepository{db: db}
}

func (r *blockchainRepository) CreateTransaction(tx *domain.BlockchainTransaction) error {
	return r.db.Create(tx).Error
}

func (r *blockchainRepository) FindByPurchaseID(purchaseID uuid.UUID) (*domain.BlockchainTransaction, error) {
	var tx domain.BlockchainTransaction
	err := r.db.Where("purchase_id = ?", purchaseID).First(&tx).Error
	return &tx, err
}

func (r *blockchainRepository) UpdateTransactionStatus(txHash string, status string, blockNumber int64) error {
	return r.db.Model(&domain.BlockchainTransaction{}).
		Where("transaction_hash = ?", txHash).
		Updates(map[string]interface{}{
			"status":       status,
			"block_number": blockNumber,
		}).Error
}

type nftRepository struct {
	db *gorm.DB
}

func NewNFTRepository(db *gorm.DB) domain.NFTRepository {
	return &nftRepository{db: db}
}

func (r *nftRepository) Create(nft *domain.NFTOwnership) error {
	return r.db.Create(nft).Error
}

func (r *nftRepository) FindByProductID(productID uuid.UUID) (*domain.NFTOwnership, error) {
	var nft domain.NFTOwnership
	err := r.db.Preload("Product").Preload("Owner").
		Where("product_id = ?", productID).First(&nft).Error
	return &nft, err
}

func (r *nftRepository) FindByOwnerID(ownerID uuid.UUID) ([]*domain.NFTOwnership, error) {
	var nfts []*domain.NFTOwnership
	err := r.db.Preload("Product").Preload("Product.Images").
		Where("owner_id = ?", ownerID).Find(&nfts).Error
	return nfts, err
}

func (r *nftRepository) Transfer(nftID, newOwnerID uuid.UUID) error {
	return r.db.Model(&domain.NFTOwnership{}).
		Where("id = ?", nftID).
		Update("owner_id", newOwnerID).Error
}
