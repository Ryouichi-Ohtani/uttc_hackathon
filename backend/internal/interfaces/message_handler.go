package interfaces

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/yourusername/ecomate/backend/internal/domain"
	"github.com/yourusername/ecomate/backend/internal/usecase"
)

type MessageHandler struct {
	messageUseCase usecase.MessageUseCase
	authUseCase    *usecase.AuthUseCase
	upgrader       websocket.Upgrader
	connections    map[uuid.UUID]map[*websocket.Conn]bool // conversationID -> connections
	mu             sync.RWMutex
}

func NewMessageHandler(messageUseCase usecase.MessageUseCase, authUseCase *usecase.AuthUseCase) *MessageHandler {
	return &MessageHandler{
		messageUseCase: messageUseCase,
		authUseCase:    authUseCase,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // In production, check origin properly
			},
		},
		connections: make(map[uuid.UUID]map[*websocket.Conn]bool),
	}
}

func (h *MessageHandler) CreateConversation(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req domain.CreateConversationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	conversation, err := h.messageUseCase.CreateConversation(userID.(uuid.UUID), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, conversation)
}

func (h *MessageHandler) GetOrCreateConversation(c *gin.Context) {
	userID, _ := c.Get("user_id")

	productID, err := uuid.Parse(c.Param("productId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
		return
	}

	sellerID, err := uuid.Parse(c.Param("sellerId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid seller id"})
		return
	}

	conversation, err := h.messageUseCase.GetOrCreateConversation(productID, userID.(uuid.UUID), sellerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, conversation)
}

func (h *MessageHandler) ListConversations(c *gin.Context) {
	userID, _ := c.Get("user_id")

	conversations, err := h.messageUseCase.GetUserConversations(userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"conversations": conversations})
}

func (h *MessageHandler) GetMessages(c *gin.Context) {
	conversationID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid conversation id"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))

	messages, pagination, err := h.messageUseCase.GetMessages(conversationID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"messages":   messages,
		"pagination": pagination,
	})
}

func (h *MessageHandler) SendMessage(c *gin.Context) {
	userID, _ := c.Get("user_id")

	conversationID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid conversation id"})
		return
	}

	var req domain.SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	message, err := h.messageUseCase.SendMessage(conversationID, userID.(uuid.UUID), req.Content)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Broadcast to WebSocket connections
	h.broadcast(conversationID, domain.WSMessage{
		Type: domain.WSMessageTypeMessage,
		Data: message,
	})

	c.JSON(http.StatusCreated, message)
}

func (h *MessageHandler) SuggestMessage(c *gin.Context) {
	userID, _ := c.Get("user_id")

	conversationID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid conversation id"})
		return
	}

	suggestion, err := h.messageUseCase.SuggestMessage(conversationID, userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"suggestion": suggestion})
}

func (h *MessageHandler) WebSocketHandler(c *gin.Context) {
	conversationID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid conversation id"})
		return
	}

	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}

	var userID uuid.UUID
	authenticated := false

	// Wait for auth message
	var authMsg domain.WSMessage
	if err := conn.ReadJSON(&authMsg); err != nil {
		conn.Close()
		return
	}

	if authMsg.Type == domain.WSMessageTypeAuth && authMsg.Token != "" {
		validatedUserID, err := h.authUseCase.ValidateToken(authMsg.Token)
		if err == nil {
			userID = validatedUserID
			authenticated = true

			// Register connection
			h.mu.Lock()
			if h.connections[conversationID] == nil {
				h.connections[conversationID] = make(map[*websocket.Conn]bool)
			}
			h.connections[conversationID][conn] = true
			h.mu.Unlock()

			// Send success
			conn.WriteJSON(domain.WSMessage{
				Type: domain.WSMessageTypeAuth,
				Data: gin.H{"status": "authenticated"},
			})
		}
	}

	if !authenticated {
		conn.WriteJSON(domain.WSMessage{
			Type: domain.WSMessageTypeAuth,
			Data: gin.H{"error": "authentication failed"},
		})
		conn.Close()
		return
	}

	// Clean up on disconnect
	defer func() {
		h.mu.Lock()
		delete(h.connections[conversationID], conn)
		if len(h.connections[conversationID]) == 0 {
			delete(h.connections, conversationID)
		}
		h.mu.Unlock()
		conn.Close()
	}()

	// Handle incoming messages
	for {
		var msg domain.WSMessage
		if err := conn.ReadJSON(&msg); err != nil {
			break
		}

		switch msg.Type {
		case domain.WSMessageTypeSend:
			if msg.Content != "" {
				message, err := h.messageUseCase.SendMessage(conversationID, userID, msg.Content)
				if err == nil {
					h.broadcast(conversationID, domain.WSMessage{
						Type: domain.WSMessageTypeMessage,
						Data: message,
					})
				}
			}
		case domain.WSMessageTypeTyping:
			h.broadcast(conversationID, domain.WSMessage{
				Type: domain.WSMessageTypeTyping,
				Data: gin.H{"user_id": userID},
			})
		}
	}
}

func (h *MessageHandler) broadcast(conversationID uuid.UUID, msg domain.WSMessage) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	msgBytes, _ := json.Marshal(msg)

	for conn := range h.connections[conversationID] {
		if err := conn.WriteMessage(websocket.TextMessage, msgBytes); err != nil {
			log.Printf("Error broadcasting: %v", err)
		}
	}
}
