package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Task represents a todo item in our application
type Task struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id" doc:"Unique identifier for the task"`
	Title       string             `json:"title" doc:"Title of the task" minLength:"1" maxLength:"200"`
	Description string             `json:"description,omitempty" doc:"Detailed description of the task" maxLength:"1000"`
	Completed   bool               `json:"completed" doc:"Whether the task is completed"`
}

// CreateTaskInput is the input for creating a new task
type CreateTaskInput struct {
	Body struct {
		Title       string `json:"title" doc:"Title of the task" minLength:"1" maxLength:"200" example:"Buy groceries"`
		Description string `json:"description,omitempty" doc:"Detailed description" maxLength:"1000" example:"Buy milk, eggs, and bread"`
	}
}

// CreateTaskOutput is the response for creating a task
type CreateTaskOutput struct {
	Body Task
}

// GetTasksOutput is the response for getting all tasks
type GetTasksOutput struct {
	Body []Task
}

// GetTaskInput is the input for getting a single task
type GetTaskInput struct {
	ID string `path:"id" doc:"Task ID" minLength:"24" maxLength:"24"`
}

// GetTaskOutput is the response for getting a single task
type GetTaskOutput struct {
	Body Task
}

// UpdateTaskInput is the input for updating a task
type UpdateTaskInput struct {
	ID   string `path:"id" doc:"Task ID" minLength:"24" maxLength:"24"`
	Body struct {
		Title       *string `json:"title,omitempty" doc:"Title of the task" minLength:"1" maxLength:"200"`
		Description *string `json:"description,omitempty" doc:"Detailed description" maxLength:"1000"`
		Completed   *bool   `json:"completed,omitempty" doc:"Whether the task is completed"`
	}
}

// UpdateTaskOutput is the response for updating a task
type UpdateTaskOutput struct {
	Body Task
}

// DeleteTaskInput is the input for deleting a task
type DeleteTaskInput struct {
	ID string `path:"id" doc:"Task ID" minLength:"24" maxLength:"24"`
}

// DeleteTaskOutput is the response for deleting a task
type DeleteTaskOutput struct {
	Body struct {
		Message string `json:"message" doc:"Success message"`
		ID      string `json:"id" doc:"Deleted task ID"`
	}
}

// HealthOutput is the response for the health check
type HealthOutput struct {
	Body struct {
		Status  string `json:"status" doc:"Health status" example:"healthy"`
		Message string `json:"message" doc:"Health message" example:"Server is running with MongoDB!"`
	}
}
