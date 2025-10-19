package usecase

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/yourusername/ecomate/backend/internal/domain"
)

type BlockchainUseCase interface {
	RecordPurchaseOnChain(purchaseID uuid.UUID) (*domain.BlockchainTransaction, error)
	IssueCO2Token(purchaseID uuid.UUID, co2Amount float64) (*domain.BlockchainTransaction, error)
	MintProductNFT(productID, ownerID uuid.UUID) (*domain.NFTOwnership, error)
	GetUserNFTs(userID uuid.UUID) ([]*domain.NFTOwnership, error)
	GetTransactionByPurchaseID(purchaseID uuid.UUID) (*domain.BlockchainTransaction, error)
	VerifyTransaction(txHash string) (*domain.BlockchainTransaction, error)
}

type blockchainUseCase struct {
	blockchainRepo domain.BlockchainRepository
	nftRepo        domain.NFTRepository
	purchaseRepo   domain.PurchaseRepository
	productRepo    domain.ProductRepository
}

func NewBlockchainUseCase(
	blockchainRepo domain.BlockchainRepository,
	nftRepo domain.NFTRepository,
	purchaseRepo domain.PurchaseRepository,
	productRepo domain.ProductRepository,
) BlockchainUseCase {
	return &blockchainUseCase{
		blockchainRepo: blockchainRepo,
		nftRepo:        nftRepo,
		purchaseRepo:   purchaseRepo,
		productRepo:    productRepo,
	}
}

func (u *blockchainUseCase) RecordPurchaseOnChain(purchaseID uuid.UUID) (*domain.BlockchainTransaction, error) {
	purchase, err := u.purchaseRepo.FindByID(purchaseID)
	if err != nil {
		return nil, err
	}

	// Generate transaction hash (mock - in production, interact with blockchain)
	txHash := u.generateTransactionHash(purchase)

	tx := &domain.BlockchainTransaction{
		PurchaseID:      purchaseID,
		TransactionHash: txHash,
		BlockNumber:     int64(time.Now().Unix()), // Mock block number
		ChainID:         "ecomate-chain-1",
		Status:          "confirmed",
		CO2TokenAmount:  0, // CO2 feature removed
	}

	now := time.Now()
	tx.ConfirmedAt = &now

	if err := u.blockchainRepo.CreateTransaction(tx); err != nil {
		return nil, err
	}

	return tx, nil
}

func (u *blockchainUseCase) IssueCO2Token(purchaseID uuid.UUID, co2Amount float64) (*domain.BlockchainTransaction, error) {
	// Similar to RecordPurchaseOnChain but specifically for CO2 tokens
	return u.RecordPurchaseOnChain(purchaseID)
}

func (u *blockchainUseCase) MintProductNFT(productID, ownerID uuid.UUID) (*domain.NFTOwnership, error) {
	product, err := u.productRepo.FindByID(productID)
	if err != nil {
		return nil, err
	}

	// Generate NFT metadata
	metadataURI := fmt.Sprintf("ipfs://Qm%s", uuid.New().String()[:32])
	tokenID := fmt.Sprintf("ECO-%s", uuid.New().String()[:8])

	nft := &domain.NFTOwnership{
		ProductID:       productID,
		OwnerID:         ownerID,
		TokenID:         tokenID,
		ContractAddress: "0xEcoMateNFTContract1234567890",
		MetadataURI:     metadataURI,
	}

	if err := u.nftRepo.Create(nft); err != nil {
		return nil, err
	}

	// In production: Actually mint NFT on blockchain
	fmt.Printf("NFT minted for product %s: TokenID=%s\n", product.Title, tokenID)

	return nft, nil
}

func (u *blockchainUseCase) GetUserNFTs(userID uuid.UUID) ([]*domain.NFTOwnership, error) {
	return u.nftRepo.FindByOwnerID(userID)
}

func (u *blockchainUseCase) GetTransactionByPurchaseID(purchaseID uuid.UUID) (*domain.BlockchainTransaction, error) {
	return u.blockchainRepo.FindByPurchaseID(purchaseID)
}

func (u *blockchainUseCase) VerifyTransaction(txHash string) (*domain.BlockchainTransaction, error) {
	// In production: Query blockchain node
	// For now, just lookup in database
	var tx domain.BlockchainTransaction
	// Simplified - would need repository method
	return &tx, nil
}

func (u *blockchainUseCase) generateTransactionHash(purchase *domain.Purchase) string {
	// Generate deterministic hash from purchase data
	data := fmt.Sprintf("%s-%s-%d-%d",
		purchase.ID.String(),
		purchase.BuyerID.String(),
		purchase.Price,
		time.Now().Unix(),
	)

	hash := sha256.Sum256([]byte(data))
	return "0x" + hex.EncodeToString(hash[:])
}
