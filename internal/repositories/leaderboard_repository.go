package repositories

import (
	"code-valley-api/internal/database"
	"code-valley-api/internal/models"

	"gorm.io/gorm"
)

type LeaderboardRepository struct {
	db *gorm.DB
}

func NewLeaderboardRepository() *LeaderboardRepository {
	return &LeaderboardRepository{
		db: database.GetDB(),
	}
}

func (r *LeaderboardRepository) GetTopUsersByCoins(limit int) ([]models.User, error) {
	var users []models.User
	err := r.db.Order("coins DESC").Limit(limit).Find(&users).Error
	return users, err
}

func (r *LeaderboardRepository) GetTopUsersByEXP(limit int) ([]models.User, error) {
	var users []models.User
	err := r.db.Order("exp DESC").Limit(limit).Find(&users).Error
	return users, err
}

func (r *LeaderboardRepository) GetTopUsersByTasksCompleted(limit int) ([]models.UserStatistics, error) {
	var stats []models.UserStatistics
	err := r.db.Preload("User").Order("tasks_completed DESC").Limit(limit).Find(&stats).Error
	return stats, err
}