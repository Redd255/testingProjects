package main

import (
	"fmt"
	"html/template"
	"net/http"
	"sync"
)

// Define the ToDo item struct
type ToDo struct {
	ID     int
	Task   string
	Status bool
}

var (
	todoList  = []ToDo{}
	mu        sync.Mutex // To handle concurrency
	idCounter int        = 1
)

// Load HTML templates
var templates = template.Must(template.ParseFiles("./templates/index.html"))

// Handler to display the todo list
func indexHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	// Render the template with the todo list
	templates.ExecuteTemplate(w, "index.html", todoList)
}

// Handler to add a new todo item
func addHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		task := r.FormValue("task")

		if task != "" {
			mu.Lock()
			todo := ToDo{
				ID:     idCounter,
				Task:   task,
				Status: false,
			}
			idCounter++
			todoList = append(todoList, todo)
			mu.Unlock()
		}
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

// Handler to mark a task as complete
func completeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		id := r.FormValue("id")

		mu.Lock()
		for i, todo := range todoList {
			if fmt.Sprint(todo.ID) == id {
				todoList[i].Status = true
				break
			}
		}
		mu.Unlock()
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func main() {
	// Route handlers
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/add", addHandler)
	http.HandleFunc("/complete", completeHandler)

	// Serve static assets (e.g., CSS)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Start the server
	fmt.Println("Server started at :8080")
	http.ListenAndServe(":8080", nil)
}
