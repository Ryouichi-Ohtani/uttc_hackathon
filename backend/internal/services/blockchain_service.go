package services

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Block represents a single block in the blockchain
type Block struct {
	Index        int                    `json:"index"`
	Timestamp    time.Time              `json:"timestamp"`
	Data         map[string]interface{} `json:"data"`
	PreviousHash string                 `json:"previous_hash"`
	Hash         string                 `json:"hash"`
	Nonce        int                    `json:"nonce"`
}

// Blockchain represents the entire chain
type Blockchain struct {
	Chain      []Block `json:"chain"`
	Difficulty int     `json:"difficulty"`
}

// TransactionRecord represents a transaction on the blockchain
type TransactionRecord struct {
	TransactionID uuid.UUID `json:"transaction_id"`
	ProductID     uuid.UUID `json:"product_id"`
	SellerID      uuid.UUID `json:"seller_id"`
	BuyerID       uuid.UUID `json:"buyer_id"`
	Price         int       `json:"price"`
	CO2Saved      float64   `json:"co2_saved_kg"`
	Timestamp     time.Time `json:"timestamp"`
	Verified      bool      `json:"verified"`
}

// BlockchainService handles blockchain operations
type BlockchainService struct {
	blockchain *Blockchain
}

func NewBlockchainService() *BlockchainService {
	service := &BlockchainService{
		blockchain: &Blockchain{
			Chain:      []Block{},
			Difficulty: 4, // Number of leading zeros required in hash
		},
	}

	// Create genesis block
	genesisBlock := Block{
		Index:        0,
		Timestamp:    time.Now(),
		Data:         map[string]interface{}{"type": "genesis"},
		PreviousHash: "0",
		Nonce:        0,
	}
	genesisBlock.Hash = service.calculateHash(genesisBlock)
	service.blockchain.Chain = append(service.blockchain.Chain, genesisBlock)

	return service
}

// RecordTransaction records a transaction on the blockchain
func (s *BlockchainService) RecordTransaction(tx TransactionRecord) (*Block, error) {
	data := map[string]interface{}{
		"type":           "transaction",
		"transaction_id": tx.TransactionID.String(),
		"product_id":     tx.ProductID.String(),
		"seller_id":      tx.SellerID.String(),
		"buyer_id":       tx.BuyerID.String(),
		"price":          tx.Price,
		"co2_saved":      tx.CO2Saved,
		"timestamp":      tx.Timestamp,
	}

	block := s.createBlock(data)
	s.blockchain.Chain = append(s.blockchain.Chain, *block)

	return block, nil
}

// VerifyTransaction verifies a transaction exists on the blockchain
func (s *BlockchainService) VerifyTransaction(transactionID uuid.UUID) (bool, *TransactionRecord, error) {
	for _, block := range s.blockchain.Chain {
		if txID, ok := block.Data["transaction_id"].(string); ok {
			if txID == transactionID.String() {
				// Reconstruct transaction record
				tx := &TransactionRecord{
					TransactionID: transactionID,
					Verified:      true,
				}

				if productID, ok := block.Data["product_id"].(string); ok {
					tx.ProductID, _ = uuid.Parse(productID)
				}
				if sellerID, ok := block.Data["seller_id"].(string); ok {
					tx.SellerID, _ = uuid.Parse(sellerID)
				}
				if buyerID, ok := block.Data["buyer_id"].(string); ok {
					tx.BuyerID, _ = uuid.Parse(buyerID)
				}
				if price, ok := block.Data["price"].(int); ok {
					tx.Price = price
				}
				if co2, ok := block.Data["co2_saved"].(float64); ok {
					tx.CO2Saved = co2
				}

				return true, tx, nil
			}
		}
	}

	return false, nil, fmt.Errorf("transaction not found on blockchain")
}

