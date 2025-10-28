// THIS FILE IS DEPRECATED
// The project has been refactored into a production-ready structure.
//
// New entry point: cmd/api/main.go
//
// To run the server:
//   go run cmd/api/main.go
//
// Or with hot-reload:
//   air
//
// Project structure:
//   cmd/api/          - Application entry point
//   internal/models/  - Data structures
//   internal/handlers/ - HTTP handlers
//   internal/middleware/ - Middleware functions
//   internal/database/ - Database connections
//
// See Learning files/API_FILE_STRUCTURE.md for more details.

package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("⚠️  This main.go is deprecated!")
	fmt.Println("")
	fmt.Println("The project has been refactored into a production-ready structure.")
	fmt.Println("")
	fmt.Println("To run the server, use:")
	fmt.Println("  go run cmd/api/main.go")
	fmt.Println("")
	fmt.Println("Or with hot-reload:")
	fmt.Println("  air")
	fmt.Println("")
	fmt.Println("See cmd/api/main.go for the new entry point.")
	os.Exit(1)
}
