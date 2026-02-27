package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Username  string         `json:"username" gorm:"uniqueIndex;not null"`
	Password  string         `json:"-" gorm:"not null"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=4"`
	Password string `json:"password" validate:"required,min=6"`
}

// Response
type AuthResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

type UserResponse struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
}

// Mock database with UUIDs
var Users = []User{
	{
		ID:       uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
		Username: "john",
		Password: "password123", // Hash this in production!
	},
	{
		ID:       uuid.MustParse("6ba7b810-9dad-11d1-80b4-00c04fd430c8"),
		Username: "test",
		Password: "test123",
	},
}

// Helper function to create a new user with generated UUID
func NewUser(username, password string) User {

	return User{
		ID:       uuid.New(),
		Username: username,
		Password: password, // Remember to hash this in production!
	}
}

// Find user by username (helper function)
func FindUserByUsername(username string) *User {
	for _, user := range Users {
		if user.Username == username {
			return &user
		}
	}
	return nil
}