// GetTransactionHistory gets all transactions for a user
func (s *BlockchainService) GetTransactionHistory(userID uuid.UUID) []TransactionRecord {
	var transactions []TransactionRecord

	for _, block := range s.blockchain.Chain {
		if block.Data["type"] != "transaction" {
			continue
		}

		sellerID, _ := block.Data["seller_id"].(string)
		buyerID, _ := block.Data["buyer_id"].(string)

		if sellerID == userID.String() || buyerID == userID.String() {
			tx := TransactionRecord{
				Verified: true,
			}

			if txID, ok := block.Data["transaction_id"].(string); ok {
				tx.TransactionID, _ = uuid.Parse(txID)
			}
			if productID, ok := block.Data["product_id"].(string); ok {
				tx.ProductID, _ = uuid.Parse(productID)
			}
			if seller, ok := block.Data["seller_id"].(string); ok {
				tx.SellerID, _ = uuid.Parse(seller)
			}
			if buyer, ok := block.Data["buyer_id"].(string); ok {
				tx.BuyerID, _ = uuid.Parse(buyer)
			}
			if price, ok := block.Data["price"].(int); ok {
				tx.Price = price
			}
			if co2, ok := block.Data["co2_saved"].(float64); ok {
				tx.CO2Saved = co2
			}

			transactions = append(transactions, tx)
		}
	}

	return transactions
}

// ValidateChain validates the entire blockchain
func (s *BlockchainService) ValidateChain() bool {
	for i := 1; i < len(s.blockchain.Chain); i++ {
		currentBlock := s.blockchain.Chain[i]
		previousBlock := s.blockchain.Chain[i-1]

		// Check if current hash is correct
		if currentBlock.Hash != s.calculateHash(currentBlock) {
			return false
		}

		// Check if previous hash matches
		if currentBlock.PreviousHash != previousBlock.Hash {
			return false
		}

		// Check proof of work
		if !s.isValidProof(currentBlock) {
			return false
		}
	}

	return true
}

// GetCO2Certificate generates a CO2 reduction certificate
func (s *BlockchainService) GetCO2Certificate(userID uuid.UUID) (*CO2Certificate, error) {
	transactions := s.GetTransactionHistory(userID)

	totalCO2 := 0.0
	transactionCount := 0

	for _, tx := range transactions {
		totalCO2 += tx.CO2Saved
		transactionCount++
	}

	certificate := &CO2Certificate{
		UserID:           userID,
		TotalCO2Saved:    totalCO2,
		TransactionCount: transactionCount,
		IssuedAt:         time.Now(),
		CertificateHash:  s.generateCertificateHash(userID, totalCO2, transactionCount),
		Verified:         s.ValidateChain(),
	}

	return certificate, nil
}

type CO2Certificate struct {
	UserID           uuid.UUID `json:"user_id"`
	TotalCO2Saved    float64   `json:"total_co2_saved_kg"`
	TransactionCount int       `json:"transaction_count"`
	IssuedAt         time.Time `json:"issued_at"`
	CertificateHash  string    `json:"certificate_hash"`
	Verified         bool      `json:"verified"`
}

func (s *BlockchainService) createBlock(data map[string]interface{}) *Block {
	prevBlock := s.blockchain.Chain[len(s.blockchain.Chain)-1]

	block := &Block{
		Index:        prevBlock.Index + 1,
		Timestamp:    time.Now(),
		Data:         data,
		PreviousHash: prevBlock.Hash,
		Nonce:        0,
	}

	// Proof of work
	for {
		block.Hash = s.calculateHash(*block)
		if s.isValidProof(*block) {
			break
		}
		block.Nonce++
	}

	return block
}

func (s *BlockchainService) calculateHash(block Block) string {
	data, _ := json.Marshal(map[string]interface{}{
		"index":         block.Index,
		"timestamp":     block.Timestamp,
		"data":          block.Data,
		"previous_hash": block.PreviousHash,
		"nonce":         block.Nonce,
	})

	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}

func (s *BlockchainService) isValidProof(block Block) bool {
	prefix := ""
	for i := 0; i < s.blockchain.Difficulty; i++ {
		prefix += "0"
	}

	return block.Hash[:s.blockchain.Difficulty] == prefix
}

func (s *BlockchainService) generateCertificateHash(userID uuid.UUID, totalCO2 float64, txCount int) string {
	data := fmt.Sprintf("%s-%f-%d-%s", userID.String(), totalCO2, txCount, time.Now().String())
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}
