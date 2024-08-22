package slips

import (
	"sailormoon/backend/utils"

	"github.com/gofiber/fiber/v2"
)

type SlipsController struct {
	Service *SlipsService
}

func (uc *SlipsController) InitializeRoutes(router fiber.Router) {
	router.Get(
		"/slips",
		uc.getSlips,
	)
}

func (uc *SlipsController) getSlips(c *fiber.Ctx) error {
	users, err := uc.Service.GetSlips()
	if err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Failed to fetch users"})
	}
	return c.Status(fiber.StatusOK).JSON(utils.FormatSuccessResponse(users, fiber.StatusOK))
}
