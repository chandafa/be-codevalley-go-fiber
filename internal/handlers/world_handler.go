package handlers

import (
	"code-valley-api/internal/models"
	"code-valley-api/internal/services"
	"code-valley-api/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type WorldHandler struct {
	worldService *services.WorldService
}

func NewWorldHandler(worldService *services.WorldService) *WorldHandler {
	return &WorldHandler{
		worldService: worldService,
	}
}

func (h *WorldHandler) GetMapState(c *fiber.Ctx) error {
	user := c.Locals("user").(*utils.Claims)
	mapName := c.Params("map_name")

	mapState, err := h.worldService.GetMapState(mapName, user.UserID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse(err.Error()))
	}

	return c.JSON(models.SuccessResponse("Map state retrieved successfully", mapState))
}

func (h *WorldHandler) GetPlayerPosition(c *fiber.Ctx) error {
	user := c.Locals("user").(*utils.Claims)

	position, err := h.worldService.GetPlayerPosition(user.UserID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse("Player position not found"))
	}

	return c.JSON(models.SuccessResponse("Player position retrieved successfully", position))
}

func (h *WorldHandler) TeleportPlayer(c *fiber.Ctx) error {
	user := c.Locals("user").(*utils.Claims)

	var req services.TeleportRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse("Invalid request body"))
	}

	position, err := h.worldService.TeleportPlayer(user.UserID, req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse(err.Error()))
	}

	return c.JSON(models.SuccessResponse("Player teleported successfully", position))
}

func (h *WorldHandler) GetGameTime(c *fiber.Ctx) error {
	gameTime, err := h.worldService.GetGameTime()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse("Failed to get game time"))
	}

	return c.JSON(models.SuccessResponse("Game time retrieved successfully", gameTime))
}

func (h *WorldHandler) GetCodeFarms(c *fiber.Ctx) error {
	user := c.Locals("user").(*utils.Claims)

	farms, err := h.worldService.GetUserCodeFarms(user.UserID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse("Failed to get code farms"))
	}

	return c.JSON(models.SuccessResponse("Code farms retrieved successfully", farms))
}

func (h *WorldHandler) PlantCode(c *fiber.Ctx) error {
	user := c.Locals("user").(*utils.Claims)

	var req services.PlantCodeRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse("Invalid request body"))
	}

	farm, err := h.worldService.PlantCode(user.UserID, req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse(err.Error()))
	}

	return c.JSON(models.SuccessResponse("Code planted successfully", farm))
}

func (h *WorldHandler) WaterCode(c *fiber.Ctx) error {
	user := c.Locals("user").(*utils.Claims)

	idParam := c.Params("id")
	farmID, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse("Invalid farm ID"))
	}

	farm, err := h.worldService.WaterCode(user.UserID, farmID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse(err.Error()))
	}

	return c.JSON(models.SuccessResponse("Code watered successfully", farm))
}

func (h *WorldHandler) HarvestCode(c *fiber.Ctx) error {
	user := c.Locals("user").(*utils.Claims)

	idParam := c.Params("id")
	farmID, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse("Invalid farm ID"))
	}

	result, err := h.worldService.HarvestCode(user.UserID, farmID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse(err.Error()))
	}

	return c.JSON(models.SuccessResponse("Code harvested successfully", result))
}