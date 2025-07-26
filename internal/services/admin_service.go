package services

import (
	"errors"
	"time"

	"code-valley-api/internal/models"
	"code-valley-api/internal/repositories"
	"code-valley-api/internal/utils"

	"github.com/google/uuid"
)

type AdminService struct {
	userRepo *repositories.UserRepository
}

func NewAdminService() *AdminService {
	return &AdminService{
		userRepo: repositories.NewUserRepository(),
	}
}

type BanUserRequest struct {
	Reason   string `json:"reason" validate:"required"`
	Duration int    `json:"duration"` // days, 0 for permanent
}

type ChangeRoleRequest struct {
	Role string `json:"role" validate:"required,oneof=player admin"`
}

func (s *AdminService) GetAllUsers(pagination utils.PaginationParams) (*models.PaginatedResponse, error) {
	users, total, err := s.userRepo.GetAll(pagination)
	if err != nil {
		return nil, err
	}

	data := make([]interface{}, len(users))
	for i, user := range users {
		data[i] = user.ToResponse()
	}

	totalPages := int(total) / pagination.PerPage
	if int(total)%pagination.PerPage > 0 {
		totalPages++
	}

	return &models.PaginatedResponse{
		Data: data,
		Meta: models.PaginationMeta{
			CurrentPage: pagination.Page,
			PerPage:     pagination.PerPage,
			Total:       int(total),
			TotalPages:  totalPages,
		},
	}, nil
}

func (s *AdminService) BanUser(userID uuid.UUID, req BanUserRequest) error {
	if err := utils.ValidateStruct(req); err != nil {
		return err
	}

	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	if user.Role == models.RoleAdmin {
		return errors.New("cannot ban admin users")
	}

	// In a real implementation, you'd have a bans table
	// For now, we'll just mark the user as banned in a comment or separate field
	user.Bio = "BANNED: " + req.Reason
	
	return s.userRepo.Update(user)
}

func (s *AdminService) ChangeUserRole(userID uuid.UUID, req ChangeRoleRequest) (*models.UserResponse, error) {
	if err := utils.ValidateStruct(req); err != nil {
		return nil, err
	}

	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	user.Role = models.UserRole(req.Role)
	
	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	response := user.ToResponse()
	return &response, nil
}

func (s *AdminService) GetSystemStats() (map[string]interface{}, error) {
	// In a real implementation, you'd query various tables for statistics
	stats := map[string]interface{}{
		"total_users":        1000,
		"active_users":       750,
		"total_quests":       50,
		"completed_quests":   2500,
		"total_guilds":       25,
		"online_users":       150,
		"server_uptime":      "5 days, 12 hours",
		"last_updated":       time.Now(),
	}

	return stats, nil
}

func (s *AdminService) GetAuditLogs(pagination utils.PaginationParams) (*models.PaginatedResponse, error) {
	// In a real implementation, you'd have an audit_logs table
	// For now, we'll return mock data
	logs := []map[string]interface{}{
		{
			"id":        uuid.New(),
			"action":    "user_login",
			"user_id":   uuid.New(),
			"timestamp": time.Now().Add(-1 * time.Hour),
			"details":   "User logged in successfully",
		},
		{
			"id":        uuid.New(),
			"action":    "quest_completed",
			"user_id":   uuid.New(),
			"timestamp": time.Now().Add(-2 * time.Hour),
			"details":   "Quest 'First Steps' completed",
		},
	}

	data := make([]interface{}, len(logs))
	for i, log := range logs {
		data[i] = log
	}

	return &models.PaginatedResponse{
		Data: data,
		Meta: models.PaginationMeta{
			CurrentPage: pagination.Page,
			PerPage:     pagination.PerPage,
			Total:       len(logs),
			TotalPages:  1,
		},
	}, nil
}