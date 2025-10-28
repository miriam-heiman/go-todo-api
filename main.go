// package main - This declares that this file belongs to the "main" package
// In Go, "main" is special - it tells Go this is a program you can run, not a library
// Think of it like "this is the entry point" of your application
package main

// import - This is how you bring in code from other packages (like importing in JavaScript or Python)
// The parentheses create a "block" where you list all the imports
import (
	"context"       // context = context package - used for cancellation and timeouts in Go
	"encoding/json" // encoding/json = JSON package - for converting Go data structures to/from JSON
	"fmt"           // fmt = format package - for printing text and formatting strings
	"log"           // log = logging package - for writing messages to the console
	"net/http"      // net/http = HTTP package - for creating web servers (like Express.js)
	"os"            // os = operating system package - for accessing environment variables and system functions
	"time"          // time = time package - for handling time-related operations

	// Third-party packages
	"github.com/joho/godotenv" // godotenv = package for loading .env files (like dotenv in Node.js)

	// MongoDB packages - from the official MongoDB Go driver
	"go.mongodb.org/mongo-driver/bson"           // bson = Binary JSON - MongoDB's data format
	"go.mongodb.org/mongo-driver/bson/primitive" // primitive = MongoDB primitive types like ObjectID
	"go.mongodb.org/mongo-driver/mongo"          // mongo = Main MongoDB driver package
	"go.mongodb.org/mongo-driver/mongo/options"  // options = Configuration options for MongoDB operations
)

// type Task struct - This creates a custom data type called "Task"
// "struct" is Go's way of grouping related data together (like a class without methods, or a JavaScript object with a fixed shape)
// Think of it as defining a template: "Every Task must have these exact fields"
type Task struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"` // MongoDB's ObjectID type for unique IDs, bson tag tells MongoDB how to store it
	Title       string             `json:"title"`                   // Title field of type string (text), json tag says use "title" in JSON
	Description string             `json:"description"`             // Description field of type string, json tag says use "description" in JSON
	Completed   bool               `json:"completed"`               // Completed field of type bool (boolean: true or false), json tag says use "completed" in JSON
}

// MongoDB connection variables - these will hold our database connection
var (
	client     *mongo.Client     // client = MongoDB client connection
	collection *mongo.Collection // collection = reference to our tasks collection in the database
	ctx        context.Context   // ctx = context for database operations (handles timeouts and cancellation)
)

// ========================================
// MIDDLEWARE FUNCTIONS
// ========================================
// Middleware is code that runs BEFORE your handlers (like filters or interceptors)
// They can modify requests/responses, log information, add headers, check authentication, etc.
// In Go, middleware is a function that takes a handler and returns a new handler

// loggingMiddleware - Logs information about each HTTP request
// This helps with debugging and monitoring by showing: method, path, and how long the request took
func loggingMiddleware(next http.Handler) http.Handler {
	// http.HandlerFunc converts a function into an http.Handler
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Record the start time of the request
		start := time.Now()

		// Call the next handler in the chain (the actual route handler)
		next.ServeHTTP(w, r)

		// After the handler finishes, log the request details
		// %s = string, %v = any value
		log.Printf(
			"%s %s %s",
			r.Method,           // HTTP method (GET, POST, PUT, DELETE)
			r.URL.Path,         // The URL path that was requested
			time.Since(start),  // How long the request took to process
		)
	})
}

// corsMiddleware - Enables Cross-Origin Resource Sharing (CORS)
// CORS allows your API to be accessed from web browsers on different domains
// Without CORS, browsers block requests from other websites for security
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers to allow requests from any origin
		// Access-Control-Allow-Origin: which domains can access this API (* = all domains)
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// Access-Control-Allow-Methods: which HTTP methods are allowed
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		// Access-Control-Allow-Headers: which headers the client can send
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight requests (OPTIONS method)
		// Browsers send OPTIONS requests before actual requests to check if CORS is allowed
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

// chainMiddleware - Helper function to apply multiple middleware in order
// This makes it easy to wrap a handler with multiple middleware functions
// Usage: chainMiddleware(handler, middleware1, middleware2, middleware3)
func chainMiddleware(h http.Handler, middleware ...func(http.Handler) http.Handler) http.Handler {
	// Apply middleware in reverse order so they execute in the order provided
	// If you pass [logging, cors], it will execute: logging -> cors -> handler
	for i := len(middleware) - 1; i >= 0; i-- {
		h = middleware[i](h)
	}
	return h
}

