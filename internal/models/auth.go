package models

// ============================================================================
// AUTHENTICATION MODELS
// ============================================================================
// These models define the structure for login and JWT authentication

// LoginInput is what the user sends when logging in
type LoginInput struct {
	Body struct {
		Username string `json:"username" doc:"Username" minLength:"3" maxLength:"50" example:"john"`
		Password string `json:"password" doc:"Password" minLength:"6" maxLength:"100" example:"secret123"`
	}
}

// LoginOutput is what the server sends back after a successful login
// Contains the JWT token and when it expires
type LoginOutput struct {
	Body struct {
		Token     string `json:"token" doc:"JWT token for authentication" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
		ExpiresIn int64  `json:"expires_in" doc:"Token expiration time in seconds" example:"3600"`
	}
}

// JWTClaims is the data stored inside the JWT token
// This is what the token "claims" about the user
// The server will check this data when validating the token
type JWTClaims struct {
	UserID   string `json:"user_id"`  // Unique identifier for the user
	Username string `json:"username"` // Username (for logging/display)
	Exp      int64  `json:"exp"`      // Expiration time (Unix timestamp)
	Iat      int64  `json:"iat"`      // Issued at time (when token was created)
}
