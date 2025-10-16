package interfaces

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/yourusername/ecomate/backend/internal/usecase"
)

type SustainabilityHandler struct {
	sustainabilityUseCase usecase.SustainabilityUseCase
}

func NewSustainabilityHandler(sustainabilityUseCase usecase.SustainabilityUseCase) *SustainabilityHandler {
	return &SustainabilityHandler{
		sustainabilityUseCase: sustainabilityUseCase,
	}
}

func (h *SustainabilityHandler) GetDashboard(c *gin.Context) {
	userID, _ := c.Get("user_id")

	dashboard, err := h.sustainabilityUseCase.GetDashboard(userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dashboard)
}

func (h *SustainabilityHandler) GetLeaderboard(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	period := c.DefaultQuery("period", "all")

	leaderboard, err := h.sustainabilityUseCase.GetLeaderboard(limit, period)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"leaderboard": leaderboard})
}

func (h *SustainabilityHandler) GetUserFavorites(c *gin.Context) {
	userID, _ := c.Get("user_id")

	favorites, err := h.sustainabilityUseCase.GetUserFavorites(userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"favorites": favorites})
}
