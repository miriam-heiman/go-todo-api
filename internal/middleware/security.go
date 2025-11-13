// This middleware adds security headers to protect against common attacks

package middleware

import (
	"net/http"
)

// SecurityHeaders adds HTTP security headers to all responses
// These headers protect against common web vulnerabilities
func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// X-Content-Type-Options: Prevents MIME type sniffing
		// Stops browsers from guessing content types
		// Prevents execution of JavaScript disguised as images
		w.Header().Set("X-Content-Type-Options", "nosniff")

		// X-Frame-Options: Prevents clickjacking attacks
		// Stops your API from being embedded in iframes on malicious sites
		// DENY = never allow framing
		w.Header().Set("X-Frame-Options", "DENY")

		// X-XSS-Protection: Enables browser XSS filters
		// mode=block = stop page from loading if XSS is detected
		w.Header().Set("X-XSS-Protection", "1; mode=block")

		// Content-Security-Policy: Controls what resources can be loaded
		// Modern approach for XSS filtering (better than X-XSS-Protection but doesn't work in all browsers)
		w.Header().Set("Content-Security-Policy", "default-src 'none'; frame-ancestors 'none'")

		// Referrer-Policy: Controls how much referrer information is sent
		// no-referrer = don't send referrer header (protects user privacy)
		w.Header().Set("Referrer-Policy", "no-referrer")

		// Strict-Transport-Security (HSTS): Forces HTTPS
		// max-age=31536000 = enforce HTTPS for 1 year
		// includeSubDomains = apply to all subdomains
		w.Header().Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")

		// Permissions-Policy: Controls browser features
		// Disables geolocation, microphone, camera access
		w.Header().Set("Permissions-Policy", "geolocation=(), microphone=(), camera=()")

		// Prevent caching of API responses with sensitive data:
		w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, private")

		// Hide server/framework version info:
		w.Header().Set("Server", "") // Remove server identification

		// Call next handler
		next.ServeHTTP(w, r)
	})
}

// SecurityHeadersChi is the Chi-compatible version
func SecurityHeadersChi(next http.Handler) http.Handler {
	return SecurityHeaders(next)
}
