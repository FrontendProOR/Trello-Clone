package main

import (
	"log"
	"notifications-service/config"
	"notifications-service/controllers"
	"notifications-service/repositories"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
)

func main() {
	// Retry logika
	var session *gocql.Session
	var err error
	for i := 0; i < 10; i++ {
		session, err = config.ConnectToCassandra()
		if err == nil {
			break
		}
		log.Printf("Cassandra is not ready yet. Retrying in 5 seconds... (%d/10)\n", i+1)
		time.Sleep(5 * time.Second)
	}
	if err != nil {
		log.Fatalf("Failed to connect to Cassandra after retries: %v", err)
		return
	}
	defer session.Close()
	log.Println("Connected to Cassandra successfully.")

	// Repo i tabela
	repo := repositories.NewNotificationsRepository(session)
	err = repo.CreateTable()
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}
	log.Println("Table created successfully (if not already present).")

	controller := controllers.NotificationsController{Repo: repo}

	// Gin router + CORS
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:4200"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Rute
	router.GET("/notifications/:user_id", controller.GetNotifications)
	router.GET("/notifications/unread/:user_id", controller.GetUnreadNotifications)
	router.POST("/notifications", controller.CreateNotification)
	router.PUT("/notifications/read", controller.MarkNotificationAsRead)
	router.DELETE("/notifications", controller.DeleteNotification)
	router.GET("/health", func(ctx *gin.Context) {
		ctx.String(200, "Service is running")
	})

	log.Println("Server is running on port 8084")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
