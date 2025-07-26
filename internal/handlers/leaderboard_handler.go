package handlers

import (
	"code-valley-api/internal/models"
	"code-valley-api/internal/services"

	"github.com/gofiber/fiber/v2"
)

type LeaderboardHandler struct {
	leaderboardService *services.LeaderboardService
}

func NewLeaderboardHandler(leaderboardService *services.LeaderboardService) *LeaderboardHandler {
	return &LeaderboardHandler{
		leaderboardService: leaderboardService,
	}
}

func (h *LeaderboardHandler) GetCoinLeaderboard(c *fiber.Ctx) error {
	leaderboard, err := h.leaderboardService.GetTopUsersByCoins(50)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse("Failed to fetch coin leaderboard"))
	}

	return c.JSON(models.SuccessResponse("Coin leaderboard retrieved successfully", leaderboard))
}

func (h *LeaderboardHandler) GetEXPLeaderboard(c *fiber.Ctx) error {
	leaderboard, err := h.leaderboardService.GetTopUsersByEXP(50)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse("Failed to fetch EXP leaderboard"))
	}

	return c.JSON(models.SuccessResponse("EXP leaderboard retrieved successfully", leaderboard))
}

func (h *LeaderboardHandler) GetTaskLeaderboard(c *fiber.Ctx) error {
	leaderboard, err := h.leaderboardService.GetTopUsersByTasksCompleted(50)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse("Failed to fetch task leaderboard"))
	}

	return c.JSON(models.SuccessResponse("Task leaderboard retrieved successfully", leaderboard))
}