package interfaces

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/yourusername/ecomate/backend/internal/usecase"
)

type SalesPredictionHandler struct {
	salesPredictionUseCase usecase.SalesPredictionUseCase
}

func NewSalesPredictionHandler(salesPredictionUseCase usecase.SalesPredictionUseCase) *SalesPredictionHandler {
	return &SalesPredictionHandler{
		salesPredictionUseCase: salesPredictionUseCase,
	}
}

// PredictProductPrice handles GET /predictions/products/:id/price
func (h *SalesPredictionHandler) PredictProductPrice(c *gin.Context) {
	productID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	prediction, err := h.salesPredictionUseCase.PredictProductSalePrice(productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, prediction)
}

// PredictSellerRevenue handles GET /predictions/revenue
func (h *SalesPredictionHandler) PredictSellerRevenue(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	days, _ := strconv.Atoi(c.DefaultQuery("days", "30"))

	prediction, err := h.salesPredictionUseCase.PredictSellerRevenue(userID.(uuid.UUID), days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, prediction)
}

// GetMarketTrends handles GET /predictions/market-trends
func (h *SalesPredictionHandler) GetMarketTrends(c *gin.Context) {
	category := c.Query("category")
	if category == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Category is required"})
		return
	}

	days, _ := strconv.Atoi(c.DefaultQuery("days", "30"))

	trends, err := h.salesPredictionUseCase.GetMarketTrends(category, days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, trends)
}
