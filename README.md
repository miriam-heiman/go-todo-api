# Go To-Do REST API

A simple REST API built with Go to learn the language. This project demonstrates building a CRUD API with MongoDB database.

## 🚀 Features

- **HTTP Server** - Runs on port 8080
- **CRUD Operations** - Create, Read, Update, Delete tasks
- **JSON API** - Returns data in JSON format
- **Error Handling** - Proper HTTP status codes
- **MongoDB Integration** - Cloud database storage with MongoDB Atlas

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
   # Copy the example file
   cp .env.example .env
   
   # Edit .env and add your MongoDB connection string
   # Get it from: https://cloud.mongodb.com
   ```

4. Run the server:
   ```bash
   go run main.go
   ```

## 📖 Usage

### Start the Server
```bash
go run main.go
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

## 📚 Learning Resources

Check out the `Learning files/` directory for detailed explanations:
- `CODE_STRUCTURE.md` - Understanding Go file organization
- `DEPENDENCIES.md` - Understanding go.mod and go.sum
- `MONGODB_SETUP.md` - How to add MongoDB
- `SUMMARY.md` - Complete project overview
- `TESTING.md` - How to test the API

## 🛠 Tech Stack

- **Go** - Programming language
- **net/http** - HTTP server and client
- **encoding/json** - JSON encoding/decoding
- **strconv** - String conversion utilities

## 📝 Project Structure

```
go-todo-api/
├── main.go              # Main application code
├── go.mod               # Go module dependencies
├── go.sum               # Dependency checksums
├── README.md            # This file
└── Learning files/      # Learning resources
    ├── CODE_STRUCTURE.md
    ├── DEPENDENCIES.md
    ├── MONGODB_SETUP.md
    ├── SUMMARY.md
    └── TESTING.md
```

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

- [ ] Add MongoDB for persistent storage
- [ ] Add user authentication
- [ ] Add task categories/tags
- [ ] Add due dates
- [ ] Add filtering and sorting
- [ ] Add pagination
- [ ] Add middleware (logging, CORS)
- [ ] Add unit tests

## 📄 License

This project is for educational purposes.

