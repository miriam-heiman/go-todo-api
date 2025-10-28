package handlers

import (
	"context"

	"go-todo-api/internal/models"
)

// Health handles the health check endpoint
func Health(ctx context.Context, input *struct{}) (*models.HealthOutput, error) {
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
