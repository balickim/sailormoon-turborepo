package users

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

var users = []User{}
var lastID = 0

type UserService struct{}

func (s *UserService) CreateUser(name, email string) (User, error) {
	lastID++
	user := User{ID: lastID, Name: name, Email: email}
	users = append(users, user)
	return user, nil
}

func (s *UserService) GetAllUsers() ([]User, error) {
	return users, nil
}
