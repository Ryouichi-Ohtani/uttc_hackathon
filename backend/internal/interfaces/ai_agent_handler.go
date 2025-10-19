package interfaces

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/yourusername/ecomate/backend/internal/domain"
	"github.com/yourusername/ecomate/backend/internal/usecase"
)

type AIAgentHandler struct {
	aiAgentUC *usecase.AIAgentUseCase
}

func NewAIAgentHandler(aiAgentUC *usecase.AIAgentUseCase) *AIAgentHandler {
	return &AIAgentHandler{
		aiAgentUC: aiAgentUC,
	}
}

// ========== AI Listing Agent Endpoints ==========

// POST /api/v1/ai-agent/listing/generate
// AI出品エージェント: 画像からすべてを自動生成
func (h *AIAgentHandler) GenerateListing(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req domain.AIListingGenerationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.aiAgentUC.GenerateListingFromImages(c.Request.Context(), userID.(uuid.UUID), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// POST /api/v1/ai-agent/listing/:id/approve
// AI出品承認: ユーザーが承認画面で修正した内容を反映
func (h *AIAgentHandler) ApproveListing(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	productIDStr := c.Param("id")
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product ID"})
		return
	}

	var modifications map[string]interface{}
	if err := c.ShouldBindJSON(&modifications); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.aiAgentUC.ApproveAndModifyListing(c.Request.Context(), userID.(uuid.UUID), productID, modifications); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "listing approved and updated"})
}

// GET /api/v1/ai-agent/listing/:id/data
// AI生成データ取得
func (h *AIAgentHandler) GetListingData(c *gin.Context) {
	productIDStr := c.Param("id")
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product ID"})
		return
	}

	data, err := h.aiAgentUC.GetListingData(productID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "listing data not found"})
		return
	}

	c.JSON(http.StatusOK, data)
}

// ========== AI Negotiation Agent Endpoints ==========

// POST /api/v1/ai-agent/negotiation/enable
// AI交渉エージェントを有効化/設定
func (h *AIAgentHandler) EnableNegotiation(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req domain.ToggleAINegotiationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.aiAgentUC.EnableAINegotiation(c.Request.Context(), userID.(uuid.UUID), &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "AI negotiation enabled"})
}

// GET /api/v1/ai-agent/negotiation/:product_id
// AI交渉設定を取得
func (h *AIAgentHandler) GetNegotiationSettings(c *gin.Context) {
	productIDStr := c.Param("product_id")
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product ID"})
		return
	}

	settings, err := h.aiAgentUC.GetNegotiationSettings(productID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "negotiation settings not found"})
		return
	}

	c.JSON(http.StatusOK, settings)
}

// DELETE /api/v1/ai-agent/negotiation/:product_id
// AI交渉を無効化
func (h *AIAgentHandler) DisableNegotiation(c *gin.Context) {
	_, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	productIDStr := c.Param("product_id")
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product ID"})
		return
	}

	// TODO: 所有権確認

	if err := h.aiAgentUC.DisableNegotiation(productID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "AI negotiation disabled"})
}

// ========== AI Shipping Agent Endpoints ==========

// POST /api/v1/ai-agent/shipping/prepare
// AI配送準備エージェント: 購入後に配送情報を自動準備
func (h *AIAgentHandler) PrepareShipping(c *gin.Context) {
	var req domain.AIShippingPreparationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	prep, err := h.aiAgentUC.PrepareShipping(c.Request.Context(), req.PurchaseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, prep)
}

// GET /api/v1/ai-agent/shipping/:purchase_id
// 配送準備情報を取得
func (h *AIAgentHandler) GetShippingPreparation(c *gin.Context) {
	purchaseIDStr := c.Param("purchase_id")
	purchaseID, err := uuid.Parse(purchaseIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid purchase ID"})
		return
	}

	prep, err := h.aiAgentUC.GetShippingPreparation(purchaseID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "shipping preparation not found"})
		return
	}

	c.JSON(http.StatusOK, prep)
}

// POST /api/v1/ai-agent/shipping/:purchase_id/approve
// 配送情報を承認
func (h *AIAgentHandler) ApproveShipping(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	purchaseIDStr := c.Param("purchase_id")
	purchaseID, err := uuid.Parse(purchaseIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid purchase ID"})
		return
	}

	var req domain.ApproveShippingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.aiAgentUC.ApproveShipping(c.Request.Context(), userID.(uuid.UUID), purchaseID, &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "shipping approved"})
}

// ========== Statistics ==========

// GET /api/v1/ai-agent/stats
// AIエージェント利用統計
func (h *AIAgentHandler) GetAgentStats(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	stats, err := h.aiAgentUC.GetAgentStats(c.Request.Context(), userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// ========== Utility ==========

// RegisterRoutes registers all AI agent routes
func (h *AIAgentHandler) RegisterRoutes(router *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	aiAgent := router.Group("/ai-agent")
	aiAgent.Use(authMiddleware)
	{
		// Listing Agent
		aiAgent.POST("/listing/generate", h.GenerateListing)
		aiAgent.POST("/listing/:id/approve", h.ApproveListing)
		aiAgent.GET("/listing/:id/data", h.GetListingData)

		// Negotiation Agent
		aiAgent.POST("/negotiation/enable", h.EnableNegotiation)
		aiAgent.GET("/negotiation/:product_id", h.GetNegotiationSettings)
		aiAgent.DELETE("/negotiation/:product_id", h.DisableNegotiation)

		// Shipping Agent
		aiAgent.POST("/shipping/prepare", h.PrepareShipping)
		aiAgent.GET("/shipping/:purchase_id", h.GetShippingPreparation)
		aiAgent.POST("/shipping/:purchase_id/approve", h.ApproveShipping)

		// Stats
		aiAgent.GET("/stats", h.GetAgentStats)
	}
}
