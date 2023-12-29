package user

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/umtdemr/go-todo/server"
	"net/http"
)

type APIRoute struct {
	Route      string
	Repository Repository
}

func NewAPIRoute(repository Repository) *APIRoute {
	return &APIRoute{Route: "user", Repository: repository}
}

func (s *APIRoute) RegisterAPIRoutes(r *mux.Router) {
	r.HandleFunc("/user/register", s.handleCreateUser)
	r.HandleFunc("/user/login", s.handleLogin)
}

func (s *APIRoute) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		server.RespondWithError(w, "only post requests are allowed", http.StatusBadRequest)
		return
	}

	var userCreateData CreateUserData

	decoder := json.NewDecoder(r.Body)

	decodeErr := decoder.Decode(&userCreateData)
	if decodeErr != nil {
		server.RespondWithError(w, "make sure you provided all the necessary values", http.StatusBadRequest)
		return
	}

	createErr := s.Repository.CreateUser(&userCreateData)
	if createErr != nil {
		server.RespondWithError(w, fmt.Sprintf("error while creating user: %v", createErr), http.StatusBadRequest)
		return
	}

	response := make(map[string]string)
	response["message"] = "success"
	server.Respond(w, response)
	return
}

func (s *APIRoute) handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		server.RespondWithError(w, "only post requests are allowed", http.StatusBadRequest)
		return
	}

	var userLoginData LoginUserData

	decoder := json.NewDecoder(r.Body)

	decodeErr := decoder.Decode(&userLoginData)
	if decodeErr != nil ||
		(userLoginData.Password == nil || (userLoginData.Email == nil && userLoginData.Username == nil)) {
		server.RespondWithError(w, "make sure you provided all the necessary values", http.StatusBadRequest)
		return
	}

	isLoggedIn := s.Repository.Login(&userLoginData)
	response := make(map[string]string)
	message := "username or password is invalid"
	if isLoggedIn {
		message = "successfully logged in"
	}
	response["message"] = message
	server.Respond(w, response)
}
