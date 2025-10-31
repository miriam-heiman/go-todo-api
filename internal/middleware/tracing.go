package middleware

import (
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// Tracing wraps the handler with OpenTelemetry tracing
// This automatically creates a span for every HTTP request
// The span includes: method, path, status code, duration
func Tracing(next http.Handler) http.Handler {
	// otelhttp.NewHandler wraps our handler and:
	// 1. Creates a span when request starts
	// 2. Adds HTTP attributes (method, path, status)
	// 3. Ends the span when request finishes
	// 4. Records errors if they occur
	return otelhttp.NewHandler(next, "http-server",
		otelhttp.WithSpanNameFormatter(func(operation string, r *http.Request) string {
			// Custom span name: "GET /tasks" instead of just "http-server"
			return r.Method + " " + r.URL.Path
		}),
	)
}

// TracingChi is the Chi-compatible version
func TracingChi(next http.Handler) http.Handler {
	return Tracing(next)
}
