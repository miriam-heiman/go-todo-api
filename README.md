# Go To-Do REST API

A production-ready REST API built with Go using the Huma framework. This project demonstrates modern API development with automatic OpenAPI documentation, request validation, and MongoDB database storage.

## 🚀 Features

- **Huma Framework** - Modern REST API framework with automatic OpenAPI 3.1 documentation
- **Interactive API Docs** - Swagger-like UI at `/docs`
- **Automatic Validation** - Request/response validation using struct tags
- **CRUD Operations** - Create, Read, Update, Delete tasks
- **JSON Schema** - Automatic schema generation for all types
- **MongoDB Integration** - Persistent database storage with local MongoDB
- **Middleware** - Logging and CORS support
- **Hot Reload** - Air for automatic server restart on code changes
- **Production Structure** - Clean `cmd/` and `internal/` package organization
- **RFC 7807 Errors** - Standard problem details for errors

## 📦 Installation

1. Install Go (if not already installed):
   ```bash
   brew install go
   ```

2. Clone this repository:
   ```bash
   git clone <repository-url>
   cd go-todo-api
   ```

3. Set up environment variables:
   ```bash
   # Copy the example .env file
   cp .env.example .env
   
   # Edit .env and add your MongoDB connection string
   # Get it from MongoDB Atlas: https://www.mongodb.com/cloud/atlas
   ```

4. Run the server:
   ```bash
   # With hot-reload (recommended for development)
   air

   # Or manually
   go run cmd/api/main.go
   ```

## 📖 Usage

### Start the Server
```bash
# With hot-reload (automatically restarts on code changes)
air

# Or manually
go run cmd/api/main.go
```

The server will start on `http://localhost:8080`

### API Endpoints

#### Get All Tasks
```bash
curl http://localhost:8080/tasks
```

#### Get Task by ID
```bash
curl http://localhost:8080/tasks?id=1
```

#### Create a Task
```bash
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{"title": "My Task", "description": "Task description"}'
```

#### Update a Task
```bash
curl -X PUT http://localhost:8080/tasks?id=1 \
  -H "Content-Type: application/json" \
  -d '{"title": "Updated Task", "completed": true}'
```

#### Delete a Task
```bash
curl -X DELETE http://localhost:8080/tasks?id=1
```

#### Health Check
```bash
curl http://localhost:8080/health
```

## 📚 API Documentation

This API includes automatic interactive documentation:

- **Interactive Docs:** http://localhost:8080/docs
- **OpenAPI JSON:** http://localhost:8080/openapi.json
- **OpenAPI YAML:** http://localhost:8080/openapi.yaml

The documentation is generated automatically from code and includes:
- Request/response schemas
- Validation rules
- Example values
- Try-it-out functionality

## 📚 Learning Resources

Check out the `Learning files/` directory for detailed explanations:
- `CODE_STRUCTURE.md` - Understanding Go file organization
- `DEPENDENCIES.md` - Understanding go.mod and go.sum
- `MONGODB_SETUP.md` - How to add MongoDB
- `MIDDLEWARE_EXPLAINED.md` - Understanding middleware
- `API_FILE_STRUCTURE.md` - Production API organization
- **`HUMA_FRAMEWORK.md` - Modern REST API framework guide**
- `SUMMARY.md` - Complete project overview
- `TESTING.md` - How to test the API

## 🛠 Tech Stack

- **Go** - Programming language
- **Huma v2** - Modern REST API framework with OpenAPI generation
- **Chi Router** - HTTP router for Huma
- **MongoDB** - Database for persistent storage
- **godotenv** - Environment variable management
- **Air** - Hot-reload development tool

## 📝 Project Structure

```
go-todo-api/
├── cmd/
│   └── api/
│       └── main.go          # Application entry point
├── internal/                # Private application code
│   ├── handlers/            # HTTP request handlers
│   │   ├── home.go
│   │   ├── health.go
│   │   └── tasks.go
│   ├── middleware/          # Middleware functions
│   │   ├── logging.go
│   │   ├── cors.go
│   │   └── chain.go
│   ├── models/              # Data structures
│   │   └── task.go
│   ├── database/            # Database connections
│   │   └── mongo.go
│   └── config/              # Configuration
├── Learning files/          # Learning resources
│   ├── CODE_STRUCTURE.md
│   ├── DEPENDENCIES.md
│   ├── MONGODB_SETUP.md
│   ├── MIDDLEWARE_EXPLAINED.md
│   ├── API_FILE_STRUCTURE.md
│   ├── SUMMARY.md
│   └── TESTING.md
├── go.mod                   # Go module dependencies
├── go.sum                   # Dependency checksums
├── .env                     # Environment variables (gitignored)
├── .env.example             # Example environment file
├── .air.toml                # Hot-reload configuration
└── README.md                # This file
```

**Production-Ready Structure:**
- `cmd/` - Application entry points
- `internal/` - Private application code (can't be imported by other projects)
- Clean separation of concerns (handlers, models, middleware, database)

## 🎓 What I Learned

This project taught me:
- Go package structure and organization
- HTTP handlers and routing
- JSON encoding/decoding
- Error handling patterns
- Struct definitions
- Slice manipulation
- Switch statements
- URL query parameters

## 🔮 Future Improvements

- [x] Add MongoDB for persistent storage
- [x] Add environment variable support
- [ ] Add user authentication
- [ ] Add task categories/tags
- [ ] Add due dates
- [ ] Add filtering and sorting
- [ ] Add pagination
- [ ] Add middleware (logging, CORS)
- [ ] Add unit tests

## 📄 License

This project is for educational purposes.

