package services

import (
	"errors"

	"code-valley-api/internal/models"
	"code-valley-api/internal/repositories"

	"github.com/google/uuid"
)

type NotificationService struct {
	notificationRepo *repositories.NotificationRepository
}

func NewNotificationService() *NotificationService {
	return &NotificationService{
		notificationRepo: repositories.NewNotificationRepository(),
	}
}

func (s *NotificationService) GetUserNotifications(userID uuid.UUID) ([]models.Notification, error) {
	return s.notificationRepo.GetUserNotifications(userID)
}

func (s *NotificationService) MarkAsRead(userID, notificationID uuid.UUID) error {
	// Verify notification belongs to user
	notifications, err := s.notificationRepo.GetUserNotifications(userID)
	if err != nil {
		return err
	}

	found := false
	for _, notification := range notifications {
		if notification.ID == notificationID {
			found = true
			break
		}
	}

	if !found {
		return errors.New("notification not found")
	}

	return s.notificationRepo.MarkAsRead(notificationID)
}

func (s *NotificationService) MarkAllAsRead(userID uuid.UUID) error {
	return s.notificationRepo.MarkAllAsRead(userID)
}

func (s *NotificationService) CreateNotification(userID uuid.UUID, notificationType models.NotificationType, title, message string, data models.NotificationData) error {
	notification := &models.Notification{
		UserID:  userID,
		Type:    notificationType,
		Title:   title,
		Message: message,
		Data:    data,
	}

	return s.notificationRepo.Create(notification)
}