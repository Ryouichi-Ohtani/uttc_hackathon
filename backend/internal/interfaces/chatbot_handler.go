package interfaces

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/ecomate/backend/internal/domain"
	"github.com/yourusername/ecomate/backend/internal/infrastructure"
)

type ChatbotHandler struct {
	aiClient    *infrastructure.AIClient
	productRepo domain.ProductRepository
}

func NewChatbotHandler(aiClient *infrastructure.AIClient, productRepo domain.ProductRepository) *ChatbotHandler {
	return &ChatbotHandler{
		aiClient:    aiClient,
		productRepo: productRepo,
	}
}

type ChatRequest struct {
	Messages []ChatMessage `json:"messages" binding:"required"`
	Context  string        `json:"context"`
}

type ChatMessage struct {
	Role    string `json:"role" binding:"required,oneof=user assistant"`
	Content string `json:"content" binding:"required"`
}

type ChatResponse struct {
	Message string `json:"message"`
}

// Chat handles AI chatbot requests
func (h *ChatbotHandler) Chat(c *gin.Context) {
	var req ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get products data from database
	ctx := context.Background()
	products, _, err := h.productRepo.List(&domain.ProductFilters{
		Limit: 20,
	})
	if err != nil {
		// Continue without product data
		products = []*domain.Product{}
	}

	// Build product context for AI
	productContext := h.buildProductContext(products)

	// Convert to infrastructure types
	messages := make([]infrastructure.ChatMessage, len(req.Messages))
	for i, msg := range req.Messages {
		messages[i] = infrastructure.ChatMessage{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}

	chatReq := &infrastructure.ChatRequest{
		Messages: messages,
		Context:  productContext,
	}

	// Get AI response
	response, err := h.aiClient.Chat(ctx, chatReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate response"})
		return
	}

	c.JSON(http.StatusOK, ChatResponse{
		Message: response,
	})
}

func (h *ChatbotHandler) buildProductContext(products []*domain.Product) string {
	if len(products) == 0 {
		return ""
	}

	type ProductInfo struct {
		ID          string  `json:"id"`
		Title       string  `json:"title"`
		Description string  `json:"description"`
		Price       int     `json:"price"`
		Category    string  `json:"category"`
		Condition   string  `json:"condition"`
		CO2Impact   float64 `json:"co2_impact_kg"`
		ImageURL    string  `json:"image_url"`
	}

	productList := make([]ProductInfo, 0, len(products))
	for _, p := range products {
		imageURL := ""
		if len(p.Images) > 0 {
			for _, img := range p.Images {
				if img.IsPrimary {
					if img.CDNURL != "" {
						imageURL = img.CDNURL
					} else {
						imageURL = img.ImageURL
					}
					break
				}
			}
			if imageURL == "" && len(p.Images) > 0 {
				if p.Images[0].CDNURL != "" {
					imageURL = p.Images[0].CDNURL
				} else {
					imageURL = p.Images[0].ImageURL
				}
			}
		}

		productList = append(productList, ProductInfo{
			ID:          p.ID.String(),
			Title:       p.Title,
			Description: p.Description,
			Price:       p.Price,
			Category:    p.Category,
			Condition:   string(p.Condition),
			CO2Impact:   0, // CO2 feature removed
			ImageURL:    imageURL,
		})
	}

	jsonData, err := json.Marshal(productList)
	if err != nil {
		return ""
	}

	context := fmt.Sprintf(`
利用可能な商品リスト（JSON形式）:
%s

商品を推薦する時は、以下のフォーマットを使用してください：
[PRODUCT:商品ID:商品名:商品画像URL]

例: [PRODUCT:123e4567-e89b-12d3-a456-426614174000:ハイトップスニーカー:https://example.com/image.jpg]

このフォーマットを使うと、フロントエンドで商品カードとして表示されます。
複数の商品を推薦する場合は、改行して複数の[PRODUCT]タグを含めてください。
`, string(jsonData))

	return context
}
