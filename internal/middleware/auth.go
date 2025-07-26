package middleware

import (
	"strings"

	"code-valley-api/internal/config"
	"code-valley-api/internal/models"
	"code-valley-api/internal/utils"

	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse("Authorization header required"))
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse("Invalid authorization format"))
		}

		claims, err := utils.ValidateJWT(tokenString, cfg.JWT.Secret)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse("Invalid token"))
		}

		c.Locals("user", claims)
		return c.Next()
	}
}

func RequireRole(roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := c.Locals("user").(*utils.Claims)
		
		for _, role := range roles {
			if user.Role == role {
				return c.Next()
			}
		}
		
		return c.Status(fiber.StatusForbidden).JSON(models.ErrorResponse("Insufficient permissions"))
	}
}