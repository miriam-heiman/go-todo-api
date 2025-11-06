// ============================================================================
// PACKAGE DECLARATION
// ============================================================================
// Package handlers contains all the HTTP handler functions for our API
package handlers

// ============================================================================
// IMPORTS
// ============================================================================
import (
	// STANDARD LIBRARY PACKAGE
	"context" // context = for managing request context

	// OUR OWN PACKAGE
	"go-todo-api/internal/models" // Our data structures (HealthOutput)
)

// ============================================================================
// HEALTH CHECK ENDPOINT
// ============================================================================
// Health handles the health check endpoint
// This is called when someone makes a GET request to /health
//
// Purpose: Health checks are used by monitoring tools, load balancers, and
// orchestration systems (like Kubernetes) to verify the service is running.
// If this endpoint returns successfully, it means:
// - The server is up and responding to requests
// - The Go application is running without crashing
//
// Huma Handler Signature:
// - Input: context.Context + *struct{} (no parameters needed)
// - Output: *models.HealthOutput (contains status and message) + error
//
// Example request:  GET /health
// Example response: {"status": "healthy", "message": "Server is running with MongoDB!"}
func Health(ctx context.Context, input *models.HealthInput) (*models.HealthOutput, error) {

	// Return a simple success response
	return &models.HealthOutput{
		Body: struct {
			Status  string `json:"status" doc:"Health status" example:"healthy"`
			Message string `json:"message" doc:"Health message" example:"Server is running with MongoDB!"`
		}{
			Status:  "healthy",
			Message: "Server is running with MongoDB!",
		},
	}, nil
}

// ============================================================================
// WHY HEALTH CHECKS MATTER
// ============================================================================
//
// Health checks are a standard practice in production systems:
//
// 1. **Monitoring Systems**: Tools like Datadog, New Relic, or Prometheus
//    regularly hit /health to detect if the service goes down
//
// 2. **Load Balancers**: AWS ELB, nginx, etc. use health checks to know
//    which servers are available to receive traffic
//
// 3. **Container Orchestration**: Kubernetes uses health checks to:
//    - Know when a container is ready to receive traffic
//    - Automatically restart unhealthy containers
//    - Route traffic away from failing instances
//
// 4. **Deployment Systems**: CI/CD pipelines check health after deployment
//    to verify the new version started successfully
//
// Advanced health checks might also:
// - Check database connectivity (ping MongoDB)
// - Check external dependencies (APIs, Redis, etc.)
// - Return degraded status if some features are down
//
// For now, our simple health check just confirms the server is running.
//
// ============================================================================
