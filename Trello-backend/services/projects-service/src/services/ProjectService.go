package services

import (
	"errors"
	"fmt"
	"log"
	"projects-service/src/clients"
	"projects-service/src/models"
	"projects-service/src/repositories"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetAllProjects() ([]models.Project, error) {
	return repositories.GetAllProjects()
}

func GetProjectById(id string) (models.Project, error) {
    log.Printf("Received project ID: %s\n", id) // Ispisuje ID koji se prosleđuje
    objectId, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return models.Project{}, fmt.Errorf("invalid project ID format for ID %s: %v", id, err)
    }
    return repositories.GetProjectById(objectId)
}

func GetProjectsByManagerId(id string) ([]models.Project, error) {
    return repositories.GetProjectsByManagerId(id)
}

func GetProjectsByMemberId(id string) ([]models.Project, error) {
    return repositories.GetProjectsByMemberId(id)
}

func InsertProject(project models.Project) (interface{}, error) {
	return repositories.InsertProject(project)
}

func UpdateProject(id string, project models.Project) (interface{}, error) {
	objectId, _ := primitive.ObjectIDFromHex(id)
	return repositories.UpdateProject(objectId, project)
}

func DeleteProject(id string) (interface{}, error) {
	objectId, _ := primitive.ObjectIDFromHex(id)
	return repositories.DeleteProject(objectId)
}
// AddTaskToProjectService adds a task to a project by calling the repository function
func AddTaskToProjectService(projectID string, task models.Task) (interface{}, error) {
	// Convert the project ID from string to ObjectID
	objectID, err := primitive.ObjectIDFromHex(projectID)
	if err != nil {
		return nil, err
	}

    //set assigne to managerId of project get project with projectId then use managerId from project to set assigneTo of the task
    project, err := repositories.GetProjectById(objectID)
    if err != nil {
        return nil, err
    }
    task.AssignedTo = project.ManagerID

	// Call the repository function to add the task
	return repositories.AddTaskToProject(objectID, task)
}
func IsTaskInProject(projectId, taskId string) bool {
    projectObjectID, err := primitive.ObjectIDFromHex(projectId)
    if err != nil {
        log.Printf("Invalid project ID format: %s", projectId)
        return false
    }

    taskObjectID, err := primitive.ObjectIDFromHex(taskId)
    if err != nil {
        log.Printf("Invalid task ID format: %s", taskId)
        return false
    }

    project, err := repositories.GetProjectById(projectObjectID)
    if err != nil {
        log.Printf("Error fetching project: %v", err)
        return false
    }

    for _, id := range project.Tasks {
        if id == taskObjectID {
            return true
        }
    }

    log.Printf("Task %s not found in project %s", taskId, projectId)
    return false
}


// AddUserToTask poziva odgovarajući klijent za dodavanje korisnika na zadatak
func AddUserToTask(taskId, userId string) error {
    return clients.AddUserToTask(taskId, userId)
}

func RemoveUserFromTask(taskId, userId string) error {
    // Pozivamo klijenta za uklanjanje korisnika sa zadatka
    return clients.RemoveUserFromTask(taskId, userId)
}

func IsUserManagerOfProject(projectId, userID primitive.ObjectID) (bool, error) {
    project, err := repositories.GetProjectById(projectId)
    if err != nil {
        return false , errors.New("project not found")
    }
    if project.ManagerID == userID {
        return true, nil
    }
    return false, nil
}


func IsUserMemberOfProject(projectID, userID primitive.ObjectID) (bool, error) {
    // Dohvatamo projekat
    project, err := repositories.GetProjectById(projectID)
    if err != nil {
        log.Printf("Error fetching project: %v", err)
        return false, err
    }

    // Proveravamo da li je korisnik u listi članova projekta
    for _, member := range project.Members {
        if member == userID {
            return true, nil
        }
    }

    return false, nil
}
func GetTasksByProjectId(projectId primitive.ObjectID) ([]models.Task, error) {
    return repositories.GetTasksByProjectId(projectId)
}