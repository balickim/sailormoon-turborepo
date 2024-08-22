package database

import (
	"log"

	"github.com/bxcodec/faker/v4"
	"gorm.io/gorm"
)

func Seed(db *gorm.DB) {
	var HOW_MANY_SLIPS = 10

	// Seed Slips with faker
	slips := []SlipsEntity{}
	for i := 0; i < HOW_MANY_SLIPS; i++ {
		slip := SlipsEntity{
			Number:     i,
			IsOccupied: true,
			Notes:      faker.Sentence(),
		}
		slips = append(slips, slip)
	}

	for _, slip := range slips {
		if err := db.Where("number = ?", slip.Number).FirstOrCreate(&slip).Error; err != nil {
			log.Fatalf("Failed to seed slips: %v", err)
		}
	}

	// Seed Users with faker
	users := []UsersEntity{}
	for i := 0; i < 2; i++ {
		user := UsersEntity{
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
		users = append(users, user)
	}

	for _, user := range users {
		if err := db.Where("email = ?", user.Email).FirstOrCreate(&user).Error; err != nil {
			log.Fatalf("Failed to seed users: %v", err)
		}
	}

	// Seed Boats with faker
	boats := []BoatsEntity{}
	for i := 0; i < 2; i++ {
		boat := BoatsEntity{
			Name:    faker.Name(),
			Type:    faker.Word(),
			Length:  faker.Word(),
			Width:   faker.Word(),
			Weight:  faker.Word(),
			Draft:   faker.Word(),
			OwnerID: users[i].ID,
			Notes:   faker.Sentence(),
		}
		boats = append(boats, boat)
	}

	for _, boat := range boats {
		if err := db.Where("name = ?", boat.Name).FirstOrCreate(&boat).Error; err != nil {
			log.Fatalf("Failed to seed boats: %v", err)
		}
	}
}
