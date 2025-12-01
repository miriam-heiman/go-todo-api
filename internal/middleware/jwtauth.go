package middleware

import (
	"net/http"
	"strings"

	"go-todo-api/internal/auth"
	"go-todo-api/internal/logger"
)

// ============================================================================
// JWT AUTHENTICATION MIDDLEWARE
// ============================================================================
// This middleware checks JWT tokens on every request
// Replaces the old API key authentication

// JWTAuth checks if the request has a valid JWT token
// Flow:
// 1. Extract token from Authorization header
// 2. Validate token (signature, expiration)
// 3. If valid: allow request
// 4. If  invalid: return 401 Unauthorized
func JWTAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Step 1: GET TOKEN FROM HEADER
		// Token should be in header: Authorization: Bearer <token>
		authHeader := r.Header.Get("Authorization")

		// Check if Authorization header exists
		if authHeader == "" {
			logger.Log.Warn("Missing Authorization header", "path", r.URL.Path)
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		// Step 2: Extract token
		// Authorization header format: "Bearer eyJhbGc..."
		// We need to extract the token part (after "Bearer ")
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			logger.Log.Warn("Invalid Authorization header format", "header", authHeader)
			http.Error(w, "Invalid Authorization header format. Use: Bearer <token>", http.StatusUnauthorized)
			return
		}

		token := parts[1]

		// Step 3: Validate token
		// Check if token is valid (signature, expiration etc.)
		claims, err := auth.ValidateToken(token)
		if err != nil {
			// Token is invalid
			logger.Log.Warn("Invalid JWT token",
				"error", err.Error(), "path", r.URL.Path)
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// Step 4: Token is valid - Allow request
		// Log successful authentication
		logger.Log.Info("JWT authentication successful",
			"user_id", claims.UserID,
			"username", claims.Username,
			"path", r.URL.Path,
		)

		// Continue to next handler
		next.ServeHTTP(w, r)

	})
}

// JWTAuthChi is the Chi-compatible version
func JWTAuthChi(next http.Handler) http.Handler {
	return JWTAuth(next)
}
