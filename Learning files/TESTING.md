# Testing Your To-Do API

Your server is running at **http://localhost:8080**

## Quick Tests with curl

Open your terminal and try these commands:

### 1. Visit the Homepage
```bash
curl http://localhost:8080/
```

### 2. Get All Tasks
```bash
curl http://localhost:8080/tasks
```

### 3. Get a Specific Task by ID
```bash
curl http://localhost:8080/tasks?id=1
```

### 4. Create a New Task (POST)
```bash
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Build my first API",
    "description": "Learn Go by building this to-do API",
    "completed": false
  }'
```

### 5. Update a Task (PUT)
```bash
curl -X PUT http://localhost:8080/tasks?id=1 \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Learn Go Basics",
    "completed": true
  }'
```

### 6. Delete a Task
```bash
curl -X DELETE http://localhost:8080/tasks?id=1
```

### 7. Health Check
```bash
curl http://localhost:8080/health
```

## Or Use Your Browser

1. Visit **http://localhost:8080/** - See the homepage
2. Visit **http://localhost:8080/tasks** - Get all tasks
3. Visit **http://localhost:8080/health** - Health check

## What You Just Built!

✅ **HTTP Server** - Listens on port 8080
✅ **GET All Tasks** - Returns all tasks as JSON
✅ **GET Task by ID** - Returns a specific task
✅ **POST Create Task** - Creates a new task
✅ **PUT Update Task** - Updates an existing task
✅ **DELETE Task** - Deletes a task
✅ **Error Handling** - Proper HTTP status codes
✅ **JSON Encoding/Decoding** - Converts Go structs to/from JSON

## Key Go Concepts You Learned

1. **HTTP Handlers** - Functions that handle web requests
2. **JSON Encoding/Decoding** - Converting data to/from JSON
3. **Error Handling** - Checking for and handling errors
4. **Switch Statements** - Routing based on HTTP method
5. **Slice Manipulation** - Adding, removing, updating items in slices
6. **URL Query Parameters** - Getting data from URLs
7. **Structs** - Defining custom data types
8. **Pointers** - Referencing variables with `&`

Congrats! You've built your first API in Go! 🎉

