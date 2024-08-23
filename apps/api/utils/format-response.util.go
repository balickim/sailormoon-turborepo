package utils

import "github.com/gofiber/fiber/v2"

func FormatSuccessResponse(data interface{}, statusCode int, meta ...map[string]interface{}) fiber.Map {
	response := fiber.Map{
		"success":    true,
		"statusCode": statusCode,
		"data":       data,
	}

	if len(meta) > 0 && meta[0] != nil {
		response["meta"] = meta[0]
	}

	return response
}
