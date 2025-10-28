package middleware

import "net/http"

// Chain applies multiple middleware to a handler in order
// This makes it easy to wrap a handler with multiple middleware functions
// Usage: Chain(handler, middleware1, middleware2, middleware3)
func Chain(h http.Handler, middleware ...func(http.Handler) http.Handler) http.Handler {
	// Apply middleware in reverse order so they execute in the order provided
	// If you pass [logging, cors], it will execute: logging -> cors -> handler
	for i := len(middleware) - 1; i >= 0; i-- {
		h = middleware[i](h)
	}
	return h
}
