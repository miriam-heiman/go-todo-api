# API File Structure - A Beginner's Guide

## Your Current Structure (Single File)

Right now, your entire API is in one file:

```
go-todo-api/
├── main.go              # Everything is here! (530+ lines)
├── go.mod
├── go.sum
├── .env
├── .air.toml
└── README.md
```

**This is perfectly fine for learning!** But as your project grows, this becomes hard to maintain.

## The Problem with One Big File

Imagine you have 1000+ lines in `main.go`:
- Hard to find specific functions
- Difficult to test individual pieces
- Multiple people can't work on different features easily
- Changes in one area might accidentally break another

**Solution:** Split your code into multiple files and folders by responsibility.

## Evolution of API Structure

### Stage 1: Single File (Your Current Stage ✅)
**When:** Learning, prototyping, small APIs (< 500 lines)

```
go-todo-api/
└── main.go              # Everything: handlers, DB, middleware, structs
```

**Pros:**
- Simple to understand
- Easy to navigate (everything in one place)
- Perfect for learning

**Cons:**
- Gets messy as it grows
- Hard to reuse code
- Difficult to test

---

### Stage 2: Multiple Files, Same Package
**When:** API grows to 500-1000 lines

```
go-todo-api/
├── main.go              # Entry point: main(), setup
├── handlers.go          # All handler functions
├── middleware.go        # Middleware functions
├── models.go            # Data structures (Task, User, etc.)
└── database.go          # Database connection and queries
```

**Pros:**
- Better organization
- Easier to find things
- Still simple (all in `package main`)

**Cons:**
- All code still tightly coupled
- Can't import parts into other projects
- Testing is still tricky

---

### Stage 3: Package-Based Structure (Recommended ⭐)
**When:** Production APIs, team projects, reusable code

```
go-todo-api/
├── main.go                      # Entry point only (30-50 lines)
├── go.mod
├── go.sum
├── .env
│
├── cmd/                         # Application entry points
│   └── api/
│       └── main.go              # Server startup
│
├── internal/                    # Private application code
│   ├── handlers/                # HTTP handlers
│   │   ├── tasks.go             # Task-related handlers
│   │   ├── health.go            # Health check handler
│   │   └── home.go              # Home page handler
│   │
│   ├── middleware/              # Middleware functions
│   │   ├── logging.go           # Logging middleware
│   │   ├── cors.go              # CORS middleware
│   │   └── auth.go              # Authentication middleware
│   │
│   ├── models/                  # Data structures
│   │   ├── task.go              # Task struct and methods
│   │   └── response.go          # API response structs
│   │
│   ├── database/                # Database layer
│   │   ├── mongo.go             # MongoDB connection
│   │   └── queries.go           # Database queries
│   │
│   └── config/                  # Configuration
│       └── config.go            # App configuration
│
├── pkg/                         # Public, reusable code
│   └── utils/                   # Utility functions
│       ├── validator.go         # Input validation
│       └── logger.go            # Custom logger
│
└── tests/                       # Test files
    ├── handlers_test.go
    └── middleware_test.go
```

**Pros:**
- Very organized and scalable
- Easy to find and modify code
- Testable (each package can be tested independently)
- Team-friendly (different people work on different packages)
- Reusable (pkg/ can be imported by other projects)

**Cons:**
- More complex initially
- More files to navigate
- Requires understanding of Go packages

---

## Detailed Breakdown of Package Structure

### 1. `cmd/` Directory
**Purpose:** Application entry points

```go
// cmd/api/main.go
package main

import (
    "log"
    "net/http"
    "github.com/yourusername/go-todo-api/internal/handlers"
    "github.com/yourusername/go-todo-api/internal/database"
)

func main() {
    // Initialize database
    database.Connect()

    // Setup routes
    mux := http.NewServeMux()
    handlers.RegisterRoutes(mux)

    // Start server
    log.Fatal(http.ListenAndServe(":8080", mux))
}
```

**Why:**
- Keeps main.go small and focused
- Easy to create multiple entry points (API server, CLI tools, workers)

---

