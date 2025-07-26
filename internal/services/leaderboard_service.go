package services

import (
	"code-valley-api/internal/models"
	"code-valley-api/internal/repositories"
)

type LeaderboardService struct {
	leaderboardRepo *repositories.LeaderboardRepository
}

func NewLeaderboardService() *LeaderboardService {
	return &LeaderboardService{
		leaderboardRepo: repositories.NewLeaderboardRepository(),
	}
}

func (s *LeaderboardService) GetTopUsersByCoins(limit int) ([]models.User, error) {
	return s.leaderboardRepo.GetTopUsersByCoins(limit)
}

func (s *LeaderboardService) GetTopUsersByEXP(limit int) ([]models.User, error) {
	return s.leaderboardRepo.GetTopUsersByEXP(limit)
}

func (s *LeaderboardService) GetTopUsersByTasksCompleted(limit int) ([]models.UserStatistics, error) {
	return s.leaderboardRepo.GetTopUsersByTasksCompleted(limit)
}