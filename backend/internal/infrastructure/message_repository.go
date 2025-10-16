package infrastructure

import (
	"github.com/google/uuid"
	"github.com/yourusername/ecomate/backend/internal/domain"
	"gorm.io/gorm"
)

type messageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) domain.MessageRepository {
	return &messageRepository{db: db}
}

func (r *messageRepository) CreateConversation(conversation *domain.Conversation) error {
	return r.db.Create(conversation).Error
}

func (r *messageRepository) FindConversationByID(id uuid.UUID) (*domain.Conversation, error) {
	var conversation domain.Conversation
	if err := r.db.
		Preload("Product").
		Preload("Participants").
		Preload("Participants.User").
		Where("id = ?", id).
		First(&conversation).Error; err != nil {
		return nil, err
	}
	return &conversation, nil
}

func (r *messageRepository) FindConversationByParticipants(productID uuid.UUID, userID1, userID2 uuid.UUID) (*domain.Conversation, error) {
	var conversation domain.Conversation

	// Find conversation with product and both users as participants
	err := r.db.
		Joins("JOIN conversation_participants cp1 ON cp1.conversation_id = conversations.id AND cp1.user_id = ?", userID1).
		Joins("JOIN conversation_participants cp2 ON cp2.conversation_id = conversations.id AND cp2.user_id = ?", userID2).
		Where("conversations.product_id = ?", productID).
		Preload("Product").
		Preload("Participants").
		Preload("Participants.User").
		First(&conversation).Error

	if err != nil {
		return nil, err
	}

	return &conversation, nil
}

func (r *messageRepository) GetUserConversations(userID uuid.UUID) ([]*domain.Conversation, error) {
	var conversations []*domain.Conversation

	err := r.db.
		Joins("JOIN conversation_participants ON conversation_participants.conversation_id = conversations.id").
		Where("conversation_participants.user_id = ?", userID).
		Preload("Product").
		Preload("Product.Images", "is_primary = true").
		Preload("Participants").
		Preload("Participants.User").
		Preload("Messages", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at DESC").Limit(1)
		}).
		Order("last_message_at DESC").
		Find(&conversations).Error

	if err != nil {
		return nil, err
	}

	return conversations, nil
}

func (r *messageRepository) CreateMessage(message *domain.Message) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Create message
		if err := tx.Create(message).Error; err != nil {
			return err
		}

		// Update conversation last_message_at
		if err := tx.Model(&domain.Conversation{}).
			Where("id = ?", message.ConversationID).
			Update("last_message_at", message.CreatedAt).
			Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *messageRepository) GetMessages(conversationID uuid.UUID, page, limit int) ([]*domain.Message, *domain.PaginationResponse, error) {
	var messages []*domain.Message
	var total int64

	query := r.db.Model(&domain.Message{}).
		Where("conversation_id = ?", conversationID).
		Preload("Sender")

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, nil, err
	}

	// Pagination
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 50
	}

	offset := (page - 1) * limit
	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&messages).Error; err != nil {
		return nil, nil, err
	}

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	pagination := &domain.PaginationResponse{
		Page:       page,
		Limit:      limit,
		Total:      int(total),
		TotalPages: totalPages,
	}

	return messages, pagination, nil
}

func (r *messageRepository) MarkAsRead(messageID, userID uuid.UUID) error {
	return r.db.Model(&domain.Message{}).
		Where("id = ? AND sender_id != ?", messageID, userID).
		Update("is_read", true).
		Error
}
