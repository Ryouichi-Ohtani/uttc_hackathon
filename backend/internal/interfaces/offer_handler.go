package interfaces

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/yourusername/ecomate/backend/internal/domain"
	"github.com/yourusername/ecomate/backend/internal/usecase"
)

type OfferHandler struct {
	offerUseCase usecase.OfferUseCase
}

func NewOfferHandler(offerUseCase usecase.OfferUseCase) *OfferHandler {
	return &OfferHandler{
		offerUseCase: offerUseCase,
	}
}

// CreateOffer handles POST /offers
func (h *OfferHandler) CreateOffer(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req domain.CreateOfferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	productID, err := uuid.Parse(req.ProductID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	offer, err := h.offerUseCase.CreateOffer(
		userID.(uuid.UUID),
		productID,
		req.OfferPrice,
		req.Message,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, offer)
}

// RespondOffer handles PATCH /offers/:id/respond
func (h *OfferHandler) RespondOffer(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	offerID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offer ID"})
		return
	}

	var req domain.RespondOfferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	offer, err := h.offerUseCase.RespondOffer(
		offerID,
		userID.(uuid.UUID),
		req.Accept,
		req.Message,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, offer)
}

// GetMyOffers handles GET /offers/my
func (h *OfferHandler) GetMyOffers(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	role := c.DefaultQuery("role", "buyer") // buyer or seller

	var offers []*domain.Offer
	var err error

	if role == "seller" {
		offers, err = h.offerUseCase.GetSellerOffers(userID.(uuid.UUID))
	} else {
		offers, err = h.offerUseCase.GetBuyerOffers(userID.(uuid.UUID))
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"offers": offers})
}

// GetProductOffers handles GET /offers/products/:id
func (h *OfferHandler) GetProductOffers(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	productID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	offers, err := h.offerUseCase.GetProductOffers(productID, userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"offers": offers})
}

// GetNegotiationSuggestion handles GET /offers/products/:id/ai-suggestion
func (h *OfferHandler) GetNegotiationSuggestion(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	productID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	isBuyer := c.DefaultQuery("role", "buyer") == "buyer"

	suggestion, err := h.offerUseCase.GetNegotiationSuggestion(productID, userID.(uuid.UUID), isBuyer)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, suggestion)
}
