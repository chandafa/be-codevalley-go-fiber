package services

import (
	"errors"
	// "time"

	"code-valley-api/internal/models"
	"code-valley-api/internal/repositories"
	"code-valley-api/internal/websocket"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FriendService struct {
	friendRepo *repositories.FriendRepository
	userRepo   *repositories.UserRepository
}

func NewFriendService() *FriendService {
	return &FriendService{
		friendRepo: repositories.NewFriendRepository(),
		userRepo:   repositories.NewUserRepository(),
	}
}

func (s *FriendService) GetUserFriends(userID uuid.UUID) ([]models.User, error) {
	return s.friendRepo.GetUserFriends(userID)
}

func (s *FriendService) SendFriendRequest(requesterID uuid.UUID, username string) error {
	// Get addressee by username
	addressee, err := s.userRepo.GetByUsername(username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	if requesterID == addressee.ID {
		return errors.New("cannot send friend request to yourself")
	}

	// Check if friendship already exists
	existing, err := s.friendRepo.GetFriendship(requesterID, addressee.ID)
	if err == nil {
		switch existing.Status {
		case models.FriendshipStatusAccepted:
			return errors.New("already friends")
		case models.FriendshipStatusPending:
			return errors.New("friend request already sent")
		case models.FriendshipStatusBlocked:
			return errors.New("cannot send friend request")
		}
	}

	// Create friendship request
	friendship := &models.Friendship{
		RequesterID: requesterID,
		AddresseeID: addressee.ID,
		Status:      models.FriendshipStatusPending,
	}

	if err := s.friendRepo.CreateFriendship(friendship); err != nil {
		return err
	}

	// Send real-time notification
	websocket.NotifyFriendRequest(addressee.ID, map[string]interface{}{
		"requester_id": requesterID,
		"type":         "friend_request",
	})

	return nil
}

func (s *FriendService) AcceptFriendRequest(addresseeID uuid.UUID, username string) error {
	// Get requester by username
	requester, err := s.userRepo.GetByUsername(username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	// Get friendship
	friendship, err := s.friendRepo.GetFriendship(requester.ID, addresseeID)
	if err != nil {
		return errors.New("friend request not found")
	}

	if friendship.Status != models.FriendshipStatusPending {
		return errors.New("friend request is not pending")
	}

	if friendship.AddresseeID != addresseeID {
		return errors.New("you cannot accept this friend request")
	}

	// Update friendship status
	friendship.Status = models.FriendshipStatusAccepted
	if err := s.friendRepo.UpdateFriendship(friendship); err != nil {
		return err
	}

	// Send real-time notification to requester
	websocket.NotifyFriendRequest(requester.ID, map[string]interface{}{
		"addressee_id": addresseeID,
		"type":         "friend_request_accepted",
	})

	return nil
}

func (s *FriendService) RemoveFriend(userID uuid.UUID, username string) error {
	// Get friend by username
	friend, err := s.userRepo.GetByUsername(username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	// Delete friendship
	if err := s.friendRepo.DeleteFriendship(userID, friend.ID); err != nil {
		return err
	}

	return nil
}

func (s *FriendService) GetOnlineFriends(userID uuid.UUID) ([]models.User, error) {
	return s.friendRepo.GetOnlineFriends(userID)
}
