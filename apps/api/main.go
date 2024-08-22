package main

import (
	"log"
	"os"
	"sailormoon/backend/database"
	"sailormoon/backend/modules/slips"
	"sailormoon/backend/modules/users"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	if err := database.InitDB(); err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}

	goEnv := os.Getenv("GO_ENV")
	if goEnv == "development" {
		if err := database.DB.AutoMigrate(
			&database.UsersEntity{},
			&database.SlipsEntity{},
			&database.BoatsEntity{},
		); err != nil {
			log.Fatalf("Auto-migration failed: %v", err)
		}

		database.Seed(database.DB)
		log.Println("Database seeding and auto-migration completed.")
	} else {
		log.Println("Skipping auto migrations and database seeding in non-development environment.")
	}

	app := fiber.New()

	app.Use(cors.New())

	userService := &users.UserService{}
	userController := &users.UserController{Service: userService}
	slipsService := &slips.SlipsService{}
	slipsController := &slips.SlipsController{Service: slipsService}

	api := app.Group("/api")
	userController.InitializeRoutes(api)
	slipsController.InitializeRoutes(api)

	if err := app.Listen("127.0.0.1:3000"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
