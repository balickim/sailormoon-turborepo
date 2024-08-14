package middlewares

import (
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate = validator.New()

func ValidationMiddleware(dtoType reflect.Type) fiber.Handler {
	return func(c *fiber.Ctx) error {
		dto := reflect.New(dtoType).Interface()
		if err := c.BodyParser(dto); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
		}
		if err := validate.Struct(dto); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		c.Locals("validatedData", dto)
		return c.Next()
	}
}
