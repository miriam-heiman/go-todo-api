package tracing

import (
	// STANDARD LIBRARY PACKAGES
	"context" // Manages request lifecycles, timeouts and cancellation
	"log"     // Logging with timestamps and error handling
	"os"
	"time" // Working with the time durations and delays

	// OUR OWN PACKAGES
	"go-todo-api/internal/logger" // Our structured logger

	// THIRD-PARTY LIBRARY PACKAGES
	"go.opentelemetry.io/otel" // Exporter: Sends traces via HTTP to Jaeger/Tempo

	// OpenTelemetry core: Main OTel packages - gives access to the global tracer
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"            // Resource: Service metada
	sdktrace "go.opentelemetry.io/otel/sdk/trace"      // Trace provider: Core tracing functionality, creates spans
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0" // Semantic conventions: Standard attribute names for service.name, etc.
)

// Initialises the ServiceName variable
const ServiceName = "go-todo-api"

// Init initializes OpenTelemetry tracing
// This sets up the global tracer that the entire app will use
func Init(serviceName string) func() {
	// Step 1: Create an OTLP HTTP exporter
	// This sends traces to Jaeger (or an OTLP-compatible backend)
	ctx := context.Background()

	// Read OTLP endpoint from environment, default to localhost for local dev
	otlpEndpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if otlpEndpoint == "" {
		otlpEndpoint = "http://localhost:4318"
	}
	// Strip http:// prefix if present (the library adds it)
	if len(otlpEndpoint) > 7 && otlpEndpoint[:7] == "http://" {
		otlpEndpoint = otlpEndpoint[7:]
	}

	exporter, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpoint(otlpEndpoint),
		otlptracehttp.WithInsecure(),
	)

	if err != nil {
		logger.Log.Error("Failed to create OTLP trace exporter", "error", err)
		log.Fatal("Failed to create OTLP trace exporter:")
	}

	// Step 2: Create a resource (describes this service)
	// This adds metadata to all traces: service name, version, etc.
	res, err := resource.New(ctx, resource.WithAttributes(
		semconv.ServiceName(ServiceName),
		semconv.ServiceVersion("1.0.0"),
	),
	)
	if err != nil {
		logger.Log.Error("Failed to create resource", "error", err)
		log.Fatal("Failed to create resource")
	}

	// Step 3: Create a trace provider
	// This is the core of OpenTelemetry - it creates and manages spans
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),                // Send traces in batches (efficient)
		sdktrace.WithResource(res),                    // Attach our service metadata
		sdktrace.WithSampler(sdktrace.AlwaysSample()), // Sample 100% of traces (for learning)
	)

	// Step 4: Set as a global tracer provider
	// This makes it available everywhere in your app via otel.Tracer()
	otel.SetTracerProvider(tp)

	logger.Log.Info("OpenTelemetry tracing initialized", "endpoint", otlpEndpoint, "backend", "Jaeger")
	// Return a cleanup function
	// Call this when the server shuts down to flush any remaining traces
	return func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := tp.Shutdown(ctx); err != nil {
			logger.Log.Error("Error shutting down tracer provider", "error", err)
		}
	}
}
