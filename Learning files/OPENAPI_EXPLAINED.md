# OpenAPI Documentation Explained (For Beginners)

## What is OpenAPI?

**OpenAPI** (formerly called "Swagger") is a **standard format for describing REST APIs**. Think of it as an "instruction manual" or "contract" for your API that both humans and computers can understand.

### Real-World Analogy

Imagine you're giving someone directions to your house:
- **Without OpenAPI**: You verbally tell them the directions (easy to misunderstand, forget, or get wrong)
- **With OpenAPI**: You give them a GPS coordinate, address, and detailed map (precise, standardized, machine-readable)

OpenAPI does the same thing for your API - it provides a precise, standardized description that tools and people can use.

---

## Why OpenAPI is Amazing

### 1. **Interactive Documentation** 📚
Users can read about your API **AND** test it directly in their browser (no Postman needed!)

**Your API's Interactive Docs**: http://localhost:8080/docs

Try it now:
1. Open http://localhost:8080/docs in your browser
2. You'll see all your endpoints listed
3. Click "GET /tasks" → "Try it out" → "Execute"
4. The documentation runs the request and shows the response!

### 2. **Automatic Generation** 🤖
Huma creates OpenAPI documentation **automatically** from your code. You don't write it manually!

**How?** Huma reads your Go code:
```go
type CreateTaskInput struct {
    Body struct {
        Title       string `json:"title" minLength:"1" maxLength:"200" doc:"Title of the task"`
        Description string `json:"description,omitempty" maxLength:"1000" doc:"Detailed description"`
    }
}
```

And generates OpenAPI spec:
```json
{
  "title": {
    "description": "Title of the task",
    "minLength": 1,
    "maxLength": 200,
    "type": "string"
  }
}
```

### 3. **Client Code Generation** 🛠️
Tools can read your OpenAPI spec and **auto-generate client libraries** in any programming language!

Example: Generate a TypeScript client for your frontend:
```bash
npx @openapitools/openapi-generator-cli generate \
  -i http://localhost:8080/openapi.json \
  -g typescript-fetch \
  -o ./generated-client
```

Now your frontend has type-safe functions like:
```typescript
import { TasksApi } from './generated-client';

const api = new TasksApi();
const tasks = await api.listTasks(); // Fully typed!
```

### 4. **Testing Tools** 🧪
Import your OpenAPI spec into tools like:
- **Postman**: Click "Import" → Paste http://localhost:8080/openapi.json
- **Insomnia**: Same process
- **Bruno**: Import from URL
- **Thunder Client** (VS Code): Import collection

All your endpoints are instantly available with examples!

### 5. **Contract-First Development** 🤝
Frontend and backend teams can work **in parallel**:
- Backend publishes OpenAPI spec
- Frontend generates client and builds UI
- Both teams work independently using the "contract"

---

## Your API's OpenAPI Endpoints

Your TODO API has **3 OpenAPI endpoints** automatically created by Huma:

### 1. Interactive Documentation (Swagger UI)
**URL**: http://localhost:8080/docs

**What you see**:
- Beautiful web interface
- All endpoints organized by tags (System, Tasks)
- Click any endpoint to see details
- "Try it out" buttons to test endpoints live
- Request/response examples

**When to use**:
- Learning the API
- Manual testing
- Sharing with teammates
- API demos

### 2. OpenAPI JSON Format
**URL**: http://localhost:8080/openapi.json

**What you get**: Machine-readable JSON describing your entire API

**When to use**:
- Importing into Postman/Insomnia
- Generating client libraries
- CI/CD pipelines that validate APIs
- API gateways (Kong, AWS API Gateway)

### 3. OpenAPI YAML Format
**URL**: http://localhost:8080/openapi.yaml

**What you get**: Same as JSON, but in YAML format (more human-readable)

**When to use**:
- When you prefer YAML over JSON
- Some tools prefer YAML (like Stoplight)
- Committing to git (YAML diffs are easier to read)

---

## Understanding Your API's OpenAPI Spec

Let's break down what Huma generated for your TODO API:

### Basic Info

