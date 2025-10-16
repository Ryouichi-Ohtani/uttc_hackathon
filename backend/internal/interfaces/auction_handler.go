package interfaces

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/yourusername/ecomate/backend/internal/domain"
	"github.com/yourusername/ecomate/backend/internal/usecase"
)

type AuctionHandler struct {
	auctionUseCase usecase.AuctionUseCase
	connections    map[uuid.UUID]map[*websocket.Conn]bool
	mu             sync.RWMutex
	upgrader       websocket.Upgrader
}

func NewAuctionHandler(auctionUseCase usecase.AuctionUseCase) *AuctionHandler {
	return &AuctionHandler{
		auctionUseCase: auctionUseCase,
		connections:    make(map[uuid.UUID]map[*websocket.Conn]bool),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

// CreateAuction handles POST /auctions
func (h *AuctionHandler) CreateAuction(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req domain.CreateAuctionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	productID, err := uuid.Parse(req.ProductID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	auction, err := h.auctionUseCase.CreateAuction(
		userID.(uuid.UUID),
		productID,
		req.StartPrice,
		req.MinBidIncrement,
		req.DurationMinutes,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, auction)
}

// PlaceBid handles POST /auctions/:id/bids
func (h *AuctionHandler) PlaceBid(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	auctionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid auction ID"})
		return
	}

	var req domain.PlaceBidRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	bid, err := h.auctionUseCase.PlaceBid(auctionID, userID.(uuid.UUID), req.Amount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Broadcast bid to all connected clients
	h.broadcastBid(auctionID, bid)

	c.JSON(http.StatusCreated, bid)
}

// GetAuction handles GET /auctions/:id
func (h *AuctionHandler) GetAuction(c *gin.Context) {
	auctionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid auction ID"})
		return
	}

	auction, err := h.auctionUseCase.GetAuction(auctionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Auction not found"})
		return
	}

	c.JSON(http.StatusOK, auction)
}

// GetActiveAuctions handles GET /auctions
func (h *AuctionHandler) GetActiveAuctions(c *gin.Context) {
	auctions, err := h.auctionUseCase.GetActiveAuctions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"auctions": auctions})
}

// GetAuctionBids handles GET /auctions/:id/bids
func (h *AuctionHandler) GetAuctionBids(c *gin.Context) {
	auctionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid auction ID"})
		return
	}

	bids, err := h.auctionUseCase.GetAuctionBids(auctionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"bids": bids})
}

// WebSocketHandler handles WebSocket connections for real-time bidding
func (h *AuctionHandler) WebSocketHandler(c *gin.Context) {
	auctionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid auction ID"})
		return
	}

	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	// Add connection
	h.mu.Lock()
	if h.connections[auctionID] == nil {
		h.connections[auctionID] = make(map[*websocket.Conn]bool)
	}
	h.connections[auctionID][conn] = true
	h.mu.Unlock()

	// Remove connection on close
	defer func() {
		h.mu.Lock()
		delete(h.connections[auctionID], conn)
		h.mu.Unlock()
	}()

	// Keep connection alive
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}
}

func (h *AuctionHandler) broadcastBid(auctionID uuid.UUID, bid *domain.Bid) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	connections, exists := h.connections[auctionID]
	if !exists {
		return
	}

	for conn := range connections {
		conn.WriteJSON(map[string]interface{}{
			"type": "new_bid",
			"data": bid,
		})
	}
}
