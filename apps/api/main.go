package main

import (
	"log"
	"sailormoon/backend/database"
	"sailormoon/backend/modules/boats"
	"sailormoon/backend/modules/slips"
	"sailormoon/backend/modules/users"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found: %v", err)
	}

	if err := database.InitDB(); err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}

	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173",
		AllowMethods:     "GET,POST,PUT,DELETE",
		AllowHeaders:     "Origin,Content-Type,Accept",
		AllowCredentials: true,
	}))

	userService := &users.UserService{}
	userController := &users.UserController{Service: userService}
	slipsService := &slips.SlipsService{}
	slipsController := &slips.SlipsController{Service: slipsService}
	boatsService := &boats.BoatsService{}
	boatsController := &boats.BoatsController{Service: boatsService}

	userController.InitializeRoutes(app.Group("/users"))
	slipsController.InitializeRoutes(app.Group("/slips"))
	boatsController.InitializeRoutes(app.Group("/boats"))

	if err := app.Listen("127.0.0.1:3000"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
