package users

import (
	"sailormoon/backend/database"
	"sailormoon/backend/modules/users/entities"
	"sailormoon/backend/utils"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct{}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (s *UserService) CreateUser(name, email, password string) (entities.UsersEntity, error) {
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return entities.UsersEntity{}, err
	}

	user := entities.UsersEntity{Name: name, Email: email, Password: hashedPassword}
	err = database.DB.Create(&user).Error
	if err != nil {
		return entities.UsersEntity{}, utils.HandleDBError(err)
	}

	return user, nil
}

func (s *UserService) GetAllUsers() ([]entities.UsersEntity, error) {
	var users []entities.UsersEntity
	if err := database.DB.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
