package services

import (
	"errors"
	"time"
	"users-service/src/models"
	"users-service/src/repositories"
	"users-service/src/utils"

	"golang.org/x/crypto/bcrypt"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// InsertUser inserts a new user into the repository.
func InsertUser(user models.User) (*mongo.InsertOneResult, error) {
	return repositories.InsertUser(user)
}

// GetAllUsers retrieves all users from the repository.
func GetAllUsers() ([]models.User, error) {
	return repositories.GetAllUsers()
}

// GetUserById retrieves a user by ID from the repository.
func GetUserById(id string) (models.User, error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.User{}, err
	}
	return repositories.GetUserById(objectId)
}

func UpdateUser(id string, user models.User) (interface{}, error) {
	objectId, _ := primitive.ObjectIDFromHex(id)
	return repositories.UpdateUser(objectId, user)
}

// DeleteUser deletes a user from the repository.
func DeleteUser(id string) error {
	return repositories.DeleteUser(id)
}

// GetUserByEmail retrieves a user by email from the repository.
func GetUserByEmail(email string) (models.User, error) {
	return repositories.GetUserByEmail(email)
}

// GetUserByUsername retrieves a user by username from the repository.
func GetUserByUsername(username string) (models.User, error) {
	return repositories.GetUserByUsername(username)
}

// AddUserToProject adds a user to a project.
func AddUserToProject(projectId string, userId primitive.ObjectID) error {
	objectId, err := primitive.ObjectIDFromHex(projectId)
	if err != nil {
		return err
	}
	return repositories.AddUserToProject(objectId, userId)
}

// RemoveUserFromProject removes a user from a project.
func RemoveUserFromProject(projectId string, userId primitive.ObjectID) error {
	objectId, err := primitive.ObjectIDFromHex(projectId)
	if err != nil {
		return err
	}
	return repositories.RemoveUserFromProject(objectId, userId)
}

// RegisterUser handles the registration of a new user.
func RegisterUser(firstName, lastName, username, email string, password string, role models.UserRole) error {
	existingUser, err := repositories.GetUserByEmail(email)
	if err == nil && existingUser.Email != "" {
		return errors.New("user with this email already exists")
	}
	// existingUsername, err := repositories.GetUserByUsername(username)
	// if err == nil && existingUsername.Username != "" {
	// 	return errors.New("user with this username already exists")
	// }

	var code = utils.GenerateCode()
	err = SendVerificationEmail(email, code)
	if err != nil {
		return err
	}

	user := models.User{
		ID:        [12]byte{},
		Username:  username,
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Password:  password,
		Role:      role,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		IsActive:  false,
		Code:      code,
	}

	_, err = repositories.InsertUser(user)
	if err != nil {
		return err
	}

	err = repositories.SaveVerificationCode(email, code)
	if err != nil {
		return err
	}
	return nil
}

// VerifyUser verifies a user's email.
func VerifyUser(email, code string) error {
	user, err := repositories.GetUserByEmail(email)
	if err != nil {
		return err
	}

	if user.IsActive {
		return errors.New("user is already verified and active")
	}

	err = repositories.ActivateUser(email, code)
	if err != nil {
		return err
	}
	return nil
}

// LoginUser authenticates a user and returns a JWT token
func LoginUser(email, password string) (string, error) {
	// Authenticate user
	user, err := repositories.LoginUser(email, password)
	if err != nil {
		return "", err
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(user.ID.Hex(), user.Email, string(user.Role))
	if err != nil {
		return "", err
	}

	return token, nil
}

func ChangePassword(email, oldPassword, newPassword string) error {
	//all logic before repository 
	if email == "" || oldPassword == "" || newPassword == "" {
		return errors.New("email, oldPassword and newPassword are required")
	}
	
	user, err := repositories.GetUserByEmail(email)
	if err != nil {
		return err
	}
	
	if !user.IsActive {
		return errors.New("user is not active")
	}
	
	// err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword))
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword))
	if err != nil {
		return errors.New("invalid old password")
	}
	
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}
	
	user.Password = string(hashedPassword)
	
	err = repositories.ChangePassword(user)
	if err != nil {
		return err
	}
	
	return nil
}