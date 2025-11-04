// ============================================================================
// PACKAGE DECLARATION
// ============================================================================
// Package database handles all MongoDB connection and management
// This package is responsible for connecting to MongoDB and providing
// access to the database collection throughout the application
package database

// ============================================================================
// IMPORTS
// ============================================================================
import (
	// STANDARD LIBRARY PACKAGES
	"context" // context = for managing timeouts and cancellation
	"log"     // log = for error logging and fatal errors
	"os"      // os = for reading environment variables
	"time"    // time = for creating timeouts

	// OUR OWN PACKAGE
	logger "go-todo-api/internal/logger" // Our structured logger

	// THIRD-PARTY PACKAGES
	"github.com/joho/godotenv"                  // godotenv = loads .env file into environment
	"go.mongodb.org/mongo-driver/mongo"         // mongo = MongoDB driver for Go
	"go.mongodb.org/mongo-driver/mongo/options" // options = MongoDB connection options
)

// ============================================================================
// PACKAGE-LEVEL VARIABLES (SHARED ACROSS ALL FILES IN THIS PACKAGE)
// ============================================================================
// var() = declares multiple variables at once
// These are package-level variables (not inside a function)
// They're accessible to all functions in this package
var (
	// client is the MongoDB client connection
	// *mongo.Client = pointer to a Client (the * means it can be nil)
	// This holds the connection to MongoDB server
	client *mongo.Client

	// collection is the specific MongoDB collection we're working with
	// In MongoDB: Database → Collection → Documents
	// Like SQL:   Database → Table → Rows
	// Our collection is: "todoapi" database, "tasks" collection
	collection *mongo.Collection
)

// ============================================================================
// CONNECT TO MONGODB
// ============================================================================
// Connect initializes the MongoDB connection
// This function is called once at server startup (in main.go)
// It:
// 1. Loads environment variables from .env file
// 2. Connects to MongoDB using connection string
// 3. Pings MongoDB to verify connection works
// 4. Sets up the collection we'll use for all operations
func Connect() {
	// ----------------------------------------------------------------------------
	// STEP 1: LOAD ENVIRONMENT VARIABLES FROM .env FILE
	// ----------------------------------------------------------------------------
	// godotenv.Load() reads the .env file and puts variables into environment
	// Example .env file:
	//   MONGO_URI=mongodb://localhost:27017/todoapi
	if err := godotenv.Load(); err != nil {
		// If .env file doesn't exist, that's okay in production
		// (production systems usually set environment variables directly)
		logger.Log.Warn("No .env file found", "fallback", "environment variables")
	}

	// ----------------------------------------------------------------------------
	// STEP 2: CREATE CONTEXT WITH TIMEOUT
	// ----------------------------------------------------------------------------
	// Create a context that will automatically timeout after 10 seconds
	// This prevents the connection attempt from hanging forever
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel() // Clean up context when function exits

	// ----------------------------------------------------------------------------
	// STEP 3: GET MONGODB CONNECTION STRING FROM ENVIRONMENT
	// ----------------------------------------------------------------------------
	// os.Getenv() reads an environment variable
	// MONGO_URI format:
	//   Local:  mongodb://localhost:27017/todoapi
	//   Atlas:  mongodb+srv://username:password@cluster.mongodb.net/todoapi
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		// If MONGO_URI is not set, we can't connect to database
		// logger.log.error uses the structured logger we set up
		// log.Fatal() prints error and exits the program (like a crash)
		// We first log the error with structured logging, then we use log.fatal() to exit the program.
		logger.Log.Error("MONGO_URI not found", "required", true)
		log.Fatal("MONGO_URI not found. Please set it in your .env file")
	}

	// ----------------------------------------------------------------------------
	// STEP 4: CREATE MONGODB CLIENT WITH CONNECTION OPTIONS
	// ----------------------------------------------------------------------------
	// options.Client() creates a ClientOptions object
	// .ApplyURI() tells it to use our connection string
	clientOptions := options.Client().ApplyURI(mongoURI)

	// ----------------------------------------------------------------------------
	// STEP 5: ACTUALLY CONNECT TO MONGODB
	// ----------------------------------------------------------------------------
	// mongo.Connect() establishes the connection to MongoDB server
	// Note: We're assigning to the package-level "client" variable (not creating a new one)
	// That's why we use "var err error" first - to avoid shadowing with :=
	var err error
	client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		// If connection fails (wrong URI, MongoDB not running, network issue)
		logger.Log.Error("Failed to connect to MongoDB", "error", err)
		log.Fatal("Failed to connect to MongoDB:")
	}

	// ----------------------------------------------------------------------------
	// STEP 6: PING MONGODB TO VERIFY CONNECTION WORKS
	// ----------------------------------------------------------------------------
	// Just because Connect() succeeded doesn't mean we can actually talk to MongoDB
	// Ping() sends a test message to verify the connection is working
	err = client.Ping(ctx, nil)
	if err != nil {
		// If ping fails, the connection isn't working properly
		logger.Log.Error("Failed to ping MongoDB", "error", err)
		log.Fatal("Failed to ping MongoDB:")
	}

	// ----------------------------------------------------------------------------
	// STEP 7: SELECT DATABASE AND COLLECTION
	// ----------------------------------------------------------------------------
	// MongoDB structure: Server → Database → Collection → Documents
	// client.Database("todoapi") = selects the "todoapi" database
	// .Collection("tasks") = selects the "tasks" collection within that database
	//
	// Note: MongoDB will automatically create the database and collection
	// the first time we insert a document - we don't need to create them manually!
	collection = client.Database("todoapi").Collection("tasks")

	// ----------------------------------------------------------------------------
	// STEP 8: LOG SUCCESS
	// ----------------------------------------------------------------------------
	logger.Log.Info("Connected to MongoDB", "database", "todoapi", "collection", "tasks")
}

