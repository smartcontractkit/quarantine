package main

import (
	"fmt"
	"log"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Email string `json:"email"`
}

type UserService struct {
	users map[int]*User
	nextID int
}

func NewUserService() *UserService {
	return &UserService{
		users:  make(map[int]*User),
		nextID: 1,
	}
}

func (us *UserService) CreateUser(name, email string) *User {
	user := &User{
		ID:    us.nextID,
		Name:  name,
		Email: email,
	}
	us.users[user.ID] = user
	us.nextID++
	return user
}

func (us *UserService) GetUser(id int) (*User, bool) {
	user, exists := us.users[id]
	return user, exists
}

func (us *UserService) DeleteUser(id int) bool {
	if _, exists := us.users[id]; exists {
		delete(us.users, id)
		return true
	}
	return false
}

func (us *UserService) ListUsers() []*User {
	users := make([]*User, 0, len(us.users))
	for _, user := range us.users {
		users = append(users, user)
	}
	return users
}

func main() {
	service := NewUserService()
	
	user1 := service.CreateUser("Alice", "alice@example.com")
	user2 := service.CreateUser("Bob", "bob@example.com")
	
	fmt.Printf("Created users: %+v, %+v\n", user1, user2)
	
	if user, exists := service.GetUser(1); exists {
		fmt.Printf("Found user: %+v\n", user)
	}
	
	users := service.ListUsers()
	fmt.Printf("All users: %+v\n", users)
	
	if service.DeleteUser(1) {
		log.Println("User 1 deleted")
	}
	
	users = service.ListUsers()
	fmt.Printf("Remaining users: %+v\n", users)
}
