package handlers

import (
	"code-valley-api/internal/models"
	"code-valley-api/internal/services"
	"code-valley-api/internal/utils"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req services.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse("Invalid request body"))
	}

	response, err := h.authService.Register(req)
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

	return c.Status(fiber.StatusCreated).JSON(models.SuccessResponse("User registered successfully", response))
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req services.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse("Invalid request body"))
	}

	response, err := h.authService.Login(req)
	if err != nil {
		if validationErrors := utils.FormatValidationErrors(err); len(validationErrors) > 0 {
			return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
				Success: false,
				Message: "Validation failed",
				Data:    validationErrors,
			})
		}
		return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse(err.Error()))
	}

	// return c.JSON(models.SuccessResponse("Login successful", response))
	return c.JSON(models.SuccessResponse("Login successful", map[string]interface{}{
		"token": response.Token,
		"user":  response.User,
	}))

}

func (h *AuthHandler) GetProfile(c *fiber.Ctx) error {
	user := c.Locals("user").(*utils.Claims)

	profile, err := h.authService.GetProfile(user.UserID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse("User not found"))
	}

	return c.JSON(models.SuccessResponse("Profile retrieved successfully", profile))
}

func (h *AuthHandler) UpdateProfile(c *fiber.Ctx) error {
	user := c.Locals("user").(*utils.Claims)

	var req services.UpdateProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse("Invalid request body"))
	}

	profile, err := h.authService.UpdateProfile(user.UserID, req)
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

	return c.JSON(models.SuccessResponse("Profile updated successfully", profile))
}

func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	// For simplicity, we'll just return the current user info
	// In a real application, you might want to implement proper refresh token logic
	user := c.Locals("user").(*utils.Claims)

	profile, err := h.authService.GetProfile(user.UserID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse("User not found"))
	}

	return c.JSON(models.SuccessResponse("Token refreshed successfully", map[string]interface{}{
		"user": profile,
	}))
}
