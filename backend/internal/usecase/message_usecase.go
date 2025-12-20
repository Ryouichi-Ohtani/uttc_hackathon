package usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/yourusername/ecomate/backend/internal/domain"
	"github.com/yourusername/ecomate/backend/internal/infrastructure"
)

type MessageUseCase interface {
	CreateConversation(userID uuid.UUID, req *domain.CreateConversationRequest) (*domain.Conversation, error)
	GetOrCreateConversation(productID, buyerID, sellerID uuid.UUID) (*domain.Conversation, error)
	GetConversation(conversationID uuid.UUID) (*domain.Conversation, error)
	GetUserConversations(userID uuid.UUID) ([]*domain.Conversation, error)
	SendMessage(conversationID, senderID uuid.UUID, content string) (*domain.Message, error)
	GetMessages(conversationID uuid.UUID, page, limit int) ([]*domain.Message, *domain.PaginationResponse, error)
	SuggestMessage(conversationID, senderID uuid.UUID) (string, error)
}

type messageUseCase struct {
	messageRepo domain.MessageRepository
	productRepo domain.ProductRepository
	aiClient    *infrastructure.AIClient
}

func NewMessageUseCase(
	messageRepo domain.MessageRepository,
	productRepo domain.ProductRepository,
	aiClient *infrastructure.AIClient,
) MessageUseCase {
	return &messageUseCase{
		messageRepo: messageRepo,
		productRepo: productRepo,
		aiClient:    aiClient,
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

func (u *messageUseCase) SuggestMessage(conversationID, senderID uuid.UUID) (string, error) {
	if u.aiClient == nil {
		return "", errors.New("AIサービスが利用できません")
	}

	conversation, err := u.messageRepo.FindConversationByID(conversationID)
	if err != nil {
		return "", err
	}

	if conversation.Product == nil {
		return "", errors.New("商品情報が不足しています")
	}

	messages, _, err := u.messageRepo.GetMessages(conversationID, 1, 50)
	if err != nil {
		return "", err
	}

	var history []string
	for i := len(messages) - 1; i >= 0; i-- {
		msg := messages[i]
		roleLabel := "購入者"
		if conversation.Product != nil && msg.SenderID == conversation.Product.SellerID {
			roleLabel = "出品者"
		}
		history = append(history, fmt.Sprintf("%s: %s", roleLabel, msg.Content))
	}

	var imageURLs []string
	for _, img := range conversation.Product.Images {
		if img.CDNURL != "" {
			imageURLs = append(imageURLs, img.CDNURL)
		} else {
			imageURLs = append(imageURLs, img.ImageURL)
		}
		if len(imageURLs) >= 3 {
			break
		}
	}

	role := "購入者"
	if conversation.Product != nil && senderID == conversation.Product.SellerID {
		role = "出品者"
	}

	var promptBuilder strings.Builder
	fmt.Fprintf(&promptBuilder, "あなたはEcoMateのAIメッセージアシスタントです。\n")
	fmt.Fprintf(&promptBuilder, "以下の情報をもとに、次に送る%sとして丁寧な日本語メッセージを1件だけ出力してください。\n\n", role)
	fmt.Fprintf(&promptBuilder, "商品名: %s\n", conversation.Product.Title)
	fmt.Fprintf(&promptBuilder, "説明: %s\n", conversation.Product.Description)
	fmt.Fprintf(&promptBuilder, "価格: ¥%d\n", conversation.Product.Price)
	fmt.Fprintf(&promptBuilder, "カテゴリ: %s\n", conversation.Product.Category)
	fmt.Fprintf(&promptBuilder, "状態: %s\n", conversation.Product.Condition)
	if len(imageURLs) > 0 {
		fmt.Fprintf(&promptBuilder, "画像URL: %s\n", strings.Join(imageURLs, ", "))
	}
	if len(history) > 0 {
		fmt.Fprintf(&promptBuilder, "\n直近の交渉履歴:\n%s\n", strings.Join(history, "\n"))
	}
	fmt.Fprintf(&promptBuilder, "\n%sとして次のメッセージを送りたいことを踏まえて、具体的かつ丁寧な文章を作成してください。", role)

	req := &infrastructure.MessageSuggestionRequest{
		Role:               role,
		ProductTitle:       conversation.Product.Title,
		ProductDescription: conversation.Product.Description,
		ProductPrice:       conversation.Product.Price,
		ProductCategory:    conversation.Product.Category,
		ProductCondition:   string(conversation.Product.Condition),
		ImageURLs:          imageURLs,
		History:            history,
	}
	result, err := u.aiClient.GenerateMessageSuggestion(context.Background(), req)
	if err != nil {
		return "", fmt.Errorf("AIメッセージ提案に失敗しました: %w", err)
	}

	return strings.TrimSpace(result), nil
}
