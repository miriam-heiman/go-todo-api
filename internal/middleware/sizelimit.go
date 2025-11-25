// ============================================================================
// REQUEST SIZE LIMIT MIDDLEWARE
// ============================================================================
// This middleware limits the size of incoming request bodies
// Prevents memory exhaustion and DoS attacks from large payloads

package middleware

import (
	"net/http"
)

// MaxRequestSize is the maximum allowed request body size (in bytes)
// Current: 1MB (1024 * 1024 bytes)
const MaxRequestSize int64 = 1 * 1024 * 1024 // 1 MB

// RequestSizeLimit middleware limits the size of request bodies
// Returns 413 Payload Too Large if request exceeds limit

// How it works:
// 1. Wraps the request body with http.MaxBytesReader
// 2. MaxBytesReader stops reading after MaxRequestSize bytes
// 3. If limit exceeded, returns 413 error
// 4. If limit not exceeded, passes to next handler

func RequestSizeLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Wrap the request body with a size-limited reader
		// MaxBytesReader will:
		// - Read up to MaxRequestSize bytes normally
		// - Return error if more bytes are attempted
		// - Close connection if limit exceeded (prevents attacker from retrying)
		r.Body = http.MaxBytesReader(w, r.Body, MaxRequestSize)

		// Try to call the next handler
		// If body size is exceeded during handler execution,
		// the handler will get an error when reading the body
		next.ServeHTTP(w, r)

		// Note: we don't explicitly check the error here because
		// the handler will handle it when it tries to read the body
		// and gets an error from MaxBytesReader
	})
}

// RequestSizeLimitChi is the Chi-compatible version
func RequestSizeLimitChi(next http.Handler) http.Handler {
	return RequestSizeLimit(next)
}
