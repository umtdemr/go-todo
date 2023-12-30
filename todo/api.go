package todo

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/umtdemr/go-todo/server"
	"github.com/umtdemr/go-todo/user"
	"net/http"
	"os"
	"strconv"
)

type APIRoute struct {
	Route      string
	Repository Repository
}

func NewTodoAPIRoute(repository Repository) *APIRoute {
	return &APIRoute{Route: "todo", Repository: repository}
}

func (s *APIRoute) RegisterRoutes(router *mux.Router) {
	router.Handle("/todo", user.AuthMiddleware(http.HandlerFunc(s.handleList)))
	router.Handle("/todo/list", user.AuthMiddleware(http.HandlerFunc(s.handleList)))
	router.Handle("/todo/create", user.AuthMiddleware(http.HandlerFunc(s.handleAdd)))
	router.Handle("/todo/update", user.AuthMiddleware(http.HandlerFunc(s.handleUpdate)))
	router.Handle("/todo/{id}", user.AuthMiddleware(http.HandlerFunc(s.handleFetchAndDelete)))
}

func (s *APIRoute) handleList(w http.ResponseWriter, r *http.Request) {
	todos, err := s.Repository.GetAllTodos()

	if err != nil {
		server.RespondWithError(w, fmt.Sprintf("error while getting list: %s", err), http.StatusBadRequest)
		return
	}
	server.Respond(w, todos)
}

func (s *APIRoute) handleAdd(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		server.RespondWithError(w, "not valid", http.StatusBadRequest)
		return
	}

	var createTodoType CreateTodoData

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&createTodoType)

	if err != nil {
		server.RespondWithError(w, fmt.Sprintf("parsing error: %v", err), http.StatusBadRequest)
		return
	}

	if createTodoType.Title == "" {
		server.RespondWithError(w, "Title need to be sent", http.StatusBadRequest)
		return
	}

	createTodoData := NewTodo(createTodoType.Title)
	createdTodo, createErr := s.Repository.CreateTodo(createTodoData)
	if createErr != nil {
		fmt.Fprintf(os.Stderr, "error while generating the todo: %s\n", createErr)
	}

	server.Respond(w, createdTodo)
}

func (s *APIRoute) handleUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		server.RespondWithError(w, "Only POST methods are allowed", http.StatusBadRequest)
		return
	}

	var updateData UpdateTodoData

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&updateData)

	if err != nil {
		server.RespondWithError(w, fmt.Sprintf("error while parsing: %s", err), http.StatusBadRequest)
		return
	}

	if updateData.Id == nil {
		server.RespondWithError(w, "ID is required for updating", http.StatusBadRequest)
		return
	}

	updatedTodo, updateErr := s.Repository.UpdateTodo(&updateData)
	if updateErr != nil {
		server.RespondWithError(w, fmt.Sprintf("Error while updating: %s", updateErr), http.StatusBadRequest)
		return
	}

	server.Respond(w, updatedTodo)
}

func (s *APIRoute) handleFetchAndDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodDelete {
		server.RespondWithError(w, "Only GET and DELETE requests are allowed", http.StatusBadRequest)
		return
	}

	pathVars := mux.Vars(r)
	todoId, ok := pathVars["id"]
	if !ok {
		server.RespondWithError(w, "couldn't find the id", http.StatusBadRequest)
		return
	}
	todoIdInt, parseErr := strconv.Atoi(todoId)

	if parseErr != nil {
		server.RespondWithError(w, "need and integer value as id", http.StatusBadRequest)
		return
	}

	if r.Method == http.MethodDelete {
		removedTodo, removeErr := s.Repository.RemoveTodo(todoIdInt)
		if removeErr != nil {
			server.RespondWithError(w, removeErr.Error(), http.StatusBadRequest)
			return
		} else {
			server.Respond(w, removedTodo)
			return
		}
	} else {
		fetchedTodo, err := s.Repository.GetTodo(todoIdInt)
		if err != nil {
			server.RespondWithError(w, err.Error(), http.StatusBadRequest)
			return
		}
		server.Respond(w, fetchedTodo)
		return
	}
}
