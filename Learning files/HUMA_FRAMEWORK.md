# Huma Framework - Modern REST API Development

## What is Huma?

Huma is a modern REST API framework for Go that eliminates boilerplate code and provides automatic OpenAPI documentation generation. It's designed for building production-ready APIs with less code and better developer experience.

**Official Site:** https://huma.rocks
**GitHub:** https://github.com/danielgtaylor/huma

## Why Use Huma?

### Before Huma (Standard Library):

```go
func CreateTask(w http.ResponseWriter, r *http.Request) {
    var task Task

    // Manual JSON decoding
    if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }

    // Manual validation
    if task.Title == "" {
        http.Error(w, "Title is required", http.StatusBadRequest)
        return
    }
    if len(task.Title) > 200 {
        http.Error(w, "Title too long", http.StatusBadRequest)
        return
    }

    // Business logic
    result, err := db.Insert(task)
    if err != nil {
        http.Error(w, "Failed to create", http.StatusInternalServerError)
        return
    }

    // Manual JSON encoding
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(result)
}
```

### After Huma:

```go
func CreateTask(ctx context.Context, input *CreateTaskInput) (*CreateTaskOutput, error) {
    // Input already validated and decoded!
    newTask := Task{
        Title:       input.Body.Title,
        Description: input.Body.Description,
        Completed:   false,
    }

    // Business logic
    result, err := db.Insert(newTask)
    if err != nil {
        return nil, huma.Error500InternalServerError("Failed to create")
    }

    // Return output (auto-encoded to JSON)
    return &CreateTaskOutput{Body: result}, nil
}
```

**What Huma Does Automatically:**
- ✅ Validates input against struct tags
- ✅ Decodes request body to struct
- ✅ Encodes response to JSON
- ✅ Generates OpenAPI documentation
- ✅ Provides JSON schemas
- ✅ Handles errors properly

## Key Features

### 1. Automatic OpenAPI 3.1 Documentation

Huma generates complete OpenAPI documentation from your code:

```go
// This single registration...
huma.Register(api, huma.Operation{
    OperationID: "create-task",
    Method:      http.MethodPost,
    Path:        "/tasks",
    Summary:     "Create a new task",
    Description: "Add a new TODO task to the database",
    Tags:        []string{"Tasks"},
}, handlers.CreateTask)

// ...generates OpenAPI spec with:
// - Request/response schemas
// - Validation rules
// - Example values
// - Error responses
```

**Access OpenAPI Docs:**
- **Interactive UI:** http://localhost:8080/docs
- **JSON spec:** http://localhost:8080/openapi.json
- **YAML spec:** http://localhost:8080/openapi.yaml

### 2. Automatic Request Validation

Define validation in struct tags - Huma handles enforcement:

```go
type CreateTaskInput struct {
    Body struct {
        Title       string `json:"title" minLength:"1" maxLength:"200" doc:"Title of the task"`
        Description string `json:"description,omitempty" maxLength:"1000"`
    }
}
```

**Validation Rules:**
- `minLength` / `maxLength` - String length constraints
- `minimum` / `maximum` - Number constraints
- `pattern` - Regex validation
- `enum` - Allowed values
- `required` - Required fields (fields without `omitempty`)
- `format` - Special formats (email, uri, date-time, etc.)

### 3. JSON Schema Generation

Every type gets a JSON Schema automatically:

```json
{
    "$schema": "http://localhost:8080/schemas/Task.json",
    "id": "6900f242845ebd7da9c8e24f",
    "title": "Test Huma API",
    "completed": false
}
```

Clients can fetch schemas for validation and code generation.

### 4. RFC 7807 Problem Details

Errors follow the standard Problem Details format:

```json
{
    "type": "about:blank",
    "title": "Bad Request",
    "status": 400,
    "detail": "Validation failed: title must be at least 1 characters",
    "instance": "/tasks",
    "errors": [
        {
            "location": "body.title",
            "message": "must be at least 1 characters",
            "value": ""
        }
    ]
}
```

### 5. Type-Safe Handlers

Handlers use Go's type system for safety:

```go
// Input and output types are explicit
func GetTask(ctx context.Context, input *GetTaskInput) (*GetTaskOutput, error) {
    // input.ID is a string (from path parameter)
    // Return value must match GetTaskOutput
    // Errors are returned, not written to ResponseWriter
}
```

## Huma Handler Pattern

### Standard Signature

```go
func HandlerName(
    ctx context.Context,    // Request context
    input *InputType,       // Parsed and validated input
) (*OutputType, error) {    // Response or error
    // Your logic here
}
```

### Input Structure

