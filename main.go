package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/umtdemr/go-todo/respond"
	"github.com/umtdemr/go-todo/todo"
	"net/http"
	"slices"
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

func main() {
	http.HandleFunc("/list", handleList)
	http.HandleFunc("/create", handleAdd)
	http.HandleFunc("/update", handleUpdate)
	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		fmt.Printf("got an error: %s\n", err)
	}
}
