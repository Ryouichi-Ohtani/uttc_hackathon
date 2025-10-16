package usecase

import (
	"time"

	"github.com/google/uuid"
	"github.com/yourusername/ecomate/backend/internal/domain"
)

type NotificationUseCase interface {
	Create(userID uuid.UUID, notifType domain.NotificationType, title, message, link string) error
	GetUserNotifications(userID uuid.UUID, limit int) ([]*domain.Notification, error)
	MarkAsRead(id uuid.UUID) error
	MarkAllAsRead(userID uuid.UUID) error
	GetUnreadCount(userID uuid.UUID) (int64, error)
}

type notificationUseCase struct {
	notificationRepo domain.NotificationRepository
}

func NewNotificationUseCase(notificationRepo domain.NotificationRepository) NotificationUseCase {
	return &notificationUseCase{
		notificationRepo: notificationRepo,
	}
}

func (u *notificationUseCase) Create(userID uuid.UUID, notifType domain.NotificationType, title, message, link string) error {
	notification := &domain.Notification{
		UserID:    userID,
		Type:      notifType,
		Title:     title,
		Message:   message,
		Link:      link,
		IsRead:    false,
		CreatedAt: time.Now(),
	}
	return u.notificationRepo.Create(notification)
}

func (u *notificationUseCase) GetUserNotifications(userID uuid.UUID, limit int) ([]*domain.Notification, error) {
	return u.notificationRepo.FindByUserID(userID, limit)
}

func (u *notificationUseCase) MarkAsRead(id uuid.UUID) error {
	return u.notificationRepo.MarkAsRead(id)
}

func (u *notificationUseCase) MarkAllAsRead(userID uuid.UUID) error {
	return u.notificationRepo.MarkAllAsRead(userID)
}

func (u *notificationUseCase) GetUnreadCount(userID uuid.UUID) (int64, error) {
	return u.notificationRepo.GetUnreadCount(userID)
}
