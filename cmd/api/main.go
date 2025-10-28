// Package main is the entry point for the TODO API server
package main

import (
	"fmt"
	"log"
	"net/http"

	"go-todo-api/internal/database"
	"go-todo-api/internal/handlers"
	"go-todo-api/internal/middleware"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
)

func main() {
	// Initialize database connection
	database.Connect()

	// Create a new Chi router
	router := chi.NewMux()

	// Apply custom middleware
	router.Use(middleware.LoggingChi)
	router.Use(middleware.CORSChi)

	// Create Huma API with OpenAPI documentation
	api := humachi.New(router, huma.DefaultConfig("TODO API", "1.0.0"))

	// Configure API metadata
	api.OpenAPI().Info.Description = "A production-ready REST API for managing TODO tasks"
	api.OpenAPI().Info.Contact = &huma.Contact{
		Name: "Your Name",
		URL:  "https://github.com/yourusername/go-todo-api",
	}

	// Register operations with OpenAPI documentation
	huma.Register(api, huma.Operation{
		OperationID: "get-health",
		Method:      http.MethodGet,
		Path:        "/health",
		Summary:     "Health check",
		Description: "Check if the API server is running and healthy",
		Tags:        []string{"System"},
	}, handlers.Health)

	huma.Register(api, huma.Operation{
		OperationID: "list-tasks",
		Method:      http.MethodGet,
		Path:        "/tasks",
		Summary:     "List all tasks",
		Description: "Retrieve all TODO tasks from the database",
		Tags:        []string{"Tasks"},
	}, handlers.GetAllTasks)

	huma.Register(api, huma.Operation{
		OperationID: "get-task",
		Method:      http.MethodGet,
		Path:        "/tasks/{id}",
		Summary:     "Get a task by ID",
		Description: "Retrieve a specific task using its unique identifier",
		Tags:        []string{"Tasks"},
	}, handlers.GetTaskByID)

	huma.Register(api, huma.Operation{
		OperationID:   "create-task",
		Method:        http.MethodPost,
		Path:          "/tasks",
		Summary:       "Create a new task",
		Description:   "Add a new TODO task to the database",
		Tags:          []string{"Tasks"},
		DefaultStatus: http.StatusCreated,
	}, handlers.CreateTask)

	huma.Register(api, huma.Operation{
		OperationID: "update-task",
		Method:      http.MethodPut,
		Path:        "/tasks/{id}",
		Summary:     "Update a task",
		Description: "Update an existing task's title, description, or completion status",
		Tags:        []string{"Tasks"},
	}, handlers.UpdateTask)

	huma.Register(api, huma.Operation{
		OperationID: "delete-task",
		Method:      http.MethodDelete,
		Path:        "/tasks/{id}",
		Summary:     "Delete a task",
		Description: "Remove a task from the database",
		Tags:        []string{"Tasks"},
	}, handlers.DeleteTask)

	// Print startup information
	fmt.Println("🚀 Server starting on http://localhost:8080")
	fmt.Println("✨ Framework: Huma v2 with Chi router")
	fmt.Println("✨ Middleware enabled: Logging, CORS")
	fmt.Println("📁 Production structure: cmd/ and internal/ packages")
	fmt.Println("📚 OpenAPI Documentation available at:")
	fmt.Println("  - http://localhost:8080/docs (Interactive API docs)")
	fmt.Println("  - http://localhost:8080/openapi.json (OpenAPI spec)")
	fmt.Println("  - http://localhost:8080/openapi.yaml (OpenAPI spec)")
	fmt.Println("\n🎯 Try these endpoints:")
	fmt.Println("  - GET    /health")
	fmt.Println("  - GET    /tasks")
	fmt.Println("  - POST   /tasks")
	fmt.Println("  - GET    /tasks/{id}")
	fmt.Println("  - PUT    /tasks/{id}")
	fmt.Println("  - DELETE /tasks/{id}")

	// Start the HTTP server
	port := ":8080"
	log.Fatal(http.ListenAndServe(port, router))
}
