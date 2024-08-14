package utils

import "github.com/gofiber/fiber/v2"

func FormatSuccessResponse(data interface{}, statusCode int) fiber.Map {
	return fiber.Map{
		"success":    true,
		"statusCode": statusCode,
		"data":       data,
	}
}
