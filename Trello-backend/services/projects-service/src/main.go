package main

import (
	"log"
	"projects-service/src/config"
	"projects-service/src/controllers"
	"projects-service/src/repositories"

	"projects-service/src/middlewares"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	config.ConnectDB()
	repositories.Init()
	router := gin.Default()
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// corsConfig := cors.Config{
	// 	AllowOrigins:     []string{"http://localhost:4200"},
	// 	AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	// 	AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
	// 	ExposeHeaders:    []string{"Content-Length","Authorization"},
	// 	AllowCredentials: true,
	// }
	newCorsConfig := cors.Config{
		AllowOrigins:     []string{"http://localhost:4200"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization","Content-Type", "Accept"},
		AllowCredentials: true,
	}

	router.Use(cors.New(newCorsConfig))

	private := router.Group("/projects")
	private.Use(middlewares.AuthMiddleware())
	{
		private.GET("/", controllers.GetAllProjects)
		private.GET("/:id", controllers.GetProjectById)
		private.POST("", controllers.InsertProject)
		private.PUT("/:id", controllers.UpdateProject)
		private.DELETE("/:id", controllers.DeleteProject)
		private.PUT("/:id/tasks", controllers.AddTaskToProject)
		
		private.PUT("/:id/tasks/:taskId/users/:userId", controllers.AddUserToTaskInProject)
		private.DELETE("/:id/tasks/:taskId/users/:userId", controllers.RemoveUserFromTask)

		private.GET("/manager/:id", controllers.GetProjectsByManagerId)
		private.GET("/member/:id", controllers.GetProjectsByMemberId)

		private.GET("/:id/members", controllers.GetAllMembersInProject)

		private.GET("/:id/tasks", controllers.GetTasksByProjectId)
	}

	router.Run(":8080")
}
