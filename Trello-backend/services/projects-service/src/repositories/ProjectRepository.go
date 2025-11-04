package repositories

import (
	"context"
	"fmt"
	"log"
	"projects-service/src/clients"
	"projects-service/src/config"
	"projects-service/src/models"

	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// var projectCollection *mongo.Collection = config.DB.Collection("projects")
var projectCollection *mongo.Collection
var taskCollection *mongo.Collection
var userCollection *mongo.Collection

func Init() {
	projectCollection = config.DB.Collection("projects")
	taskCollection = config.DB.Collection("tasks")
	userCollection = config.DB.Collection("users")
}
func GetAllProjects() ([]models.Project, error) {
	var projects []models.Project
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := projectCollection.Find(ctx, bson.M{})
	if err != nil {
		return projects, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var project models.Project
		if err := cursor.Decode(&project); err != nil {
			return projects, err
		}
		projects = append(projects, project)
	}
	return projects, nil
}

func GetProjectById(id primitive.ObjectID) (models.Project, error) {
	var project models.Project
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := projectCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&project)
	return project, err
}

func GetProjectsByManagerId(id string) ([]models.Project, error) {
	var projects []models.Project
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	cursor, err := projectCollection.Find(ctx, bson.M{"manager_id": objectID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var project models.Project
		if err := cursor.Decode(&project); err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return projects, nil
}

func GetProjectsByMemberId(id string) ([]models.Project, error) {
	var projects []models.Project
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Convert the string id to ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	// Query MongoDB to find projects where the member's id is in the members array
	cursor, err := projectCollection.Find(ctx, bson.M{"members": objectID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// Iterate through the cursor and decode each project
	for cursor.Next(ctx) {
		var project models.Project
		if err := cursor.Decode(&project); err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}

	// Check if there was an error during iteration
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	// Return the projects
	return projects, nil
}

func InsertProject(project models.Project) (*mongo.InsertOneResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	project.CreatedAt = time.Now()
	result, err := projectCollection.InsertOne(ctx, project)
	return result, err
}

func UpdateProject(id primitive.ObjectID, project models.Project) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	update := bson.M{"$set": project}
	result, err := projectCollection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return result, err
}

func DeleteProject(id primitive.ObjectID) (*mongo.DeleteResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := projectCollection.DeleteOne(ctx, bson.M{"_id": id})
	return result, err
}

// AddTaskToProject adds a new task to an existing project
func AddTaskToProject(projectID primitive.ObjectID, task models.Task) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create task in tasks-service and get its ID directly as an ObjectID
	taskID, err := clients.CreateTaskInTasksService(task, projectID.Hex())
	if err != nil {
		log.Printf("Error sending task to tasks-service: %v", err)
		return nil, err
	}

	// Push the ObjectID of the task directly into the `tasks` array
	update := bson.M{"$push": bson.M{"tasks": taskID}}
	result, err := projectCollection.UpdateOne(ctx, bson.M{"_id": projectID}, update)

	if err != nil || result.MatchedCount == 0 {
		log.Printf("Failed to add task ID %s to project %s", taskID.Hex(), projectID.Hex())
	}

	return result, err
}
func GetTasksByProjectId(projectID primitive.ObjectID) ([]models.Task, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    // Dohvatanje projekta
    var project models.Project
    err := projectCollection.FindOne(ctx, bson.M{"_id": projectID}).Decode(&project)
    if err != nil {
        log.Printf("Error fetching project with ID %v: %v", projectID, err)
        return nil, fmt.Errorf("Project not found: %v", err)
    }

    // Provera da li projekat ima taskove
    if len(project.Tasks) == 0 {
        log.Printf("Project with ID %v has no tasks", projectID)
        return []models.Task{}, nil // Projekat nema taskove
    }

    // Dohvatanje taskova
    filter := bson.M{"_id": bson.M{"$in": project.Tasks}}
    log.Printf("Filter for tasks: %+v", filter)  // Proveri filter koji se koristi za pretragu taskova

    cursor, err := taskCollection.Find(ctx, filter)
    if err != nil {
        log.Printf("Error finding tasks for project %v: %v", projectID, err)
        return nil, fmt.Errorf("Error finding tasks: %v", err)
    }
    defer cursor.Close(ctx)

    // Provera da li je cursor prazan
    if !cursor.Next(ctx) {
        log.Printf("No tasks found for project with ID %v", projectID)
        return nil, fmt.Errorf("No tasks found for project with ID %v", projectID)
    }

    var tasks []models.Task
    // Dohvatanje svih taskova
    if err := cursor.All(ctx, &tasks); err != nil {
        log.Printf("Error decoding tasks for project with ID %v: %v", projectID, err)
        return nil, fmt.Errorf("Error decoding tasks: %v", err)
    }

    // Provera da li su tasks inicijalizovani
    if len(tasks) == 0 {
        log.Printf("No tasks decoded for project with ID %v", projectID)
    }

    return tasks, nil
}

func GetAllMembersInProject(projectID primitive.ObjectID) ([]models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Dohvatanje projekta
	var project models.Project
	err := projectCollection.FindOne(ctx, bson.M{"_id": projectID}).Decode(&project)
	if err != nil {
		log.Printf("Error fetching project with ID %v: %v", projectID, err)
		return nil, fmt.Errorf("Project not found: %v", err)
	}

	// Dohvatanje članova projekta
	filter := bson.M{"_id": bson.M{"$in": project.Members}}
	cursor, err := userCollection.Find(ctx, filter)
	if err != nil {
		log.Printf("Error finding members for project %v: %v", projectID, err)
		return nil, fmt.Errorf("Error finding members: %v", err)
	}
	defer cursor.Close(ctx)

	// Provera da li je cursor prazan
	if !cursor.Next(ctx) {
		log.Printf("No members found for project with ID %v", projectID)
		return nil, fmt.Errorf("No members found for project with ID %v", projectID)
	}

	var members []models.User
	// Dohvatanje svih članova projekta
	if err := cursor.All(ctx, &members); err != nil {
		log.Printf("Error decoding members for project with ID %v: %v", projectID, err)
		return nil, fmt.Errorf("Error decoding members: %v", err)
	}

	// Provera da li su članovi inicijalizovani
	if len(members) == 0 {
		log.Printf("No members decoded for project with ID %v", projectID)
	}

	return members, nil
	}