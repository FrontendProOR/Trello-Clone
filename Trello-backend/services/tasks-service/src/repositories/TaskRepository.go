package repositories

import (
	"context"
	"fmt"
	"log"
	"tasks-service/src/config"
	"tasks-service/src/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var taskCollection *mongo.Collection

func Init() {
	taskCollection = config.DB.Collection("tasks")
}

func GetAllTasks() ([]models.Task, error) {
	var tasks []models.Task
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := taskCollection.Find(ctx, bson.M{})
	if err != nil {
		return tasks, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var task models.Task
		if err := cursor.Decode(&task); err != nil {
			return tasks, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func GetTaskById(id primitive.ObjectID) (models.Task, error) {
	var task models.Task
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := taskCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&task)
	return task, err
}

func InsertTask(task models.Task) (*mongo.InsertOneResult, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    // Proverite da li je `project_id` validan pre postavljanja
    if task.ProjectID.IsZero() {
        return nil, fmt.Errorf("invalid project ID for task")
    }
	task.Status = string(models.Pending)
    task.CreatedAt = time.Now()
    result, err := taskCollection.InsertOne(ctx, task)
    return result, err
}

func UpdateTask(id primitive.ObjectID, task models.Task) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	update := bson.M{"$set": task}
	result, err := taskCollection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return result, err
}

func DeleteTask(id primitive.ObjectID) (*mongo.DeleteResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := taskCollection.DeleteOne(ctx, bson.M{"_id": id})
	return result, err
}
func AddUserToTask(taskID primitive.ObjectID, userID primitive.ObjectID) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    log.Printf("Updating task %s with assigned user %s", taskID.Hex(), userID.Hex())

    // Ažuriramo AssignedTo polje u zadatku
    update := bson.M{"$set": bson.M{"assigned_to": userID}}
    result, err := taskCollection.UpdateOne(ctx, bson.M{"_id": taskID}, update)

    if err != nil {
        log.Printf("Failed to update task %s: %v", taskID.Hex(), err)
        return err
    }

    log.Printf("Update result: %v", result)
    if result.ModifiedCount == 0 {
        log.Printf("No task was updated. Verify that task %s exists.", taskID.Hex())
    }

    return nil
}
func RemoveUserFromTask(taskID primitive.ObjectID, userID string) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
	
    // Proveravamo da li korisnik već nije dodeljen
    update := bson.M{"$unset": bson.M{"assigned_to": ""}}

    _, err := taskCollection.UpdateOne(ctx, bson.M{"_id": taskID}, update)
    return err
}
