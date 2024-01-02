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
	Route   string
	Service *Service
}

func NewTodoAPIRoute(service *Service) *APIRoute {
	return &APIRoute{Route: "todo", Service: service}
}

func (s *APIRoute) RegisterRoutes(router *mux.Router, userService user.Service) {
	router.Handle("/todo", userService.AuthMiddleware(http.HandlerFunc(s.handleList)))
	router.Handle("/todo/list", userService.AuthMiddleware(http.HandlerFunc(s.handleList)))
	router.Handle("/todo/create", userService.AuthMiddleware(http.HandlerFunc(s.handleAdd)))
	router.Handle("/todo/update", userService.AuthMiddleware(http.HandlerFunc(s.handleUpdate)))
	router.Handle("/todo/{id}", userService.AuthMiddleware(http.HandlerFunc(s.handleFetchAndDelete)))
}

func (s *APIRoute) handleList(w http.ResponseWriter, r *http.Request) {
	authenticatedUser := r.Context().Value("user").(*user.VisibleUser)
	todos, err := s.Service.GetAllTodos(authenticatedUser.Id)

	if err != nil {
		server.RespondWithError(w, fmt.Sprintf("error while getting list: %s", err), http.StatusBadRequest)
		return
	}
	server.Respond(w, todos)
}

func (s *APIRoute) handleAdd(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		err := server.ErrNotValidMethod.With("only POST methods are allowed")
		server.RespondWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	var createTodoType CreateTodoData

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&createTodoType)

	if err != nil {
		server.RespondWithError(w, fmt.Sprintf("parsing error: %v", err), http.StatusBadRequest)
		return
	}

	authenticatedUser := r.Context().Value("user").(*user.VisibleUser)

	createdTodo, createErr := s.Service.CreateTodo(&createTodoType, authenticatedUser.Id)
	if createErr != nil {
		fmt.Fprintf(os.Stderr, "error while generating the todo: %s\n", createErr)
	}

	server.Respond(w, createdTodo)
}

func (s *APIRoute) handleUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		err := server.ErrNotValidMethod.With("only POST methods are allowed")
		server.RespondWithError(w, err.Error(), http.StatusBadRequest)
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

	authenticatedUser := r.Context().Value("user").(*user.VisibleUser)

	updatedTodo, updateErr := s.Service.UpdateTodo(&updateData, authenticatedUser.Id)
	if updateErr != nil {
		server.RespondWithError(w, fmt.Sprintf("Error while updating: %s", updateErr), http.StatusBadRequest)
		return
	}

	server.Respond(w, updatedTodo)
}

func (s *APIRoute) handleFetchAndDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodDelete {
		err := server.ErrNotValidMethod.With("only GET and DELETE methods are allowed")
		server.RespondWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	pathVars := mux.Vars(r)
	todoId, ok := pathVars["id"]
	if !ok {
		err := server.ErrInvalidRequest.With("id is not sent")
		server.RespondWithError(w, err.Error(), http.StatusBadRequest)
		return
	}
	todoIdInt, parseErr := strconv.Atoi(todoId)

	if parseErr != nil {
		err := server.ErrInvalidRequest.With("need a numeric value for the id")
		server.RespondWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	authenticatedUser := r.Context().Value("user").(*user.VisibleUser)
	if r.Method == http.MethodDelete {
		removedTodo, removeErr := s.Service.RemoveTodo(todoIdInt, authenticatedUser.Id)
		if removeErr != nil {
			server.RespondWithError(w, removeErr.Error(), http.StatusBadRequest)
			return
		} else {
			server.Respond(w, removedTodo)
			return
		}
	} else {
		fetchedTodo, err := s.Service.GetTodo(todoIdInt, authenticatedUser.Id)
		if err != nil {
			server.RespondWithError(w, err.Error(), http.StatusBadRequest)
			return
		}
		server.Respond(w, fetchedTodo)
		return
	}
}
