package main

import (
	"sailormoon/backend/modules/users"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	app := fiber.New()

	// Enable CORS with default settings
	app.Use(cors.New())

	userService := &users.UserService{}
	userController := &users.UserController{Service: userService}

	api := app.Group("/api")
	userController.InitializeRoutes(api)

	app.Listen("127.0.0.1:3000")
}
