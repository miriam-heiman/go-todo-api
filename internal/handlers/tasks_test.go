package handlers

import (
	"context"
	"os"
	"testing"

	"go-todo-api/internal/database"
	"go-todo-api/internal/logger"
	"go-todo-api/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TestMain runs before all tests and handles setup/teardown
func TestMain(m *testing.M) {

	// Setup: Initialise logger first (database.Connect needs it)
	logger.Init()

	// Setup: Connect to MongoDB before running tests
	database.Connect()

	// Run all tests
	code := m.Run()

	// Teardown: Close connection after all tests
	database.Close()

	// Exit with test result code
	os.Exit(code)
}

// Note: These tests require MongoDB to be running
// Run with: go test /internal handlers -v

// ============================================================================
// GET ALL TASKS - EMPTY DATABASE
// ============================================================================

// TestGetAllTasks_EmptyDatabase tests getting tasks when database is empty
func TestGetAllTasks_EmptyDatabase(t *testing.T) {
	// Skip this MongoDB integration function in short mode
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange: Clean database
	ctx := context.Background()
	collection := database.GetCollection()
	collection.DeleteMany(ctx, bson.M{}) // Clear all tasks

	// Act: Get all tasks
	input := &models.GetTasksInput{}
	output, err := GetAllTasks(ctx, input)

	// Assert
	if err != nil {
		t.Fatalf("GetAllTasks returned error: %v", err)
	}

	if output == nil {
		t.Fatal("Output is nil")
	}

	if len(output.Body) != 0 {
		t.Errorf("Expected 0 tasks, got %d", len(output.Body))
	}

	t.Log("✅ GetAllTasks with empty database passed")
}

// ============================================================================
// GETALLTASKS - WITH TASKS
// ============================================================================

// TestGetAllTasks_WithTasks tests getting tasks when some exist
func TestGetAllTasks_WithTasks(t *testing.T) {
	// Skip this MongoDB integration function in short mode
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange: Clean database and insert test tasks
	ctx := context.Background()
	collection := database.GetCollection()
	collection.DeleteMany(ctx, bson.M{})

	// Insert 2 test tasks
	testTasks := []interface{}{
		models.Task{
			ID:          primitive.NewObjectID(),
			Title:       "Test Task 1",
			Description: "First test task",
			Completed:   false,
		},
		models.Task{
			ID:          primitive.NewObjectID(),
			Title:       "Test Task 2",
			Description: "Second test task",
			Completed:   true,
		},
	}

	_, err := collection.InsertMany(ctx, testTasks)
	if err != nil {
		t.Fatalf("Failed to insert test tasks: %v", err)
	}

	// Act: Get all tasks
	input := &models.GetTasksInput{}
	output, err := GetAllTasks(ctx, input)

	// Assert
	if err != nil {
		t.Fatalf("GetAllTasks returned error: %v", err)
	}

	if len(output.Body) != 2 {
		t.Errorf("Expected 2 tasks, got %d", len(output.Body))
	}

	// Cleanup
	collection.DeleteMany(ctx, bson.M{})
	t.Log("✅ GetAllTasks with tasks passed")
}

// ============================================================================
// GETALLTASKS - FILTER BY COMPLETED
// ============================================================================

// TestGetAllTasks_FilteredCompleted tests filtering by completed status
func TestGetAllTasks_FilterCompleted(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange
	ctx := context.Background()
	collection := database.GetCollection()
	collection.DeleteMany(ctx, bson.M{})

	// Insert mix of completed and incomplete tasks
	testTasks := []interface{}{
		models.Task{ID: primitive.NewObjectID(), Title: "Task 1", Completed: false},
		models.Task{ID: primitive.NewObjectID(), Title: "Task 2", Completed: true},
		models.Task{ID: primitive.NewObjectID(), Title: "Task 3", Completed: false},
	}
	collection.InsertMany(ctx, testTasks)

	// Act: Get only completed tasks
	input := &models.GetTasksInput{Completed: "true"}
	output, err := GetAllTasks(ctx, input)

	// Assert
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	if len(output.Body) != 1 {
		t.Errorf("Expected 1 completed tasks, got %d", len(output.Body))
	}

	if output.Body[0].Title != "Task 2" {
		t.Errorf("Expected 'Task 2', got '%v'", output.Body[0])
	}

	// Cleanup
	collection.DeleteMany(ctx, bson.M{})
	t.Log("✅ Filter by completed passed")
}

// ============================================================================
// TEST CREATETASK
// ============================================================================
// TestCreateTask tests creating a new task
func TestCreateTask(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange
	ctx := context.Background()
	collection := database.GetCollection()
	collection.DeleteMany(ctx, bson.M{})

	input := &models.CreateTaskInput{
		Body: struct {
			Title       string `json:"title" doc:"Title of the task" minLength:"1" maxLength:"200" example:"Buy groceries"`
			Description string `json:"description,omitempty" doc:"Detailed description" maxLength:"1000" example:"Buy milk, eggs, and bread"`
		}{
			Title:       "New Test Task",
			Description: "Testing task creation",
		},
	}

	// Act
	output, err := CreateTask(ctx, input)

	// Assert
	if err != nil {
		t.Fatalf("CreateTask returned error: %v", err)
	}

	if output.Body.Title != "New Test Task" {
		t.Error("Expected non-zero ID")
	}

	if output.Body.Completed != false {
		t.Error("Expected new task to be incomplete")
	}

	// Cleanup
	collection.DeleteMany(ctx, bson.M{})
	t.Logf("✅ CreateTask passed. Created task with ID: %s", output.Body.ID.Hex())
}

// ============================================================================
// TEST GETTASKBHID
// ============================================================================
// TestGetTaskByID tests retrieving a specific task
func TestGetTaskByID(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange: Create a task first
	ctx := context.Background()
	collection := database.GetCollection()
	collection.DeleteMany(ctx, bson.M{})

	testTask := models.Task{
		ID:          primitive.NewObjectID(),
		Title:       "Find Me",
		Description: "Test finding by ID",
		Completed:   false,
	}
	_, err := collection.InsertOne(ctx, testTask)
	if err != nil {
		t.Fatalf("Failed to insert test task: %v", err)
	}

	// Act: Get the task by ID
	input := &models.GetTaskInput{
		ID: testTask.ID.Hex(),
	}
	output, err := GetTaskByID(ctx, input)

	// Assert
	if err != nil {
		t.Fatalf("GetTaskByID returned error: %v", err)
	}

	if output.Body.Title != "Find Me" {
		t.Errorf("Expected 'Find Me', got '%s'", output.Body.Title)
	}

	if output.Body.ID != testTask.ID {
		t.Error("ID mismatch")
	}

	// Cleanup
	collection.DeleteMany(ctx, bson.M{})
	t.Log("✅ GetTaskByID passed")
}

// ============================================================================
// TEST GETTASKBHID - INVALID ID
// ============================================================================
// TestGetTaskByID tests retrieving a specific task
func TestGetTaskByID_InvalidID(t *testing.T) {
	// No database needed, so no testing.Short() check
	ctx := context.Background()

	// Create input with bad ID
	input := &models.GetTaskInput{
		ID: "invalid-id-format", // This ID will flag as invalid as it is not a valid 24 character hex string
	}

	// Call handler - will try to parse "invalid-id-format"
	_, err := GetTaskByID(ctx, input) // We expect an error here and don't care about any other output

	// Assert: Check that we got an error
	if err == nil {
		// This is not the outcome we want, we want an invalid ID to return an error
		t.Fatalf("Expected error for invalid ID, got nil")
	}

	// If we reach here, err is not nil, which is correct
	t.Log("✅ Invalid ID error handling passed")
}

// ============================================================================
// TEST UPDATE TASK
// ============================================================================
// TestUpdateTask tests updating an existing task
func TestUpdateTask(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange: Create a task first
	ctx := context.Background()
	collection := database.GetCollection()

	// First, ensure database is completely clean
	collection.DeleteMany(ctx, bson.M{})

	testTask := models.Task{
		ID:          primitive.NewObjectID(),
		Title:       "Original Title",
		Description: "Original Description",
		Completed:   false,
	}

	_, err := collection.InsertOne(ctx, testTask)
	if err != nil {
		t.Fatalf("Failed to insert test task: %v", err)
	}

	// Act: Update the task
	title := "Updated Title"
	description := "Updated Description"
	completed := true

	input := &models.UpdateTaskInput{
		ID: testTask.ID.Hex(),
	}

	input.Body.Title = &title
	input.Body.Description = &description
	input.Body.Completed = &completed

	output, err := UpdateTask(ctx, input)

	// Assert
	if err != nil {
		t.Fatalf("UpdateTask returned error: %v", err)
	}

	if output.Body.Title != "Updated Title" {
		t.Errorf("Expected 'Updated Title', got '%s'", output.Body.Title)
	}

	if output.Body.Completed != true {
		t.Error("Expected task to be completed")
	}

	// Cleanup
	collection.DeleteMany(ctx, bson.M{})
	t.Log("✅ UpdateTask passed")
}

// ============================================================================
// TEST DELETETASK
// ============================================================================
// TestDeleteTask tests deleting a task
func TestDeleteTask(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange: Create a task first
	ctx := context.Background()
	collection := database.GetCollection()
	collection.DeleteMany(ctx, bson.M{})

	testTask := models.Task{
		ID:          primitive.NewObjectID(),
		Title:       "Delete Me",
		Description: "This task will be deleted",
		Completed:   false,
	}
	collection.InsertOne(ctx, testTask)

	// Verify task exists
	count, _ := collection.CountDocuments(ctx, bson.M{})
	if count != 1 {
		t.Fatalf("Expected 1 task before delete, got %d", count)
	}

	// Act: Delete the task
	input := &models.DeleteTaskInput{
		ID: testTask.ID.Hex(),
	}
	output, err := DeleteTask(ctx, input)

	// Assert
	if err != nil {
		t.Fatalf("DeleteTask returned error: %v", err)
	}

	if output.Body.Message == "" {
		t.Error("Expected success message")
	}

	// Verify task is deleted
	count, _ = collection.CountDocuments(ctx, bson.M{})
	if count != 0 {
		t.Errorf("Expected 0 tasks after delete, got %d", count)
	}

	t.Log("✅ DeleteTask passed")
}

// ============================================================================
// TEST DELETETASK - Not found
// ============================================================================

// TestDeleteTask_NotFound tests deleting non-existent task
func TestDeleteTask_NotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	ctx := context.Background()
	collection := database.GetCollection()
	collection.DeleteMany(ctx, bson.M{})

	// Try to delete task that doesn't exist
	input := &models.DeleteTaskInput{
		ID: primitive.NewObjectID().Hex(),
	}

	_, err := DeleteTask(ctx, input)

	// Assert: Should return error
	if err == nil {
		t.Error("Expected error for non-existent task, got nil")
	}

	t.Log("✅ DeleteTask not found error handling passed")
}
