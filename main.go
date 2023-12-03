package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/umtdemr/go-todo/respond"
	"github.com/umtdemr/go-todo/todo"
	"net/http"
	"slices"
	"strings"
	"time"
)

type Database struct {
	Todos *[]*todo.Todo
}

var db = Database{
	Todos: &[]*todo.Todo{todo.NewTodo("initial todo")},
}

func handleList(w http.ResponseWriter, r *http.Request) {
	respond.Respond(w, db.Todos)
}

func handleAdd(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respond.WithError(w, "not valid", http.StatusBadRequest)
		return
	}

	var createTodoType todo.CreateTodoData

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&createTodoType)

	if err != nil {
		respond.WithError(w, fmt.Sprintf("parsing error: %v", err), http.StatusBadRequest)
		return
	}

	if createTodoType.Title == "" {
		respond.WithError(w, "Title need to be sent", http.StatusBadRequest)
		return
	}

	createdTodo := todo.NewTodo(createTodoType.Title)
	*db.Todos = append(*db.Todos, createdTodo)

	respond.Respond(w, createdTodo)
}

func handleUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respond.WithError(w, "Only POST methods are allowed", http.StatusBadRequest)
		return
	}

	var updateData todo.UpdateTodoData

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&updateData)

	if err != nil {
		respond.WithError(w, fmt.Sprintf("error while parsing: %s", err), http.StatusBadRequest)
		return
	}

	id, uuidParseErr := uuid.Parse(updateData.Id)

	if uuidParseErr != nil {
		respond.WithError(w, "the id you sent is invalid", http.StatusBadRequest)
		return
	}

	todoIndex := slices.IndexFunc(*db.Todos, func(c *todo.Todo) bool { return c.Id == id })
	todoItem := (*db.Todos)[todoIndex]

	if updateData.Title != nil {
		todoItem.Title = *updateData.Title
	}

	if updateData.Done != nil {
		todoItem.Done = *updateData.Done
	}

	todoItem.UpdatedAt = time.Now()

	respond.Respond(w, todoItem)
}

func handleFetchAndDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodDelete {
		respond.WithError(w, "Only GET and DELETE requests are allowed", http.StatusBadRequest)
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 2 {
		respond.WithError(w, "Invalid request", http.StatusBadRequest)
		return
	}
	todoId := parts[1]

	parsedTodoId, parseUuidErr := uuid.Parse(todoId)

	if parseUuidErr != nil {
		respond.WithError(w, "the id you sent is invalid", http.StatusBadRequest)
		return
	}

	todoIndex := slices.IndexFunc(*db.Todos, func(c *todo.Todo) bool { return c.Id == parsedTodoId })
	if r.Method == http.MethodDelete {
		*(db.Todos) = slices.Delete(*db.Todos, todoIndex, 1)
		respond.Respond(w, (*db.Todos)[todoIndex])
	} else {
		respond.Respond(w, (*db.Todos)[todoIndex])
	}
}

func main() {
	http.HandleFunc("/list", handleList)
	http.HandleFunc("/create", handleAdd)
	http.HandleFunc("/update", handleUpdate)
	http.HandleFunc("/", handleFetchAndDelete)
	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		fmt.Printf("got an error: %s\n", err)
	}
}
