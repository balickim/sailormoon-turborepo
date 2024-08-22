package users

import (
	"sailormoon/backend/database"
	"sailormoon/backend/utils"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct{}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (s *UserService) CreateUser(name, email, password string) (database.UsersEntity, error) {
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return database.UsersEntity{}, err
	}

	user := database.UsersEntity{FirstName: name, Email: email, Password: hashedPassword}
	err = database.DB.Create(&user).Error
	if err != nil {
		return database.UsersEntity{}, utils.HandleDBError(err)
	}

	return user, nil
}

func (s *UserService) GetAllUsers() ([]database.UsersEntity, error) {
	var users []database.UsersEntity
	if err := database.DB.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
