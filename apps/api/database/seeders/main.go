package main

import (
	"log"
	"math/rand"
	"sailormoon/backend/database"
	"time"

	"github.com/bxcodec/faker/v4"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

const (
	HOW_MANY_USERS     = 100
	HOW_MANY_SLIPS     = 1000
	HOW_MANY_BOATS     = 100
	MIN_USERS_PER_BOAT = 20
	MAX_USERS_PER_BOAT = 50
	MIN_SLIPS_PER_BOAT = 10
	MAX_SLIPS_PER_BOAT = 30
)

func Seed(db *gorm.DB) {
	tx := db.Begin()

	rand.Seed(time.Now().UnixNano())

	var users []database.UsersEntity
	for i := 0; i < HOW_MANY_USERS; i++ {
		user := database.UsersEntity{
			Email:       faker.Email(),
			Password:    faker.Password(),
			LastName:    faker.LastName(),
			FirstName:   faker.FirstName(),
			Phone:       faker.Phonenumber(),
			Address:     faker.Word(),
			NIP:         faker.Word(),
			CompanyName: faker.Word(),
			Notes:       faker.Sentence(),
		}
		if err := tx.Create(&user).Error; err != nil {
			tx.Rollback()
			log.Fatalf("Failed to seed users: %v", err)
		}
		users = append(users, user)
		log.Printf("Created user: %s, ID: %d", user.Email, user.ID)
	}

	var slips []database.SlipsEntity
	for i := 0; i < HOW_MANY_SLIPS; i++ {
		slip := database.SlipsEntity{
			Number:     i + 1,
			IsOccupied: true,
			Notes:      faker.Sentence(),
		}
		if err := tx.Create(&slip).Error; err != nil {
			tx.Rollback()
			log.Fatalf("Failed to seed slips: %v", err)
		}
		slips = append(slips, slip)
	}

	for i := 0; i < HOW_MANY_BOATS; i++ {
		// Select a random number of users for this boat
		numUsers := rand.Intn(MAX_USERS_PER_BOAT-MIN_USERS_PER_BOAT+1) + MIN_USERS_PER_BOAT
		boatUsers := []*database.UsersEntity{}
		for j := 0; j < numUsers; j++ {
			boatUsers = append(boatUsers, &users[rand.Intn(len(users))])
		}

		// Select a random number of slips for this boat
		numSlips := rand.Intn(MAX_SLIPS_PER_BOAT-MIN_SLIPS_PER_BOAT+1) + MIN_SLIPS_PER_BOAT
		boatSlips := []*database.SlipsEntity{}
		for j := 0; j < numSlips; j++ {
			boatSlips = append(boatSlips, &slips[rand.Intn(len(slips))])
		}

		boat := database.BoatsEntity{
			Name:   faker.Name(),
			Type:   faker.Word(),
			Length: faker.Word(),
			Width:  faker.Word(),
			Weight: faker.Word(),
			Draft:  faker.Word(),
			Owners: boatUsers,
			Slips:  boatSlips,
			Notes:  faker.Sentence(),
		}
		if err := tx.Create(&boat).Error; err != nil {
			tx.Rollback()
			log.Fatalf("Failed to seed boats: %v", err)
		}
	}

	if err := tx.Commit().Error; err != nil {
		log.Fatalf("Transaction commit failed: %v", err)
	}

	log.Println("Database seeding completed successfully.")
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	if err := database.InitDB(); err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}

	db := database.DB

	Seed(db)
}
