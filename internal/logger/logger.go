package logger

import (
	"log/slog"
	"os"
)

// Global logger instance
// All parts of the app will use this single logger
var Log *slog.Logger

// Init initialises the structured logger
// Call lthis once at startup before using log
func Init() {
	// Create a JSON handler that writes to stdout (console)
	// JSON format makes it easy for Loki to parse
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo, // Log Info, Warn, Error (skip Debug in production)
	})

	// Create the logger with our handler
	Log = slog.New(handler)

	Log.Info("Logger initialised", "format", "json")
}
