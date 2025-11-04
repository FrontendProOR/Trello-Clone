package clients

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func CheckIfManager(projectId, userId string, authHeader string) (bool, error) {
    projectServiceURL := fmt.Sprintf("http://projects-service:8080/projects/%s/isManager/%s", projectId, userId)

    log.Printf("CheckIfManager: Sending request to %s", projectServiceURL)

    req, err := http.NewRequest(http.MethodGet, projectServiceURL, nil)
    if err != nil {
        log.Printf("CheckIfManager: Error creating request: %v", err)
        return false, err
    }

    // Dodajte Authorization header
    req.Header.Set("Authorization", authHeader)

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        log.Printf("CheckIfManager: Error making request to project-service: %v", err)
        return false, err
    }
    defer resp.Body.Close()

    log.Printf("CheckIfManager: Response status code: %d", resp.StatusCode)

    if resp.StatusCode != http.StatusOK {
        log.Printf("CheckIfManager: Project-service returned status code %d", resp.StatusCode)
        return false, nil
    }

    var response struct {
        IsManager bool `json:"isManager"`
    }
    if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
        log.Printf("CheckIfManager: Error decoding response: %v", err)
        return false, err
    }

    log.Printf("CheckIfManager: Manager status: %t", response.IsManager)
    return response.IsManager, nil
}

