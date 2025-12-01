package handlers

import (
	"context"

	"go-todo-api/internal/auth"
	"go-todo-api/internal/logger"
	"go-todo-api/internal/models"

	"github.com/danielgtaylor/huma/v2"
)

// ============================================================================
// LOGIN HANDLER
// ============================================================================
// Login authenticates a user and returns a JWT token
// This is called when someone makes a POST request to /auth/login77

func Login(ctx context.Context, input *models.LoginInput) (*models.LoginOutput, error) {
	// Step 1: VALIDATE CREDENTIALS
	// Check if username and password are correct
	userID, err := auth.ValidateCredentials(input.Body.Username, input.Body.Password)
	if err != nil {
		// Invalid credentials - log the attempt and return error
		logger.Log.Warn("Failed login attempt", "username", input.Body.Username)
		return nil, huma.Error401Unauthorized("Invalid username or password")
	}

	// Step 2: GENERATE JWT TOKEN
	// Credentials are valid - create a token for this user
	token, err := auth.GenerateToken(userID, input.Body.Username)
	if err != nil {
		// Token generation failed (shouldn't happen)
		logger.Log.Error("Failed to generate token", "error", err, "user_id", userID)
		return nil, huma.Error500InternalServerError("Failed to generate token")
	}

	// Step 3: LOG SUCCESSFUL LOGIN
	logger.Log.Info("User logged in successfully", "username", input.Body.Username, "user_id", userID)

	// Step 4: RETURN TOKEN TO USER
	// Return the token and expiration time
	return &models.LoginOutput{
		Body: struct {
			Token     string `json:"token" doc:"JWT token for authentication" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
			ExpiresIn int64  `json:"expires_in" doc:"Token expiration time in seconds" example:"3600"`
		}{
			Token:     token,
			ExpiresIn: int64(auth.TokenExpiration.Seconds()), // 3600 seconds (1 hour)
		},
	}, nil
}
