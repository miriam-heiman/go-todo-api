package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Task represents a todo item in our application
type Task struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"` // MongoDB's ObjectID type for unique IDs
	Title       string             `json:"title"`                   // Title of the task
	Description string             `json:"description"`             // Description of the task
	Completed   bool               `json:"completed"`               // Whether the task is completed
}