```go
type GetTaskInput struct {
    // Path parameters
    ID string `path:"id" doc:"Task ID" minLength:"24" maxLength:"24"`

    // Query parameters
    Filter string `query:"filter" doc:"Filter tasks"`

    // Headers
    Authorization string `header:"Authorization"`

    // Request body
    Body struct {
        Title string `json:"title" minLength:"1"`
    }
}
```

### Output Structure

```go
type GetTaskOutput struct {
    // Response body
    Body Task

    // Custom headers (optional)
    ETag string `header:"ETag"`
}
```

## Migration Guide

### Step 1: Add Dependencies

```bash
go get github.com/danielgtaylor/huma/v2
go get github.com/danielgtaylor/huma/v2/adapters/humachi
go get github.com/go-chi/chi/v5
```

### Step 2: Create Input/Output Types

```go
// Before (no types needed)
func GetTask(w http.ResponseWriter, r *http.Request) {...}

// After (define types)
type GetTaskInput struct {
    ID string `path:"id" minLength:"24" maxLength:"24"`
}

type GetTaskOutput struct {
    Body Task
}

func GetTask(ctx context.Context, input *GetTaskInput) (*GetTaskOutput, error) {...}
```

### Step 3: Update Handler Logic

```go
// Before
func GetTask(w http.ResponseWriter, r *http.Request) {
    id := r.URL.Query().Get("id")
    if id == "" {
        http.Error(w, "ID required", http.StatusBadRequest)
        return
    }

    task, err := db.FindByID(id)
    if err != nil {
        http.Error(w, "Not found", http.StatusNotFound)
        return
    }

    json.NewEncoder(w).Encode(task)
}

// After
func GetTask(ctx context.Context, input *GetTaskInput) (*GetTaskOutput, error) {
    // input.ID is already validated!
    task, err := db.FindByID(input.ID)
    if err != nil {
        return nil, huma.Error404NotFound("Task not found")
    }

    return &GetTaskOutput{Body: task}, nil
}
```

### Step 4: Register Operations

```go
// Create Huma API
router := chi.NewMux()
api := humachi.New(router, huma.DefaultConfig("My API", "1.0.0"))

// Register operations
huma.Register(api, huma.Operation{
    OperationID: "get-task",
    Method:      http.MethodGet,
    Path:        "/tasks/{id}",
    Summary:     "Get a task by ID",
    Tags:        []string{"Tasks"},
}, handlers.GetTask)
```

## Common Patterns

### Pagination

```go
type ListTasksInput struct {
    Page  int `query:"page" minimum:"1" default:"1"`
    Limit int `query:"limit" minimum:"1" maximum:"100" default:"20"`
}

type ListTasksOutput struct {
    Body struct {
        Tasks []Task `json:"tasks"`
        Total int    `json:"total"`
        Page  int    `json:"page"`
    }
}
```

### Authentication

```go
type AuthenticatedInput struct {
    Token string `header:"Authorization" doc:"Bearer token"`
}

func ProtectedHandler(ctx context.Context, input *AuthenticatedInput) (*Output, error) {
    if !validateToken(input.Token) {
        return nil, huma.Error401Unauthorized("Invalid token")
    }
    // ...
}
```

### File Uploads

```go
type UploadInput struct {
    Body struct {
        File []byte `json:"file" contentEncoding:"base64"`
    }
}
```

### Conditional Responses

```go
type GetTaskOutput struct {
    Body  Task
    ETag  string `header:"ETag"`
}

func GetTask(ctx context.Context, input *GetTaskInput) (*GetTaskOutput, error) {
    task, err := db.FindByID(input.ID)
    if err != nil {
        return nil, huma.Error404NotFound("Task not found")
    }

    etag := generateETag(task)
    return &GetTaskOutput{
        Body: task,
        ETag: etag,
    }, nil
}
```

## Error Handling

### Standard Errors

```go
// 400 Bad Request
return nil, huma.Error400BadRequest("Invalid input")

// 401 Unauthorized
return nil, huma.Error401Unauthorized("Authentication required")

// 403 Forbidden
return nil, huma.Error403Forbidden("Access denied")

// 404 Not Found
return nil, huma.Error404NotFound("Resource not found")

// 500 Internal Server Error
return nil, huma.Error500InternalServerError("Database error")
```

### Custom Errors

```go
return nil, &huma.ErrorModel{
    Status: 422,
    Title:  "Unprocessable Entity",
    Detail: "Task cannot be completed because...",
    Errors: []huma.ErrorDetail{
        {
            Location: "body.status",
            Message:  "Invalid status transition",
            Value:    "pending",
        },
    },
}
```

## Documentation Tips

### Use Descriptive Tags