```json
{
  "openapi": "3.1.0",
  "info": {
    "title": "TODO API",
    "description": "A production-ready REST API for managing TODO tasks",
    "version": "1.0.0",
    "contact": {
      "name": "Your Name",
      "url": "https://github.com/yourusername/go-todo-api"
    }
  }
}
```

**Where this comes from**: Your `main.go` file:
```go
config := huma.DefaultConfig("TODO API", "1.0.0")
config.Info.Description = "A production-ready REST API for managing TODO tasks"
```

### Endpoints (Paths)

```json
{
  "paths": {
    "/tasks": {
      "get": {
        "summary": "List all tasks",
        "description": "Retrieve all TODO tasks from the database",
        "operationId": "list-tasks",
        "tags": ["Tasks"],
        "responses": {
          "200": {
            "description": "OK",
            "content": {
              "application/json": {
                "schema": {
                  "type": "array",
                  "items": { "$ref": "#/components/schemas/Task" }
                }
              }
            }
          }
        }
      }
    }
  }
}
```

**Where this comes from**: Your `main.go` registration:
```go
huma.Register(api, huma.Operation{
    OperationID: "list-tasks",
    Method:      http.MethodGet,
    Path:        "/tasks",
    Summary:     "List all tasks",
    Description: "Retrieve all TODO tasks from the database",
    Tags:        []string{"Tasks"},
}, handlers.GetAllTasks)
```

### Data Models (Schemas)

```json
{
  "components": {
    "schemas": {
      "Task": {
        "type": "object",
        "required": ["id", "title", "completed"],
        "properties": {
          "id": {
            "type": "string",
            "description": "Unique identifier for the task"
          },
          "title": {
            "type": "string",
            "minLength": 1,
            "maxLength": 200,
            "description": "Title of the task"
          },
          "description": {
            "type": "string",
            "maxLength": 1000,
            "description": "Detailed description of the task"
          },
          "completed": {
            "type": "boolean",
            "description": "Whether the task is completed"
          }
        }
      }
    }
  }
}
```

**Where this comes from**: Your `models/task.go` struct:
```go
type Task struct {
    ID          primitive.ObjectID `json:"id" doc:"Unique identifier for the task"`
    Title       string             `json:"title" minLength:"1" maxLength:"200" doc:"Title of the task"`
    Description string             `json:"description,omitempty" maxLength:"1000" doc:"Detailed description of the task"`
    Completed   bool               `json:"completed" doc:"Whether the task is completed"`
}
```

**See the magic?** Huma automatically converts:
- `doc:"..."` → `"description": "..."`
- `minLength:"1"` → `"minLength": 1`
- `maxLength:"200"` → `"maxLength": 200`
- `omitempty` → makes field optional in schema

---

## How to Use OpenAPI Documentation

### For API Consumers (Users of Your API)

1. **Explore the API**
   ```bash
   # Open in browser
   open http://localhost:8080/docs
   ```

2. **Test Endpoints**
   - Click any endpoint (e.g., "GET /tasks")
   - Click "Try it out"
   - Click "Execute"
   - See the real response!

3. **Copy curl Commands**
   - After executing a request, scroll down
   - Copy the "curl" command
   - Run it in your terminal

4. **See Request/Response Examples**
   - Every endpoint shows example requests
   - Example responses for success and errors
   - All possible status codes documented

### For Developers (Building Clients)

1. **Generate Client Library**
   ```bash
   # TypeScript client
   npx @openapitools/openapi-generator-cli generate \
     -i http://localhost:8080/openapi.json \
     -g typescript-fetch \
     -o ./sdk

   # Python client
   npx @openapitools/openapi-generator-cli generate \
     -i http://localhost:8080/openapi.json \
     -g python \
     -o ./sdk

   # Go client
   npx @openapitools/openapi-generator-cli generate \
     -i http://localhost:8080/openapi.json \
     -g go \
     -o ./sdk
   ```

2. **Import into Testing Tools**
   - **Postman**: Import → Link → http://localhost:8080/openapi.json
   - **Insomnia**: Create → Import from URL
   - **Thunder Client**: Collections → Import → From URL

3. **Validate API Changes**
   ```bash
   # Install openapi-diff
   npm install -g openapi-diff

   # Compare old vs new spec
   openapi-diff http://localhost:8080/openapi.json ./old-spec.json
   ```

