package services

import (
	"errors"

	"code-valley-api/internal/config"
	"code-valley-api/internal/models"
	"code-valley-api/internal/repositories"
	"code-valley-api/internal/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuthService struct {
	userRepo *repositories.UserRepository
	cfg      *config.Config
}

func NewAuthService(cfg *config.Config) *AuthService {
	return &AuthService{
		userRepo: repositories.NewUserRepository(),
		cfg:      cfg,
	}
}

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required,min=3,max=30"`
	Password string `json:"password" validate:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
	Token        string              `json:"token"`
	User         models.UserResponse `json:"user"`
	AccessToken  string              `json:"access_token"`
	RefreshToken string              `json:"refresh_token"`
}

func (s *AuthService) Register(req RegisterRequest) (*AuthResponse, error) {
	// Validate input
	if err := utils.ValidateStruct(req); err != nil {
		return nil, err
	}

	// Check if user already exists
	if _, err := s.userRepo.GetByEmail(req.Email); err == nil {
		return nil, errors.New("email already registered")
	}

	if _, err := s.userRepo.GetByUsername(req.Username); err == nil {
		return nil, errors.New("username already taken")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &models.User{
		Email:        req.Email,
		Username:     req.Username,
		PasswordHash: hashedPassword,
		Role:         models.RolePlayer,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	// Generate token
	token, err := utils.GenerateJWT(
		user.ID,
		user.Email,
		string(user.Role),
		s.cfg.JWT.Secret,
		s.cfg.JWT.ExpireHours,
	)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		Token: token,
		User:  user.ToResponse(),
	}, nil
}

func (s *AuthService) Login(req LoginRequest) (*AuthResponse, error) {
	// Validate input
	if err := utils.ValidateStruct(req); err != nil {
		return nil, err
	}

	// Get user by email
	user, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid credentials")
		}
		return nil, err
	}

	// Check password
	if err := utils.CheckPassword(user.PasswordHash, req.Password); err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Generate token
	token, err := utils.GenerateJWT(
		user.ID,
		user.Email,
		string(user.Role),
		s.cfg.JWT.Secret,
		s.cfg.JWT.ExpireHours,
	)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		Token: token,
		User:  user.ToResponse(),
	}, nil
}

func (s *AuthService) GetProfile(userID uuid.UUID) (*models.UserResponse, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	response := user.ToResponse()
	return &response, nil
}

type UpdateProfileRequest struct {
	Username  string `json:"username" validate:"omitempty,min=3,max=30"`
	Bio       string `json:"bio" validate:"omitempty,max=500"`
	AvatarURL string `json:"avatar_url" validate:"omitempty,url"`
}

func (s *AuthService) UpdateProfile(userID uuid.UUID, req UpdateProfileRequest) (*models.UserResponse, error) {
	// Validate input
	if err := utils.ValidateStruct(req); err != nil {
		return nil, err
	}

	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	// Check if username is taken by another user
	if req.Username != "" && req.Username != user.Username {
		if _, err := s.userRepo.GetByUsername(req.Username); err == nil {
			return nil, errors.New("username already taken")
		}
	}

	// Update fields
	if req.Username != "" {
		user.Username = req.Username
	}
	if req.Bio != "" {
		user.Bio = req.Bio
	}
	if req.AvatarURL != "" {
		user.AvatarURL = req.AvatarURL
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	response := user.ToResponse()
	return &response, nil
}
