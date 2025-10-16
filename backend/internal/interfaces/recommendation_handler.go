package interfaces

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/yourusername/ecomate/backend/internal/usecase"
)

type RecommendationHandler struct {
	recommendationUseCase usecase.RecommendationUseCase
}

func NewRecommendationHandler(recommendationUseCase usecase.RecommendationUseCase) *RecommendationHandler {
	return &RecommendationHandler{
		recommendationUseCase: recommendationUseCase,
	}
}

// GetRecommendations handles GET /recommendations
func (h *RecommendationHandler) GetRecommendations(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	recommendations, err := h.recommendationUseCase.GetRecommendations(userID.(uuid.UUID), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"recommendations": recommendations})
}
