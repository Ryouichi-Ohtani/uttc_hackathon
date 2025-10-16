package interfaces

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/yourusername/ecomate/backend/internal/usecase"
)

type BlockchainHandler struct {
	blockchainUseCase usecase.BlockchainUseCase
}

func NewBlockchainHandler(blockchainUseCase usecase.BlockchainUseCase) *BlockchainHandler {
	return &BlockchainHandler{
		blockchainUseCase: blockchainUseCase,
	}
}

type RecordPurchaseRequest struct {
	PurchaseID uuid.UUID `json:"purchase_id" binding:"required"`
}

func (h *BlockchainHandler) RecordPurchase(c *gin.Context) {
	var req RecordPurchaseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx, err := h.blockchainUseCase.RecordPurchaseOnChain(req.PurchaseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tx)
}

type MintNFTRequest struct {
	ProductID uuid.UUID `json:"product_id" binding:"required"`
}

func (h *BlockchainHandler) MintNFT(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req MintNFTRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	nft, err := h.blockchainUseCase.MintProductNFT(req.ProductID, userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, nft)
}

func (h *BlockchainHandler) GetMyNFTs(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	nfts, err := h.blockchainUseCase.GetUserNFTs(userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"nfts": nfts,
	})
}

func (h *BlockchainHandler) GetPurchaseTransaction(c *gin.Context) {
	purchaseID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid purchase ID"})
		return
	}

	tx, err := h.blockchainUseCase.GetTransactionByPurchaseID(purchaseID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "transaction not found"})
		return
	}

	c.JSON(http.StatusOK, tx)
}
