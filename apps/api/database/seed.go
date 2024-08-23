package database

import (
	"log"

	"github.com/bxcodec/faker/v4"
	"gorm.io/gorm"
)

func Seed(db *gorm.DB) {
	var HOW_MANY_SLIPS = 1000

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
		if err := db.Create(&user).Error; err != nil {
			log.Fatalf("Failed to seed users: %v", err)
		}
		users = append(users, user)
	}

	for _, user := range users {
		log.Printf("Created user: %s, ID: %d", user.Email, user.ID)
	}

	slips := []SlipsEntity{}
	for i := 0; i < HOW_MANY_SLIPS; i++ {
		slip := SlipsEntity{
			Number:     i + 1,
			IsOccupied: true,
			Notes:      faker.Sentence(),
		}
		if err := db.Create(&slip).Error; err != nil {
			log.Fatalf("Failed to seed slips: %v", err)
		}
		slips = append(slips, slip)
	}

	for i := 0; i < HOW_MANY_SLIPS; i++ {
		boat := BoatsEntity{
			Name:   faker.Name(),
			Type:   faker.Word(),
			Length: faker.Word(),
			Width:  faker.Word(),
			Weight: faker.Word(),
			Draft:  faker.Word(),
			Owners: []*UsersEntity{&users[i%len(users)]},
			Slips:  []*SlipsEntity{&slips[i%len(slips)]},
			Notes:  faker.Sentence(),
		}
		if err := db.Create(&boat).Error; err != nil {
			log.Fatalf("Failed to seed boats: %v", err)
		}
	}
}
