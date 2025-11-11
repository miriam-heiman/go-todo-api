package handlers

import (
	// STANDARD LIBARIES
	"context"
	"strings"
	"testing"

	// THIRD-PARTY LIBRARIES
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/humatest"

	//  OUR OWN PACKAGES
	"go-todo-api/internal/models"
)

// TestHealthHandler tests the health check endpoint
func TestHealthHandler(t *testing.T) {

	// Act: Call the Health handler directly
	ctx := context.Background()
	input := &models.HealthInput{}

	output, err := Health(ctx, input)

	// Assert: Check the results
	if err != nil {
		t.Fatalf("Health handler returned error: %v", err)
	}

	if output.Body.Status != "healthy" {
		t.Errorf("Expected status 'healthy', got '%s'", output.Body.Status)
	}

	if output.Body.Message == "" {
		t.Error("Expected non-empty message")
	}

	// Log success
	t.Logf("✅ Health check passed: %s - %s", output.Body.Status, output.Body.Message)
}

// TestHealthHandler_Integration tests the health endpoint via HTTP
func TestHealthHandler_Integration(t *testing.T) {
	// This tests the FULL HTTP flow (more realistic)

	// Arrange: Create test API and register endpoint
	_, api := humatest.New(t)

	// Register the health endpoint like in main.go
	huma.Register(api, huma.Operation{
		OperationID: "health-check",
		Method:      "GET",
		Path:        "/health",
		Summary:     "Health check",
	}, Health)

	// Act: Make HTTP request to /health
	resp := api.Get("/health")

	// Assert: Check HTTP response
	if resp.Code != 200 {
		t.Errorf("Expected status 200, got %d", resp.Code)
	}

	// Check response body contains expected JSON
	body := resp.Body.String()
	if !strings.Contains(body, "healthy") {
		t.Errorf("Response doesn't contain 'healthy': %s", body)
	}

	if !strings.Contains(body, "status") {
		t.Errorf("Response doesn't contain 'status' field: %s", body)
	}

	t.Logf("✅ Integration test passed. Response: %s", body)
}
