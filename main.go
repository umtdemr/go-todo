package main

import (
	"fmt"
	"github.com/umtdemr/go-todo/respond"
	"net/http"
	"time"
)

type todo struct {
	Title     string    `json:"title"`
	Done      bool      `json:"done"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	initialTodoS := []todo{todo{Title: "initial", CreatedAt: time.Now()}}
	respond.Respond(w, initialTodoS)
}

func main() {
	http.HandleFunc("/", handleHome)
	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		fmt.Printf("got an error: %s\n", err)
	}
}
