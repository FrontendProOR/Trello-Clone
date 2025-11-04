package clients

import (
	"bytes"
	"encoding/json"
	"projects-service/src/models"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"fmt"
	"log"
	"net/http"
)

// AddUserToTask šalje HTTP PUT zahtev tasks-service-u da doda korisnika na zadatak
func AddUserToTask(taskId string, userId string) error {
    tasksServiceUrl := fmt.Sprintf("http://tasks-service:8080/tasks/%s/users/%s", taskId, userId)
    req, err := http.NewRequest(http.MethodPut, tasksServiceUrl, nil)
    if err != nil {
        log.Printf("Error creating HTTP request: %v", err)
        return err
    }

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        log.Printf("Error making HTTP request: %v", err)
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        log.Printf("Failed to add user to task, status code: %d", resp.StatusCode)
        return fmt.Errorf("failed to add user to task, status code: %d", resp.StatusCode)
    }

    log.Println("User successfully added to task.")
    return nil
}
func CreateTaskInTasksService(task models.Task, projectId string) (primitive.ObjectID, error) {
    id, err := primitive.ObjectIDFromHex(projectId)
    if err != nil {
        // Handle the error
        return primitive.NilObjectID, err
    }
    task.ProjectID = id
    // Kreiranje JSON tela sa zadatkom
    jsonTask, err := json.Marshal(task)
    if err != nil {
        return primitive.NilObjectID, err
    }

    // API poziv ka tasks-service
    tasksServiceUrl := "http://tasks-service:8080/tasks"
    req, err := http.NewRequest(http.MethodPost, tasksServiceUrl, bytes.NewBuffer(jsonTask))
    if err != nil {
        return primitive.NilObjectID, err
    }
    //get token
    token := req.Header.Get("Authorization")

    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Accept", "application/json")
    req.Header.Set("Authorization", "Bearer "+token)

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return primitive.NilObjectID, err
    }
    defer resp.Body.Close()

    // Provera uspešnosti zahteva
    if resp.StatusCode != http.StatusCreated {
        return primitive.NilObjectID, fmt.Errorf("failed to create task in tasks-service, status code: %d", resp.StatusCode)
    }

    // Dekodiramo odgovor i dobijamo ID zadatka
    var result struct {
        ID primitive.ObjectID `json:"id"`
    }
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return primitive.NilObjectID, err
    }

    log.Printf("Task created in tasks-service with ID: %s", result.ID.Hex())
    return result.ID, nil
}
func RemoveUserFromTask(taskId, userId string) error {
    tasksServiceUrl := fmt.Sprintf("http://tasks-service:8080/tasks/%s/users/%s", taskId, userId)
    
    req, err := http.NewRequest(http.MethodDelete, tasksServiceUrl, nil)
    if err != nil {
        log.Printf("Error creating HTTP request: %v", err)
        return err
    }

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        log.Printf("Error making HTTP request: %v", err)
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        log.Printf("Failed to remove user from task, status code: %d", resp.StatusCode)
        return fmt.Errorf("failed to remove user from task, status code: %d", resp.StatusCode)
    }

    log.Println("User successfully removed from task.")
    return nil
}
