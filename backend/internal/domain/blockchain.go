package domain

import (
	"time"

	"github.com/google/uuid"
)

type BlockchainTransaction struct {
	ID              uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	PurchaseID      uuid.UUID `json:"purchase_id" gorm:"type:uuid;not null;index"`
	Purchase        *Purchase `json:"purchase,omitempty" gorm:"foreignKey:PurchaseID"`
	TransactionHash string    `json:"transaction_hash" gorm:"uniqueIndex"`
	BlockNumber     int64     `json:"block_number"`
	ChainID         string    `json:"chain_id"`
	Status          string    `json:"status"` // pending, confirmed, failed
	CO2TokenAmount  float64   `json:"co2_token_amount" gorm:"type:decimal(10,2)"`
	CreatedAt       time.Time `json:"created_at"`
	ConfirmedAt     *time.Time `json:"confirmed_at"`
}

type NFTOwnership struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ProductID   uuid.UUID `json:"product_id" gorm:"type:uuid;not null;index"`
	Product     *Product  `json:"product,omitempty" gorm:"foreignKey:ProductID"`
	OwnerID     uuid.UUID `json:"owner_id" gorm:"type:uuid;not null;index"`
	Owner       *User     `json:"owner,omitempty" gorm:"foreignKey:OwnerID"`
	TokenID     string    `json:"token_id" gorm:"uniqueIndex"`
	ContractAddress string `json:"contract_address"`
	MetadataURI string    `json:"metadata_uri"`
	CreatedAt   time.Time `json:"created_at"`
}

type BlockchainRepository interface {
	CreateTransaction(tx *BlockchainTransaction) error
	FindByPurchaseID(purchaseID uuid.UUID) (*BlockchainTransaction, error)
	UpdateTransactionStatus(txHash string, status string, blockNumber int64) error
}

type NFTRepository interface {
	Create(nft *NFTOwnership) error
	FindByProductID(productID uuid.UUID) (*NFTOwnership, error)
	FindByOwnerID(ownerID uuid.UUID) ([]*NFTOwnership, error)
	Transfer(nftID, newOwnerID uuid.UUID) error
}
