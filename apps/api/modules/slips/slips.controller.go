package slips

import (
	"sailormoon/backend/auth"
	"sailormoon/backend/modules/users"
	"sailormoon/backend/utils"

	"github.com/gofiber/fiber/v2"
)

type SlipsController struct {
	Service      *SlipsService
	UsersService *users.UsersService
}

func (controller *SlipsController) InitializeRoutes(router fiber.Router) {
	router.Get(
		"/",
		auth.SessionAuthMiddleware(controller.UsersService),
		controller.getSlips,
	)
}

func (controller *SlipsController) getSlips(c *fiber.Ctx) error {
	params, meta, err := utils.ParseQueryParams(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid query parameters"})
	}

	slips, total, err := controller.Service.GetSlips(params)
	if err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Failed to fetch slips"})
	}

	utils.UpdateMetaWithTotal(meta, total, params.PageSize)
	return c.Status(fiber.StatusOK).JSON(utils.FormatSuccessResponse(slips, fiber.StatusOK, meta))
}
