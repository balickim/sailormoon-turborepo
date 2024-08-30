package users

import (
	"encoding/json"
	"errors"
	"sailormoon/backend/database"
	"sailormoon/backend/utils"

	"github.com/allegro/bigcache/v3"
	"golang.org/x/crypto/bcrypt"
)

type UsersService struct {
	Cache *bigcache.BigCache
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (s *UsersService) CreateUser(name, email, password string) (database.UsersEntity, error) {
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return database.UsersEntity{}, err
	}

	user := database.UsersEntity{FirstName: name, Email: email, Password: hashedPassword}
	err = database.DB.Create(&user).Error
	if err != nil {
		return database.UsersEntity{}, utils.HandleDBError(err)
	}

	sessionID := utils.GenerateSessionID()
	userJSON, _ := json.Marshal(user)
	err = s.Cache.Set(sessionID, userJSON)
	if err != nil {
		return database.UsersEntity{}, err
	}

	return user, nil
}

func (s *UsersService) Login(email, password string) (string, error) {
	var user database.UsersEntity
	if err := database.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return "", errors.New("invalid email or password")
	}

	if err := CheckPassword(user.Password, password); err != nil {
		return "", errors.New("invalid email or password")
	}

	sessionID := utils.GenerateSessionID()
	userJSON, _ := json.Marshal(user)
	err := s.Cache.Set(sessionID, userJSON)
	if err != nil {
		return "", err
	}

	return sessionID, nil
}

func (s *UsersService) GetAllUsers() ([]database.UsersEntity, error) {
	var users []database.UsersEntity
	if err := database.DB.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (s *UsersService) GetSession(sessionID string) (database.UsersEntity, error) {
	if s.Cache == nil {
		return database.UsersEntity{}, errors.New("cache is not initialized")
	}

	userData, err := s.Cache.Get(sessionID)
	if err != nil {
		return database.UsersEntity{}, err
	}

	var user database.UsersEntity
	if err := json.Unmarshal(userData, &user); err != nil {
		return database.UsersEntity{}, err
	}
	return user, nil
}
