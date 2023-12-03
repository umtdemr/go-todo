package main

import (
	"fmt"
	"github.com/umtdemr/go-todo/respond"
	"github.com/umtdemr/go-todo/todo"
	"net/http"
)

type Database struct {
	Todos []todo.Todo
}

var db = Database{
	Todos: []todo.Todo{todo.NewTodo("initial")},
}

func handleList(w http.ResponseWriter, r *http.Request) {
	respond.Respond(w, db.Todos)
}

func main() {
	http.HandleFunc("/list", handleList)
	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		fmt.Printf("got an error: %s\n", err)
	}
}
