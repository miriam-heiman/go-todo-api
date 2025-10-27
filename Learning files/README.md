# Go To-Do API - Learning Project

A simple REST API built with Go to learn the language!

## What This Does

- **HTTP Server**: Listens on port 8080
- **Tasks Management**: Store and retrieve tasks
- **JSON API**: Returns data in JSON format

## How to Run

1. Make sure you're in this directory:
   ```bash
   cd ~/Documents/GitHub/go-todo-api
   ```

2. Run the server:
   ```bash
   go run main.go
   ```

3. Visit in your browser:
   - http://localhost:8080/ (homepage)
   - http://localhost:8080/tasks (get all tasks)
   - http://localhost:8080/health (health check)

## Understanding the Code

### Package Declaration
```go
package main
```
This means "this is a runnable program" (not a library).

### Imports
```go
import (
    "fmt"    // For printing and formatting
    "log"    // For logging
    "net/http"  // For HTTP server
)
```

### Struct (Custom Data Type)
```go
type Task struct {
    ID          int    `json:"id"`
    Title       string `json:"title"`
    Description string `json:"description"`
    Completed   bool   `json:"completed"`
}
```
This defines what a Task looks like. The `json:"..."` tags tell Go how to convert it to JSON.

### The `main()` Function
This runs when you start the program. It:
1. Creates a sample task
2. Sets up URL routes
3. Starts listening on port 8080

### Handler Functions
Functions that handle HTTP requests. They receive:
- `w http.ResponseWriter` - Used to send response back to client
- `r *http.Request` - Contains info about the incoming request

## Next Steps

We'll expand this to handle:
- Creating new tasks (POST)
- Updating tasks (PUT)
- Deleting tasks (DELETE)
- Proper JSON responses

