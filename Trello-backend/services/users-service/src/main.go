package main

import (
	"os"
	"users-service/src/config"
	"users-service/src/controllers"
	"users-service/src/repositories"

	"github.com/joho/godotenv"

	"log"
	"users-service/src/middlewares"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	
	config.Init()
	config.ConnectDB()
	repositories.Init()
	
	// httpClient := &http.Client{
	// 	Timeout: 5 * time.Second,
	// }

	router := gin.Default()
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, proceeding with environment variables")
	}

	// Dobijanje putanje blacklist fajla iz varijable okruženja
	blacklistPath := os.Getenv("PASSWORD_BLACKLIST_PATH")
	if blacklistPath == "" {
		log.Fatalf("PASSWORD_BLACKLIST_PATH is not set in the environment")
	}

	// Učitajte blacklist
	err = config.LoadPasswordBlacklist(blacklistPath)
	if err != nil {
		log.Fatalf("Failed to load password blacklist: %v", err)
	}

	newCorsConfig := cors.Config{
		AllowOrigins:     []string{"http://localhost:4200"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization","Content-Type", "Accept"},
		AllowCredentials: true,
	}

	router.Use(cors.New(newCorsConfig))

	router.POST("/users/register", controllers.RegisterUser)
	router.POST("/users/verify", controllers.VerifyUser)
	router.POST("/users/login", controllers.LoginUser)
	router.POST("/users/check-password", controllers.CheckPasswordBlacklist)
	router.POST("/users/magic-link", controllers.GenerateMagicLink)
	router.GET("/users/magic-link/validate", controllers.ConfirmMagicLink)
	router.GET("/users/captcha", controllers.GenerateCaptcha)       // Generisanje CAPTCHA
	router.POST("/users/captcha/validate", controllers.ValidateCaptcha) // Validacija CAPTCHA

	private := router.Group("/users")
	private.Use(middlewares.AuthMiddleware())
	{
		private.GET("/", controllers.GetAllUsers)
		private.GET("/:id", controllers.GetUserById)
		private.GET("/username/:username", controllers.GetUserByUsername)
		private.GET("/email/:email", controllers.GetUserByEmail)
		private.POST("/", controllers.InsertUser)
		private.DELETE("/:id", controllers.DeleteUser)
		private.PUT("/change-password", controllers.ChangePassword)
		private.PUT("/:id", controllers.UpdateUser)

	}
	
	projects := router.Group("/projects")
	projects.Use(middlewares.AuthMiddleware())
	{
		// projects.PUT("/:projectId/users/:userId", controllers.AddUserToProject)
		projects.PUT("/:projectId/users/:email", controllers.AddUserToProjectByEmail)
		// projects.DELETE("/:projectId/users/:userId", controllers.RemoveUserFromProject)
		projects.DELETE("/:projectId/users/:email", controllers.RemoveUserFromProjectByEmail)
	}
	
	
	router.Run(":8080")
}
