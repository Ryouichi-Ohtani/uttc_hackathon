package interfaces

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/yourusername/ecomate/backend/internal/domain"
	"github.com/yourusername/ecomate/backend/internal/usecase"
)

type PurchaseHandler struct {
	purchaseUseCase usecase.PurchaseUseCase
}

func NewPurchaseHandler(purchaseUseCase usecase.PurchaseUseCase) *PurchaseHandler {
	return &PurchaseHandler{
		purchaseUseCase: purchaseUseCase,
	}
}

func (h *PurchaseHandler) Create(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req domain.CreatePurchaseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	purchase, err := h.purchaseUseCase.Create(userID.(uuid.UUID), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, purchase)
}

func (h *PurchaseHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	purchase, err := h.purchaseUseCase.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "purchase not found"})
		return
	}

	c.JSON(http.StatusOK, purchase)
}

func (h *PurchaseHandler) List(c *gin.Context) {
	userID, _ := c.Get("user_id")

	role := c.Query("role") // "buyer" or "seller" or ""
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	purchases, pagination, err := h.purchaseUseCase.ListByUser(userID.(uuid.UUID), role, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"purchases":  purchases,
		"pagination": pagination,
	})
}

func (h *PurchaseHandler) Complete(c *gin.Context) {
	userID, _ := c.Get("user_id")

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.purchaseUseCase.CompletePurchase(id, userID.(uuid.UUID)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "purchase completed successfully"})
}

func (h *PurchaseHandler) GetShippingLabel(c *gin.Context) {
	userID, _ := c.Get("user_id")

	purchaseID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid purchase id"})
		return
	}

	label, err := h.purchaseUseCase.GetShippingLabel(purchaseID, userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, label)
}

func (h *PurchaseHandler) GenerateShippingLabel(c *gin.Context) {
	userID, _ := c.Get("user_id")

	purchaseID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid purchase id"})
		return
	}

	label, err := h.purchaseUseCase.GenerateShippingLabel(purchaseID, userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, label)
}
