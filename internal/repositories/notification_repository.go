package repositories

import (
	"code-valley-api/internal/database"
	"code-valley-api/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type NotificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository() *NotificationRepository {
	return &NotificationRepository{
		db: database.GetDB(),
	}
}

func (r *NotificationRepository) Create(notification *models.Notification) error {
	return r.db.Create(notification).Error
}

func (r *NotificationRepository) GetUserNotifications(userID uuid.UUID) ([]models.Notification, error) {
	var notifications []models.Notification
	err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&notifications).Error
	return notifications, err
}

func (r *NotificationRepository) MarkAsRead(notificationID uuid.UUID) error {
	return r.db.Model(&models.Notification{}).Where("id = ?", notificationID).Update("is_read", true).Error
}

func (r *NotificationRepository) MarkAllAsRead(userID uuid.UUID) error {
	return r.db.Model(&models.Notification{}).Where("user_id = ? AND is_read = ?", userID, false).Update("is_read", true).Error
}

func (r *NotificationRepository) GetUnreadCount(userID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.Model(&models.Notification{}).Where("user_id = ? AND is_read = ?", userID, false).Count(&count).Error
	return count, err
}