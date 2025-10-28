# Understanding Middleware - A Beginner's Guide

## What is Middleware?

Imagine you're entering a building. Before you can reach your destination (a specific office), you might pass through:
1. A security checkpoint (checks your ID)
2. A reception desk (logs your visit)
3. An elevator (gets you to the right floor)

**Middleware works the same way in web applications!**

Middleware is code that runs **between** receiving a request and sending a response. Every HTTP request passes through your middleware before reaching your actual handler functions.

## Visual Flow

```
Client Request
     ↓
[Logging Middleware]     ← Logs: "GET /tasks started"
     ↓
[CORS Middleware]        ← Adds: "Access-Control-Allow-Origin: *"
     ↓
[Your Handler]           ← Executes: tasksHandler()
     ↓
[CORS Middleware]        ← Returns through middleware
     ↓
[Logging Middleware]     ← Logs: "GET /tasks took 500µs"
     ↓
Client Response
```

## Why Use Middleware?

### 1. DRY Principle (Don't Repeat Yourself)
Without middleware, you'd have to add the same code to every handler:

```go
// WITHOUT middleware - repetitive! ❌
func tasksHandler(w http.ResponseWriter, r *http.Request) {
    log.Printf("Request: %s %s", r.Method, r.URL.Path)  // Repeated
    w.Header().Set("Access-Control-Allow-Origin", "*")   // Repeated

    // Your actual logic here...
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
    log.Printf("Request: %s %s", r.Method, r.URL.Path)  // Repeated again!
    w.Header().Set("Access-Control-Allow-Origin", "*")   // Repeated again!

    // Your actual logic here...
}
```

With middleware, you write it once:

```go
// WITH middleware - clean! ✅
func tasksHandler(w http.ResponseWriter, r *http.Request) {
    // Just your business logic - middleware handles the rest!
    // Your actual logic here...
}
```

### 2. Separation of Concerns
Each middleware has one job:
- **Logging middleware**: Only worries about logging
- **CORS middleware**: Only worries about CORS headers
- **Auth middleware**: Only worries about authentication

This makes code easier to understand, test, and maintain.

### 3. Easy to Add/Remove Features
Want to add logging to all routes? Just add one middleware.
Want to remove it? Remove one line.

## How Middleware Works in Go

### The Pattern

Middleware in Go follows this pattern:

```go
func middlewareName(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Code BEFORE calling the handler

        next.ServeHTTP(w, r)  // Call the next handler in the chain

        // Code AFTER calling the handler
    })
}
```

Let's break this down:

1. **Takes a handler**: `next http.Handler` - the next thing to call
2. **Returns a handler**: `http.Handler` - a new wrapped handler
3. **Calls the next handler**: `next.ServeHTTP(w, r)` - passes control to the next middleware or handler

### Example: Logging Middleware

```go
func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // BEFORE: Record start time
        start := time.Now()

        // Call the next handler (your actual route handler)
        next.ServeHTTP(w, r)

        // AFTER: Log how long it took
        log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
    })
}
```

**What happens:**
1. Request arrives
2. Middleware records the start time
3. Middleware calls your handler (`next.ServeHTTP`)
4. Your handler does its work and returns
5. Middleware calculates elapsed time
6. Middleware logs the request details
7. Response sent to client

## Our Middleware Examples

### 1. Logging Middleware

**Purpose**: Track every request with timing information

```go
func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        next.ServeHTTP(w, r)
        log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
    })
}
```

**Output:**
```
2025/10/28 15:03:15 GET /health 9.333µs
2025/10/28 15:04:04 GET /tasks 732.75µs
2025/10/28 15:04:16 POST /tasks 948.958µs
```

**Benefits:**
- See every request your server receives
- Track performance (slow requests stand out)
- Debug issues (did the request even arrive?)
- Monitor usage patterns

### 2. CORS Middleware

**Purpose**: Allow browsers from other domains to access your API

```go
func corsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Add CORS headers before calling handler
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

        // Handle preflight requests
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return  // Don't call next handler for OPTIONS
        }

        next.ServeHTTP(w, r)
    })
}
```

**What is CORS?**
- **C**ross-**O**rigin **R**esource **S**haring
- Browser security feature that blocks requests from different domains
- Example: A website at `example.com` trying to call your API at `localhost:8080`
- Without CORS, the browser blocks this
- With CORS headers, you tell the browser "it's okay, allow this"

**What happens:**
1. Browser makes request to your API
2. CORS middleware adds special headers
3. Browser sees the headers and allows the request

### 3. Chaining Middleware

**Purpose**: Apply multiple middleware in a specific order

```go
func chainMiddleware(h http.Handler, middleware ...func(http.Handler) http.Handler) http.Handler {
    for i := len(middleware) - 1; i >= 0; i-- {
        h = middleware[i](h)
    }
    return h
}
```

**Usage:**
```go
handler := chainMiddleware(
    mux,
    loggingMiddleware,  // Executes first
    corsMiddleware,     // Executes second
)
```

**Why reverse order?**
We apply middleware from right to left (reverse) so they execute left to right:
- `loggingMiddleware(corsMiddleware(yourHandler))`
- Request flows: logging → cors → handler
- Response flows back: handler → cors → logging

## Applying Middleware in main()

**Old way (no middleware):**
```go
func main() {
    http.HandleFunc("/tasks", tasksHandler)
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

**New way (with middleware):**
```go
func main() {
    // Create a router
    mux := http.NewServeMux()
    mux.Handle("/tasks", http.HandlerFunc(tasksHandler))

    // Wrap with middleware
    handler := chainMiddleware(
        mux,
        loggingMiddleware,
        corsMiddleware,
    )

    // Use the wrapped handler
    log.Fatal(http.ListenAndServe(":8080", handler))
}
```

## Common Middleware Use Cases

Here are other middleware you might build:

### 1. Authentication Middleware
```go
func authMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        token := r.Header.Get("Authorization")
        if token == "" {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return  // Stop here, don't call next
        }
        next.ServeHTTP(w, r)
    })
}
```

### 2. Rate Limiting Middleware
```go
func rateLimitMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if tooManyRequests() {
            http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
            return
        }
        next.ServeHTTP(w, r)
    })
}
```

### 3. Request ID Middleware
```go
func requestIDMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        requestID := generateID()
        w.Header().Set("X-Request-ID", requestID)
        next.ServeHTTP(w, r)
    })
}
```

## Key Takeaways

1. **Middleware runs before handlers** - Like checkpoints before reaching your destination
2. **Middleware wraps handlers** - Like layers of an onion around your core logic
3. **Order matters** - First middleware applied = first to execute
4. **Reusable** - Write once, apply to all routes
5. **Composable** - Stack multiple middleware together

## Testing Your Middleware

You can test your middleware is working:

1. **Check logs** - Look for request logs in your terminal
2. **Check headers** - Use `curl -i` to see response headers
3. **Check browser** - Open developer tools and look at network requests

Example:
```bash
# See CORS headers in response
curl -i http://localhost:8080/tasks

# Output includes:
# Access-Control-Allow-Origin: *
# Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS
```

## Next Steps

Now that you understand middleware, you can:
- Add authentication middleware to protect routes
- Add validation middleware to check request data
- Add compression middleware to reduce response sizes
- Add metrics middleware to track API usage

Middleware is a powerful pattern that keeps your code clean and maintainable!
