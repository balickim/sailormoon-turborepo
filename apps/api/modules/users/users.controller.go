package users

import (
	"reflect"
	"sailormoon/backend/middlewares"
	"sailormoon/backend/modules/users/dtos"
	"sailormoon/backend/utils"

	"github.com/gofiber/fiber/v2"
)

type UserController struct {
	Service *UserService
}

func (uc *UserController) InitializeRoutes(router fiber.Router) {
	router.Post("/users", middlewares.ValidationMiddleware(reflect.TypeOf(dtos.CreateUserDto{})), uc.createUser)
	router.Get("/users", uc.getAllUsers)
}

func (uc *UserController) createUser(c *fiber.Ctx) error {
	dto := c.Locals("validatedData").(*dtos.CreateUserDto)
	user, err := uc.Service.CreateUser(dto.Name, dto.Email)
	if err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Failed to create user"})
	}
	return c.Status(fiber.StatusCreated).JSON(utils.FormatSuccessResponse(user, fiber.StatusCreated))
}

func (uc *UserController) getAllUsers(c *fiber.Ctx) error {
	users, err := uc.Service.GetAllUsers()
	if err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Failed to fetch users"})
	}
	return c.Status(fiber.StatusOK).JSON(utils.FormatSuccessResponse(users, fiber.StatusOK))
}
