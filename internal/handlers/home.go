package handlers

import (
	"fmt"
	"net/http"
)

// Home handles requests to the homepage
func Home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, `
		<h1>Welcome to my Go To-Do API!</h1>
		<p>Now powered by MongoDB!</p>
		<p>Your tasks are now saved in the cloud </p>
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
