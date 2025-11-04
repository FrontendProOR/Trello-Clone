package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Project struct {
	ID                   primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name                 string             `bson:"name" json:"name"`
	Description          string             `bson:"description" json:"description"`
	ManagerID            primitive.ObjectID             `bson:"manager_id" json:"managerId"`
	CreatedAt            time.Time          `bson:"created_at" json:"createdAt"`
	ExpectedEndDate      time.Time          `bson:"expected_end_date" json:"expectedEndDate"`
	MinMembers           int                `bson:"min_members" json:"minMembers"`
	MaxMembers           int                `bson:"max_members" json:"maxMembers"`
	Status               string             `bson:"status" json:"status"`
	Members              []primitive.ObjectID           `bson:"members" json:"members"`
	Tasks                []primitive.ObjectID             `bson:"tasks" json:"tasks"`
	WorkflowGraphID      string             `bson:"workflow_graph_id" json:"workflowGraphId"`
	Documents            []Document         `bson:"documents" json:"documents"`
	NotificationsEnabled bool               `bson:"notifications_enabled" json:"notificationsEnabled"`
	History              []Activity         `bson:"history" json:"history"`
	Analytics            ProjectAnalytics   `bson:"analytics" json:"analytics"`
}
