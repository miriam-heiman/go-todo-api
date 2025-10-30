package middleware

import (
	"net/http"
	"os"
)

// Auth checks if the request has a valid API key
// This protects endpoints from unauthorised access
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Step 1: Get the API key from environment variable
		// In production, this would come from secure storage
		validAPIKey := os.Getenv("API_KEY")

		// Step 2: Get the API key from the request header
		// Client must send: X-API-Key: their-key-here
		requestAPIKey := r.Header.Get("X-API-Key")

		// Step 3: Check if API key is missing
		if requestAPIKey == "" {
			// Return 401 Unauthorised
			http.Error(w, "API key required", http.StatusUnauthorized)
			return
		}

		// Step 4: Check if API key is invalid
		if requestAPIKey != validAPIKey {
			// Return 403 Forbidden
			http.Error(w, "Invalid API key", http.StatusForbidden)
			return
		}

		// Step 5: API key is valid - allow request to continue
		next.ServeHTTP(w, r)
	})
}

// AuthChi is the Chi-compatible version
func AuthChi(next http.Handler) http.Handler {
	return Auth(next)
}
