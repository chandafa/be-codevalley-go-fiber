package handlers

import (
	"code-valley-api/internal/models"
	"code-valley-api/internal/services"
	"code-valley-api/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type QuestHandler struct {
	questService *services.QuestService
}

func NewQuestHandler(questService *services.QuestService) *QuestHandler {
	return &QuestHandler{
		questService: questService,
	}
}

func (h *QuestHandler) GetQuests(c *fiber.Ctx) error {
	pagination := utils.GetPaginationParams(c)

	response, err := h.questService.GetQuests(pagination)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse("Failed to fetch quests"))
	}

	return c.JSON(models.SuccessResponse("Quests retrieved successfully", response))
}

func (h *QuestHandler) GetQuestByID(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse("Invalid quest ID"))
	}

	quest, err := h.questService.GetQuestByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse("Quest not found"))
	}

	return c.JSON(models.SuccessResponse("Quest retrieved successfully", quest))
}

func (h *QuestHandler) StartQuest(c *fiber.Ctx) error {
	user := c.Locals("user").(*utils.Claims)
	
	idParam := c.Params("id")
	questID, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse("Invalid quest ID"))
	}

	progress, err := h.questService.StartQuest(user.UserID, questID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse(err.Error()))
	}

	return c.JSON(models.SuccessResponse("Quest started successfully", progress))
}

func (h *QuestHandler) CompleteQuest(c *fiber.Ctx) error {
	user := c.Locals("user").(*utils.Claims)
	
	idParam := c.Params("id")
	questID, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse("Invalid quest ID"))
	}

	var req services.CompleteQuestRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse("Invalid request body"))
	}

	progress, err := h.questService.CompleteQuest(user.UserID, questID, req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse(err.Error()))
	}

	return c.JSON(models.SuccessResponse("Quest completed successfully", progress))
}

func (h *QuestHandler) GetUserProgress(c *fiber.Ctx) error {
	user := c.Locals("user").(*utils.Claims)

	progress, err := h.questService.GetUserProgress(user.UserID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse("Failed to fetch progress"))
	}

	return c.JSON(models.SuccessResponse("Progress retrieved successfully", progress))
}

// Admin handlers
func (h *QuestHandler) CreateQuest(c *fiber.Ctx) error {
	var req services.CreateQuestRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse("Invalid request body"))
	}

	quest, err := h.questService.CreateQuest(req)
	if err != nil {
		if validationErrors := utils.FormatValidationErrors(err); len(validationErrors) > 0 {
			return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
				Success: false,
				Message: "Validation failed",
				Data:    validationErrors,
			})
		}
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse(err.Error()))
	}

	return c.Status(fiber.StatusCreated).JSON(models.SuccessResponse("Quest created successfully", quest))
}

func (h *QuestHandler) UpdateQuest(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse("Invalid quest ID"))
	}

	var req services.CreateQuestRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse("Invalid request body"))
	}

	quest, err := h.questService.UpdateQuest(id, req)
	if err != nil {
		if validationErrors := utils.FormatValidationErrors(err); len(validationErrors) > 0 {
			return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
				Success: false,
				Message: "Validation failed",
				Data:    validationErrors,
			})
		}
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse(err.Error()))
	}

	return c.JSON(models.SuccessResponse("Quest updated successfully", quest))
}

func (h *QuestHandler) DeleteQuest(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse("Invalid quest ID"))
	}

	if err := h.questService.DeleteQuest(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse("Failed to delete quest"))
	}

	return c.JSON(models.SuccessResponse("Quest deleted successfully", nil))
}