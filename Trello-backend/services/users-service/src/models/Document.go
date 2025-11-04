package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Document struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ProjectID  string             `bson:"project_id" json:"projectId"`
	TaskID     string             `bson:"task_id,omitempty" json:"taskId"`
	Name       string             `bson:"name" json:"name"`
	URL        string             `bson:"url" json:"url"`
	UploadedAt time.Time          `bson:"uploaded_at" json:"uploadedAt"`
	UploadedBy string             `bson:"uploaded_by" json:"uploadedBy"`
}
