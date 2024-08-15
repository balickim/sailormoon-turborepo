package main

import (
	"log"
	"sailormoon/backend/database"
	"sailormoon/backend/modules/users"
	"sailormoon/backend/modules/users/entities"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	if err := database.InitDB(); err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}

	if err := database.DB.AutoMigrate(
		&entities.UsersEntity{},
	); err != nil {
		log.Fatalf("Auto-migration failed: %v", err)
	}

	app := fiber.New()
	app.Use(cors.New())

	userService := &users.UserService{}
	userController := &users.UserController{Service: userService}

	api := app.Group("/api")
	userController.InitializeRoutes(api)

	if err := app.Listen("127.0.0.1:3000"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
