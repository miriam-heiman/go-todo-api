// package main - This declares that this file belongs to the "main" package
// In Go, "main" is special - it tells Go this is a program you can run, not a library
// Think of it like "this is the entry point" of your application
package main

// import - This is how you bring in code from other packages (like importing in JavaScript or Python)
// The parentheses create a "block" where you list all the imports
import (
	"encoding/json" // encoding/json = JSON package - for converting Go data structures to/from JSON
	"fmt"           // fmt = format package - for printing text and formatting strings
	"log"           // log = logging package - for writing messages to the console
	"net/http"      // net/http = HTTP package - for creating web servers (like Express.js)
	"strconv"       // strconv = string conversion package - for converting strings to numbers and vice versa
)

// type Task struct - This creates a custom data type called "Task"
// "struct" is Go's way of grouping related data together (like a class without methods, or a JavaScript object with a fixed shape)
// Think of it as defining a template: "Every Task must have these exact fields"
type Task struct {
	ID          int    `json:"id"`          // ID field of type int (integer/number), json tag tells Go to use "id" when converting to JSON
	Title       string `json:"title"`       // Title field of type string (text), json tag says use "title" in JSON
	Description string `json:"description"` // Description field of type string, json tag says use "description" in JSON
	Completed   bool   `json:"completed"`   // Completed field of type bool (boolean: true or false), json tag says use "completed" in JSON
}

// var tasks []Task - Declares a variable called "tasks"
// []Task means "a slice of Task" - think of it as an array or list that can grow/shrink
// It starts empty - we can add Task objects to it
var tasks []Task

// var nextID = 1 - Declares a variable called nextID and sets it to 1
// We'll use this to give each new task a unique ID number
var nextID = 1

