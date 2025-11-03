// ============================================================================
// PACKAGE DECLARATION
// ============================================================================
// Package handlers contains all the HTTP handler functions for our API
// Handlers are the "business logic" - they process requests and return responses
// Think of handlers as the "controllers" in MVC pattern
package handlers

// ============================================================================
// IMPORTS
// ============================================================================
import (
	// STANDARD LIBRARY PACKAGES
	"context" // context = for managing request timeouts and cancellation
	"fmt"     // fmt = for printing formatted output to console
	"time"    // time = for working with time durations and timeouts

	// OUR OWN PACKAGES
	"go-todo-api/internal/database" // Our database connection code
	"go-todo-api/internal/models"   // Our data structures (Task, Input/Output types)

	// THIRD-PARTY PACKAGES
	"github.com/danielgtaylor/huma/v2"           // Huma = REST API framework with error helpers
	"go.mongodb.org/mongo-driver/bson"           // bson = MongoDB's query language (like SQL)
	"go.mongodb.org/mongo-driver/bson/primitive" // primitive = MongoDB types (ObjectID)
	"go.mongodb.org/mongo-driver/mongo"          // mongo = MongoDB driver for Go

	// OPEN TELEMETRY SPAN PACKAGES
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

// ============================================================================
// GET ALL TASKS - LIST OPERATION (WITH FILTERING)
// ============================================================================
// GetAllTasks retrieves all tasks from the database, with optional filtering
// This is called when someone makes a GET request to /tasks
//
// Huma Handler Signature:
// - Input: context.Context (for timeouts) + *models.GetTasksInput (query parameters)
// - Output: *models.GetTasksOutput (contains array of tasks) + error
//
// Example requests:
// GET /tasks                    → Returns all tasks
// GET /tasks?completed=true     → Returns only completed tasks
// GET /tasks?completed=false    → Returns only incomplete tasks
func GetAllTasks(ctx context.Context, input *models.GetTasksInput) (*models.GetTasksOutput, error) {
	// ----------------------------------------------------------------------------
	// STEP 1: CREATE A TRACER
	// ----------------------------------------------------------------------------
	// Get a tracer object named 'handlers'.
	// The name 'handlers' groups all spans from this file together in Jaeger
	tracer := otel.Tracer("handlers")

	// ----------------------------------------------------------------------------
	// STEP 2: CREATE THE MAIN HANDLER SPAN
	// ----------------------------------------------------------------------------
	// Create a new span called 'GetAllTasks'. This is the name that will show in Jaeger.
	// Returns the updated context with the span attached and handlerSpan, the span object.
	// Defer stops the span whne the function exits
	ctx, handlerSpan := tracer.Start(ctx, "GetAllTasks")
	defer handlerSpan.End()

	// ----------------------------------------------------------------------------
	// STEP 3: BUILD FILTER AND ADD ATTRIBUTES
	// ----------------------------------------------------------------------------
	// Build the MongoDB filter
	// SetAttributes adds metadata to the span
	filter := bson.M{}
	switch input.Completed {
	case "true":
		filter["completed"] = true
		handlerSpan.SetAttributes(attribute.String("filter.completed", input.Completed))
	case "false":
		filter["completed"] = false
		handlerSpan.SetAttributes(attribute.String("filter.completed", input.Completed))
	}

	// ----------------------------------------------------------------------------
	// STEP 4: CREATE DATABASE SPAN
	// ----------------------------------------------------------------------------
	// Create a child span for the database query
	collection := database.GetCollection()
	ctx, dbSpan := tracer.Start(ctx, "MongoDB.Find")
	dbSpan.SetAttributes( // adds 3 tags: the database type, which collection and what operation
		attribute.String("db.system", "mongodb"),
		attribute.String("db.collection", "tasks"),
		attribute.String("db.operation", "find"),
	)

	// Create database timeout context from the SPAN context
	dbCtx, cancel := context.WithTimeout(ctx, 5*time.Second) // Use span's ctx
	defer cancel()

	// ----------------------------------------------------------------------------
	// STEP 5: EXECUTE QUERY AND END SPAN
	// ----------------------------------------------------------------------------
	cursor, err := collection.Find(dbCtx, filter)
	dbSpan.End() // Stop the database timer immediately. We manually end it here (not defer) because we want precise timing.

	// ----------------------------------------------------------------------------
	// STEP 6: RECORD ERRORS
	// ----------------------------------------------------------------------------
	// If there's an error, RecordError() marks the span as failed.
	// The span will show up red in Jaeger and an error message is attached to the span.
	if err != nil {
		handlerSpan.RecordError(err) // Record error on span
		return nil, huma.Error500InternalServerError("Failed to fetch tasks from the database")
	}
	defer cursor.Close(dbCtx)

	// Decode results
	var tasks []models.Task
	if err = cursor.All(dbCtx, &tasks); err != nil {
		handlerSpan.RecordError(err)
		return nil, huma.Error500InternalServerError("Failed to decode tasks")
	}

	if tasks == nil {
		tasks = []models.Task{}
	}

	// ----------------------------------------------------------------------------
	// STEP 7: ADD RESULT METRICS
	// ----------------------------------------------------------------------------
	// Add result count to span
	handlerSpan.SetAttributes(attribute.Int("result.count", len(tasks)))

	if input.Completed != "" {
		fmt.Printf("✅ Retrieved %d tasks from MongoDB (filtered by completed=%s)\n", len(tasks), input.Completed)
	} else {
		fmt.Printf("✅ Retrieved %d tasks from MongoDB (no filter)\n", len(tasks))
	}

	return &models.GetTasksOutput{Body: tasks}, nil
}

// ============================================================================
// GET TASK BY ID - SPECIFIC TASK FILTERING
// ============================================================================

func GetTaskByID(ctx context.Context, input *models.GetTaskInput) (*models.GetTaskOutput, error) {
	// ----------------------------------------------------------------------------
	// STEP 1: CONVERT STRING ID TO MONGODB OBJECTID
	// ----------------------------------------------------------------------------
	// The ID comes from the URL as a string like "6900d436e231fdbb964c3c1c"
	// MongoDB needs this as an ObjectID type, so we convert it
	// primitive.ObjectIDFromHex() converts 24-character hex string → ObjectID
	objectID, err := primitive.ObjectIDFromHex(input.ID)
	if err != nil {
		// If the ID is not a valid 24-character hex string, return HTTP 400 error
		// Example invalid IDs: "123", "abc", "6900d436e231fdbb964c3c1" (too short)
		return nil, huma.Error400BadRequest("Invalid task ID format")
	}

	// ----------------------------------------------------------------------------
	// STEP 2: CREATE DATABASE CONTEXT WITH TIMEOUT
	// ----------------------------------------------------------------------------
	dbCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// ----------------------------------------------------------------------------
	// STEP 3: QUERY DATABASE FOR THE SPECIFIC TASK
	// ----------------------------------------------------------------------------
	// Create a variable to hold the result
	var task models.Task

	// Get the collection and find one document that matches the ID
	collection := database.GetCollection()
	// bson.M{"_id": objectID} = filter that matches documents where _id field equals objectID
	// This is like: SELECT * FROM tasks WHERE _id = objectID (in SQL)
	// .Decode(&task) = put the result into our task variable
	err = collection.FindOne(dbCtx, bson.M{"_id": objectID}).Decode(&task)

	// ----------------------------------------------------------------------------
	// STEP 4: HANDLE ERRORS
	// ----------------------------------------------------------------------------
	if err != nil {
		// Check if the error is "no documents found"
		if err == mongo.ErrNoDocuments {
			// Task with this ID doesn't exist → return HTTP 404 error
			return nil, huma.Error404NotFound("Task not found")
		}
		// Any other error (database connection issue, etc.) → HTTP 500 error
		return nil, huma.Error500InternalServerError("Failed to fetch task")
	}

	// ----------------------------------------------------------------------------
	// STEP 5: LOG SUCCESS AND RETURN RESULT
	// ----------------------------------------------------------------------------
	// .Hex() converts ObjectID back to string for logging
	fmt.Printf("✅ Retrieved task with ID %s\n", objectID.Hex())

	// Return the output struct with the task we found
	return &models.GetTaskOutput{Body: task}, nil
}

// ============================================================================
// CREATE TASK - CREATE OPERATION
// ============================================================================
// CreateTask creates a new task in the database
// This is called when someone makes a POST request to /tasks
//
// Huma Handler Signature:
// - Input: context.Context + *models.CreateTaskInput (contains title & description)
// - Output: *models.CreateTaskOutput (contains newly created task with ID) + error
//
// Example request:  POST /tasks with body: {"title": "Buy milk", "description": "From the store"}
// Example response: {"id": "6900d436e231fdbb964c3c1c", "title": "Buy milk", "description": "From the store", "completed": false}
func CreateTask(ctx context.Context, input *models.CreateTaskInput) (*models.CreateTaskOutput, error) {
	// Create tracer and handler span
	tracer := otel.Tracer("handlers")
	ctx, handlerSpan := tracer.Start(ctx, "CreateTask")
	defer handlerSpan.End()

	// ----------------------------------------------------------------------------
	// STEP 1: CREATE NEW TASK STRUCT FROM INPUT
	// ----------------------------------------------------------------------------
	// Take the data from the request body and create a Task struct
	// Note: We're NOT setting the ID here - MongoDB will generate it for us
	// Note: Completed defaults to false for new tasks
	newTask := models.Task{
		Title:       input.Body.Title,       // From request body
		Description: input.Body.Description, // From request body (can be empty)
		Completed:   false,                  // Always starts as not completed
	}

	// Add task attributes to span
	handlerSpan.SetAttributes(
		attribute.String("task.title", input.Body.Title),
		attribute.Bool("task.completed", false),
	)

	// ----------------------------------------------------------------------------
	// STEP 2: CREATE DATABASE CONTEXT WITH TIMEOUT
	// ----------------------------------------------------------------------------
	dbCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// ----------------------------------------------------------------------------
	// STEP 3: INSERT THE NEW TASK INTO MONGODB
	// ----------------------------------------------------------------------------
	// Create database span
	_, dbSpan := tracer.Start(ctx, "MongoDB.InsertOne")
	dbSpan.SetAttributes(
		attribute.String("db.system", "mongodb"),
		attribute.String("db.collection", "tasks"),
	)

	collection := database.GetCollection()
	// InsertOne() adds the newTask to the database
	// It returns:
	//   - result.InsertedID = the auto-generated MongoDB ID for this document
	//   - err = any error that occurred during insertion
	result, err := collection.InsertOne(dbCtx, newTask)

	// Error recorded and will be visible in Jaeger
	if err != nil {
		handlerSpan.RecordError(err)
		dbSpan.End()
		// If insertion fails (database down, disk full, etc.) → HTTP 500 error
		return nil, huma.Error500InternalServerError("Failed to create task in database")
	}
	// End the span once the task has been added to the db
	dbSpan.End()

	// ----------------------------------------------------------------------------
	// STEP 4: SET THE AUTO-GENERATED ID ON OUR TASK
	// ----------------------------------------------------------------------------
	// MongoDB generated an ID and put it in result.InsertedID
	// result.InsertedID is type interface{}, so we need to convert it
	// .(primitive.ObjectID) = type assertion (like casting in other languages)
	// This says: "I know this is an ObjectID, treat it as one"
	newTask.ID = result.InsertedID.(primitive.ObjectID)

	// Record the generated ID in the span
	handlerSpan.SetAttributes(attribute.String("task.id", newTask.ID.Hex()))

	// ----------------------------------------------------------------------------
	// STEP 5: LOG SUCCESS AND RETURN THE NEW TASK
	// ----------------------------------------------------------------------------
	// Print success message with the task title and new ID
	fmt.Printf("✅ Created new task: %s with ID %s\n", newTask.Title, newTask.ID.Hex())

	// Return the complete task (now with its ID) to the client
	// HTTP status will be 201 Created (set in main.go with DefaultStatus)
	return &models.CreateTaskOutput{Body: newTask}, nil
}

// ============================================================================
// UPDATE TASK - UPDATE OPERATION
// ============================================================================
// UpdateTask updates an existing task in the database
// This is called when someone makes a PUT request to /tasks/{id}
//
// Huma Handler Signature:
// - Input: context.Context + *models.UpdateTaskInput (contains ID + optional fields)
// - Output: *models.UpdateTaskOutput (contains updated task) + error
//
// Example request:  PUT /tasks/6900d436e231fdbb964c3c1c with body: {"completed": true}
// Example response: {"id": "6900d436e231fdbb964c3c1c", "title": "Buy milk", "completed": true}
//
// IMPORTANT: This is a PARTIAL update (also called PATCH-like behavior)
// - Client only sends fields they want to change
// - Fields not sent remain unchanged
// - We use pointers (*string, *bool) to distinguish "not sent" from "sent but empty"
func UpdateTask(ctx context.Context, input *models.UpdateTaskInput) (*models.UpdateTaskOutput, error) {
	// Create tracer and handler span
	tracer := otel.Tracer("handlers")
	ctx, handlerSpan := tracer.Start(ctx, "UpdateTask")
	defer handlerSpan.End()

	// Add task ID to span attributes
	handlerSpan.SetAttributes(attribute.String("task.id", input.ID))

	// ----------------------------------------------------------------------------
	// STEP 1: CONVERT STRING ID TO MONGODB OBJECTID
	// ----------------------------------------------------------------------------
	objectID, err := primitive.ObjectIDFromHex(input.ID)
	if err != nil {
		return nil, huma.Error400BadRequest("Invalid task ID format")
	}

	// ----------------------------------------------------------------------------
	// STEP 2: CREATE DATABASE CONTEXT WITH TIMEOUT
	// ----------------------------------------------------------------------------
	dbCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	collection := database.GetCollection()

	// ----------------------------------------------------------------------------
	// STEP 3: CHECK IF TASK EXISTS (OPTIONAL BUT GOOD PRACTICE)
	// ----------------------------------------------------------------------------
	// Create span for FindOne operation
	_, findSpan := tracer.Start(ctx, "MongoDB.FindONe")
	findSpan.SetAttributes(
		attribute.String("db.system", "mongodb"),
		attribute.String("db.collection", "tasks"),
		attribute.String("db.operation", "findOne"),
	)

	// Find the existing task first to verify it exists
	// This gives us a better error message if the task doesn't exist
	var existingTask models.Task
	err = collection.FindOne(dbCtx, bson.M{"_id": objectID}).Decode(&existingTask)
	if err != nil {
		findSpan.End()
		handlerSpan.RecordError(err)
		if err == mongo.ErrNoDocuments {
			return nil, huma.Error404NotFound("Task not found")
		}
		return nil, huma.Error500InternalServerError("Failed to fetch task")
	}

	findSpan.End()

	// ----------------------------------------------------------------------------
	// STEP 4: BUILD UPDATE DOCUMENT WITH ONLY PROVIDED FIELDS
	// ----------------------------------------------------------------------------
	// MongoDB update format: { "$set": { "field1": "value1", "field2": "value2" } }
	// $set = MongoDB operator that updates specific fields without replacing entire document
	update := bson.M{"$set": bson.M{}} // Create empty update document

	// Check each field to see if it was provided in the request
	// Remember: input.Body.Title is a *string (pointer)
	// If pointer is nil, field was not sent in request
	// If pointer is not nil, field was sent (even if empty string)

	if input.Body.Title != nil {
		// *input.Body.Title = dereference the pointer to get actual string value
		update["$set"].(bson.M)["title"] = *input.Body.Title
	}
	if input.Body.Description != nil {
		update["$set"].(bson.M)["description"] = *input.Body.Description
	}
	if input.Body.Completed != nil {
		// *input.Body.Completed = dereference the pointer to get actual bool value
		update["$set"].(bson.M)["completed"] = *input.Body.Completed
	}

	// ----------------------------------------------------------------------------
	// STEP 5: VALIDATE THAT AT LEAST ONE FIELD WAS PROVIDED
	// ----------------------------------------------------------------------------
	// If client sent empty body {}, there's nothing to update
	if len(update["$set"].(bson.M)) == 0 {
		return nil, huma.Error400BadRequest("No fields to update")
	}

	// ----------------------------------------------------------------------------
	// STEP 6: PERFORM THE UPDATE IN MONGODB
	// ----------------------------------------------------------------------------
	// Create span for UpdateOne operation
	_, updateSpan := tracer.Start(ctx, "MongoDB.UpdateOne")
	updateSpan.SetAttributes(
		attribute.String("db.system", "mongodb"),
		attribute.String("db.collection", "tasks"),
		attribute.String("db.operation", "updateOne"),
	)

	// UpdateOne(filter, update) updates the first document matching the filter
	// Returns result with MatchedCount (how many docs matched) and ModifiedCount
	result, err := collection.UpdateOne(dbCtx, bson.M{"_id": objectID}, update)
	if err != nil {
		updateSpan.End()
		handlerSpan.RecordError(err)
		return nil, huma.Error500InternalServerError("Failed to update task")
	}
	updateSpan.End()

	// Add modified count to span
	handlerSpan.SetAttributes(attribute.Int64("result.modifiedCount", result.ModifiedCount))

	// Double-check that a document was actually matched (should always be 1)
	if result.MatchedCount == 0 {
		return nil, huma.Error404NotFound("Task not found")
	}

	// ----------------------------------------------------------------------------
	// STEP 7: FETCH THE UPDATED TASK TO RETURN IT
	// ----------------------------------------------------------------------------
	// After updating, get the latest version of the task from database
	// This ensures we return the complete, up-to-date task to the client
	var updatedTask models.Task
	collection.FindOne(dbCtx, bson.M{"_id": objectID}).Decode(&updatedTask)

	// ----------------------------------------------------------------------------
	// STEP 8: LOG SUCCESS AND RETURN UPDATED TASK
	// ----------------------------------------------------------------------------
	fmt.Printf("✅ Updated task with ID %s\n", objectID.Hex())
	return &models.UpdateTaskOutput{Body: updatedTask}, nil
}

// ============================================================================
// DELETE TASK - DELETE OPERATION
// ============================================================================
// DeleteTask deletes a task from the database
// This is called when someone makes a DELETE request to /tasks/{id}
//
// Huma Handler Signature:
// - Input: context.Context + *models.DeleteTaskInput (contains ID from URL path)
// - Output: *models.DeleteTaskOutput (contains success message) + error
//
// Example request:  DELETE /tasks/6900d436e231fdbb964c3c1c
// Example response: {"message": "Task deleted successfully", "id": "6900d436e231fdbb964c3c1c"}
func DeleteTask(ctx context.Context, input *models.DeleteTaskInput) (*models.DeleteTaskOutput, error) {
	// ----------------------------------------------------------------------------
	// STEP 1: CONVERT STRING ID TO MONGODB OBJECTID
	// ----------------------------------------------------------------------------
	objectID, err := primitive.ObjectIDFromHex(input.ID)
	if err != nil {
		return nil, huma.Error400BadRequest("Invalid task ID format")
	}

	// ----------------------------------------------------------------------------
	// STEP 2: CREATE DATABASE CONTEXT WITH TIMEOUT
	// ----------------------------------------------------------------------------
	dbCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// ----------------------------------------------------------------------------
	// STEP 3: DELETE THE TASK FROM MONGODB
	// ----------------------------------------------------------------------------
	collection := database.GetCollection()
	// DeleteOne(filter) removes the first document that matches the filter
	// Returns result with DeletedCount (how many documents were deleted)
	// Should be either 0 (not found) or 1 (successfully deleted)
	result, err := collection.DeleteOne(dbCtx, bson.M{"_id": objectID})
	if err != nil {
		// Database error during deletion → HTTP 500 error
		return nil, huma.Error500InternalServerError("Failed to delete task")
	}

	// ----------------------------------------------------------------------------
	// STEP 4: CHECK IF TASK WAS ACTUALLY DELETED
	// ----------------------------------------------------------------------------
	// If DeletedCount is 0, no document with that ID existed
	if result.DeletedCount == 0 {
		return nil, huma.Error404NotFound("Task not found")
	}

	// ----------------------------------------------------------------------------
	// STEP 5: LOG SUCCESS AND RETURN CONFIRMATION
	// ----------------------------------------------------------------------------
	fmt.Printf("✅ Deleted task with ID %s\n", objectID.Hex())

	// Return a success message with the deleted task's ID
	// This uses an anonymous struct (defined inline without a type name)
	// The struct is defined in models.DeleteTaskOutput, but we create it here
	return &models.DeleteTaskOutput{
		Body: struct {
			Message string `json:"message" doc:"Success message"`
			ID      string `json:"id" doc:"Deleted task ID"`
		}{
			Message: "Task deleted successfully", // Success message
			ID:      input.ID,                    // Echo back the ID that was deleted
		},
	}, nil
}

// ============================================================================
// HOW THESE HANDLERS WORK WITH HUMA
// ============================================================================
//
// Each handler follows the same pattern:
// 1. Validate input (convert IDs, check formats)
// 2. Create database context with timeout (prevents hanging)
// 3. Perform database operation (Find, Insert, Update, Delete)
// 4. Handle errors (404, 400, 500)
// 5. Log success and return result
//
// Huma automatically:
// - Validates request against struct tags (minLength, maxLength)
// - Converts JSON request body to Input structs
// - Converts Output structs to JSON response
// - Handles error responses (uses RFC 7807 Problem Details)
// - Generates OpenAPI documentation
//
// Error codes used:
// - 400 Bad Request: Invalid input (bad ID format, validation failed)
// - 404 Not Found: Task doesn't exist
// - 500 Internal Server Error: Database or server error
//
// ============================================================================
