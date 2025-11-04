package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)


type Status string
const (
	Pending Status = "Pending"
	Finished  Status = "Finished"
	InProgress  Status = "In_Progress"
)

type Task struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ProjectID   primitive.ObjectID `bson:"project_id" json:"projectId"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description" json:"description"`
	AssignedTo  primitive.ObjectID             `bson:"assigned_to" json:"assignedTo"`
	CreatedAt   time.Time          `bson:"created_at" json:"createdAt"`
	DueDate     time.Time          `bson:"due_date" json:"dueDate"`
	Status      Status             `bson:"status" json:"status"`
	Priority    string             `bson:"priority" json:"priority"`
	// SubTasks        []SubTask          `bson:"sub_tasks" json:"subTasks"`
	// Comments        []Comment          `bson:"comments" json:"comments"`
	DocumentLinks   []string `bson:"document_links" json:"documentLinks"`
	WorkflowStateID string   `bson:"workflow_state_id" json:"workflowStateId"`
}