// func main() - This is THE main function that runs when you start your program
// It's like the entry point - Go automatically calls this function when you run your program
// The empty () means it takes no parameters
func main() {
	// tasks = append(tasks, Task{...}) - This adds a new Task to our tasks list
	// append() adds an item to the end of a slice (like push() in JavaScript arrays)
	// Task{...} creates a new Task object (called a "struct literal")
	// Inside the braces, we're setting each field using "FieldName: value" syntax
	tasks = append(tasks, Task{
		ID:          1,                               // Set ID field to 1
		Title:       "Learn Go Basics",               // Set Title to this text
		Description: "Understand Go's main concepts", // Set Description to this text
		Completed:   false,                           // Set Completed to false (not done yet)
	})

	// nextID = 2 - Set the nextID counter to 2 since we just used ID 1
	nextID = 2

	// http.HandleFunc() - This tells the server "when someone visits this URL, run this function"
	// It maps URL paths to handler functions
	// The first parameter is the URL path (like "/" for homepage)
	// The second parameter is the name of the function to call when that URL is visited
	http.HandleFunc("/", homeHandler)         // Visit "/" (homepage) → call homeHandler function
	http.HandleFunc("/tasks", tasksHandler)   // Visit "/tasks" → call tasksHandler function
	http.HandleFunc("/health", healthHandler) // Visit "/health" → call healthHandler function

	// fmt.Println() - This prints text to the console (like console.log() in JavaScript)
	// Each line prints a different message to help you know the server started
	fmt.Println("🚀 Server starting on http://localhost:8080")
	fmt.Println("Try visiting:")
	fmt.Println("  - http://localhost:8080/ (homepage)")
	fmt.Println("  - http://localhost:8080/tasks (get all tasks)")
	fmt.Println("  - http://localhost:8080/health (health check)")

	// log.Fatal(http.ListenAndServe(":8080", nil)) - This starts the actual HTTP server
	// http.ListenAndServe() starts listening for incoming HTTP requests on port 8080
	// ":8080" means "listen on all network interfaces on port 8080"
	// nil means "use default settings" (we could pass custom settings here, but we're not)
	// log.Fatal() means "if this fails, print the error and exit the program"
	// This line blocks (waits forever) until the server stops or crashes
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// func homeHandler(...) - This is a function that handles requests to the homepage
// It takes two parameters:
//   - w http.ResponseWriter - This is what we use to write data back to the client (like res in Express.js)
//   - r *http.Request - This contains information about the incoming request (like req in Express.js)
//
// The * means it's a pointer (reference to the request object, not a copy of it)
func homeHandler(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set() - This sets HTTP headers (metadata about our response)
	// "Content-Type" tells the browser/client what kind of data we're sending
	// "text/html" means "this is HTML that you should render as a web page"
	w.Header().Set("Content-Type", "text/html")

	// fmt.Fprintf(w, ...) - This writes formatted text directly to the response
	// w is the response writer (where we're sending data to the client)
	// The backticks ` create a "raw string literal" - everything inside is treated literally (including newlines)
	// This means we can write multi-line HTML and it preserves the formatting
	fmt.Fprintf(w, `
		<h1>Welcome to my Go To-Do API!</h1>
		<p>This is running on a Go HTTP server!</p>
		<h2>Available endpoints:</h2>
		<ul>
			<li><a href="/tasks">GET /tasks</a> - Get all tasks</li>
			<li>POST /tasks - Create a new task</li>
			<li>GET /tasks?id=X - Get task by ID</li>
			<li>PUT /tasks?id=X - Update a task</li>
			<li>DELETE /tasks?id=X - Delete a task</li>
		</ul>
		<p><a href="/health">Health Check</a></p>
	`)
}

// func healthHandler(...) - This function handles the health check endpoint
// Health checks are used to verify the server is running properly (monitoring tools use this)
func healthHandler(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set() - Set the content type to JSON
	// "application/json" tells the client "this is JSON data"
	w.Header().Set("Content-Type", "application/json")

	// fmt.Fprintf(w, ...) - Write a JSON string directly to the response
	// The backtick creates a raw string literal containing JSON
	fmt.Fprintf(w, `{"status": "healthy", "message": "Server is running!"}`)
}

// func tasksHandler(...) - This function handles requests to the /tasks endpoint
// It handles different HTTP methods: GET (read), POST (create), PUT (update), DELETE (delete)
func tasksHandler(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set() - Tell the client we're sending JSON data in our response
	w.Header().Set("Content-Type", "application/json")

	// r.Method - This gives us the HTTP method used in the request (GET, POST, PUT, DELETE, etc.)
	// We'll check which method was used and call different code for each one
	switch r.Method {

	// case "GET" - Handle requests to read/fetch tasks
	case "GET":
		// Check if an ID parameter was provided in the URL query (like ?id=1)
		// r.URL.Query() gets the query parameters from the URL
		idParam := r.URL.Query().Get("id") // Get the value of the "id" parameter, or "" if not present

		// if idParam == "" - If no ID was provided, return all tasks
		if idParam == "" {
			// Call a helper function to get and return all tasks
			getAllTasksHandler(w)
		} else {
			// strconv.Atoi(idParam) - Convert the ID parameter from string to integer
			// Atoi returns TWO values: the integer and an error (this is Go's way of error handling)
			id, err := strconv.Atoi(idParam)

			// if err != nil - If there was an error converting the ID (like if it's not a number)
			if err != nil {
				// http.Error() - Send an error response back to the client
				// Status 400 means "Bad Request" (the client sent invalid data)
				http.Error(w, "Invalid ID", http.StatusBadRequest)
				return // Exit the function early
			}

			// Call a helper function to get and return a specific task by ID
			getTaskByIDHandler(w, id)
		}

	// case "POST" - Handle requests to create a new task
	case "POST":
		createTaskHandler(w, r)

	// case "PUT" - Handle requests to update an existing task
	case "PUT":
		updateTaskHandler(w, r)

	// case "DELETE" - Handle requests to delete a task
	case "DELETE":
		deleteTaskHandler(w, r)

	// default - Handle any HTTP method we don't support
	default:
		// http.Error() - Send an error response (405 = Method Not Allowed)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// func getAllTasksHandler(...) - Helper function to get and return all tasks
func getAllTasksHandler(w http.ResponseWriter) {
	// json.NewEncoder(w) - Create a JSON encoder that writes to the response writer
	// .Encode(tasks) - Convert our tasks slice to JSON and write it to the response
	// If anything goes wrong, Encode returns an error
	err := json.NewEncoder(w).Encode(tasks)

	// if err != nil - Check if encoding failed
	if err != nil {
		// http.Error() - Send an error response back to the client
		// Status 500 means "Internal Server Error" (something went wrong on our end)
		http.Error(w, "Failed to encode tasks", http.StatusInternalServerError)
		return
	}

	// fmt.Println() - Print a success message to the console
	fmt.Println("✅ Retrieved all tasks")
}

// func getTaskByIDHandler(...) - Helper function to get a specific task by its ID
func getTaskByIDHandler(w http.ResponseWriter, id int) {
	// Loop through all tasks to find one with a matching ID
	// for range - This is Go's way of looping through a slice (like forEach in JavaScript)
	// _ means "ignore this value" (we don't need the index)
	for _, task := range tasks {
		// if task.ID == id - Check if this task's ID matches the one we're looking for
		if task.ID == id {
			// json.NewEncoder(w).Encode(task) - Convert this single task to JSON and send it
			json.NewEncoder(w).Encode(task)
			fmt.Printf("✅ Retrieved task with ID %d\n", id) // Print with the ID substituted in
			return                                          // Exit the function since we found the task
		}
	}

	// If we got here, the task wasn't found
	// http.Error() - Send a "404 Not Found" error
	http.Error(w, "Task not found", http.StatusNotFound)
}

// func createTaskHandler(...) - Helper function to create a new task
func createTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Declare a variable to hold the incoming task data
	var newTask Task

	// json.NewDecoder(r.Body) - Create a JSON decoder that reads from the request body
	// .Decode(&newTask) - Convert the JSON from the request body into a Task struct
	// &newTask means "pass a reference to the newTask variable" (so Decode can modify it)
	err := json.NewDecoder(r.Body).Decode(&newTask)

	// if err != nil - Check if decoding failed (maybe the JSON was malformed)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Check that required fields are provided
	// newTask.Title == "" - Check if title is empty
	if newTask.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}

	// Assign the next available ID to this new task
	// We use nextID global variable and then increment it for the next task
	newTask.ID = nextID
	nextID++ // Increment for next time

	// Set completed to false by default if not specified
	newTask.Completed = false

	// tasks = append(tasks, newTask) - Add the new task to our tasks slice
	tasks = append(tasks, newTask)

	// json.NewEncoder(w).Encode(newTask) - Successfully return the created task as JSON
	json.NewEncoder(w).Encode(newTask)
	fmt.Printf("✅ Created new task: %s\n", newTask.Title)
}

// func updateTaskHandler(...) - Helper function to update an existing task
func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Get the ID parameter from the URL query
	idParam := r.URL.Query().Get("id")

	// if idParam == "" - Check if ID was provided
	if idParam == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	// Convert ID from string to integer
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Declare a variable to hold the updated task data
	var updatedTask Task

	// Decode the JSON from the request body
	err = json.NewDecoder(r.Body).Decode(&updatedTask)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Loop through tasks to find the one to update
	// We need the index (i) this time so we can modify the task in place
	for i, task := range tasks {
		if task.ID == id {
			// Keep the original ID
			updatedTask.ID = id

			// If title is provided, update it; otherwise keep the old title
			if updatedTask.Title == "" {
				updatedTask.Title = task.Title
			}

			// If description is provided, update it; otherwise keep the old one
			if updatedTask.Description == "" {
				updatedTask.Description = task.Description
			}

			// Replace the task at index i with the updated version
			tasks[i] = updatedTask

			// Return the updated task as JSON
			json.NewEncoder(w).Encode(updatedTask)
			fmt.Printf("✅ Updated task with ID %d\n", id)
			return
		}
	}

	// Task wasn't found
	http.Error(w, "Task not found", http.StatusNotFound)
}

// func deleteTaskHandler(...) - Helper function to delete a task
func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Get the ID parameter from the URL query
	idParam := r.URL.Query().Get("id")

	// if idParam == "" - Check if ID was provided
	if idParam == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	// Convert ID from string to integer
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Loop through tasks to find the one to delete
	// We need the index (i) to remove the task
	for i, task := range tasks {
		if task.ID == id {
			// Delete the task from the slice using a trick:
			// tasks[:i] gets all tasks before index i
			// tasks[i+1:] gets all tasks after index i
			// ... unpacks the second slice (syntax for appending multiple items)
			tasks = append(tasks[:i], tasks[i+1:]...)

			// Send a success response
			fmt.Fprintf(w, `{"message": "Task deleted successfully"}`)
			fmt.Printf("✅ Deleted task with ID %d\n", id)
			return
		}
	}

	// Task wasn't found
	http.Error(w, "Task not found", http.StatusNotFound)
}
