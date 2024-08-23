package main

import (
	"log"
	"sailormoon/backend/database"

	"github.com/joho/godotenv"
)

func RunMigrations() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	if err := database.InitDB(); err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}

	db := database.DB

	if err := database.Migrate(db); err != nil {
		log.Fatalf("Auto-migration failed: %v", err)
	}

	log.Println("Database migration completed.")
}

func main() {
	RunMigrations()
}
