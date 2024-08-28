package boats

import (
	"sailormoon/backend/utils"

	"github.com/gofiber/fiber/v2"
)

type BoatsController struct {
	Service *BoatsService
}

func (controller *BoatsController) InitializeRoutes(router fiber.Router) {
	router.Get(
		"/",
		controller.getBoats,
	)
}

func (controller *BoatsController) getBoats(c *fiber.Ctx) error {
	params, meta, err := utils.ParseQueryParams(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid query parameters"})
	}

	boats, total, err := controller.Service.GetBoats(params)
	if err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Failed to fetch boats"})
	}

	utils.UpdateMetaWithTotal(meta, total, params.PageSize)
	return c.Status(fiber.StatusOK).JSON(utils.FormatSuccessResponse(boats, fiber.StatusOK, meta))
}
