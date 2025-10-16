package domain

import (
	"time"

	"github.com/google/uuid"
)

type Conversation struct {
	ID            uuid.UUID                  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ProductID     *uuid.UUID                 `json:"product_id" gorm:"type:uuid"`
	Product       *Product                   `json:"product,omitempty" gorm:"foreignKey:ProductID"`
	Participants  []ConversationParticipant  `json:"participants,omitempty" gorm:"foreignKey:ConversationID"`
	Messages      []Message                  `json:"messages,omitempty" gorm:"foreignKey:ConversationID"`
	LastMessageAt time.Time                  `json:"last_message_at"`
	CreatedAt     time.Time                  `json:"created_at"`
}

type ConversationParticipant struct {
	ID             uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ConversationID uuid.UUID `json:"conversation_id" gorm:"type:uuid;not null;index"`
	UserID         uuid.UUID `json:"user_id" gorm:"type:uuid;not null;index"`
	User           *User     `json:"user,omitempty" gorm:"foreignKey:UserID"`
	LastReadAt     time.Time `json:"last_read_at"`
	CreatedAt      time.Time `json:"created_at"`
}

type Message struct {
	ID             uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ConversationID uuid.UUID `json:"conversation_id" gorm:"type:uuid;not null;index"`
	SenderID       uuid.UUID `json:"sender_id" gorm:"type:uuid;not null"`
	Sender         *User     `json:"sender,omitempty" gorm:"foreignKey:SenderID"`
	Content        string    `json:"content" gorm:"not null"`
	IsRead         bool      `json:"is_read" gorm:"default:false"`
	CreatedAt      time.Time `json:"created_at"`
}

type MessageRepository interface {
	CreateConversation(conversation *Conversation) error
	FindConversationByID(id uuid.UUID) (*Conversation, error)
	FindConversationByParticipants(productID uuid.UUID, userID1, userID2 uuid.UUID) (*Conversation, error)
	GetUserConversations(userID uuid.UUID) ([]*Conversation, error)

	CreateMessage(message *Message) error
	GetMessages(conversationID uuid.UUID, page, limit int) ([]*Message, *PaginationResponse, error)
	MarkAsRead(messageID, userID uuid.UUID) error
}

type CreateConversationRequest struct {
	ProductID     uuid.UUID `json:"product_id" binding:"required"`
	ParticipantID uuid.UUID `json:"participant_id" binding:"required"`
}

type SendMessageRequest struct {
	Content string `json:"content" binding:"required,min=1"`
}

// WebSocket message types
type WSMessageType string

const (
	WSMessageTypeAuth    WSMessageType = "auth"
	WSMessageTypeMessage WSMessageType = "message"
	WSMessageTypeTyping  WSMessageType = "typing"
	WSMessageTypeRead    WSMessageType = "read"
	WSMessageTypeSend    WSMessageType = "send_message"
	WSMessageTypeMarkRead WSMessageType = "mark_read"
)

type WSMessage struct {
	Type WSMessageType `json:"type"`
	Data interface{}   `json:"data,omitempty"`
	Token string       `json:"token,omitempty"`
	Content string     `json:"content,omitempty"`
	MessageID *uuid.UUID `json:"message_id,omitempty"`
}
