package controllers

import (
	"bytes"
	"encoding/base64"
	"log"
	"net/http"
	"time"
	"users-service/src/config"
	"users-service/src/models"
	"users-service/src/repositories"
	"users-service/src/services"

	"github.com/steambap/captcha"

	"fmt"
	"users-service/src/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GetAllUsers handles GET requests to retrieve all users.
func GetAllUsers(c *gin.Context) {
	users, err := services.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

// GetUserById handles GET requests to retrieve a user by ID.
func GetUserById(c *gin.Context) {
	id := c.Param("id")
	user, err := services.GetUserById(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Escapiraj polja pre slanja odgovora
	escapedUser := map[string]interface{}{
		"id":        user.ID.Hex(),
		"username":  utils.EscapeHTML(user.Username),
		"first_name": utils.EscapeHTML(user.FirstName),
		"last_name":  utils.EscapeHTML(user.LastName),
		"email":      utils.EscapeHTML(user.Email),
		"role":       user.Role,
		"created_at": user.CreatedAt,
		"updated_at": user.UpdatedAt,
		"is_active":  user.IsActive,
		
	}

	c.JSON(http.StatusOK, escapedUser)
}


// GetUserByEmail handles GET requests to retrieve a user by email.
func GetUserByEmail(c *gin.Context) {
	email := c.Param("email")
	user, err := services.GetUserByEmail(email)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

// GetUserByUsername handles GET requests to retrieve a user by username.
func GetUserByUsername(c *gin.Context) {
	username := c.Param("username")
	user, err := services.GetUserByUsername(username)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

// InsertUser handles POST requests to create a new user.
func InsertUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := services.InsertUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, result)
}

func UpdateUser(c *gin.Context) {
	id := c.Param("_id")
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	result, err := services.UpdateUser(id, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, result)
}


// DeleteUser handles DELETE requests to remove a user.
func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	err := services.DeleteUser(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

func AddUserToProject(c *gin.Context) {
	projectId := c.Param("projectId")
	userId := c.Param("userId")


	log.Printf("AddUserToProject: Received projectId: %s, userId: %s", projectId, userId)

	log.Printf("AddUserToProject: Received projectId: %s, userId: %s", projectId, userId)

	// Validacija projectId
	projectObjectID, err := primitive.ObjectIDFromHex(projectId)
	if err != nil {
		log.Printf("AddUserToProject: Invalid project ID: %s, error: %v", projectId, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}
	log.Printf("AddUserToProject: Validated projectObjectID: %s", projectObjectID.Hex())

	// Validacija userId
	userObjectID, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		log.Printf("AddUserToProject: Invalid user ID: %s, error: %v", userId, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	log.Printf("AddUserToProject: Validated userObjectID: %s", userObjectID.Hex())

	requesterRole := c.GetString("role")
	if requesterRole != "manager" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied. Only managers can add users to projects."})
		return
	}
	// Dodavanje korisnika u projekat
	if err := repositories.AddUserToProject(projectObjectID, userObjectID); err != nil {
		log.Printf("AddUserToProject: Failed to add user %s to project %s, error: %v", userObjectID.Hex(), projectObjectID.Hex(), err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add user to project"})
		return
	}

	log.Printf("AddUserToProject: Successfully added user %s to project %s", userObjectID.Hex(), projectObjectID.Hex())
	c.JSON(http.StatusOK, gin.H{"message": "User added to project successfully"})
}

func AddUserToProjectByEmail(c *gin.Context) {
	projectId := c.Param("projectId")
	email := c.Param("email")

	log.Printf("AddUserToProjectByEmail: Received projectId: %s, email: %s", projectId, email)

	projectObjectID, err := primitive.ObjectIDFromHex(projectId)
	if err != nil {
		log.Printf("AddUserToProjectByEmail: Invalid project ID: %s, error: %v", projectId, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}
	log.Printf("AddUserToProjectByEmail: Validated projectObjectID: %s", projectObjectID.Hex())

	user, err := services.GetUserByEmail(email)
	if err != nil {
		log.Printf("AddUserToProjectByEmail: Failed to get user by email: %s, error: %v", email, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user by email"})
		return
	}
	log.Printf("AddUserToProjectByEmail: Got user: %v", user)	

	if err := repositories.AddUserToProject(projectObjectID, user.ID); err != nil {
		log.Printf("AddUserToProjectByEmail: Failed to add user %s to project %s, error: %v", user.ID.Hex(), projectObjectID.Hex(), err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add user to project"})
		return
	}

	log.Printf("AddUserToProjectByEmail: Successfully added user %s to project %s", user.ID.Hex(), projectObjectID.Hex())
	c.JSON(http.StatusOK, gin.H{"message": "User added to project successfully"})
}

// RemoveUserFromProject handles DELETE requests to remove a user from a project.
func RemoveUserFromProject(c *gin.Context) {
	projectId := c.Param("projectId")
	userId := c.Param("userId")


	log.Printf("AddUserToProject: Received projectId: %s, userId: %s", projectId, userId)

	projectObjectID, err := primitive.ObjectIDFromHex(projectId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}
	userObjectID, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	requesterRole := c.GetString("role")
	if requesterRole != "manager" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied. Only managers can add tasks."})
		return
	}
	

	

	if err := repositories.RemoveUserFromProject(projectObjectID, userObjectID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove user from project"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User removed from project successfully"})
}

func RemoveUserFromProjectByEmail(c *gin.Context) {
	projectId := c.Param("projectId")
	email := c.Param("email")

	log.Printf("RemoveUserFromProjectByEmail: Received projectId: %s, email: %s", projectId, email)	

	projectObjectID, err := primitive.ObjectIDFromHex(projectId)
	if err != nil {
		log.Printf("RemoveUserFromProjectByEmail: Invalid project ID: %s, error: %v", projectId, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}
	log.Printf("RemoveUserFromProjectByEmail: Validated projectObjectID: %s", projectObjectID.Hex())

	user, err := services.GetUserByEmail(email)
	if err != nil {
		log.Printf("RemoveUserFromProjectByEmail: Failed to get user by email: %s, error: %v", email, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user by email"})
		return
	}
	log.Printf("RemoveUserFromProjectByEmail: Got user: %v", user)	

	if err := repositories.RemoveUserFromProject(projectObjectID, user.ID); err != nil {
		log.Printf("RemoveUserFromProjectByEmail: Failed to remove user %s from project %s, error: %v", user.ID.Hex(), projectObjectID.Hex(), err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove user from project"})
		return
	}

	log.Printf("RemoveUserFromProjectByEmail: Successfully removed user %s from project %s", user.ID.Hex(), projectObjectID.Hex())
	c.JSON(http.StatusOK, gin.H{"message": "User removed from project successfully"})
}

// RegisterUser handles POST requests for user registration.
func RegisterUser(c *gin.Context) {
	var req struct {
		FirstName string          `json:"first_name"`
		LastName  string          `json:"last_name"`
		Username  string          `json:"username"`
		Email     string          `json:"email"`
		Password  string          `json:"password"`
		Role      models.UserRole `json:"role"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if config.IsBlacklistedPassword(req.Password) {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Password is too common. Please choose a stronger password."})
        return
    }

	req.FirstName = utils.EscapeHTML(req.FirstName)
	req.LastName = utils.EscapeHTML(req.LastName)
	req.Username = utils.EscapeHTML(req.Username)
	req.Email = utils.EscapeHTML(req.Email)

	err := services.RegisterUser(req.FirstName, req.LastName, req.Username, req.Email, req.Password, req.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registration successful. Verification email sent."})
}

// VerifyUser handles POST requests for verifying a user's email.
func VerifyUser(c *gin.Context) {
	var req struct {
		Email string `json:"email"`
		Code  string `json:"code"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := services.VerifyUser(req.Email, req.Code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User verified and account activated successfully."})
}

// LoginUser handles POST requests for user login.
func LoginUser(c *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Bind the request body to the struct
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	// Validate email and password
	if req.Email == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email and password are required"})
		return
	}

	// Call the service to log in the user and get the token
	token, err := services.LoginUser(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password !!! " + err.Error()})
		return
	}

	// Return the token
	escapedEmail := utils.EscapeHTML(req.Email)
	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"email": escapedEmail,
	})
}

type ChangePasswordRequest struct {
    Email       string `json:"email"`
    OldPassword string `json:"oldPassword"`
    NewPassword string `json:"newPassword"`
}

func ChangePassword(c *gin.Context) {
    var req ChangePasswordRequest

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    err := services.ChangePassword(req.Email, req.OldPassword, req.NewPassword)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully."})
}

func CheckPasswordBlacklist(c *gin.Context) {
	var req struct {
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	isBlacklisted := config.IsBlacklistedPassword(req.Password)
	c.JSON(http.StatusOK, gin.H{"blacklisted": isBlacklisted})
}

// GenerateMagicLink handles POST requests to generate and send a magic link
func GenerateMagicLink(c *gin.Context) {
	var req struct {
		Email string `json:"email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	// Proveri da li korisnik postoji
	user, err := repositories.GetUserByEmail(req.Email)
	if err != nil || user.Email == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Generiši token za magic link
	token, err := utils.GenerateMagicLinkToken(user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate magic link"})
		return
	}

	// Kreiraj URL magic link-a
	magicLink := fmt.Sprintf("http://localhost:4200/users/magic-login?token=%s", token)

	// Pošalji magic link putem e-pošte
	err = services.SendMagicLinkEmail(req.Email, magicLink)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send magic link email"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Magic link sent successfully"})
}

// ConfirmMagicLink handles GET requests to confirm the magic link token
func ConfirmMagicLink(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token is required"})
		return
	}

	email, err := utils.ValidateMagicLinkToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		return
	}

	// Pronađi korisnika
	user, err := repositories.GetUserByEmail(email)
	if err != nil || user.Email == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Generiši JWT token
	jwtToken, err := utils.GenerateJWT(user.ID.Hex(), user.Email, string(user.Role))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate JWT"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": jwtToken})
}
var captchaStore = make(map[string]string) // Privremena mapa za čuvanje CAPTCHA tekstova

// GenerateCaptcha generiše CAPTCHA sliku i vraća njen ID
func GenerateCaptcha(c *gin.Context) {
	// Kreiramo CAPTCHA sa matematičkim izrazom
	captchaData, err := captcha.NewMathExpr(200, 100)
	if err != nil {
		log.Printf("Error generating CAPTCHA: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate CAPTCHA"})
		return
	}

	// Kreiramo buffer za sliku
	var buf bytes.Buffer
	if err := captchaData.WriteImage(&buf); err != nil {
		log.Printf("Error writing CAPTCHA to buffer: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write CAPTCHA image"})
		return
	}

	captchaID := fmt.Sprintf("%d", time.Now().UnixNano())
	captchaStore[captchaID] = captchaData.Text

	imageBase64 := base64.StdEncoding.EncodeToString(buf.Bytes())
	c.JSON(http.StatusOK, gin.H{
		"id":    captchaID,
		"image": fmt.Sprintf("data:image/png;base64,%s", imageBase64),
	})
}



// ValidateCaptcha validira CAPTCHA na osnovu unetog teksta i ID-a
func ValidateCaptcha(c *gin.Context) {
	var request struct {
		CaptchaID   string `json:"captchaId"`
		CaptchaText string `json:"captchaText"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}
	log.Printf("Received CAPTCHA ID: %s, Text: %s", request.CaptchaID, request.CaptchaText)

	// Proveravamo CAPTCHA ID i tekst
	storedText, exists := captchaStore[request.CaptchaID]
	if !exists || storedText != request.CaptchaText {
		log.Printf("CAPTCHA ID not found: %s", request.CaptchaID)

		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid CAPTCHA"})
		return
	}
	if storedText != request.CaptchaText {
		log.Printf("CAPTCHA text mismatch: Expected %s, Got %s", storedText, request.CaptchaText)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid CAPTCHA text"})
		return
	}

	// CAPTCHA je validna
	delete(captchaStore, request.CaptchaID) // Brišemo CAPTCHA iz memorije nakon validacije
	c.JSON(http.StatusOK, gin.H{"message": "CAPTCHA validated successfully"})
}