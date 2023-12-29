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
