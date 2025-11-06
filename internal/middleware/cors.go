// ============================================================================
// PACKAGE DECLARATION
// ============================================================================
// Package middleware contains HTTP middleware functions
package middleware

// ============================================================================
// IMPORTS
// ============================================================================
import "net/http" // net/http = for HTTP types and constants

// ============================================================================
// CORS MIDDLEWARE
// ============================================================================
// CORS enables Cross-Origin Resource Sharing
// CORS allows your API to be accessed from web browsers on different domains
// Without CORS, browsers block requests from other websites for security
//
// What is CORS?
// CORS = Cross-Origin Resource Sharing
// It's a security feature built into web browsers (not in curl, Postman, etc.)
//
// The Problem CORS Solves:
// Imagine you have:
// - Your API running at:     http://localhost:8080
// - Your frontend running at: http://localhost:3000
//
// Without CORS, when your frontend tries to fetch data from your API,
// the browser blocks the request with an error like:
// "Access to fetch at 'http://localhost:8080/tasks' from origin
//  'http://localhost:3000' has been blocked by CORS policy"
//
// This is a security feature to prevent malicious websites from:
// - Stealing data from other websites
// - Making unauthorized requests on behalf of users
//
// The Solution:
// Your API needs to explicitly say "I allow requests from other domains"
// by sending special HTTP headers (Access-Control-Allow-* headers)
//
// Flow with CORS:
//   Browser (localhost:3000) → Sends request with Origin header
//   API (localhost:8080) → Returns Access-Control-Allow-Origin header
//   Browser → "OK, the API allows this origin" → Allows the request
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// --------------------------------------------------------------------
		// SET CORS HEADERS (TELL BROWSER WHAT'S ALLOWED)
		// --------------------------------------------------------------------

		// Access-Control-Allow-Origin: Which domains can access this API?
		// "*" = wildcard = allow ANY domain
		// In production, you'd typically specify your frontend domain:
		//   "https://myapp.com" or "http://localhost:3000"
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// Access-Control-Allow-Methods: Which HTTP methods are allowed?
		// This tells the browser: "You can send GET, POST, PUT, DELETE, OPTIONS"
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		// Access-Control-Allow-Headers: Which headers can be sent?
		// This allows the browser to send:
		// - Content-Type (for JSON requests)
		// - Authorization (for auth tokens like JWT)
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// --------------------------------------------------------------------
		// HANDLE PREFLIGHT REQUESTS (OPTIONS METHOD)
		// --------------------------------------------------------------------
		// Before making certain requests, browsers send a "preflight" request
		// Preflight = a pre-check to see if the actual request is allowed
		//
		// When does a browser send a preflight?
		// - For requests with custom headers (like Authorization)
		// - For methods other than GET/POST with simple headers
		// - For requests with Content-Type other than form-data or text/plain
		//
		// Preflight request format:
		//   OPTIONS /tasks
		//   Origin: http://localhost:3000
		//   Access-Control-Request-Method: DELETE
		//   Access-Control-Request-Headers: Authorization
		//
		// Our response:
		//   200 OK
		//   Access-Control-Allow-Origin: *
		//   Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS
		//   Access-Control-Allow-Headers: Content-Type, Authorization
		//
		// After getting this response, the browser knows it's safe to send
		// the actual request (DELETE /tasks with Authorization header)
		if r.Method == "OPTIONS" {
			// Return 200 OK with the CORS headers we already set above
			w.WriteHeader(http.StatusOK)
			return // Stop here, don't call next handler
		}

		// --------------------------------------------------------------------
		// CALL THE NEXT HANDLER (FOR NON-PREFLIGHT REQUESTS)
		// --------------------------------------------------------------------
		// If it's not a preflight, proceed to the actual handler
		next.ServeHTTP(w, r)
	})
}

// ============================================================================
// CHI-COMPATIBLE WRAPPER
// ============================================================================
// CORSChi is the Chi-compatible middleware version
// Chi middleware has the same signature as standard middleware,
// so this is just an alias for clarity
//
// In main.go we use:
//   router.Use(middleware.CORSChi)
func CORSChi(next http.Handler) http.Handler {
	return CORS(next) // Just call the standard CORS function
}

// ============================================================================
// CORS SECURITY CONSIDERATIONS
// ============================================================================
//
// 1. **Using "*" (Allow All Origins)**:
//    - Good for: Public APIs that anyone can use
//    - Bad for: APIs with user authentication/authorization
//    - Why: Allows any website to access your API from a browser
//
// 2. **Production Best Practice**:
//    Instead of "*", specify your actual frontend domain:
//    ```go
//    allowedOrigins := []string{
//        "https://myapp.com",
//        "https://www.myapp.com",
//        "http://localhost:3000", // for development
//    }
//    origin := r.Header.Get("Origin")
//    if slices.Contains(allowedOrigins, origin) {
//        w.Header().Set("Access-Control-Allow-Origin", origin)
//    }
//    ```
//
// 3. **CORS Only Affects Browsers**:
//    - Tools like curl, Postman, or server-to-server requests ignore CORS
//    - CORS is ONLY a browser security feature
//    - You still need authentication/authorization for security!
//
// 4. **What CORS Does NOT Protect Against**:
//    - CORS doesn't authenticate users
//    - CORS doesn't authorize requests
//    - CORS doesn't encrypt data
//    - It only controls which browser-based origins can access your API
//
// ============================================================================
