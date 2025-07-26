package repositories

import (
	"code-valley-api/internal/database"
	"code-valley-api/internal/models"
	"code-valley-api/internal/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type QuestRepository struct {
	db *gorm.DB
}

func NewQuestRepository() *QuestRepository {
	return &QuestRepository{
		db: database.GetDB(),
	}
}

func (r *QuestRepository) Create(quest *models.Quest) error {
	return r.db.Create(quest).Error
}

func (r *QuestRepository) GetByID(id uuid.UUID) (*models.Quest, error) {
	var quest models.Quest
	err := r.db.First(&quest, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &quest, nil
}

func (r *QuestRepository) GetAll(pagination utils.PaginationParams) ([]models.Quest, int64, error) {
	var quests []models.Quest
	var total int64

	query := r.db.Model(&models.Quest{}).Where("is_active = ?", true)
	
	query.Count(&total)
	
	err := query.Offset(pagination.Offset).
		Limit(pagination.PerPage).
		Find(&quests).Error

	return quests, total, err
}

func (r *QuestRepository) Update(quest *models.Quest) error {
	return r.db.Save(quest).Error
}

func (r *QuestRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Quest{}, "id = ?", id).Error
}

func (r *QuestRepository) GetUserProgress(userID uuid.UUID, questID uuid.UUID) (*models.UserQuestProgress, error) {
	var progress models.UserQuestProgress
	err := r.db.Preload("Quest").
		First(&progress, "user_id = ? AND quest_id = ?", userID, questID).Error
	if err != nil {
		return nil, err
	}
	return &progress, nil
}

func (r *QuestRepository) CreateProgress(progress *models.UserQuestProgress) error {
	return r.db.Create(progress).Error
}

func (r *QuestRepository) UpdateProgress(progress *models.UserQuestProgress) error {
	return r.db.Save(progress).Error
}

func (r *QuestRepository) GetUserAllProgress(userID uuid.UUID) ([]models.UserQuestProgress, error) {
	var progress []models.UserQuestProgress
	err := r.db.Preload("Quest").
		Find(&progress, "user_id = ?", userID).Error
	return progress, err
}