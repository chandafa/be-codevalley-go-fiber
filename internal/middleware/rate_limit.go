package middleware

import (
	"time"

	"code-valley-api/internal/config"
	"code-valley-api/internal/models"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

func RateLimitMiddleware(cfg *config.Config) fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        cfg.RateLimit.Max,
		Expiration: time.Duration(cfg.RateLimit.Expiration) * time.Minute,
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(
				models.ErrorResponse("Rate limit exceeded"),
			)
		},
	})
}