// init() - This special function runs AUTOMATICALLY before main() - it's for setup/initialization
// In Go, init() functions are called automatically when the package is loaded
func init() {
	// Load environment variables from .env file
	// godotenv.Load() reads the .env file and sets environment variables
	// If the .env file doesn't exist, it will silently fail (optional for production)
	if err := godotenv.Load(); err != nil {
		// If .env file doesn't exist, that's okay - environment variables might be set another way
		log.Println("No .env file found - using environment variables or defaults")
	}

	// Create a context with a 10-second timeout
	// context.Background() creates the root context, context.WithTimeout adds a timeout
	// This ensures database operations don't hang forever if something goes wrong
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel() // defer ensures this cleanup function runs when init() exits

	// Get MongoDB connection string from environment variable
	// os.Getenv() reads environment variables from .env file or system environment
	// This keeps your connection string secret and out of version control
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		// If MONGO_URI is not set in .env or environment variables, exit with error
		log.Fatal("MONGO_URI not found. Please set it in your .env file. See .env.example for template.")
	}

	// Create a new MongoDB client
	// clientOptions.New() creates a configuration object for the MongoDB client
	clientOptions := options.Client().ApplyURI(mongoURI)

	// Connect to MongoDB
	// mongo.Connect() establishes connection to the MongoDB server
	var err error
	client, err = mongo.Connect(ctx, clientOptions)

	// Check if connection failed
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err) // log.Fatal prints error and exits program
	}

	// Ping the MongoDB server to verify connection
	// This ensures the connection is actually working
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Failed to ping MongoDB:", err)
	}

	// Select the database and collection
	// Database name: "todoapi", Collection name: "tasks"
	collection = client.Database("todoapi").Collection("tasks")

	fmt.Println("✅ Connected to MongoDB!")
}

// func main() - This is THE main function that runs when you start your program
// It's like the entry point - Go automatically calls this function when you run your program
// The empty () means it takes no parameters
func main() {
	// Create a new ServeMux (HTTP request multiplexer/router)
	// ServeMux is more flexible than the default mux and allows us to use middleware
	mux := http.NewServeMux()

	// Register routes with their handler functions
	// http.HandlerFunc() converts our handler functions into http.Handler type
	mux.Handle("/", http.HandlerFunc(homeHandler))
	mux.Handle("/tasks", http.HandlerFunc(tasksHandler))
	mux.Handle("/health", http.HandlerFunc(healthHandler))

	// Wrap the entire mux with middleware
	// The order matters: requests flow through logging -> cors -> handlers
	// Responses flow back: handlers -> cors -> logging
	handler := chainMiddleware(
		mux,
		loggingMiddleware, // First: log all requests
		corsMiddleware,    // Second: add CORS headers
	)

	// fmt.Println() - This prints text to the console (like console.log() in JavaScript)
	// Each line prints a different message to help you know the server started
	fmt.Println("🚀 Server starting on http://localhost:8080")
	fmt.Println("✨ Middleware enabled: Logging, CORS")
	fmt.Println("Try visiting:")
	fmt.Println("  - http://localhost:8080/ (homepage)")
	fmt.Println("  - http://localhost:8080/tasks (get all tasks)")
	fmt.Println("  - http://localhost:8080/health (health check)")

	// Start the HTTP server with our middleware-wrapped handler
	// Instead of passing nil, we pass our custom handler that includes middleware
	log.Fatal(http.ListenAndServe(":8080", handler))
}

// func homeHandler(...) - This is a function that handles requests to the homepage
// It takes two parameters:
//   - w http.ResponseWriter - This is what we use to write data back to the client (like res in Express.js)
//   - r *http.Request - This contains information about the incoming request (like req in Express.js)
//
// The * means it's a pointer (reference to the request object, not a copy of it)
func homeHandler(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set() - This sets HTTP headers (metadata about our response)
	// "Content-Type" tells the browser/client what kind of data we're sending
	// "text/html" means "this is HTML that you should render as a web page"
	w.Header().Set("Content-Type", "text/html")

	// fmt.Fprintf(w, ...) - This writes formatted text directly to the response
	// w is the response writer (where we're sending data to the client)
	// The backticks ` create a "raw string literal" - everything inside is treated literally (including newlines)
	// This means we can write multi-line HTML and it preserves the formatting
	fmt.Fprintf(w, `
		<h1>Welcome to my Go To-Do API!</h1>
		<p>Now powered by MongoDB!</p>
		<p>Your tasks are now saved in the cloud </p>
		<h2>Available endpoints:</h2>
		<ul>
			<li><a href="/tasks">GET /tasks</a> - Get all tasks</li>
			<li>POST /tasks - Create a new task</li>
			<li>GET /tasks?id=X - Get task by ID</li>
			<li>PUT /tasks?id=X - Update a task</li>
			<li>DELETE /tasks?id=X - Delete a task</li>
		</ul>
		<p><a href="/health">Health Check</a></p>
	`)
}

