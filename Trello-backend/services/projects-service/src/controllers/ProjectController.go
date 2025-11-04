package controllers

import (
	"log"
	"net/http"
	"projects-service/src/models"
	"projects-service/src/repositories"
	"projects-service/src/services"
	"projects-service/src/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/gin-gonic/gin"
)

func GetAllProjects(c *gin.Context) {
	projects, err := services.GetAllProjects()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, projects)
}

func GetProjectById(c *gin.Context) {
    id := c.Param("id")
    log.Printf("ID received from request: %s", id) // Ispisuje `id` iz URL-a
    project, err := services.GetProjectById(id)
    if err != nil {
        log.Println(err) // Ispisuje detaljnu grešku
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, project)
}

func GetProjectsByManagerId(c *gin.Context) {
    id := c.Param("id")
    log.Printf("ID received from request: %s", id) // Ispisuje `id` iz URL-a
    projects, err := services.GetProjectsByManagerId(id)
    if err != nil {
        log.Println(err) // Ispisuje detaljnu grešku
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, projects)
}

func GetProjectsByMemberId(c *gin.Context) {
    id := c.Param("id")
    log.Printf("ID received from request: %s", id) // Ispisuje `id` iz URL-a
    projects, err := services.GetProjectsByMemberId(id)
    if err != nil {
        log.Println(err) // Ispisuje detaljnu grešku
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, projects)
}

func InsertProject(c *gin.Context) {
    var project models.Project
    // Dohvatanje role iz konteksta
    requesterRole, exists := c.Get("role")

    requesterRoleStr, ok := requesterRole.(string)
    if !ok {
        log.Printf("InsertProject: role is not a string: %v\n", requesterRole)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user role format"})
        return
    }

    if requesterRoleStr != "manager" {
        log.Println("InsertProject: User is not a manager")
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    log.Println("InsertProject: Received request")

    // Dohvatanje user_id iz konteksta
    requesterId, exists := c.Get("user_id")
    if !exists {
        log.Println("InsertProject: user_id not found in context")
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    log.Printf("InsertProject: Retrieved user_id from context: %v\n", requesterId)

    // Bindovanje JSON podataka na Project model
    if err := c.BindJSON(&project); err != nil {
        log.Printf("InsertProject: Failed to bind JSON to project: %v\n", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project data"})
        return
    }

    // Escapovanje unetih podataka pre umetanja u bazu
    project.Name = utils.EscapeHTML(project.Name)
    project.Description = utils.EscapeHTML(project.Description)
    project.Status = utils.EscapeHTML(project.Status)

    // Provera i dodela ManagerID
    requesterIdStr, ok := requesterId.(string)
    if !ok {
        log.Printf("InsertProject: user_id is not a string: %v\n", requesterId)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
        return
    }

    log.Printf("InsertProject: user_id as string: %s\n", requesterIdStr)

    managerID, err := primitive.ObjectIDFromHex(requesterIdStr)
    if err != nil {
        log.Printf("InsertProject: Failed to convert user_id to ObjectID: %v\n", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
        return
    }

    log.Printf("InsertProject: user_id converted to ObjectID: %s\n", managerID.Hex())
    project.ManagerID = managerID // Postavljamo ManagerID nakon bindovanja

    // Poziv servisa za umetanje projekta
    result, err := services.InsertProject(project)
    if err != nil {
        log.Printf("InsertProject: Failed to insert project: %v\n", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert project"})
        return
    }

    log.Printf("InsertProject: Project successfully inserted: %+v\n", result)

    // Vraćanje uspešnog odgovora
    c.JSON(http.StatusCreated, result)
}




func UpdateProject(c *gin.Context) {
	id := c.Param("id")
	var project models.Project
	if err := c.BindJSON(&project); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	result, err := services.UpdateProject(id, project)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, result)
}

func DeleteProject(c *gin.Context) {
	id := c.Param("id")

    requseterRole := c.GetString("role")
    if requseterRole != "manager" {
        c.JSON(http.StatusForbidden, gin.H{"error": "Access denied. Only managers can delete projects."})
        return
    }
	result, err := services.DeleteProject(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, result)
}
// AddTaskToProject handles the HTTP request to add a task to a project
func AddTaskToProject(c *gin.Context) {
    // Ekstrakcija ID projekta iz URL-a
    projectID := c.Param("id")

    // Provera user_id iz konteksta
    requesterRole, exists := c.Get("role")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    // Konverzija requesterRole u string
    requesterRoleStr, ok := requesterRole.(string)
    if !ok {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user role format"})
        return
    }

    // Provera da li je korisnik manager
    if requesterRoleStr != "manager" {
        c.JSON(http.StatusForbidden, gin.H{"error": "Access denied. Only managers can add tasks."})
        return
    }

    // Parsiranje request body-a kao Task
    var task models.Task
    if err := c.BindJSON(&task); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task data"})
        return
    }

    // Poziv servisne funkcije za dodavanje zadatka u projekat
    result, err := services.AddTaskToProjectService(projectID, task)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add task to project"})
        return
    }

    // Uspešan odgovor
    c.JSON(http.StatusOK, result)
}


func AddUserToTaskInProject(c *gin.Context) {
    projectId := c.Param("id")
    taskId := c.Param("taskId")
    userId := c.Param("userId")

    log.Printf("Project ID from request: %s", projectId)
    log.Printf("Task ID from request: %s", taskId)
    log.Printf("User ID from request: %s", userId)

    // Validacija projectId
    projectObjectID, err := primitive.ObjectIDFromHex(projectId)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
        return
    }

    // Validacija userId
    userObjectID, err := primitive.ObjectIDFromHex(userId)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
        return
    }

    // Dohvatanje requesterId iz konteksta
    requesterRole, exists := c.Get("role")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    // Konverzija requesterId u string i zatim u ObjectID
    requesterRoleStr, ok := requesterRole.(string)
    if !ok {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
        return
    }

    

    if requesterRoleStr != "manager" {
        c.JSON(http.StatusForbidden, gin.H{"error": "Access denied. Only managers can add tasks."})
        return
    }

    // Provera da li zadatak pripada projektu
	belongsToProject := services.IsTaskInProject(projectId, taskId)
	    if err != nil {
        log.Printf("Error verifying task-project relation: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify task-project relation"})
        return
    }

    if !belongsToProject {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Task does not belong to the specified project"})
        return
    }

    // Provera da li je korisnik član projekta
    isMember, err := services.IsUserMemberOfProject(projectObjectID, userObjectID)
    if err != nil {
        log.Printf("Error verifying user membership: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify user membership"})
        return
    }

    if !isMember {
        c.JSON(http.StatusForbidden, gin.H{"error": "User is not a member of the project"})
        return
    }

    // Dodavanje korisnika na zadatak
    err = services.AddUserToTask(taskId, userId)
    if err != nil {
        log.Printf("Error adding user to task: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add user to task"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "User added to task in project successfully"})
}

