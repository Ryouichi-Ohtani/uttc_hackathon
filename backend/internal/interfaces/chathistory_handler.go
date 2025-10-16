package interfaces

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/yourusername/ecomate/backend/internal/usecase"
)

type ChatHistoryHandler struct {
	chatHistoryUseCase *usecase.ChatHistoryUseCase
}

func NewChatHistoryHandler(chatHistoryUseCase *usecase.ChatHistoryUseCase) *ChatHistoryHandler {
	return &ChatHistoryHandler{chatHistoryUseCase: chatHistoryUseCase}
}

func (h *ChatHistoryHandler) GetHistory(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	limitStr := c.DefaultQuery("limit", "50")
	limit, _ := strconv.Atoi(limitStr)

	histories, err := h.chatHistoryUseCase.GetHistory(userID.(uuid.UUID), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"histories": histories})
}

func (h *ChatHistoryHandler) DeleteHistory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := h.chatHistoryUseCase.DeleteHistory(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "History deleted"})
}

type CO2GoalHandler struct {
	goalUseCase *usecase.CO2GoalUseCase
}

func NewCO2GoalHandler(goalUseCase *usecase.CO2GoalUseCase) *CO2GoalHandler {
	return &CO2GoalHandler{goalUseCase: goalUseCase}
}

type CreateGoalRequest struct {
	TargetKG   float64 `json:"target_kg" binding:"required"`
	TargetDate string  `json:"target_date" binding:"required"`
}

func (h *CO2GoalHandler) CreateGoal(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req CreateGoalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	targetDate, err := time.Parse("2006-01-02", req.TargetDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
		return
	}

	goal, err := h.goalUseCase.CreateGoal(userID.(uuid.UUID), req.TargetKG, targetDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, goal)
}

func (h *CO2GoalHandler) GetGoal(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	goal, err := h.goalUseCase.GetGoal(userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if goal == nil {
		c.JSON(http.StatusOK, gin.H{"goal": nil})
		return
	}

	// Calculate progress percentage
	progress := 0.0
	if goal.TargetKG > 0 {
		progress = (goal.CurrentKG / goal.TargetKG) * 100
		if progress > 100 {
			progress = 100
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"goal":     goal,
		"progress": progress,
	})
}

type ShippingHandler struct {
	shippingUseCase *usecase.ShippingTrackingUseCase
}

func NewShippingHandler(shippingUseCase *usecase.ShippingTrackingUseCase) *ShippingHandler {
	return &ShippingHandler{shippingUseCase: shippingUseCase}
}

type CreateShippingRequest struct {
	PurchaseID     string `json:"purchase_id" binding:"required"`
	Carrier        string `json:"carrier" binding:"required"`
	ShippingMethod string `json:"shipping_method" binding:"required"`
}

func (h *ShippingHandler) CreateShipping(c *gin.Context) {
	var req CreateShippingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	purchaseID, err := uuid.Parse(req.PurchaseID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid purchase ID"})
		return
	}

	tracking, err := h.shippingUseCase.CreateTracking(purchaseID, req.Carrier, req.ShippingMethod)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, tracking)
}

func (h *ShippingHandler) GetShipping(c *gin.Context) {
	purchaseIDStr := c.Param("purchase_id")
	purchaseID, err := uuid.Parse(purchaseIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid purchase ID"})
		return
	}

	tracking, err := h.shippingUseCase.GetTracking(purchaseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if tracking == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tracking not found"})
		return
	}

	c.JSON(http.StatusOK, tracking)
}

type UpdateShippingStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

func (h *ShippingHandler) UpdateStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var req UpdateShippingStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.shippingUseCase.UpdateStatus(id, req.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Status updated"})
}
