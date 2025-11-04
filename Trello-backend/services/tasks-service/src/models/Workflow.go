package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Workflow struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ProjectID   string             `bson:"project_id" json:"projectId"`
	GraphData   string             `bson:"graph_data" json:"graphData"` // JSON-encoded graph
	CreatedAt   time.Time          `bson:"created_at" json:"createdAt"`
	LastUpdated time.Time          `bson:"last_updated" json:"lastUpdated"`
}
