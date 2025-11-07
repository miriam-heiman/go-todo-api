package logger

import (
	// STANDARD LIBRARIES
	"context"
	"log/slog"

	// THIRD-PARTY LIBRARY
	"go.opentelemetry.io/otel/trace"
)

// WithTrace returns a logger with trace context fields added
// This links logs to traces for correlation in Grafana
func WithTrace(ctx context.Context) *slog.Logger {
	// Extract the span from context
	span := trace.SpanFromContext(ctx)
	spanContext := span.SpanContext()

	// If there's a valid span, add trace and span IDs to logger
	if spanContext.IsValid() {
		return Log.With(
			slog.String("trace_id", spanContext.TraceID().String()),
			slog.String("span_id", spanContext.SpanID().String()),
		)
	}
	// If no span, return the regular logger
	return Log
}
