package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/umtdemr/go-todo/email"
	"github.com/umtdemr/go-todo/server"
	"net/http"
)

type APIRoute struct {
	Route   string
	Service Service
}

func NewAPIRoute(userService Service) *APIRoute {
	return &APIRoute{Route: "user", Service: userService}
}

func (route *APIRoute) RegisterAPIRoutes(router *mux.Router) {
	router.HandleFunc("/user/register", route.handleCreateUser)
	router.HandleFunc("/user/login", route.handleLogin)
	router.HandleFunc("/user/reset-password-request", route.handleResetPasswordRequest)
	router.HandleFunc("/user/new-password", route.handleNewPassword)
}

func (route *APIRoute) handleCreateUser(w http.ResponseWriter, r *http.Request) {
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

	createErr := route.Service.CreateUser(&userCreateData)
	if createErr != nil {
		var e UserError
		if errors.As(createErr, &e) {
			server.RespondWithError(w, fmt.Sprintf("validation error: %v", e.Error()), http.StatusBadRequest)
			return
		}
		server.RespondWithError(w, fmt.Sprintf("error while creating user: %v", createErr), http.StatusBadRequest)
		return
	}

	response := make(map[string]string)
	response["message"] = "success"
	server.Respond(w, response)
	return
}

func (route *APIRoute) handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		server.RespondWithError(w, "only post requests are allowed", http.StatusBadRequest)
		return
	}

	var userLoginData LoginUserData

	decoder := json.NewDecoder(r.Body)

	decodeErr := decoder.Decode(&userLoginData)
	if decodeErr != nil {
		server.RespondWithError(w, "couldn't decode the body", http.StatusBadRequest)
		return
	}

	tokenString, loginError := route.Service.Login(&userLoginData)

	if loginError != nil {
		server.RespondWithError(w, loginError.Error(), http.StatusBadRequest)
		return
	}

	response := make(map[string]string)
	response["message"] = "success"
	response["token"] = tokenString
	server.Respond(w, response)
}

func (route *APIRoute) handleResetPasswordRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		err := server.ErrNotValidMethod.With("only post requests are allowed")
		server.RespondWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	var resetPasswordRequestData ResetPasswordRequest

	// TODO: I can create a common handler for decoding the body since I've repeated this so much
	decoder := json.NewDecoder(r.Body)
	decodeErr := decoder.Decode(&resetPasswordRequestData)
	if decodeErr != nil {
		server.RespondWithError(w, "couldn't decode the body", http.StatusBadRequest)
		return
	}

	tokenString, err := route.Service.SendResetPasswordToken(&resetPasswordRequestData)

	if err != nil {
		server.RespondWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	sendErr := email.Send(email.SendEmailData{
		To:      []string{resetPasswordRequestData.Email},
		Subject: "Your reset password token",
		Message: fmt.Sprintf("Your reset password token is: %s", tokenString),
	})

	message := make(map[string]string)
	message["message"] = "success"

	// if there is an error while sending the email, add the token to the response
	if sendErr != nil {
		message["token"] = tokenString
	}

	server.Respond(w, message)
	return
}

func (route *APIRoute) handleNewPassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		err := server.ErrNotValidMethod.With("only post requests are allowed")
		server.RespondWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	var newPasswordData NewPasswordRequest
	decoder := json.NewDecoder(r.Body)
	decodeErr := decoder.Decode(&newPasswordData)
	if decodeErr != nil {
		server.RespondWithError(w, "couldn't decode the body", http.StatusBadRequest)
		return
	}

	errApplyNewPassword := route.Service.ApplyNewPasswordWithToken(&newPasswordData)

	if errApplyNewPassword != nil {
		server.RespondWithError(w, errApplyNewPassword.Error(), http.StatusBadRequest)
		return
	}

	message := make(map[string]string)
	message["message"] = "success"
	server.Respond(w, message)
	return
}
