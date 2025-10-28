package middleware

import "net/http"

// CORS enables Cross-Origin Resource Sharing
// CORS allows your API to be accessed from web browsers on different domains
// Without CORS, browsers block requests from other websites for security
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers to allow requests from any origin
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight requests (OPTIONS method)
		// Browsers send OPTIONS requests before actual requests to check if CORS is allowed
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

// CORSChi is the Chi-compatible middleware version
func CORSChi(next http.Handler) http.Handler {
	return CORS(next)
}
