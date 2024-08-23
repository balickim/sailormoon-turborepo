package main

import (
	"log"
	"sailormoon/backend/database"
	"sailormoon/backend/modules/slips"
	"sailormoon/backend/modules/users"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	if err := database.InitDB(); err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}

	app := fiber.New()
	app.Use(cors.New())
	api := app.Group("/api")

	userService := &users.UserService{}
	userController := &users.UserController{Service: userService}
	slipsService := &slips.SlipsService{}
	slipsController := &slips.SlipsController{Service: slipsService}

	userController.InitializeRoutes(api)
	slipsController.InitializeRoutes(api)

	if err := app.Listen("127.0.0.1:3000"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
