// ============================================================================
// PACKAGE DECLARATION
// ============================================================================
// Package middleware contains HTTP middleware functions
// Middleware = code that runs BEFORE and/or AFTER your route handlers
// Think of middleware like security checks at an airport - everyone goes through them
package middleware

// ============================================================================
// IMPORTS
// ============================================================================
import (
	"log"        // log = for printing log messages to console
	"net/http"   // net/http = for HTTP types (Handler, ResponseWriter, Request)
	"time"       // time = for measuring request duration
)

// ============================================================================
// LOGGING MIDDLEWARE
// ============================================================================
// Logging logs information about each HTTP request
// This helps with debugging and monitoring by showing: method, path, and how long the request took
//
// What it does:
// 1. Records the start time of the request
// 2. Calls the next handler (your actual route handler)
// 3. After the handler finishes, logs the request details
//
// Output format:
//   GET /tasks 5.234ms
//   POST /tasks 12.456ms
//   PUT /tasks/123 3.789ms
//
// Middleware Pattern:
// Middleware in Go uses the "wrapper" pattern:
// - It takes a handler (next) as input
// - Returns a new handler that wraps the original
// - The new handler does something before/after calling next
//
// Flow:
//   Request → Logging Middleware → Your Handler → Response
//                ↓                       ↑
//           Log start time          Log duration
func Logging(next http.Handler) http.Handler {
	// return http.HandlerFunc() creates a new handler
	// The function inside receives every HTTP request
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// --------------------------------------------------------------------
		// BEFORE THE HANDLER RUNS
		// --------------------------------------------------------------------
		// Record the start time of the request
		// time.Now() = current time (like Date.now() in JavaScript)
		start := time.Now()

		// --------------------------------------------------------------------
		// RUN THE ACTUAL HANDLER
		// --------------------------------------------------------------------
		// next.ServeHTTP() calls the next handler in the chain
		// This is where your route handler (GetAllTasks, CreateTask, etc.) runs
		// When this returns, the request has been fully processed
		next.ServeHTTP(w, r)

		// --------------------------------------------------------------------
		// AFTER THE HANDLER RUNS
		// --------------------------------------------------------------------
		// Log the request details
		// time.Since(start) = how long since start time (duration)
		// r.Method = HTTP method (GET, POST, PUT, DELETE)
		// r.URL.Path = request path (/tasks, /tasks/123, etc.)
		//
		// Example output: GET /tasks 5.234ms
		log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}

// ============================================================================
// CHI-COMPATIBLE WRAPPER
// ============================================================================
// LoggingChi is the Chi-compatible middleware version
// Chi middleware has the exact same signature as standard middleware,
// so this is just an alias for clarity (shows we're using it with Chi)
//
// In main.go we use:
//   router.Use(middleware.LoggingChi)
func LoggingChi(next http.Handler) http.Handler {
	return Logging(next) // Just call the standard Logging function
}

// ============================================================================
// WHY LOGGING MATTERS
// ============================================================================
//
// Logging every request is essential for:
//
// 1. **Debugging**: When something goes wrong, logs show what requests came in
//    - "Why did this task disappear?" → Check logs for DELETE request
//    - "Why is the API slow?" → Check request durations in logs
//
// 2. **Monitoring**: Track API usage patterns
//    - Which endpoints are most popular?
//    - What's the average response time?
//    - Are there any slow endpoints?
//
// 3. **Security**: Detect suspicious activity
//    - Multiple failed requests from same IP
//    - Unusual access patterns
//    - Potential attacks or abuse
//
// 4. **Compliance**: Many regulations require request logs
//    - GDPR, HIPAA, SOC 2 often require audit trails
//
// In production, you'd typically:
// - Send logs to a centralized logging system (like ELK Stack, Datadog)
// - Add more details (user ID, request ID, status code)
// - Use structured logging (JSON format) for easier parsing
//
// ============================================================================
