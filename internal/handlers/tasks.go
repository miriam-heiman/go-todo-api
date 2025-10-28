package handlers

import (
	"context"
	"fmt"
	"time"

	"go-todo-api/internal/database"
	"go-todo-api/internal/models"

	"github.com/danielgtaylor/huma/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetAllTasks retrieves all tasks from the database
func GetAllTasks(ctx context.Context, input *struct{}) (*models.GetTasksOutput, error) {
	dbCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := database.GetCollection()
	cursor, err := collection.Find(dbCtx, bson.M{})
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to fetch tasks from database")
	}
	defer cursor.Close(dbCtx)

	var tasks []models.Task
	if err = cursor.All(dbCtx, &tasks); err != nil {
		return nil, huma.Error500InternalServerError("Failed to decode tasks")
	}

	// Initialize empty array instead of null if no tasks
	if tasks == nil {
		tasks = []models.Task{}
	}

	fmt.Printf("✅ Retrieved %d tasks from MongoDB\n", len(tasks))
	return &models.GetTasksOutput{Body: tasks}, nil
}

// GetTaskByID retrieves a specific task by its ID
func GetTaskByID(ctx context.Context, input *models.GetTaskInput) (*models.GetTaskOutput, error) {
	objectID, err := primitive.ObjectIDFromHex(input.ID)
	if err != nil {
		return nil, huma.Error400BadRequest("Invalid task ID format")
	}

	dbCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var task models.Task
	collection := database.GetCollection()
	err = collection.FindOne(dbCtx, bson.M{"_id": objectID}).Decode(&task)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, huma.Error404NotFound("Task not found")
		}
		return nil, huma.Error500InternalServerError("Failed to fetch task")
	}

	fmt.Printf("✅ Retrieved task with ID %s\n", objectID.Hex())
	return &models.GetTaskOutput{Body: task}, nil
}

// CreateTask creates a new task in the database
func CreateTask(ctx context.Context, input *models.CreateTaskInput) (*models.CreateTaskOutput, error) {
	newTask := models.Task{
		Title:       input.Body.Title,
		Description: input.Body.Description,
		Completed:   false,
	}

	dbCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := database.GetCollection()
	result, err := collection.InsertOne(dbCtx, newTask)
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to create task in database")
	}

	newTask.ID = result.InsertedID.(primitive.ObjectID)
	fmt.Printf("✅ Created new task: %s with ID %s\n", newTask.Title, newTask.ID.Hex())

	return &models.CreateTaskOutput{Body: newTask}, nil
}

// UpdateTask updates an existing task in the database
func UpdateTask(ctx context.Context, input *models.UpdateTaskInput) (*models.UpdateTaskOutput, error) {
	objectID, err := primitive.ObjectIDFromHex(input.ID)
	if err != nil {
		return nil, huma.Error400BadRequest("Invalid task ID format")
	}

	dbCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := database.GetCollection()

	// Find existing task first
	var existingTask models.Task
	err = collection.FindOne(dbCtx, bson.M{"_id": objectID}).Decode(&existingTask)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, huma.Error404NotFound("Task not found")
		}
		return nil, huma.Error500InternalServerError("Failed to fetch task")
	}

	// Build update document with only provided fields
	update := bson.M{"$set": bson.M{}}

	if input.Body.Title != nil {
		update["$set"].(bson.M)["title"] = *input.Body.Title
	}
	if input.Body.Description != nil {
		update["$set"].(bson.M)["description"] = *input.Body.Description
	}
	if input.Body.Completed != nil {
		update["$set"].(bson.M)["completed"] = *input.Body.Completed
	}

	// Only update if there are fields to update
	if len(update["$set"].(bson.M)) == 0 {
		return nil, huma.Error400BadRequest("No fields to update")
	}

	result, err := collection.UpdateOne(dbCtx, bson.M{"_id": objectID}, update)
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to update task")
	}

	if result.MatchedCount == 0 {
		return nil, huma.Error404NotFound("Task not found")
	}

	// Fetch and return the updated task
	var updatedTask models.Task
	collection.FindOne(dbCtx, bson.M{"_id": objectID}).Decode(&updatedTask)

	fmt.Printf("✅ Updated task with ID %s\n", objectID.Hex())
	return &models.UpdateTaskOutput{Body: updatedTask}, nil
}

// DeleteTask deletes a task from the database
func DeleteTask(ctx context.Context, input *models.DeleteTaskInput) (*models.DeleteTaskOutput, error) {
	objectID, err := primitive.ObjectIDFromHex(input.ID)
	if err != nil {
		return nil, huma.Error400BadRequest("Invalid task ID format")
	}

	dbCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := database.GetCollection()
	result, err := collection.DeleteOne(dbCtx, bson.M{"_id": objectID})
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to delete task")
	}

	if result.DeletedCount == 0 {
		return nil, huma.Error404NotFound("Task not found")
	}

	fmt.Printf("✅ Deleted task with ID %s\n", objectID.Hex())

	return &models.DeleteTaskOutput{
		Body: struct {
			Message string `json:"message" doc:"Success message"`
			ID      string `json:"id" doc:"Deleted task ID"`
		}{
			Message: "Task deleted successfully",
			ID:      input.ID,
		},
	}, nil
}
