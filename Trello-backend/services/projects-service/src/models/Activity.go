package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Activity struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID      string             `bson:"user_id" json:"userId"`
	ProjectID   string             `bson:"project_id" json:"projectId"`
	TaskID      string             `bson:"task_id,omitempty" json:"taskId"`
	Action      string             `bson:"action" json:"action"`
	Timestamp   time.Time          `bson:"timestamp" json:"timestamp"`
	Description string             `bson:"description" json:"description"`
}