// func healthHandler(...) - This function handles the health check endpoint
// Health checks are used to verify the server is running properly (monitoring tools use this)
func healthHandler(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set() - Set the content type to JSON
	// "application/json" tells the client "this is JSON data"
	w.Header().Set("Content-Type", "application/json")

	// fmt.Fprintf(w, ...) - Write a JSON string directly to the response
	// The backtick creates a raw string literal containing JSON
	fmt.Fprintf(w, `{"status": "healthy", "message": "Server is running with MongoDB!"}`)
}

// func tasksHandler(...) - This function handles requests to the /tasks endpoint
// It handles different HTTP methods: GET (read), POST (create), PUT (update), DELETE (delete)
func tasksHandler(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set() - Tell the client we're sending JSON data in our response
	w.Header().Set("Content-Type", "application/json")

	// r.Method - This gives us the HTTP method used in the request (GET, POST, PUT, DELETE, etc.)
	// We'll check which method was used and call different code for each one
	switch r.Method {

	// case "GET" - Handle requests to read/fetch tasks
	case "GET":
		// Check if an ID parameter was provided in the URL query (like ?id=...)
		// r.URL.Query() gets the query parameters from the URL
		idParam := r.URL.Query().Get("id") // Get the value of the "id" parameter, or "" if not present

		// if idParam == "" - If no ID was provided, return all tasks
		if idParam == "" {
			// Call a helper function to get and return all tasks from MongoDB
			getAllTasksHandler(w)
		} else {
			// Convert MongoDB ObjectID from string
			// primitive.ObjectIDFromHex() converts a string to MongoDB ObjectID
			// This returns TWO values: the ObjectID and an error (this is Go's way of error handling)
			objectID, err := primitive.ObjectIDFromHex(idParam)

			// if err != nil - If there was an error converting the ID
			if err != nil {
				// http.Error() - Send an error response back to the client
				// Status 400 means "Bad Request" (the client sent invalid data)
				http.Error(w, "Invalid ID", http.StatusBadRequest)
				return // Exit the function early
			}

			// Call a helper function to get and return a specific task by ID from MongoDB
			getTaskByIDHandler(w, objectID)
		}

	// case "POST" - Handle requests to create a new task
	case "POST":
		createTaskHandler(w, r)

	// case "PUT" - Handle requests to update an existing task
	case "PUT":
		updateTaskHandler(w, r)

	// case "DELETE" - Handle requests to delete a task
	case "DELETE":
		deleteTaskHandler(w, r)

	// default - Handle any HTTP method we don't support
	default:
		// http.Error() - Send an error response (405 = Method Not Allowed)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// func getAllTasksHandler(...) - Helper function to get and return all tasks from MongoDB
func getAllTasksHandler(w http.ResponseWriter) {
	// Create a context with 5-second timeout for this specific operation
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() // Cleanup when function exits

	// collection.Find() - Query MongoDB for all documents in the collection
	// The first parameter is a filter (nil = no filter, get all documents)
	// The second parameter is the cursor (placeholder for results)
	cursor, err := collection.Find(ctx, bson.M{})

	// if err != nil - Check if query failed
	if err != nil {
		http.Error(w, "Failed to fetch tasks from database", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx) // Make sure to close the cursor when done

	// Create a slice to hold all tasks
	var tasks []Task

	// cursor.All() - Read all results from the cursor into our tasks slice
	// This automatically iterates through all documents in the cursor
	err = cursor.All(ctx, &tasks)

	// if err != nil - Check if reading results failed
	if err != nil {
		http.Error(w, "Failed to decode tasks", http.StatusInternalServerError)
		return
	}

	// json.NewEncoder(w).Encode(tasks) - Convert tasks slice to JSON and send to client
	json.NewEncoder(w).Encode(tasks)

	// fmt.Println() - Print success message to console
	fmt.Printf("✅ Retrieved %d tasks from MongoDB\n", len(tasks))
}

// func getTaskByIDHandler(...) - Helper function to get a specific task by its ID from MongoDB
func getTaskByIDHandler(w http.ResponseWriter, objectID primitive.ObjectID) {
	// Create a context with 5-second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// collection.FindOne() - Find a single document matching the filter
	// bson.M{"_id": objectID} creates a filter that matches documents with this specific ID
	// Result is stored in a single Result object
	var task Task
	err := collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&task)

	// if err != nil - Check if task not found or other error
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// Task doesn't exist in database
			http.Error(w, "Task not found", http.StatusNotFound)
		} else {
			// Some other error occurred
			http.Error(w, "Failed to fetch task", http.StatusInternalServerError)
		}
		return
	}

	// Convert task to JSON and send to client
	json.NewEncoder(w).Encode(task)
	fmt.Printf("✅ Retrieved task with ID %s\n", objectID.Hex())
}

