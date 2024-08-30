package users

import (
	"reflect"
	"sailormoon/backend/middlewares"
	"sailormoon/backend/utils"
	"time"

	"github.com/gofiber/fiber/v2"
)

type UsersController struct {
	Service *UsersService
}

func (uc *UsersController) InitializeRoutes(router fiber.Router) {
	router.Post(
		"/",
		middlewares.ValidationMiddleware(reflect.TypeOf(CreateUserDto{})),
		uc.createUser,
	)
	router.Get(
		"/",
		uc.getAllUsers,
	)
	router.Post(
		"/login",
		uc.loginUser,
	)
}

func (uc *UsersController) createUser(c *fiber.Ctx) error {
	dto := c.Locals("validatedData").(*CreateUserDto)
	user, err := uc.Service.CreateUser(dto.Name, dto.Email, dto.Password)
	if err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(utils.FormatSuccessResponse(user, fiber.StatusCreated))
}

func (uc *UsersController) getAllUsers(c *fiber.Ctx) error {
	users, err := uc.Service.GetAllUsers()
	if err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Failed to fetch users"})
	}
	return c.Status(fiber.StatusOK).JSON(utils.FormatSuccessResponse(users, fiber.StatusOK))
}

func (uc *UsersController) loginUser(c *fiber.Ctx) error {
	var dto struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&dto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	sessionID, err := uc.Service.Login(dto.Email, dto.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		HTTPOnly: true,
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Login successful"})
}
