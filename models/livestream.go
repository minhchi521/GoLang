package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LiveStream struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title     string             `bson:"title" json:"title"`
	HostID    string             `bson:"host_id" json:"host_id"`
	Status    string             `bson:"status" json:"status"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
	StartedAt *time.Time         `bson:"started_at,omitempty" json:"started_at,omitempty"`
	EndedAt   *time.Time         `bson:"ended_at,omitempty" json:"ended_at,omitempty"`
	// Add other fields as needed
}
