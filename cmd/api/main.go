// ============================================================================
// PACKAGE DECLARATION
// ============================================================================
// Package main is special in Go - it tells Go "this is an executable program"
// When you run "go run main.go", Go looks for the "main" package and the "main()" function
package main

// ============================================================================
// IMPORTS
// ============================================================================
// Import statements bring in code from other packages (like "import" in Python or JavaScript)
import (
	// STANDARD LIBRARY PACKAGES (built into Go)
	"fmt"      // fmt = "format" - for printing text to the console (like console.log)
	"log"      // log = for error messages and logging
	"net/http" // net/http = for creating web servers and handling HTTP requests

	// OUR OWN PACKAGES (code we wrote in this project)
	"go-todo-api/internal/database"   // Our database connection code
	"go-todo-api/internal/handlers"   // Our API endpoint handlers (the logic for each route)
	"go-todo-api/internal/logger"     // Our structured logged setup
	"go-todo-api/internal/middleware" // Our middleware (code that runs before handlers)
	"go-todo-api/internal/tracing"    // Our tracing code setup

	// THIRD-PARTY PACKAGES (external libraries we installed)
	"github.com/danielgtaylor/huma/v2"                  // Huma = Modern REST API framework
	"github.com/danielgtaylor/huma/v2/adapters/humachi" // Adapter to use Huma with Chi router
	"github.com/go-chi/chi/v5"                          // Chi = HTTP router (handles URL routing)
)

