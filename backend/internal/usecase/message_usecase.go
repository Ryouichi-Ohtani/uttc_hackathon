package usecase

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/yourusername/ecomate/backend/internal/domain"
)

type MessageUseCase interface {
	CreateConversation(userID uuid.UUID, req *domain.CreateConversationRequest) (*domain.Conversation, error)
	GetOrCreateConversation(productID, buyerID, sellerID uuid.UUID) (*domain.Conversation, error)
	GetConversation(conversationID uuid.UUID) (*domain.Conversation, error)
	GetUserConversations(userID uuid.UUID) ([]*domain.Conversation, error)
	SendMessage(conversationID, senderID uuid.UUID, content string) (*domain.Message, error)
	GetMessages(conversationID uuid.UUID, page, limit int) ([]*domain.Message, *domain.PaginationResponse, error)
}

type messageUseCase struct {
	messageRepo domain.MessageRepository
	productRepo domain.ProductRepository
}

func NewMessageUseCase(
	messageRepo domain.MessageRepository,
	productRepo domain.ProductRepository,
) MessageUseCase {
	return &messageUseCase{
		messageRepo: messageRepo,
		productRepo: productRepo,
	}
}

func (u *messageUseCase) CreateConversation(userID uuid.UUID, req *domain.CreateConversationRequest) (*domain.Conversation, error) {
	// Get product to find seller
	product, err := u.productRepo.FindByID(req.ProductID)
	if err != nil {
		return nil, errors.New("product not found")
	}

	// Check if conversation already exists
	existing, _ := u.messageRepo.FindConversationByParticipants(req.ProductID, userID, req.ParticipantID)
	if existing != nil {
		return existing, nil
	}

	// Create new conversation
	now := time.Now()
	conversation := &domain.Conversation{
		ProductID:     &req.ProductID,
		LastMessageAt: now,
		CreatedAt:     now,
		Participants: []domain.ConversationParticipant{
			{
				UserID:     userID,
				LastReadAt: now,
				CreatedAt:  now,
			},
			{
				UserID:     product.SellerID,
				LastReadAt: now,
				CreatedAt:  now,
			},
		},
	}

	if err := u.messageRepo.CreateConversation(conversation); err != nil {
		return nil, err
	}

	return u.messageRepo.FindConversationByID(conversation.ID)
}

func (u *messageUseCase) GetOrCreateConversation(productID, buyerID, sellerID uuid.UUID) (*domain.Conversation, error) {
	// Check if conversation already exists
	existing, _ := u.messageRepo.FindConversationByParticipants(productID, buyerID, sellerID)
	if existing != nil {
		return existing, nil
	}

	// Create new conversation
	now := time.Now()
	conversation := &domain.Conversation{
		ProductID:     &productID,
		LastMessageAt: now,
		CreatedAt:     now,
		Participants: []domain.ConversationParticipant{
			{
				UserID:     buyerID,
				LastReadAt: now,
				CreatedAt:  now,
			},
			{
				UserID:     sellerID,
				LastReadAt: now,
				CreatedAt:  now,
			},
		},
	}

	if err := u.messageRepo.CreateConversation(conversation); err != nil {
		return nil, err
	}

	return u.messageRepo.FindConversationByID(conversation.ID)
}

func (u *messageUseCase) GetConversation(conversationID uuid.UUID) (*domain.Conversation, error) {
	return u.messageRepo.FindConversationByID(conversationID)
}

func (u *messageUseCase) GetUserConversations(userID uuid.UUID) ([]*domain.Conversation, error) {
	return u.messageRepo.GetUserConversations(userID)
}

func (u *messageUseCase) SendMessage(conversationID, senderID uuid.UUID, content string) (*domain.Message, error) {
	// Verify conversation exists and user is participant
	conversation, err := u.messageRepo.FindConversationByID(conversationID)
	if err != nil {
		return nil, errors.New("conversation not found")
	}

	isParticipant := false
	for _, p := range conversation.Participants {
		if p.UserID == senderID {
			isParticipant = true
			break
		}
	}

	if !isParticipant {
		return nil, errors.New("user is not a participant of this conversation")
	}

	// Create message
	message := &domain.Message{
		ConversationID: conversationID,
		SenderID:       senderID,
		Content:        content,
		IsRead:         false,
		CreatedAt:      time.Now(),
	}

	if err := u.messageRepo.CreateMessage(message); err != nil {
		return nil, err
	}

	// Reload with sender info
	message.Sender = &domain.User{}
	for _, p := range conversation.Participants {
		if p.UserID == senderID && p.User != nil {
			message.Sender = p.User
			break
		}
	}

	return message, nil
}

func (u *messageUseCase) GetMessages(conversationID uuid.UUID, page, limit int) ([]*domain.Message, *domain.PaginationResponse, error) {
	return u.messageRepo.GetMessages(conversationID, page, limit)
}
