package interfaces

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/yourusername/ecomate/backend/internal/usecase"
)

type AnalyticsHandler struct {
	analyticsUseCase usecase.AnalyticsUseCase
}

func NewAnalyticsHandler(analyticsUseCase usecase.AnalyticsUseCase) *AnalyticsHandler {
	return &AnalyticsHandler{
		analyticsUseCase: analyticsUseCase,
	}
}

// GetUserBehavior handles GET /analytics/user/behavior
func (h *AnalyticsHandler) GetUserBehavior(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	summary, err := h.analyticsUseCase.GetUserBehaviorSummary(userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, summary)
}

// GetPopularProducts handles GET /analytics/popular-products
func (h *AnalyticsHandler) GetPopularProducts(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "7"))

	products, err := h.analyticsUseCase.GetPopularProducts(days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"products": products})
}

// GetSearchTrends handles GET /analytics/search-trends
func (h *AnalyticsHandler) GetSearchTrends(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "7"))

	keywords, err := h.analyticsUseCase.GetSearchTrends(days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"keywords": keywords})
}
