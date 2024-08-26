package slips

import (
	"encoding/json"
	"math"
	"sailormoon/backend/utils"

	"github.com/gofiber/fiber/v2"
)

type SlipsController struct {
	Service *SlipsService
}

func (controller *SlipsController) InitializeRoutes(router fiber.Router) {
	router.Get(
		"/",
		controller.getSlips,
	)
}

func (controller *SlipsController) getSlips(c *fiber.Ctx) error {
	sortBy := c.Query("sort_by", "id")
	sortOrder := c.Query("sort_order", "asc")
	page := c.QueryInt("page", 1)
	pageSize := c.QueryInt("page_size", 10)
	globalFilter := c.Query("global_filter", "")

	var filters []Filter
	filtersParam := c.Query("filters", "")
	if filtersParam != "" {
		if err := json.Unmarshal([]byte(filtersParam), &filters); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid filters format"})
		}
	}

	params := GetSlipsParams{
		SortBy:       sortBy,
		SortOrder:    sortOrder,
		Page:         page,
		PageSize:     pageSize,
		Filters:      filters,
		GlobalFilter: globalFilter,
	}

	slips, total, err := controller.Service.GetSlips(params)
	if err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Failed to fetch slips"})
	}

	lastPage := int(math.Ceil(float64(total) / float64(pageSize)))
	meta := map[string]interface{}{
		"total":        total,
		"current_page": page,
		"page_size":    pageSize,
		"last_page":    lastPage,
	}

	return c.Status(fiber.StatusOK).JSON(utils.FormatSuccessResponse(slips, fiber.StatusOK, meta))
}
