// This middleware protects the API from abuse by limiting the number of requests per IP address
// This is essential for protecting Lambda deployments where excessive requests lead to high AWS bills

package middleware

import (
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"

	"go-todo-api/internal/logger"
)

// ============================================================================
// RATE LIMITER STORAGE
// ============================================================================
// 'visitor' tracks rate limit state for each IP address
type visitor struct {
	limiter  *rate.Limiter // the actual rate limiter
	lastSeen time.Time     // Last time we saw a request from this IP
}

// rateLimiter manages rate limiters for all IP addresses
type rateLimiter struct {
	visitors map[string]*visitor // Map of IP addresses to visitors
	mu       sync.RWMutex        // Lock for thread-safe access
	rate     rate.Limit          // Requests per second allowed
	burst    int                 // Maximum burst size
}

// Global rate limiter instance
var limiter *rateLimiter

// init runs when package is imported
// Sets up rate limiter with default values: 10 req/sec, burst of 20
func init() {
	limiter = &rateLimiter{
		visitors: make(map[string]*visitor),
		rate:     rate.Limit(10), // 10 requests per second
		burst:    20,             // Allow bursts up to 20 requests
	}

	// Start cleanup goroutine to remove old visitors (prevent memory leaks)
	go limiter.cleanupVisitors()
}

// ============================================================================
// RATE LIMITER METHODS
// ============================================================================

// getVisitor returns the rate limiter for an IP address
// Creates a new limiter if one doesn't exist
func (rl *rateLimiter) getVisitor(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	v, exists := rl.visitors[ip]
	if !exists {
		// Create new rate limiter for this IP
		limiter := rate.NewLimiter(rl.rate, rl.burst)
		rl.visitors[ip] = &visitor{limiter, time.Now()}
		return limiter
	}

	// Update last seen time
	v.lastSeen = time.Now()
	return v.limiter
}

// cleanupVisitors removes visitors that haven't been seen in 3 minutes
// This prevents memory leaks from accumulating stale visitors
func (rl *rateLimiter) cleanupVisitors() {
	for {
		time.Sleep(time.Minute) // Run every minute

		rl.mu.Lock()
		for ip, v := range rl.visitors {
			if time.Since(v.lastSeen) > 3*time.Minute {
				delete(rl.visitors, ip)
			}
		}
		rl.mu.Unlock()
	}
}

// ============================================================================
// MIDDLEWARE FUNCTIONS
// ============================================================================

// RateLimit middleware limits requests per IP address
// Returns 429 Too Many Requests if the limit is exceeded
func RateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract IP address from request
		ip := getIP(r)

		// Get rate limiter for this IP
		limiter := limiter.getVisitor(ip)

		// Check if request is allowed
		if !limiter.Allow() {
			// Rate limit exceeded
			logger.Log.Warn("Rate limit exceeded",
				"ip", ip,
				"path", r.URL.Path,
				"method", r.Method,
			)

			// Return 429 Too Many Requests
			http.Error(w, "Rate limit exceeded. Please try again later.", http.StatusTooManyRequests)
			return
		}

		// Request allowed - continue to next handler
		next.ServeHTTP(w, r)

	})
}

// RateLimitChi is the Chi-compatible version
func RateLimitChi(next http.Handler) http.Handler {
	return RateLimit(next)
}

// ============================================================================
// HELPER FUNCTIONS
// ============================================================================
// getIP extracts the client IP address from the request
// Handles various proxy headers and formats
func getIP(r *http.Request) string {
	// Try X-Forwarded-For header (set by proxies/load balancers)
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		return forwarded
	}

	// Try X-Real-IP header (set by some proxies)
	realIP := r.Header.Get("X-Real-IP")
	if realIP != "" {
		return realIP
	}

	// Fall back to RemoteAddr (direct connection)
	ip := r.RemoteAddr

	// Remove port if present
	// Example: "192.168.1.1:12345" should be "192.168.1.1"
	if idx := len(ip) - 1; idx >= 0 {
		for i := idx; i >= 0; i-- {
			if ip[i] == ':' {
				return ip[:i]
			}
		}
	}

	return ip
}
