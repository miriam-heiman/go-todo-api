// Package main is the entry point for the TODO API server
package main

import (
	"fmt"
	"log"
	"net/http"

	"go-todo-api/internal/database"
	"go-todo-api/internal/handlers"
	"go-todo-api/internal/middleware"
)

func main() {
	// Initialize database connection
	database.Connect()

	// Create a new HTTP router
	mux := http.NewServeMux()

	// Register routes
	mux.HandleFunc("/", handlers.Home)
	mux.HandleFunc("/tasks", handlers.Tasks)
	mux.HandleFunc("/health", handlers.Health)

	// Wrap all routes with middleware
	// Requests flow: Logging -> CORS -> Handler
	handler := middleware.Chain(
		mux,
		middleware.Logging,
		middleware.CORS,
	)

	// Print startup information
	fmt.Println("🚀 Server starting on http://localhost:8080")
	fmt.Println("✨ Middleware enabled: Logging, CORS")
	fmt.Println("📁 Production structure: cmd/ and internal/ packages")
	fmt.Println("\nTry visiting:")
	fmt.Println("  - http://localhost:8080/ (homepage)")
	fmt.Println("  - http://localhost:8080/tasks (get all tasks)")
	fmt.Println("  - http://localhost:8080/health (health check)")

	// Start the HTTP server
	port := ":8080"
	log.Fatal(http.ListenAndServe(port, handler))
}