func RemoveUserFromTask(c *gin.Context) {
    projectId := c.Param("id")
    taskId := c.Param("taskId")
    userId := c.Param("userId")

    log.Printf("Project ID from request: %s", projectId)
    log.Printf("Task ID from request: %s", taskId)
    log.Printf("User ID from request: %s", userId)

   

   

    // Dohvatanje requesterId iz konteksta
    requesterRole := c.GetString("role")
  

    if  requesterRole != "manager" {
        c.JSON(http.StatusForbidden, gin.H{"error": "Access denied. Only managers can add tasks."})
        return
    }

    

  

    // Proveravamo da li zadatak pripada projektu
    belongsToProject := services.IsTaskInProject(projectId, taskId)
   

    if !belongsToProject {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Task does not belong to the specified project"})
        return
    }

    // Uklanjamo korisnika sa zadatka
    err := services.RemoveUserFromTask(taskId, userId)
    if err != nil {
        log.Printf("Error removing user from task: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove user from task"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "User removed from task successfully"})
}
func GetTasksByProjectId(c *gin.Context) {
    projectID := c.Param("id")
    log.Printf("Fetching tasks for project ID: %s", projectID)
    
    // Konverzija string ID-a u ObjectID
    objectID, err := primitive.ObjectIDFromHex(projectID)
    if err != nil {
        log.Printf("Invalid project ID format: %s", projectID)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID format"})
        return
    }
    
    // Poziv repozitorijumske funkcije sa ObjectID-jem
    tasks, err := repositories.GetTasksByProjectId(objectID)
    if err != nil {
        log.Printf("Error fetching tasks: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tasks"})
        return
    }
    c.JSON(http.StatusOK, tasks)
}

func GetAllMembersInProject(c *gin.Context) {
    projectId := c.Param("id")
    log.Printf("Fetching members for project ID: %s", projectId)
    
    // Konverzija string ID-a u ObjectID
    objectID, err := primitive.ObjectIDFromHex(projectId)
    if err != nil {
        log.Printf("Invalid project ID format: %s", projectId)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID format"})
        return
    }
    
    // Poziv repozitorijumske funkcije sa ObjectID-jem
    members, err := repositories.GetAllMembersInProject(objectID)
    if err != nil {
        log.Printf("Error fetching members: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch members"})
        return
    }
    c.JSON(http.StatusOK, members)
}