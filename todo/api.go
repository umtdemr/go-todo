package todo

import (
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

// RegisterRoutes registers the routes for the todo API
func (s *APIRoute) RegisterRoutes(router *mux.Router, userService user.Service) {
	router.Handle("/todo", userService.AuthMiddleware(http.HandlerFunc(s.handleList)))
	router.Handle("/todo/list", userService.AuthMiddleware(http.HandlerFunc(s.handleList)))
	router.Handle("/todo/create", userService.AuthMiddleware(http.HandlerFunc(s.handleAdd)))
	router.Handle("/todo/update", userService.AuthMiddleware(http.HandlerFunc(s.handleUpdate)))
	router.Handle("/todo/{id}", userService.AuthMiddleware(http.HandlerFunc(s.handleFetchAndDelete)))
}

// handleList handles the list request
func (s *APIRoute) handleList(w http.ResponseWriter, r *http.Request) {
	authenticatedUser := r.Context().Value("user").(*user.VisibleUser)
	todos, err := s.Service.GetAllTodos(authenticatedUser.Id)

	if err != nil {
		server.RespondWithError(w, fmt.Sprintf("error while getting list: %s", err), http.StatusBadRequest)
		return
	}
	server.RespondOK(w, todos)
}

// handleAdd handles the add request
func (s *APIRoute) handleAdd(w http.ResponseWriter, r *http.Request) {
	// only POST methods are allowed
	if r.Method != http.MethodPost {
		err := server.ErrNotValidMethod.With("only POST methods are allowed")
		server.RespondWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	var createTodoType CreateTodoData

	err := server.DecodeBody(r, &createTodoType)

	if err != nil {
		server.RespondWithError(w, fmt.Sprintf("parsing error: %v", err), http.StatusBadRequest)
		return
	}

	// get the authenticated user from the context
	authenticatedUser := r.Context().Value("user").(*user.VisibleUser)

	// create the todo
	createdTodo, createErr := s.Service.CreateTodo(&createTodoType, authenticatedUser.Id)
	if createErr != nil {
		fmt.Fprintf(os.Stderr, "error while generating the todo: %s\n", createErr)
	}

	server.RespondCreated(w, createdTodo)
}

// handleUpdate handles the update request
func (s *APIRoute) handleUpdate(w http.ResponseWriter, r *http.Request) {
	// only POST methods are allowed
	if r.Method != http.MethodPost {
		err := server.ErrNotValidMethod.With("only POST methods are allowed")
		server.RespondWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	var updateData UpdateTodoData

	err := server.DecodeBody(r, &updateData)

	if err != nil {
		server.RespondWithError(w, fmt.Sprintf("error while parsing: %s", err), http.StatusBadRequest)
		return
	}

	// check if the ID is sent
	if updateData.Id == nil {
		server.RespondWithError(w, "ID is required for updating", http.StatusBadRequest)
		return
	}

	// get the authenticated user from the context
	authenticatedUser := r.Context().Value("user").(*user.VisibleUser)

	// update the todo
	updatedTodo, updateErr := s.Service.UpdateTodo(&updateData, authenticatedUser.Id)
	if updateErr != nil {
		server.RespondWithError(w, fmt.Sprintf("Error while updating: %s", updateErr), http.StatusBadRequest)
		return
	}

	server.RespondOK(w, updatedTodo)
}

// handleFetchAndDelete handles the fetch and delete requests
func (s *APIRoute) handleFetchAndDelete(w http.ResponseWriter, r *http.Request) {
	// only GET and DELETE methods are allowed
	if r.Method != http.MethodGet && r.Method != http.MethodDelete {
		err := server.ErrNotValidMethod.With("only GET and DELETE methods are allowed")
		server.RespondWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// get the ID from the path variables
	pathVars := mux.Vars(r)
	todoId, ok := pathVars["id"]

	// if the ID is not sent, respond with an error
	if !ok {
		err := server.ErrInvalidRequest.With("id is not sent")
		server.RespondWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// parse the ID to int
	todoIdInt, parseErr := strconv.Atoi(todoId)

	if parseErr != nil {
		err := server.ErrInvalidRequest.With("need a numeric value for the id")
		server.RespondWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// get the authenticated user from the context
	authenticatedUser := r.Context().Value("user").(*user.VisibleUser)

	// if the method is DELETE, remove the todo
	if r.Method == http.MethodDelete {
		removedTodo, removeErr := s.Service.RemoveTodo(todoIdInt, authenticatedUser.Id)
		if removeErr != nil {
			server.RespondWithError(w, removeErr.Error(), http.StatusBadRequest)
			return
		} else {
			server.RespondNoContent(w, removedTodo)
			return
		}
	} else {
		// if the method is GET, fetch the todo
		fetchedTodo, err := s.Service.GetTodo(todoIdInt, authenticatedUser.Id)
		if err != nil {
			server.RespondWithError(w, err.Error(), http.StatusBadRequest)
			return
		}
		server.RespondOK(w, fetchedTodo)
		return
	}
}
