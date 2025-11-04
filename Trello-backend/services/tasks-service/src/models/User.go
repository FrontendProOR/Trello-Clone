package models

import (
	"time"
)

// UserRole defines the role of a user in the system.
type UserRole string

const (
	RoleManager UserRole = "manager"
	RoleMember  UserRole = "member"
	RoleUnauth  UserRole = "unauthenticated"
)

// User represents a user in the project management platform.
type User struct {
	ID        string    `json:"id" bson:"_id"`                // Unique identifier for the user
	Username  string    `json:"username" bson:"username"`     // Unique username
	FirstName string    `json:"first_name" bson:"first_name"` // User's first name
	LastName  string    `json:"last_name" bson:"last_name"`   // User's last name
	Email     string    `json:"email" bson:"email"`           // User's email address
	Password  string    `json:"-" bson:"password"`            // Hashed password (omit in JSON)
	Role      UserRole  `json:"role" bson:"role"`             // Role of the user (manager/member)
	CreatedAt time.Time `json:"created_at" bson:"created_at"` // Timestamp for account creation
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"` // Timestamp for last update
}

// NewUser creates a new user instance.
func NewUser(username, firstName, lastName, email, password string, role UserRole) *User {
	return &User{
		ID:        generateID(), // function to generate unique ID
		Username:  username,
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Password:  hashPassword(password), // function to hash the password
		Role:      role,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// UpdateUser updates the user's information.
func (u *User) UpdateUser(firstName, lastName, email string) {
	u.FirstName = firstName
	u.LastName = lastName
	u.Email = email
	u.UpdatedAt = time.Now()
}

// HashPassword is a placeholder function for password hashing logic.
func hashPassword(password string) string {
	// Implement password hashing here (e.g., using bcrypt)
	return password // replace with hashed value
}

// GenerateID is a placeholder for generating a unique user ID.
func generateID() string {
	// Implement ID generation logic (e.g., using UUID)
	return "unique-id" // replace with generated ID
}
