package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go-todo-api/internal/database"
	"go-todo-api/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Tasks handles requests to the /tasks endpoint
// It routes different HTTP methods to appropriate handlers
func Tasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case "GET":
		// Check if an ID parameter was provided in the URL query
		idParam := r.URL.Query().Get("id")

		if idParam == "" {
			// Get all tasks
			GetAllTasks(w, r)
		} else {
			// Get specific task by ID
			GetTaskByID(w, r, idParam)
		}

	case "POST":
		CreateTask(w, r)

	case "PUT":
		UpdateTask(w, r)

	case "DELETE":
		DeleteTask(w, r)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// GetAllTasks retrieves all tasks from the database
func GetAllTasks(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := database.GetCollection()
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, "Failed to fetch tasks from database", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var tasks []models.Task
	if err = cursor.All(ctx, &tasks); err != nil {
		http.Error(w, "Failed to decode tasks", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(tasks)
	fmt.Printf("✅ Retrieved %d tasks from MongoDB\n", len(tasks))
}

// GetTaskByID retrieves a specific task by its ID
func GetTaskByID(w http.ResponseWriter, r *http.Request, idParam string) {
	objectID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var task models.Task
	collection := database.GetCollection()
	err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&task)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Task not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to fetch task", http.StatusInternalServerError)
		}
		return
	}

	json.NewEncoder(w).Encode(task)
	fmt.Printf("✅ Retrieved task with ID %s\n", objectID.Hex())
}

// CreateTask creates a new task in the database
func CreateTask(w http.ResponseWriter, r *http.Request) {
	var newTask models.Task

	if err := json.NewDecoder(r.Body).Decode(&newTask); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if newTask.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}

	newTask.Completed = false

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := database.GetCollection()
	result, err := collection.InsertOne(ctx, newTask)
	if err != nil {
		http.Error(w, "Failed to create task in database", http.StatusInternalServerError)
		return
	}

	newTask.ID = result.InsertedID.(primitive.ObjectID)
	json.NewEncoder(w).Encode(newTask)
	fmt.Printf("✅ Created new task: %s with ID %s\n", newTask.Title, newTask.ID.Hex())
}

// UpdateTask updates an existing task in the database
func UpdateTask(w http.ResponseWriter, r *http.Request) {
	idParam := r.URL.Query().Get("id")
	if idParam == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	objectID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var updatedTask models.Task
	if err := json.NewDecoder(r.Body).Decode(&updatedTask); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := database.GetCollection()

	// Find existing task first
	var existingTask models.Task
	err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&existingTask)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Task not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to fetch task", http.StatusInternalServerError)
		}
		return
	}

	// Update fields
	update := bson.M{
		"$set": bson.M{
			"title":     updatedTask.Title,
			"completed": updatedTask.Completed,
		},
	}

	// If title is empty in update, keep the existing title
	if updatedTask.Title == "" {
		update["$set"].(bson.M)["title"] = existingTask.Title
	}

	result, err := collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		http.Error(w, "Failed to update task", http.StatusInternalServerError)
		return
	}

	if result.MatchedCount == 0 {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	// Fetch and return the updated task
	collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&updatedTask)
	updatedTask.ID = objectID
	json.NewEncoder(w).Encode(updatedTask)
	fmt.Printf("✅ Updated task with ID %s\n", objectID.Hex())
}

// DeleteTask deletes a task from the database
func DeleteTask(w http.ResponseWriter, r *http.Request) {
	idParam := r.URL.Query().Get("id")
	if idParam == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	objectID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := database.GetCollection()
	result, err := collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		http.Error(w, "Failed to delete task", http.StatusInternalServerError)
		return
	}

	if result.DeletedCount == 0 {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	fmt.Fprintf(w, `{"message": "Task deleted successfully", "id": "%s"}`, idParam)
	fmt.Printf("✅ Deleted task with ID %s\n", objectID.Hex())
}