### 2. `internal/` Directory
**Purpose:** Private application code (can't be imported by other projects)

#### `internal/handlers/`
**Purpose:** HTTP request handlers

```go
// internal/handlers/tasks.go
package handlers

import (
    "net/http"
    "github.com/yourusername/go-todo-api/internal/models"
    "github.com/yourusername/go-todo-api/internal/database"
)

func GetAllTasks(w http.ResponseWriter, r *http.Request) {
    tasks, err := database.FetchAllTasks()
    if err != nil {
        http.Error(w, "Failed to fetch tasks", http.StatusInternalServerError)
        return
    }

    respondJSON(w, tasks)
}

func CreateTask(w http.ResponseWriter, r *http.Request) {
    // Handle task creation
}
```

**Benefits:**
- Each handler in its own function
- Easy to test individual handlers
- Clear separation of concerns

---

#### `internal/middleware/`
**Purpose:** Middleware functions

```go
// internal/middleware/logging.go
package middleware

import (
    "log"
    "net/http"
    "time"
)

func Logging(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        next.ServeHTTP(w, r)
        log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
    })
}
```

```go
// internal/middleware/cors.go
package middleware

import "net/http"

func CORS(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        // ... more CORS headers
        next.ServeHTTP(w, r)
    })
}
```

**Benefits:**
- Each middleware in its own file
- Easy to add/remove middleware
- Testable independently

---

#### `internal/models/`
**Purpose:** Data structures and business logic

```go
// internal/models/task.go
package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Task struct {
    ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    Title       string             `json:"title" validate:"required"`
    Description string             `json:"description"`
    Completed   bool               `json:"completed"`
    CreatedAt   time.Time          `json:"created_at"`
    UpdatedAt   time.Time          `json:"updated_at"`
}

// Validate checks if the task data is valid
func (t *Task) Validate() error {
    if t.Title == "" {
        return errors.New("title is required")
    }
    return nil
}
```

**Benefits:**
- Keep data structures separate from handlers
- Add validation methods to structs
- Easy to modify structure without touching handlers

---

#### `internal/database/`
**Purpose:** Database connection and queries

```go
// internal/database/mongo.go
package database

import (
    "context"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

var (
    client     *mongo.Client
    collection *mongo.Collection
)

func Connect() error {
    // Connection logic
}

func GetCollection() *mongo.Collection {
    return collection
}
```

```go
// internal/database/tasks.go
package database

import (
    "context"
    "github.com/yourusername/go-todo-api/internal/models"
)

func FetchAllTasks() ([]models.Task, error) {
    var tasks []models.Task
    cursor, err := collection.Find(context.Background(), bson.M{})
    if err != nil {
        return nil, err
    }
    cursor.All(context.Background(), &tasks)
    return tasks, nil
}

func CreateTask(task models.Task) error {
    _, err := collection.InsertOne(context.Background(), task)
    return err
}
```

**Benefits:**
- Database logic separate from HTTP handlers
- Easy to switch databases (MongoDB → PostgreSQL)
- Can mock database for testing

---

### 3. `pkg/` Directory
**Purpose:** Public, reusable code

```go
// pkg/utils/validator.go
package utils

import "regexp"

func IsValidEmail(email string) bool {
    pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
    matched, _ := regexp.MatchString(pattern, email)
    return matched
}
```

**When to use:**
- Code that other projects could use
- Generic utilities
- Helper functions

---

## Comparison: Before vs After

### Before (Single File)
```go
// main.go (530 lines)
package main

import (...)

type Task struct {...}

func loggingMiddleware(...) {...}
func corsMiddleware(...) {...}

func getAllTasksHandler(...) {...}
func createTaskHandler(...) {...}
// ... 15 more functions

func main() {...}
```

**Problem:** Everything is mixed together!

---

### After (Organized)

```go
// cmd/api/main.go (40 lines)
package main

func main() {
    config.Load()
    database.Connect()

    mux := http.NewServeMux()
    handlers.RegisterRoutes(mux)

    handler := middleware.Chain(mux)
    http.ListenAndServe(":8080", handler)
}
```

```go
// internal/models/task.go (30 lines)
package models

type Task struct {...}
func (t *Task) Validate() error {...}
```

```go
// internal/handlers/tasks.go (100 lines)
package handlers

func GetAllTasks(w, r) {...}
func GetTaskByID(w, r) {...}
func CreateTask(w, r) {...}
func UpdateTask(w, r) {...}
func DeleteTask(w, r) {...}
```

**Benefit:** Each file has ONE job!

---

## Industry Standard Structure (Large Projects)

For production apps at companies like Google, Uber, etc.:

```
go-todo-api/
├── cmd/
│   ├── api/                     # API server
│   │   └── main.go
│   ├── worker/                  # Background workers
│   │   └── main.go
│   └── migrate/                 # Database migrations
│       └── main.go
│
├── internal/
│   ├── api/                     # API-specific code
│   │   ├── handlers/
│   │   ├── middleware/
│   │   └── router/
│   │
│   ├── domain/                  # Business logic
│   │   ├── task/
│   │   │   ├── model.go         # Task struct
│   │   │   ├── service.go       # Business logic
│   │   │   └── repository.go    # Data access
│   │   └── user/
│   │       ├── model.go
│   │       ├── service.go
│   │       └── repository.go
│   │
│   └── infrastructure/          # External services
│       ├── database/
│       ├── cache/
│       └── email/
│
├── pkg/                         # Public libraries
│   ├── logger/
│   ├── validator/
│   └── errors/
│
├── configs/                     # Configuration files
│   ├── development.yaml
│   └── production.yaml
│
├── migrations/                  # Database migrations
│   ├── 001_create_tasks.sql
│   └── 002_add_users.sql
│
├── docs/                        # Documentation
│   ├── api.md
│   └── architecture.md
│
└── tests/                       # Integration tests
    └── integration_test.go
```

This is **advanced** - don't worry about this complexity yet!

---

## When to Refactor Your Structure

### Keep Single File When:
- Learning Go fundamentals ✅ (you're here!)
- Prototyping quickly
- Code is under 500 lines
- Solo project

### Split into Multiple Files When:
- Code reaches 500+ lines
- Hard to find functions
- Multiple people working on code
- Want to add tests

### Use Package Structure When:
- Production application
- Team of 2+ developers
- Need to test thoroughly
- Code reuse is important
- Multiple related applications

---

## Practical Example: Refactoring Your API

Let's say you wanted to refactor your current `main.go`. Here's what it would look like:

### Step 1: Create folders
```bash
mkdir -p internal/handlers
mkdir -p internal/middleware
mkdir -p internal/models
mkdir -p internal/database
```

### Step 2: Move code

**From main.go → internal/models/task.go:**
```go
package models

type Task struct {
    ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    Title       string             `json:"title"`
    Description string             `json:"description"`
    Completed   bool               `json:"completed"`
}
```

**From main.go → internal/middleware/logging.go:**
```go
package middleware

func Logging(next http.Handler) http.Handler {
    // Your logging middleware code
}
```

**From main.go → internal/handlers/tasks.go:**
```go
package handlers

func GetAllTasks(w http.ResponseWriter, r *http.Request) {
    // Your getAllTasksHandler code
}
```

### Step 3: Update main.go
```go
package main

import (
    "github.com/yourusername/go-todo-api/internal/handlers"
    "github.com/yourusername/go-todo-api/internal/middleware"
)

func main() {
    mux := http.NewServeMux()

    mux.HandleFunc("/tasks", handlers.GetAllTasks)

    handler := middleware.Logging(mux)

    http.ListenAndServe(":8080", handler)
}
```

---

## Best Practices

### 1. Group by Feature (Domain-Driven)
```
internal/
├── task/           # Everything task-related
│   ├── handler.go
│   ├── service.go
│   ├── repository.go
│   └── model.go
└── user/           # Everything user-related
    ├── handler.go
    ├── service.go
    └── model.go
```

**Benefit:** All code for a feature is together

### 2. Group by Layer (MVC-style)
```
internal/
├── handlers/       # All handlers
├── services/       # All business logic
├── repositories/   # All data access
└── models/         # All data structures
```

**Benefit:** Easy to find code by type

---

## Key Takeaways

1. **Start simple** - One file is fine for learning
2. **Refactor when needed** - Split into packages as you grow
3. **Use `internal/`** - For private application code
4. **Use `pkg/`** - For reusable libraries
5. **One responsibility per file** - Don't mix concerns
6. **Keep main.go small** - It should just wire things together
7. **Test as you go** - Good structure makes testing easier

---

## Your Next Steps

For now, **keep your single-file structure** while learning. When you're ready to refactor:

1. Start with separating into multiple files (Stage 2)
2. Then move to packages when you're comfortable (Stage 3)
3. Don't over-engineer early - simplicity is better

Remember: **The best structure is the one that makes sense for YOUR project right now!**

---

## Resources

- [Go Project Layout (Standard)](https://github.com/golang-standards/project-layout)
- [Organizing Go Code (Official Blog)](https://go.dev/blog/organizing-go-code)
- [Go by Example: Packages](https://gobyexample.com/packages)
