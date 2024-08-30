package auth

import (
	"sailormoon/backend/database"

	"github.com/gofiber/fiber/v2"
)

type SessionService interface {
	GetSession(sessionID string) (database.UsersEntity, error)
}

func SessionAuthMiddleware(service SessionService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		sessionID := c.Cookies("session_id")
		if sessionID == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		}

		user, err := service.GetSession(sessionID)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid session"})
		}

		c.Locals("user", user)
		return c.Next()
	}
}
