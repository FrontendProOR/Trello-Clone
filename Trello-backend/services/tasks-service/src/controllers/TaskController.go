package controllers

import (
	"log"
	"net/http"
	"tasks-service/src/models"
	"tasks-service/src/services"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetAllTasks(c *gin.Context) {
	tasks, err := services.GetAllTasks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, tasks)
}

func GetTaskById(c *gin.Context) {
	id := c.Param("id")
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	task, err := services.GetTaskById(objectId)
	if err != nil {
		c.JSON(http.StatusNotFound, err.Error())
		return
	}
	c.JSON(http.StatusOK, task)
}

func InsertTask(c *gin.Context) {
    var task models.Task

    // Dekodiramo JSON telo zahteva u task objekat
    if err := c.BindJSON(&task); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task data"})
        return
    }

    // Pozivamo servis za umetanje zadatka
    createdTaskID, err := services.InsertTask(task)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
        return
    }

    // Vraćamo ID kreiranog zadatka kao odgovor
    c.JSON(http.StatusCreated, gin.H{"id": createdTaskID})
}

func UpdateTask(c *gin.Context) {
	id := c.Param("id")
	var task models.Task
	if err := c.BindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	result, err := services.UpdateTask(id, task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, result)
}

func DeleteTask(c *gin.Context) {
	id := c.Param("id")
	result, err := services.DeleteTask(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, result)
}
func AddUserToTask(c *gin.Context) {
    taskId := c.Param("id")
    userId := c.Param("userId")
    log.Printf("Received request to add user %s to task %s", userId, taskId)

    // Konvertujemo taskId i userId u ObjectID
    taskObjectID, err := primitive.ObjectIDFromHex(taskId)
    if err != nil {
        log.Printf("Invalid task ID: %v", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
        return
    }

    userObjectID, err := primitive.ObjectIDFromHex(userId)
    if err != nil {
        log.Printf("Invalid user ID: %v", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
        return
    }

    // Pozivamo servis da doda korisnika na zadatak
    err = services.AddUserToTask(taskObjectID, userObjectID)
    if err != nil {
        log.Printf("Failed to add user to task: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add user to task"})
        return
    }

    log.Println("User added to task successfully.")
    c.JSON(http.StatusOK, gin.H{"message": "User added to task successfully"})
}
func RemoveUserFromTask(c *gin.Context) {
    taskId := c.Param("id")
    userId := c.Param("userId")

    log.Printf("Removing user %s from task %s", userId, taskId)

    // Konvertujemo taskId u ObjectID
    taskObjectID, err := primitive.ObjectIDFromHex(taskId)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
        return
    }
    

    // Dohvatamo status zadatka
    status, err := services.GetTaskStatus(taskId)
    if err != nil {
        log.Printf("Failed to fetch task status: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch task status"})
        return
    }

    // Proveravamo da li je zadatak završen
    if status == "Finished" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot remove user from a task with status 'Finished'"})
        return
    }

    // Uklanjamo korisnika sa zadatka
    err = services.RemoveUserFromTask(taskObjectID, userId)
    if err != nil {
        log.Printf("Failed to remove user from task: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove user from task"})
        return
    }

    log.Println("User removed from task successfully.")
    c.JSON(http.StatusOK, gin.H{"message": "User removed from task successfully"})
}

func UpdateTaskStatus(c *gin.Context) {
    id := c.Param("id")
    var statusUpdate struct {
        Status string `json:"status"`
    }

    // Parse status update request
    if err := c.BindJSON(&statusUpdate); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
        return
    }

    // Convert task ID to ObjectID
    taskObjectID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
        return
    }

    // Fetch the task by ID
    task, err := services.GetTaskById(taskObjectID)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
        return
    }

    // Get the logged-in user's ID from the context (set by AuthMiddleware)
    requesterID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized request"})
        return
    }

    // Ensure the logged-in user is assigned to the task
    if task.AssignedTo.Hex() != requesterID {
        c.JSON(http.StatusForbidden, gin.H{"error": "Only the assigned user can update the task status"})
        return
    }

    // Update the task status
    task.Status = statusUpdate.Status
    updatedTask, err := services.UpdateTask(taskObjectID.Hex(), task)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task status"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Task status updated successfully", "task": updatedTask})
}
