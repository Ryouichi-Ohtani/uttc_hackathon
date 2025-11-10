package interfaces

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/yourusername/ecomate/backend/internal/domain"
	"github.com/yourusername/ecomate/backend/internal/usecase"
)

type AutoPurchaseHandler struct {
	autoPurchaseUseCase usecase.AutoPurchaseUseCase
}

func NewAutoPurchaseHandler(autoPurchaseUseCase usecase.AutoPurchaseUseCase) *AutoPurchaseHandler {
	return &AutoPurchaseHandler{
		autoPurchaseUseCase: autoPurchaseUseCase,
	}
}

// AuthorizePayment handles payment authorization for auto-purchase
func (h *AutoPurchaseHandler) AuthorizePayment(c *gin.Context) {
	var req domain.PaymentAuthorizationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.autoPurchaseUseCase.AuthorizePayment(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// CreateWatch creates a new auto-purchase watch
func (h *AutoPurchaseHandler) CreateWatch(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req domain.CreateAutoPurchaseWatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	watch, err := h.autoPurchaseUseCase.CreateWatch(userID.(uuid.UUID), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, watch)
}

// GetUserWatches retrieves all watches for the authenticated user
func (h *AutoPurchaseHandler) GetUserWatches(c *gin.Context) {
	userID, _ := c.Get("user_id")

	watches, err := h.autoPurchaseUseCase.GetUserWatches(userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"watches": watches,
	})
}

// GetWatch retrieves a specific watch
func (h *AutoPurchaseHandler) GetWatch(c *gin.Context) {
	userID, _ := c.Get("user_id")

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	watch, err := h.autoPurchaseUseCase.GetWatchByID(id, userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, watch)
}

// CancelWatch cancels an auto-purchase watch
func (h *AutoPurchaseHandler) CancelWatch(c *gin.Context) {
	userID, _ := c.Get("user_id")

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.autoPurchaseUseCase.CancelWatch(id, userID.(uuid.UUID)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "watch cancelled successfully"})
}

// CheckAndExecute is a webhook endpoint for background jobs to trigger price checks
func (h *AutoPurchaseHandler) CheckAndExecute(c *gin.Context) {
	// In production, this should be protected by an API key or internal network restriction
	executed, err := h.autoPurchaseUseCase.CheckAndExecuteAutoPurchases()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Price check completed",
		"executed": executed,
	})
}
