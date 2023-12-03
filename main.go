package main

import (
	"encoding/json"
	"fmt"
	"github.com/umtdemr/go-todo/respond"
	"github.com/umtdemr/go-todo/todo"
	"net/http"
)

type Database struct {
	Todos *[]todo.Todo
}

var db = Database{
	Todos: &[]todo.Todo{todo.NewTodo("hey")},
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

	*db.Todos = append(*db.Todos, todo.NewTodo(createTodoType.Title))
}

func main() {
	http.HandleFunc("/list", handleList)
	http.HandleFunc("/create", handleAdd)
	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		fmt.Printf("got an error: %s\n", err)
	}
}
