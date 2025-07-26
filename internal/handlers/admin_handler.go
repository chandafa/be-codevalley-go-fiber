package handlers

import (
	"code-valley-api/internal/models"
	"code-valley-api/internal/services"
	"code-valley-api/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type AdminHandler struct {
	adminService *services.AdminService
}

func NewAdminHandler(adminService *services.AdminService) *AdminHandler {
	return &AdminHandler{
		adminService: adminService,
	}
}

func (h *AdminHandler) GetAllUsers(c *fiber.Ctx) error {
	pagination := utils.GetPaginationParams(c)

	response, err := h.adminService.GetAllUsers(pagination)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse("Failed to fetch users"))
	}

	return c.JSON(models.SuccessResponse("Users retrieved successfully", response))
}

func (h *AdminHandler) BanUser(c *fiber.Ctx) error {
	idParam := c.Params("id")
	userID, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse("Invalid user ID"))
	}

	var req services.BanUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse("Invalid request body"))
	}

	err = h.adminService.BanUser(userID, req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse(err.Error()))
	}

	return c.JSON(models.SuccessResponse("User banned successfully", nil))
}

func (h *AdminHandler) ChangeUserRole(c *fiber.Ctx) error {
	idParam := c.Params("id")
	userID, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse("Invalid user ID"))
	}

	var req services.ChangeRoleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse("Invalid request body"))
	}

	user, err := h.adminService.ChangeUserRole(userID, req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse(err.Error()))
	}

	return c.JSON(models.SuccessResponse("User role updated successfully", user))
}

func (h *AdminHandler) GetSystemStats(c *fiber.Ctx) error {
	stats, err := h.adminService.GetSystemStats()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse("Failed to fetch system stats"))
	}

	return c.JSON(models.SuccessResponse("System stats retrieved successfully", stats))
}

func (h *AdminHandler) GetAuditLogs(c *fiber.Ctx) error {
	pagination := utils.GetPaginationParams(c)

	logs, err := h.adminService.GetAuditLogs(pagination)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse("Failed to fetch audit logs"))
	}

	return c.JSON(models.SuccessResponse("Audit logs retrieved successfully", logs))
}