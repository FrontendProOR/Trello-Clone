package main

import (
	"tasks-service/src/config"
	"tasks-service/src/controllers"
	"tasks-service/src/middlewares"
	"tasks-service/src/repositories"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	config.ConnectDB()
	repositories.Init()
	router := gin.Default()

	corsConfig := cors.Config{
		AllowOrigins:     []string{"http://localhost:4200"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}

	router.Use(cors.New(corsConfig))

	tasks := router.Group("/tasks")
	tasks.Use(middlewares.AuthMiddleware())
	{
    	tasks.PUT("/:id/status", controllers.UpdateTaskStatus)
	}

	router.GET("/tasks", controllers.GetAllTasks)
	router.GET("/tasks/:id", controllers.GetTaskById)
	router.POST("/tasks", controllers.InsertTask)
	router.PUT("/tasks/:id", controllers.UpdateTask)
	router.DELETE("/tasks/:id", controllers.DeleteTask)
	router.PUT("/tasks/:id/users/:userId", controllers.AddUserToTask)
	router.DELETE("/tasks/:id/users/:userId", controllers.RemoveUserFromTask)

	//ispravniji url-ovi koji ne rade u ovom slucaju
	// router.GET("projects/:projectId/tasks", controllers.GetAllTasks)
	// router.GET("projects/:projectId/tasks/:taskId", controllers.GetTaskById)
	// router.POST("projects/:projectId/tasks", controllers.InsertTask)
	// router.PUT("projects/:projectId/tasks/:taskId", controllers.UpdateTask)
	// router.DELETE("projects/:projectId/tasks/:taskId", controllers.DeleteTask)
	// router.PUT("projects/:projectId/tasks/:taskId/users/:userId", controllers.AddUserToTask)
	// router.DELETE("projects/:projectId/tasks/:taskId/users/:userId", controllers.RemoveUserFromTask)

	router.Run(":8080")
}
