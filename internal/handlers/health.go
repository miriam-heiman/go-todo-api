package handlers

import (
	"fmt"
	"net/http"
)

// Health handles the health check endpoint
func Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"status": "healthy", "message": "Server is running with MongoDB!"}`)
}