// func createTaskHandler(...) - Helper function to create a new task in MongoDB
func createTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Declare a variable to hold the incoming task data
	var newTask Task

	// json.NewDecoder(r.Body) - Create a JSON decoder that reads from the request body
	// .Decode(&newTask) - Convert the JSON from the request body into a Task struct
	// &newTask means "pass a reference to the newTask variable" (so Decode can modify it)
	err := json.NewDecoder(r.Body).Decode(&newTask)

	// if err != nil - Check if decoding failed (maybe the JSON was malformed)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Check that required fields are provided
	// newTask.Title == "" - Check if title is empty
	if newTask.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}

	// Set completed to false by default if not specified
	newTask.Completed = false

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// collection.InsertOne() - Insert a new document into the MongoDB collection
	// This automatically generates a new ObjectID for the document
	// The result contains information about the inserted document (including its new ID)
	result, err := collection.InsertOne(ctx, newTask)

	// if err != nil - Check if insert failed
	if err != nil {
		http.Error(w, "Failed to create task in database", http.StatusInternalServerError)
		return
	}

	// Get the auto-generated ID from the insert result
	// result.InsertedID returns the ID that MongoDB generated for our document
	newTask.ID = result.InsertedID.(primitive.ObjectID) // Type assertion: "this is an ObjectID"

	// Return the created task as JSON with its new ID
	json.NewEncoder(w).Encode(newTask)
	fmt.Printf("✅ Created new task: %s with ID %s\n", newTask.Title, newTask.ID.Hex())
}

// func updateTaskHandler(...) - Helper function to update an existing task in MongoDB
func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Get the ID parameter from the URL query
	idParam := r.URL.Query().Get("id")

	// if idParam == "" - Check if ID was provided
	if idParam == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	// Convert ID from string to ObjectID
	objectID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Declare a variable to hold the updated task data
	var updatedTask Task

	// Decode the JSON from the request body
	err = json.NewDecoder(r.Body).Decode(&updatedTask)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Create context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// First, find the existing task to preserve fields that weren't provided in update
	var existingTask Task
	err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&existingTask)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Task not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to fetch task", http.StatusInternalServerError)
		}
		return
	}

	// Update fields if provided in the request, otherwise keep existing values
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

	// Update the task in MongoDB
	// collection.UpdateOne() updates a single document
	result, err := collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		http.Error(w, "Failed to update task", http.StatusInternalServerError)
		return
	}

	// Check if task was actually updated
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

// func deleteTaskHandler(...) - Helper function to delete a task from MongoDB
func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Get the ID parameter from the URL query
	idParam := r.URL.Query().Get("id")

	// if idParam == "" - Check if ID was provided
	if idParam == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	// Convert ID from string to ObjectID
	objectID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Create context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// collection.DeleteOne() - Delete a single document matching the filter
	// Returns information about the deletion (including how many documents were deleted)
	result, err := collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		http.Error(w, "Failed to delete task", http.StatusInternalServerError)
		return
	}

	// Check if task was actually deleted
	if result.DeletedCount == 0 {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	// Send success response
	fmt.Fprintf(w, `{"message": "Task deleted successfully", "id": "%s"}`, idParam)
	fmt.Printf("✅ Deleted task with ID %s\n", objectID.Hex())
}
