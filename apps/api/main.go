package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sailormoon/backend/database"
	"sailormoon/backend/modules/boats"
	"sailormoon/backend/modules/slips"
	"sailormoon/backend/modules/users"
	"syscall"
	"time"

	"github.com/allegro/bigcache/v3"
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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cacheConfig := bigcache.DefaultConfig(10 * time.Minute)
	cache, err := bigcache.New(ctx, cacheConfig)
	if err != nil {
		log.Fatalf("Failed to initialize BigCache: %v", err)
	}

	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173",
		AllowMethods:     "GET,POST,PUT,DELETE",
		AllowHeaders:     "Origin,Content-Type,Accept",
		AllowCredentials: true,
	}))

	usersService := &users.UsersService{Cache: cache}
	usersController := &users.UsersController{Service: usersService}
	slipsService := &slips.SlipsService{}
	slipsController := &slips.SlipsController{Service: slipsService, UsersService: usersService}
	boatsService := &boats.BoatsService{}
	boatsController := &boats.BoatsController{Service: boatsService}

	usersController.InitializeRoutes(app.Group("/users"))
	slipsController.InitializeRoutes(app.Group("/slips"))
	boatsController.InitializeRoutes(app.Group("/boats"))

	// Graceful shutdown on SIGINT and SIGTERM
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		if err := app.Shutdown(); err != nil {
			log.Fatalf("Failed to shutdown server: %v", err)
		}
		cancel()
	}()

	if err := app.Listen("127.0.0.1:3000"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
