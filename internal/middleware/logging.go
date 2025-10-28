package middleware

import (
	"log"
	"net/http"
	"time"
)

// Logging logs information about each HTTP request
// This helps with debugging and monitoring by showing: method, path, and how long the request took
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Record the start time of the request
		start := time.Now()

		// Call the next handler in the chain (the actual route handler)
		next.ServeHTTP(w, r)

		// After the handler finishes, log the request details
		log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}

// LoggingChi is the Chi-compatible middleware version
func LoggingChi(next http.Handler) http.Handler {
	return Logging(next)
}
