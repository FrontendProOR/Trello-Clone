package services

import (
	"tasks-service/src/models"
	"tasks-service/src/repositories"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetAllTasks() ([]models.Task, error) {
	return repositories.GetAllTasks()
}

func GetTaskById(id primitive.ObjectID) (models.Task, error) {
	return repositories.GetTaskById(id)
}

func InsertTask(task models.Task) (string, error) {
    // Postavljamo ID zadatka i trenutni datum
    task.ID = primitive.NewObjectID()
    task.CreatedAt = task.CreatedAt // postavite na trenutni timestamp ako je potrebno

    // Pozivamo funkciju u repository-u
    result, err := repositories.InsertTask(task)
    if err != nil {
        return "", err
    }
    return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

func UpdateTask(id string, task models.Task) (interface{}, error) {
	objectId, _ := primitive.ObjectIDFromHex(id)
	return repositories.UpdateTask(objectId, task)
}

func DeleteTask(id string) (interface{}, error) {
	objectId, _ := primitive.ObjectIDFromHex(id)
	return repositories.DeleteTask(objectId)
}
func AddUserToTask(taskID primitive.ObjectID, userID primitive.ObjectID) error {
    return repositories.AddUserToTask(taskID, userID)
}
func RemoveUserFromTask(taskID primitive.ObjectID, userID string) error {
    return repositories.RemoveUserFromTask(taskID, userID)
}
func GetTaskStatus(taskId string) (string, error) {
    taskObjectID, err := primitive.ObjectIDFromHex(taskId)
    if err != nil {
        return "", fmt.Errorf("invalid task ID format: %v", err)
    }

    task, err := repositories.GetTaskById(taskObjectID)
    if err != nil {
        return "", fmt.Errorf("failed to fetch task: %v", err)
    }

    return task.Status, nil
}
