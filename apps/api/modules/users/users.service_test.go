package users

import (
	"sailormoon/backend/database"
	"sailormoon/backend/modules/users/entities"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// SetupTestDB initializes a new in-memory SQLite database for testing
func SetupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	err = db.AutoMigrate(&entities.UsersEntity{})
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	database.DB = db
	return db
}

// TestHashPassword tests the HashPassword function
func TestHashPassword(t *testing.T) {
	password := "securepassword"
	hashedPassword, err := HashPassword(password)
	assert.NoError(t, err)
	assert.NotEmpty(t, hashedPassword)

	// Verify that the hashed password is different from the plain password
	assert.NotEqual(t, password, hashedPassword)

	// Verify that the hashed password can be verified
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	assert.NoError(t, err)
}

// TestCreateUser tests the CreateUser function
func TestCreateUser(t *testing.T) {
	db := SetupTestDB(t)
	defer db.Exec("DROP TABLE users")

	userService := UserService{}

	// Test creating a new user
	user, err := userService.CreateUser("John Doe", "john@example.com", "securepassword")
	assert.NoError(t, err)
	assert.NotEmpty(t, user.ID)
	assert.Equal(t, "John Doe", user.Name)
	assert.Equal(t, "john@example.com", user.Email)

	// Verify that the password was hashed
	assert.NotEqual(t, "securepassword", user.Password)

	// Test creating a user with the same email
	_, err = userService.CreateUser("Jane Doe", "john@example.com", "anotherpassword")
	assert.Error(t, err)
	assert.Equal(t, "field email is already taken", err.Error())
}

// TestGetAllUsers tests the GetAllUsers function
func TestGetAllUsers(t *testing.T) {
	db := SetupTestDB(t)
	defer db.Exec("DROP TABLE users")

	userService := UserService{}

	// Create two users
	_, err := userService.CreateUser("John Doe", "john@example.com", "securepassword")
	assert.NoError(t, err)
	_, err = userService.CreateUser("Jane Doe", "jane@example.com", "anotherpassword")
	assert.NoError(t, err)

	// Retrieve all users
	users, err := userService.GetAllUsers()
	assert.NoError(t, err)
	assert.Len(t, users, 2)

	// Verify the details of the retrieved users
	assert.Equal(t, "John Doe", users[0].Name)
	assert.Equal(t, "john@example.com", users[0].Email)
	assert.Equal(t, "Jane Doe", users[1].Name)
	assert.Equal(t, "jane@example.com", users[1].Email)
}