```go
type Task struct {
    ID          string `json:"id" doc:"Unique task identifier"`
    Title       string `json:"title" doc:"Task title (1-200 chars)" minLength:"1" maxLength:"200"`
    Description string `json:"description,omitempty" doc:"Optional detailed description" maxLength:"1000"`
    Completed   bool   `json:"completed" doc:"Whether the task is completed"`
}
```

### Add Examples

```go
type CreateTaskInput struct {
    Body struct {
        Title       string `json:"title" example:"Buy groceries"`
        Description string `json:"description" example:"Buy milk, eggs, and bread"`
    }
}
```

### Group Operations with Tags

```go
huma.Register(api, huma.Operation{
    // ...
    Tags: []string{"Tasks", "CRUD"},  // Shows up in doc groups
}, handler)
```

### Set Response Status

```go
huma.Register(api, huma.Operation{
    // ...
    DefaultStatus: http.StatusCreated,  // 201 instead of 200
}, handlers.CreateTask)
```

## Testing

### Test Operations Directly

```go
func TestCreateTask(t *testing.T) {
    input := &models.CreateTaskInput{
        Body: struct {
            Title string `json:"title"`
        }{
            Title: "Test Task",
        },
    }

    output, err := handlers.CreateTask(context.Background(), input)

    if err != nil {
        t.Fatalf("Expected no error, got %v", err)
    }
    if output.Body.Title != "Test Task" {
        t.Errorf("Expected 'Test Task', got '%s'", output.Body.Title)
    }
}
```

### Test Validation

```go
func TestValidation(t *testing.T) {
    input := &models.CreateTaskInput{
        Body: struct {
            Title string `json:"title"`
        }{
            Title: "",  // Invalid: minLength:1
        },
    }

    _, err := handlers.CreateTask(context.Background(), input)
    // Huma validates before calling the handler
    // In real usage, validation happens at the framework level
}
```

## Performance Considerations

### Huma is Fast

- Zero reflection in the hot path
- Efficient JSON parsing
- OpenAPI docs generated once at startup
- Schemas cached

### Optimization Tips

1. **Reuse structs** - Don't create new types for similar operations
2. **Use pointers for large structs** - Avoid copying
3. **Enable compression** - Use middleware
4. **Cache responses** - Add ETag support

## Comparison: Huma vs Others

| Feature | Huma | Gin | Echo | Fiber | stdlib |
|---------|------|-----|------|-------|--------|
| OpenAPI Generation | ✅ Auto | ❌ Manual | ❌ Manual | ❌ Manual | ❌ Manual |
| Request Validation | ✅ Built-in | 🔶 Plugin | 🔶 Plugin | 🔶 Plugin | ❌ Manual |
| Type Safety | ✅ Strong | 🔶 Weak | 🔶 Weak | 🔶 Weak | 🔶 Weak |
| JSON Schema | ✅ Auto | ❌ No | ❌ No | ❌ No | ❌ No |
| Error Standards | ✅ RFC 7807 | ❌ Custom | ❌ Custom | ❌ Custom | ❌ Custom |
| Learning Curve | 🔶 Medium | ✅ Easy | ✅ Easy | ✅ Easy | ✅ Easy |
| Performance | ✅ Fast | ✅ Fast | ✅ Fast | ✅ Fastest | ✅ Fast |

## Best Practices

1. **Use descriptive operation IDs** - They become function names in generated clients
2. **Tag operations logically** - Groups operations in documentation
3. **Add examples** - Makes API documentation more useful
4. **Document errors** - Explain what each error means
5. **Version your API** - Use path versioning (`/v1/tasks`)
6. **Use semantic status codes** - 201 for creation, 204 for deletion
7. **Implement pagination** - For list endpoints
8. **Add rate limiting** - Protect your API
9. **Use ETag/If-Match** - For optimistic concurrency
10. **Enable CORS properly** - Don't use `*` in production

## Resources

- **Official Docs:** https://huma.rocks/
- **GitHub:** https://github.com/danielgtaylor/huma
- **Examples:** https://github.com/danielgtaylor/huma/tree/main/examples
- **OpenAPI Spec:** https://spec.openapis.org/oas/v3.1.0

## Summary

Huma is perfect for:
- ✅ Building production APIs quickly
- ✅ Teams that value documentation
- ✅ APIs that need client generation
- ✅ Projects with strict validation requirements
- ✅ Modern RESTful services

Consider alternatives if:
- ❌ You need absolute maximum performance (use Fiber)
- ❌ You want minimal dependencies (use stdlib)
- ❌ You're building a simple prototype (any framework works)

For most production REST APIs, Huma provides the best balance of developer experience, safety, and automatic documentation.
