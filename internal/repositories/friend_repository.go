package repositories

import (
	"code-valley-api/internal/database"
	"code-valley-api/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FriendRepository struct {
	db *gorm.DB
}

func NewFriendRepository() *FriendRepository {
	return &FriendRepository{
		db: database.GetDB(),
	}
}

func (r *FriendRepository) CreateFriendship(friendship *models.Friendship) error {
	return r.db.Create(friendship).Error
}

func (r *FriendRepository) GetFriendship(requesterID, addresseeID uuid.UUID) (*models.Friendship, error) {
	var friendship models.Friendship
	err := r.db.Where("(requester_id = ? AND addressee_id = ?) OR (requester_id = ? AND addressee_id = ?)",
		requesterID, addresseeID, addresseeID, requesterID).First(&friendship).Error
	return &friendship, err
}

func (r *FriendRepository) UpdateFriendship(friendship *models.Friendship) error {
	return r.db.Save(friendship).Error
}

func (r *FriendRepository) DeleteFriendship(requesterID, addresseeID uuid.UUID) error {
	return r.db.Where("(requester_id = ? AND addressee_id = ?) OR (requester_id = ? AND addressee_id = ?)",
		requesterID, addresseeID, addresseeID, requesterID).Delete(&models.Friendship{}).Error
}

func (r *FriendRepository) GetUserFriends(userID uuid.UUID) ([]models.User, error) {
	var friends []models.User
	err := r.db.Table("users").
		Joins("JOIN friendships ON (users.id = friendships.requester_id OR users.id = friendships.addressee_id)").
		Where("(friendships.requester_id = ? OR friendships.addressee_id = ?) AND friendships.status = ? AND users.id != ?",
			userID, userID, models.FriendshipStatusAccepted, userID).
		Find(&friends).Error
	return friends, err
}

func (r *FriendRepository) GetPendingRequests(userID uuid.UUID) ([]models.Friendship, error) {
	var requests []models.Friendship
	err := r.db.Preload("Requester").Preload("Addressee").
		Where("addressee_id = ? AND status = ?", userID, models.FriendshipStatusPending).
		Find(&requests).Error
	return requests, err
}

func (r *FriendRepository) GetOnlineFriends(userID uuid.UUID) ([]models.User, error) {
	var friends []models.User
	err := r.db.Table("users").
		Joins("JOIN friendships ON (users.id = friendships.requester_id OR users.id = friendships.addressee_id)").
		Joins("JOIN online_users ON users.id = online_users.user_id").
		Where("(friendships.requester_id = ? OR friendships.addressee_id = ?) AND friendships.status = ? AND users.id != ? AND online_users.is_online = ?",
			userID, userID, models.FriendshipStatusAccepted, userID, true).
		Find(&friends).Error
	return friends, err
}