// ============================================================================
// MAIN FUNCTION
// ============================================================================
// func main() is THE entry point of the program
// When you run your program, Go automatically calls this function first
// Think of it like the "start" button of your application
func main() {

	// ------------------------------------------------------------------------
	// STEP 0: INITIALIZE STRUCTURED LOGGING
	// ------------------------------------------------------------------------
	// Set up JSON structured logging for better observability
	// This creates a global logger that all parts of the app can use
	logger.Init()

	// ------------------------------------------------------------------------
	// STEP 1: CONNECT TO DATABASE
	// ------------------------------------------------------------------------
	// Before we can handle any requests, we need to connect to MongoDB
	// This function is defined in internal/database/mongo.go
	// It reads the MONGO_URI from .env and connects to MongoDB
	database.Connect()
	// After this line, we have an active connection to MongoDB!

	// ------------------------------------------------------------------------
	// STEP 2: INITIALIZE TRACING
	// ------------------------------------------------------------------------
	// Set up OpenTelemetry tracing to track request performance
	// This returns a cleanup function that we'll call when the server shuts down
	shutdown := tracing.Init("todo-api")
	defer shutdown() // Call shutdown when main() exits to flush traces

	// ------------------------------------------------------------------------
	// STEP 3: CREATE HTTP ROUTER
	// ------------------------------------------------------------------------
	// A router decides which function (handler) to call based on the URL
	// For example: GET /tasks → calls GetAllTasks handler
	//              POST /tasks → calls CreateTask handler
	// Chi is a popular, fast router for Go
	router := chi.NewMux() // NewMux() creates a new router (Mux = "HTTP request multiplexer")

	// ------------------------------------------------------------------------
	// STEP 4: ADD MIDDLEWARE
	// ------------------------------------------------------------------------
	// Middleware is code that runs BEFORE your handlers

	// Add tracing middleware - creates spans for every request
	// This shold be first so it measures the full request duration
	router.Use(middleware.TracingChi)

	// Add logging middleware - logs every HTTP request (method, path, time)
	// Example log: "GET /tasks 2.5ms"
	router.Use(middleware.LoggingChi)

	// Add CORS middleware - allows browsers from other domains to access your API
	// CORS = Cross-Origin Resource Sharing
	// Without this, browsers block requests from other websites for security
	router.Use(middleware.CORSChi)

	// Add authentication middleware - requires valid API key for all requests
	// Every request must include header: X-API-Key: your-key-here
	router.Use(middleware.AuthChi)

	// ------------------------------------------------------------------------
	// STEP 5: CREATE HUMA API WITH OPENAPI DOCUMENTATION
	// ------------------------------------------------------------------------
	// Huma is a framework that wraps your router and adds superpowers:
	// - Automatic OpenAPI documentation generation
	// - Automatic request validation
	// - Automatic JSON encoding/decoding
	// - Better error handling

	// Create Huma config with custom context tranformer
	// This ensures OpenTelemetry spac context is passed from HTTP middleware to handlers
	config := huma.DefaultConfig("TODO API", "1.0.0")

	// Create Huma API instance with default configuration
	// "TODO API" = API name, "1.0.0" = version number
	api := humachi.New(router, config)

	// Add metadata to the API documentation
	// This shows up in the /docs page that users can see
	api.OpenAPI().Info.Description = "A production-ready REST API for managing TODO tasks"
	api.OpenAPI().Info.Contact = &huma.Contact{
		Name: "Your Name",
		URL:  "https://github.com/yourusername/go-todo-api",
	}

	// ------------------------------------------------------------------------
	// STEP 6: REGISTER API ENDPOINTS (ROUTES)
	// ------------------------------------------------------------------------
	// Each huma.Register() call tells Huma:
	// "When someone makes a [METHOD] request to [PATH], call this [HANDLER]"
	// Huma automatically generates OpenAPI documentation from these registrations

	// HEALTH CHECK ENDPOINT
	// GET /health → Returns { "status": "healthy", "message": "..." }
	// Used to check if the server is running (monitoring tools use this)
	huma.Register(api, huma.Operation{
		OperationID: "get-health",                                     // Unique ID for this operation (used in docs)
		Method:      http.MethodGet,                                   // HTTP method: GET, POST, PUT, DELETE, etc.
		Path:        "/health",                                        // URL path: http://localhost:8080/health
		Summary:     "Health check",                                   // Short description (shows in docs)
		Description: "Check if the API server is running and healthy", // Long description
		Tags:        []string{"System"},                               // Groups this endpoint under "System" in docs
	}, handlers.Health) // handlers.Health is the function that handles this request

	// GET ALL TASKS ENDPOINT
	// GET /tasks → Returns array of all tasks from database
	huma.Register(api, huma.Operation{
		OperationID: "list-tasks",
		Method:      http.MethodGet,
		Path:        "/tasks",
		Summary:     "List all tasks",
		Description: "Retrieve all TODO tasks from the database",
		Tags:        []string{"Tasks"}, // Groups under "Tasks" section in docs
	}, handlers.GetAllTasks)

	// GET SINGLE TASK BY ID ENDPOINT
	// GET /tasks/6900d436e231fdbb964c3c1c → Returns one specific task
	// {id} in the path means "this is a variable"
	// The ID from the URL is passed to the handler
	huma.Register(api, huma.Operation{
		OperationID: "get-task",
		Method:      http.MethodGet,
		Path:        "/tasks/{id}", // {id} = path parameter (captures value from URL)
		Summary:     "Get a task by ID",
		Description: "Retrieve a specific task using its unique identifier",
		Tags:        []string{"Tasks"},
	}, handlers.GetTaskByID)

	// CREATE NEW TASK ENDPOINT
	// POST /tasks with body: {"title": "Buy milk", "description": "..."}
	// Creates a new task in the database
	huma.Register(api, huma.Operation{
		OperationID:   "create-task",
		Method:        http.MethodPost, // POST = create new resource
		Path:          "/tasks",
		Summary:       "Create a new task",
		Description:   "Add a new TODO task to the database",
		Tags:          []string{"Tasks"},
		DefaultStatus: http.StatusCreated, // Return 201 Created (not 200 OK)
	}, handlers.CreateTask)

	// UPDATE EXISTING TASK ENDPOINT
	// PUT /tasks/6900d436e231fdbb964c3c1c with body: {"completed": true}
	// Updates an existing task's fields
	huma.Register(api, huma.Operation{
		OperationID: "update-task",
		Method:      http.MethodPut, // PUT = update existing resource
		Path:        "/tasks/{id}",
		Summary:     "Update a task",
		Description: "Update an existing task's title, description, or completion status",
		Tags:        []string{"Tasks"},
	}, handlers.UpdateTask)

	// DELETE TASK ENDPOINT
	// DELETE /tasks/6900d436e231fdbb964c3c1c
	// Removes a task from the database permanently
	huma.Register(api, huma.Operation{
		OperationID: "delete-task",
		Method:      http.MethodDelete, // DELETE = remove resource
		Path:        "/tasks/{id}",
		Summary:     "Delete a task",
		Description: "Remove a task from the database",
		Tags:        []string{"Tasks"},
	}, handlers.DeleteTask)

	// ------------------------------------------------------------------------
	// STEP 7: PRINT STARTUP INFORMATION
	// ------------------------------------------------------------------------
	// fmt.Println() prints text to the console (like console.log in JavaScript)
	// This helps developers know the server started successfully
	fmt.Println("🚀 Server starting on http://localhost:8080")
	fmt.Println("✨ Framework: Huma v2 with Chi router")
	fmt.Println("✨ Middleware enabled: Logging, CORS, Authentication")
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

	// ------------------------------------------------------------------------
	// STEP 8: START THE HTTP SERVER
	// ------------------------------------------------------------------------
	// This is the most important line - it actually starts the web server!

	port := ":8080" // Port 8080 = the door number your server listens on
	// :8080 means "listen on all network interfaces on port 8080"

	// http.ListenAndServe() starts the server and BLOCKS FOREVER
	// This means the program doesn't exit - it keeps running, waiting for requests
	// log.Fatal() means "if the server fails to start, print the error and exit"
	log.Fatal(http.ListenAndServe(port, router))

	// The server is now running and handling requests 24/7 until you stop it
}

// ============================================================================
// HOW THIS ALL WORKS TOGETHER
// ============================================================================
//
// 1. Program starts → main() function is called
// 2. Connect to MongoDB database
// 3. Initialize tracing requests
// 4. Create a router (Chi) to handle different URLs
// 5. Add middleware (tracing, logging, CORS) that runs before every request
// 6. Wrap router with Huma for automatic docs and validation
// 7. Register 6 endpoints (health check + 5 CRUD operations)
// 8. Print helpful startup messages
// 9. Start HTTP server on port 8080 (blocks forever, handling requests)
//
// When a request comes in:
// Request → Middleware (logging, CORS) → Router (finds matching handler)
//        → Handler (your code) → Response back to client
//
// Example flow for "GET /tasks":
// 1. Browser sends: GET http://localhost:8080/tasks
// 2. Server receives request
// 3. Tracing middleware creates a span for the request
// 4. Logging middleware logs: "GET /tasks"
// 5. Cors middleware adds Cors headers
// 6. Auth middleware checks API key
// 7. Router sees "/tasks" with GET method
// 8. Router calls handlers.GetAllTasks()
// 9. Handler queries MongoDB for all tasks
// 10. Huma converts tasks to JSON
// 11. Response sent back: [{"id": "...", "title": "..."}]
// 12. Logging middleware logs: "GET /tasks 5ms"
//
// ============================================================================
