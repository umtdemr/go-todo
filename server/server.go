package server

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

type APIServer struct {
	ListenAddr string
	Router     *mux.Router
}

func NewAPIServer(listenAddr string) *APIServer {
	return &APIServer{ListenAddr: listenAddr, Router: mux.NewRouter()}
}

func (s *APIServer) Run() {
	http.ListenAndServe(s.ListenAddr, s.Router)
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