// ============================================================================
// GET COLLECTION (GETTER FUNCTION)
// ============================================================================
// GetCollection returns the MongoDB collection for tasks
// This is called by handlers to access the database collection
//
// Why use a getter function instead of accessing "collection" directly?
// - Encapsulation: Other packages can't modify the collection variable
// - Safety: We control how the collection is accessed
// - Flexibility: We could add logic here later (logging, connection checks, etc.)
//
// Usage in handlers:
//
//	collection := database.GetCollection()
//	collection.Find(ctx, bson.M{})
func GetCollection() *mongo.Collection {
	return collection // Return the package-level collection variable
}

// ============================================================================
// CLOSE CONNECTION (CLEANUP FUNCTION)
// ============================================================================
// Close closes the MongoDB connection
// This should be called when the server is shutting down to:
// - Close all open connections gracefully
// - Release system resources (file descriptors, memory)
// - Allow pending operations to complete
//
// Note: In our current main.go, we don't call this because the server
// runs forever (until killed). In production, you'd call this in a
// shutdown handler that runs when the server receives a stop signal.
//
// Example usage (not currently used):
//
//	defer database.Close() // Call Close when main() exits
func Close() {
	// Only try to disconnect if client was actually created
	if client != nil {
		// Create context with timeout for disconnect operation
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Attempt to disconnect from MongoDB
		if err := client.Disconnect(ctx); err != nil {
			// Log error but don't crash (we're shutting down anyway)
			logger.Log.Error("Error disconnecting from MongoDB", "error", err)
		}
	}
}

// ============================================================================
// HOW THIS PACKAGE IS USED
// ============================================================================
//
// 1. Server Startup (main.go):
//    database.Connect()  ← Connects to MongoDB once
//
// 2. During Request Handling (handlers):
//    collection := database.GetCollection()  ← Get collection reference
//    collection.Find(...)                    ← Query database
//    collection.InsertOne(...)               ← Insert data
//    collection.UpdateOne(...)               ← Update data
//    collection.DeleteOne(...)               ← Delete data
//
// 3. Server Shutdown (optional, not currently implemented):
//    database.Close()  ← Disconnect from MongoDB
//
// The key insight: We connect ONCE at startup and reuse that connection
// for all requests. This is much more efficient than connecting/disconnecting
// for each request (which would be extremely slow).
//
// ============================================================================
