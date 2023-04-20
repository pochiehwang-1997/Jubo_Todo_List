package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	TodoModel struct {
		ID          primitive.ObjectID `bson:"_id,omitempty"`
		Title       string             `bson:"title"`
		Description string             `bson:"description"`
		Completed   bool               `bson:"completed"`
		CreatedAt   time.Time          `bson:"createdAt"`
	}

	Todo struct {
		ID          string    `json:"id"`
		Title       string    `json:"title"`
		Description string    `json:"description"`
		Completed   bool      `json:"completed"`
		CreatedAt   time.Time `json:"createdAt"`
	}
)
