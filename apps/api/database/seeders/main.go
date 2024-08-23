package main

import (
	"log"
	"sailormoon/backend/database"

	"github.com/bxcodec/faker/v4"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

const (
	HOW_MANY_USERS = 2
	HOW_MANY_SLIPS = 1000
)

func Seed(db *gorm.DB) {
	tx := db.Begin()

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

	for i := 0; i < HOW_MANY_SLIPS; i++ {
		boat := database.BoatsEntity{
			Name:   faker.Name(),
			Type:   faker.Word(),
			Length: faker.Word(),
			Width:  faker.Word(),
			Weight: faker.Word(),
			Draft:  faker.Word(),
			Owners: []*database.UsersEntity{&users[i%len(users)]},
			Slips:  []*database.SlipsEntity{&slips[i%len(slips)]},
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
