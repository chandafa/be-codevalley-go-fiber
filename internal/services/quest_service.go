package services

import (
	"errors"
	"time"

	"code-valley-api/internal/models"
	"code-valley-api/internal/repositories"
	"code-valley-api/internal/utils"

	"github.com/google/uuid"
	// "gorm.io/gorm"
)

type QuestService struct {
	questRepo *repositories.QuestRepository
	userRepo  *repositories.UserRepository
}

func NewQuestService() *QuestService {
	return &QuestService{
		questRepo: repositories.NewQuestRepository(),
		userRepo:  repositories.NewUserRepository(),
	}
}

func (s *QuestService) GetQuests(pagination utils.PaginationParams) (*models.PaginatedResponse, error) {
	quests, total, err := s.questRepo.GetAll(pagination)
	if err != nil {
		return nil, err
	}

	data := make([]interface{}, len(quests))
	for i, quest := range quests {
		data[i] = quest
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

func (s *QuestService) GetQuestByID(id uuid.UUID) (*models.Quest, error) {
	return s.questRepo.GetByID(id)
}

func (s *QuestService) StartQuest(userID, questID uuid.UUID) (*models.UserQuestProgress, error) {
	// Check if quest exists
	quest, err := s.questRepo.GetByID(questID)
	if err != nil {
		return nil, err
	}

	if !quest.IsActive {
		return nil, errors.New("quest is not active")
	}

	// Check if user already has progress for this quest
	existingProgress, err := s.questRepo.GetUserProgress(userID, questID)
	if err == nil {
		if existingProgress.Status == models.QuestStatusCompleted && !quest.IsRepeatable {
			return nil, errors.New("quest already completed and is not repeatable")
		}
		if existingProgress.Status == models.QuestStatusInProgress {
			return existingProgress, nil
		}
	}

	// Create new progress
	now := time.Now()
	progress := &models.UserQuestProgress{
		UserID:       userID,
		QuestID:      questID,
		Status:       models.QuestStatusInProgress,
		ProgressData: make(models.ProgressData),
		StartedAt:    &now,
	}

	if err := s.questRepo.CreateProgress(progress); err != nil {
		return nil, err
	}

	return progress, nil
}

type CompleteQuestRequest struct {
	SubmittedItems map[string]int `json:"submitted_items"`
}

func (s *QuestService) CompleteQuest(userID, questID uuid.UUID, req CompleteQuestRequest) (*models.UserQuestProgress, error) {
	// Get quest and progress
	quest, err := s.questRepo.GetByID(questID)
	if err != nil {
		return nil, err
	}

	progress, err := s.questRepo.GetUserProgress(userID, questID)
	if err != nil {
		return nil, errors.New("quest not started")
	}

	if progress.Status == models.QuestStatusCompleted {
		return nil, errors.New("quest already completed")
	}

	// Validate required items (simplified validation)
	for item, required := range quest.RequiredItems {
		if submitted, ok := req.SubmittedItems[item]; !ok || submitted < required {
			return nil, errors.New("insufficient items to complete quest")
		}
	}

	// Update progress
	now := time.Now()
	progress.Status = models.QuestStatusCompleted
	progress.CompletedAt = &now

	if err := s.questRepo.UpdateProgress(progress); err != nil {
		return nil, err
	}

	// Reward user
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	user.Coins += quest.RewardCoins
	user.EXP += quest.RewardEXP

	// Simple level calculation
	if user.EXP >= user.Level*100 {
		user.Level++
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return progress, nil
}

func (s *QuestService) GetUserProgress(userID uuid.UUID) ([]models.UserQuestProgress, error) {
	return s.questRepo.GetUserAllProgress(userID)
}

type CreateQuestRequest struct {
	Title         string               `json:"title" validate:"required"`
	Description   string               `json:"description" validate:"required"`
	RewardCoins   int                  `json:"reward_coins"`
	RewardEXP     int                  `json:"reward_exp"`
	RequiredItems models.RequiredItems `json:"required_items"`
	IsRepeatable  bool                 `json:"is_repeatable"`
	IsActive      bool                 `json:"is_active"`
}

func (s *QuestService) CreateQuest(req CreateQuestRequest) (*models.Quest, error) {
	if err := utils.ValidateStruct(req); err != nil {
		return nil, err
	}

	quest := &models.Quest{
		Title:         req.Title,
		Description:   req.Description,
		RewardCoins:   req.RewardCoins,
		RewardEXP:     req.RewardEXP,
		RequiredItems: req.RequiredItems,
		IsRepeatable:  req.IsRepeatable,
		IsActive:      req.IsActive,
	}

	if err := s.questRepo.Create(quest); err != nil {
		return nil, err
	}

	return quest, nil
}

func (s *QuestService) UpdateQuest(id uuid.UUID, req CreateQuestRequest) (*models.Quest, error) {
	if err := utils.ValidateStruct(req); err != nil {
		return nil, err
	}

	quest, err := s.questRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	quest.Title = req.Title
	quest.Description = req.Description
	quest.RewardCoins = req.RewardCoins
	quest.RewardEXP = req.RewardEXP
	quest.RequiredItems = req.RequiredItems
	quest.IsRepeatable = req.IsRepeatable
	quest.IsActive = req.IsActive

	if err := s.questRepo.Update(quest); err != nil {
		return nil, err
	}

	return quest, nil
}

func (s *QuestService) DeleteQuest(id uuid.UUID) error {
	return s.questRepo.Delete(id)
}