---

## OpenAPI Spec Structure (Simplified)

Here's the basic structure of an OpenAPI 3.1 specification:

```yaml
openapi: 3.1.0                    # OpenAPI version
info:                              # API metadata
  title: TODO API
  version: 1.0.0
  description: A REST API for tasks

servers:                           # Where the API is hosted
  - url: http://localhost:8080

paths:                             # All endpoints
  /tasks:
    get:                           # GET /tasks
      summary: List tasks
      responses:
        200:
          description: Success
          content:
            application/json:
              schema:
                type: array
                items:
                  $ref: '#/components/schemas/Task'

    post:                          # POST /tasks
      summary: Create task
      requestBody:
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/CreateTaskInput'
      responses:
        201:
          description: Created

  /tasks/{id}:
    get:                           # GET /tasks/123
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: string
      responses:
        200:
          description: Success
        404:
          description: Not found

components:                        # Reusable schemas
  schemas:
    Task:
      type: object
      required: [id, title]
      properties:
        id:
          type: string
        title:
          type: string
        completed:
          type: boolean
```

---

## Real-World Example: Your Task Model

Let's trace how your Task model becomes OpenAPI documentation:

### 1. You Write Go Code

```go
// internal/models/task.go
type Task struct {
    ID          primitive.ObjectID `json:"id" doc:"Unique identifier for the task"`
    Title       string             `json:"title" minLength:"1" maxLength:"200" doc:"Title of the task"`
    Description string             `json:"description,omitempty" maxLength:"1000" doc:"Detailed description of the task"`
    Completed   bool               `json:"completed" doc:"Whether the task is completed"`
}
```

### 2. Huma Reads Your Struct Tags

- `json:"id"` → Field name in JSON will be "id"
- `doc:"..."` → Description in OpenAPI
- `minLength:"1"` → Validation rule in OpenAPI
- `maxLength:"200"` → Maximum length in OpenAPI
- `omitempty` → Field is optional (not required)

### 3. Huma Generates OpenAPI Schema

```json
{
  "Task": {
    "type": "object",
    "required": ["id", "title", "completed"],
    "properties": {
      "id": {
        "type": "string",
        "description": "Unique identifier for the task"
      },
      "title": {
        "type": "string",
        "minLength": 1,
        "maxLength": 200,
        "description": "Title of the task"
      },
      "description": {
        "type": "string",
        "maxLength": 1000,
        "description": "Detailed description of the task"
      },
      "completed": {
        "type": "boolean",
        "description": "Whether the task is completed"
      }
    }
  }
}
```

### 4. OpenAPI Tools Use This Schema

**Postman**: Shows field types, validation rules, and descriptions

**Code Generators**: Create type-safe clients:
```typescript
// Generated TypeScript interface
interface Task {
  id: string;
  title: string; // min: 1, max: 200
  description?: string; // optional, max: 1000
  completed: boolean;
}
```

**Validators**: Check if requests match the schema

---

## Common OpenAPI Terms Explained

| Term | What It Means | Example |
|------|---------------|---------|
| **Path** | An API endpoint URL | `/tasks`, `/tasks/{id}` |
| **Operation** | An HTTP method on a path | `GET /tasks`, `POST /tasks` |
| **Parameter** | Data sent in URL or headers | `id` in `/tasks/{id}` |
| **Request Body** | Data sent in POST/PUT request | Task JSON in `POST /tasks` |
| **Response** | What the API returns | `200: Task`, `404: Not Found` |
| **Schema** | Data structure definition | The Task object structure |
| **Component** | Reusable definition | Task schema used multiple times |
| **Tag** | Group of related operations | All task operations tagged "Tasks" |
| **operationId** | Unique identifier for operation | `"list-tasks"`, `"create-task"` |

---

## Benefits of Huma's Auto-Generated OpenAPI

### What You Get For Free:

✅ **Automatic Documentation**: Every endpoint documented
✅ **Request Validation**: Huma validates against the schema
✅ **Type Safety**: Schemas match your Go structs exactly
✅ **Error Responses**: Standard RFC 7807 error format
✅ **Interactive Testing**: Swagger UI built-in
✅ **Always Up-to-Date**: Docs update when code changes
✅ **No Manual Writing**: Never write OpenAPI YAML manually

