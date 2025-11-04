package controllers

import (
	"net/http"
	"notifications-service/models"
	"notifications-service/repositories"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
)

type NotificationsController struct {
	Repo *repositories.NotificationsRepository
}

// Preuzimanje svih notifikacija za korisnika (userID)
func (c *NotificationsController) GetNotifications(ctx *gin.Context) {
	userID := ctx.Param("user_id")

	notifications, err := c.Repo.GetNotificationsByUserID(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching notifications"})
		return
	}

	ctx.JSON(http.StatusOK, notifications)
}

// Preuzimanje nepročitanih notifikacija
func (c *NotificationsController) GetUnreadNotifications(ctx *gin.Context) {
	userID := ctx.Param("user_id")

	notifications, err := c.Repo.GetUnreadNotifications(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch unread notifications"})
		return
	}

	ctx.JSON(http.StatusOK, notifications)
}

// Kreiranje notifikacije
func (c *NotificationsController) CreateNotification(ctx *gin.Context) {
	var req struct {
		UserID  string `json:"userId"`
		Message string `json:"message"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	notification := &models.Notification{
		ID:        gocql.TimeUUID().String(),
		UserID:    req.UserID,
		Message:   req.Message,
		CreatedAt: time.Now(),
		IsRead:    false,
	}

	err := c.Repo.CreateNotification(notification)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create notification"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Notification created successfully"})
}

// Označavanje notifikacije kao pročitane
func (c *NotificationsController) MarkNotificationAsRead(ctx *gin.Context) {
	var req struct {
		UserID         string `json:"userId"`
		NotificationID string `json:"notificationId"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	err := c.Repo.MarkAsRead(req.UserID, req.NotificationID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark notification as read"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Notification marked as read"})
}

// Brisanje notifikacije
func (c *NotificationsController) DeleteNotification(ctx *gin.Context) {
	var req struct {
		UserID         string `json:"userId"`
		NotificationID string `json:"notificationId"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	err := c.Repo.DeleteNotification(req.UserID, req.NotificationID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete notification"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Notification deleted successfully"})
}
