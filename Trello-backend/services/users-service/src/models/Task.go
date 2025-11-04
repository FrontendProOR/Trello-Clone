package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Task struct {

	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ProjectID       string             `bson:"project_id" json:"projectId"`
	Name            string             `bson:"name" json:"name"`
	Description     string             `bson:"description" json:"description"`
	AssignedTo      string             `bson:"assigned_to" json:"assignedTo"`
	CreatedAt       time.Time          `bson:"created_at" json:"createdAt"`
	DueDate         time.Time          `bson:"due_date" json:"dueDate"`
	Status          string             `bson:"status" json:"status"`
	Priority        string             `bson:"priority" json:"priority"`
	DocumentLinks   []string           `bson:"document_links" json:"documentLinks"`
	WorkflowStateID string             `bson:"workflow_state_id" json:"workflowStateId"`
}
