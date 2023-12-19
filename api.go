package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/umtdemr/go-todo/todo"
	"net/http"
	"os"
	"strconv"
	"sync"
)

type APIServer struct {
	listenAddr string
	repository Repository
}

func NewAPIServer(listenAddr string, repository Repository) *APIServer {
	return &APIServer{listenAddr: listenAddr, repository: repository}
}

func (s *APIServer) Run() {
	r := mux.NewRouter()
	r.HandleFunc("/list", s.handleList)
	r.HandleFunc("/create", s.handleAdd)
	r.HandleFunc("/update", s.handleUpdate)
	r.HandleFunc("/{id}", s.handleFetchAndDelete)
	http.ListenAndServe(s.listenAddr, r)
}

func Respond(w http.ResponseWriter, data interface{}) {
	isThereError := false
	jsonData, err := json.Marshal(data)

	if err != nil {
		isThereError = true
	}

	w.Header().Set("Content-Type", "application/json")
	_, writeError := w.Write(jsonData)

	if writeError != nil {
		isThereError = true
	}

	if isThereError {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func RespondWithError(w http.ResponseWriter, msg string, errCode int) {
	if msg == "" {
		msg = "An error has occurred while processing"
	}
	if errCode == 0 {
		errCode = http.StatusBadRequest
	}
	w.WriteHeader(errCode)

	resp := make(map[string]string)
	resp["message"] = msg

	encoder := json.NewEncoder(w)
	err := encoder.Encode(resp)

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"error": %v}`, err)
	}
}

type Database struct {
	Todos *[]*todo.Todo
	mutex sync.Mutex
}

var db = Database{
	Todos: &[]*todo.Todo{todo.NewTodo("initial todo")},
}

func (s *APIServer) handleList(w http.ResponseWriter, r *http.Request) {
	todos, err := s.repository.GetAllTodos()

	if err != nil {
		RespondWithError(w, fmt.Sprintf("error while getting list: %s", err), http.StatusBadRequest)
		return
	}
	Respond(w, todos)
}

func (s *APIServer) handleAdd(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		RespondWithError(w, "not valid", http.StatusBadRequest)
		return
	}

	var createTodoType todo.CreateTodoData

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&createTodoType)

	if err != nil {
		RespondWithError(w, fmt.Sprintf("parsing error: %v", err), http.StatusBadRequest)
		return
	}

	if createTodoType.Title == "" {
		RespondWithError(w, "Title need to be sent", http.StatusBadRequest)
		return
	}

	createTodoData := todo.NewTodo(createTodoType.Title)
	createdTodo, createErr := s.repository.CreateTodo(createTodoData)
	if createErr != nil {
		fmt.Fprintf(os.Stderr, "error while generating the todo: %s\n", createErr)
	}

	Respond(w, createdTodo)
}

func (s *APIServer) handleUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		RespondWithError(w, "Only POST methods are allowed", http.StatusBadRequest)
		return
	}

	var updateData todo.UpdateTodoData

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&updateData)

	if err != nil {
		RespondWithError(w, fmt.Sprintf("error while parsing: %s", err), http.StatusBadRequest)
		return
	}

	if updateData.Id == nil {
		RespondWithError(w, "ID is required for updating", http.StatusBadRequest)
		return
	}

	dbErr := s.repository.UpdateTodo(&updateData)
	if dbErr != nil {
		RespondWithError(w, fmt.Sprintf("Error while updating: %s", dbErr), http.StatusBadRequest)
		return
	}

	Respond(w, updateData)
}

func (s *APIServer) handleFetchAndDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodDelete {
		RespondWithError(w, "Only GET and DELETE requests are allowed", http.StatusBadRequest)
		return
	}

	pathVars := mux.Vars(r)
	todoId, ok := pathVars["id"]
	if !ok {
		RespondWithError(w, "couldn't find the id", http.StatusBadRequest)
		return
	}
	todoIdInt, parseErr := strconv.Atoi(todoId)

	if parseErr != nil {
		RespondWithError(w, "need and integer value as id", http.StatusBadRequest)
		return
	}

	if r.Method == http.MethodDelete {
		err := s.repository.RemoveTodo(todoIdInt)
		if err != nil {
			RespondWithError(w, err.Error(), http.StatusBadRequest)
			return
		} else {
			messageMap := make(map[string]string)
			messageMap["message"] = "removed"
			Respond(w, messageMap)
			return
		}
	} else {
		fetchedTodo, err := s.repository.GetTodo(todoIdInt)
		if err != nil {
			RespondWithError(w, err.Error(), http.StatusBadRequest)
			return
		}
		Respond(w, fetchedTodo)
		return
	}
}
