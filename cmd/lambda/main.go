// ============================================================================
// LAMBDA ENTRY POINT
// ============================================================================
// This file is the entry point for AWS Lambda deployment
// It wraps the HTTP server to work with API Gateway events

package main

import (
	"context"
	"net/http"
	"os"

	// AWS Lambda libraries
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"

	// Huma framework
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"

	// Our packages
	"go-todo-api/internal/database"
	"go-todo-api/internal/handlers"
	"go-todo-api/internal/logger"
	"go-todo-api/internal/middleware"
	"go-todo-api/internal/tracing"
)

var (
	// httpHandler is initialized once and reused across Lambda invocations
	httpHandler http.Handler
)

// init runs once when Lambda container starts (cold start)
// This is where we do expensive initialization
func init() {
	// Initialize logger
	logger.Init()
	logger.Log.Info("Lambda: Initializing...")

	// Connect to MongoDB (reused across invocations)
	database.Connect()
	logger.Log.Info("Lambda: Connected to MongoDB")

	// Initialize OpenTelemetry tracing
	shutdown := tracing.Init(tracing.ServiceName)
	defer shutdown()

	// Set up HTTP router (same as regular server)
	router := chi.NewRouter()

	// Add middleware
	router.Use(middleware.TracingChi)
	router.Use(middleware.LoggingChi)
	router.Use(middleware.RateLimitChi)
	router.Use(middleware.SecurityHeadersChi)
	router.Use(middleware.CORSChi)

	// Create Huma API
	config := huma.DefaultConfig("Go TODO API", "1.0.0")
	config.Servers = []*huma.Server{
		{URL: os.Getenv("API_BASE_URL")},
	}
	api := humachi.New(router, config)

	// Register all endpoints
	registerEndpoints(api)

	// Store the handler for reuse
	httpHandler = router

	logger.Log.Info("Lambda: Initialization complete")
}

// registerEndpoints registers all API endpoints
func registerEndpoints(api huma.API) {
	// Health check
	huma.Register(api, huma.Operation{
		OperationID: "health-check",
		Method:      "GET",
		Path:        "/health",
		Summary:     "Health check",
		Description: "Check if the API is running",
		Tags:        []string{"Health"},
	}, handlers.Health)

	// Get all tasks
	huma.Register(api, huma.Operation{
		OperationID: "get-all-tasks",
		Method:      "GET",
		Path:        "/tasks",
		Summary:     "Get all tasks",
		Description: "Retrieve all tasks with optional filtering",
		Tags:        []string{"Tasks"},
	}, handlers.GetAllTasks)

	// Create task
	huma.Register(api, huma.Operation{
		OperationID: "create-task",
		Method:      "POST",
		Path:        "/tasks",
		Summary:     "Create a new task",
		Tags:        []string{"Tasks"},
	}, handlers.CreateTask)

	// Get task by ID
	huma.Register(api, huma.Operation{
		OperationID: "get-task-by-id",
		Method:      "GET",
		Path:        "/tasks/{id}",
		Summary:     "Get a task by ID",
		Tags:        []string{"Tasks"},
	}, handlers.GetTaskByID)

	// Update task
	huma.Register(api, huma.Operation{
		OperationID: "update-task",
		Method:      "PUT",
		Path:        "/tasks/{id}",
		Summary:     "Update a task",
		Tags:        []string{"Tasks"},
	}, handlers.UpdateTask)

	// Delete task
	huma.Register(api, huma.Operation{
		OperationID: "delete-task",
		Method:      "DELETE",
		Path:        "/tasks/{id}",
		Summary:     "Delete a task",
		Tags:        []string{"Tasks"},
	}, handlers.DeleteTask)
}

// handler is called for each Lambda invocation
// It reuses the httpHandler initialized in init()
func handler(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	return httpadapter.NewV2(httpHandler).ProxyWithContext(ctx, req)
}

func main() {
	// Start Lambda runtime
	lambda.Start(handler)
}