### What You DON'T Need to Do:

❌ Write OpenAPI spec manually
❌ Keep docs in sync with code
❌ Set up separate documentation site
❌ Write API testing tools
❌ Explain request/response formats to users

---

## Try It Yourself!

### Exercise 1: Explore Your API Documentation

1. Open http://localhost:8080/docs
2. Click "GET /tasks"
3. Click "Try it out"
4. Click "Execute"
5. See the response!

### Exercise 2: View the Raw OpenAPI Spec

1. Open http://localhost:8080/openapi.json
2. Search for `"Task"` - see your Task schema
3. Search for `"minLength"` - see your validation rules
4. Search for `"description"` - see your doc tags

### Exercise 3: Import into Postman

1. Open Postman
2. Click "Import" → "Link"
3. Paste: http://localhost:8080/openapi.json
4. Click "Continue" → "Import"
5. All your endpoints are now in Postman!

### Exercise 4: Test a Request

Copy this curl command (generated from OpenAPI):
```bash
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{"title": "Test from OpenAPI", "description": "Learning OpenAPI!"}'
```

---

## OpenAPI vs Other API Documentation

| Feature | OpenAPI | Manual Docs | Postman Collections |
|---------|---------|-------------|---------------------|
| Auto-generated | ✅ | ❌ | ❌ |
| Machine-readable | ✅ | ❌ | ⚠️ Postman format only |
| Interactive testing | ✅ | ❌ | ✅ |
| Code generation | ✅ | ❌ | ❌ |
| Always in sync | ✅ | ❌ | ❌ |
| Industry standard | ✅ | ❌ | ❌ |
| Tool support | ✅ Massive | ⚠️ Limited | ⚠️ Postman only |

---

## Advanced: OpenAPI Features You're Using

Your API already uses these OpenAPI 3.1 features:

### 1. **Path Parameters**
```go
// In main.go
Path: "/tasks/{id}"
```
→ OpenAPI knows `{id}` is a path parameter

### 2. **Request Bodies**
```go
// Your CreateTaskInput
type CreateTaskInput struct {
    Body struct {
        Title string `json:"title" minLength:"1"`
    }
}
```
→ OpenAPI generates request body schema

### 3. **Response Types**
```go
// Your handlers return
func CreateTask(...) (*models.CreateTaskOutput, error)
```
→ OpenAPI documents the response type

### 4. **Validation Rules**
```go
Title string `json:"title" minLength:"1" maxLength:"200"`
```
→ OpenAPI includes validation in schema

### 5. **Error Responses**
```go
return nil, huma.Error404NotFound("Task not found")
```
→ OpenAPI documents error responses (RFC 7807 format)

### 6. **Tags for Organization**
```go
Tags: []string{"Tasks"}
```
→ OpenAPI groups endpoints by tags

---

## Summary

**OpenAPI = Standard format for describing REST APIs**

**Why it's awesome:**
- 📚 Interactive documentation (try endpoints in browser)
- 🤖 Auto-generated from your code
- 🛠️ Generate client libraries in any language
- 🧪 Import into testing tools (Postman, Insomnia)
- 🤝 Frontend/backend can work in parallel

**Your API has 3 OpenAPI endpoints:**
- http://localhost:8080/docs (interactive docs)
- http://localhost:8080/openapi.json (JSON spec)
- http://localhost:8080/openapi.yaml (YAML spec)

**Huma makes OpenAPI easy:**
- Write Go code with struct tags
- Huma generates OpenAPI automatically
- Docs always match your code
- No manual YAML writing needed!

**Next time someone asks "How do I use your API?"**
→ Send them: http://localhost:8080/docs
→ They can read AND test it immediately!

---

## Further Reading

- **OpenAPI Specification**: https://spec.openapis.org/oas/latest.html
- **Huma Documentation**: https://huma.rocks
- **Swagger UI**: https://swagger.io/tools/swagger-ui/
- **OpenAPI Generator**: https://openapi-generator.tech

---

**Remember**: With Huma, you get world-class API documentation for free just by writing good Go code with proper struct tags! 🎉
