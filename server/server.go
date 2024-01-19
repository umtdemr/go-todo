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
	http.ListenAndServe(s.ListenAddr, requestLogger(s.Router))
}

func Respond(w http.ResponseWriter, data interface{}, statusCode int) {
	jsonData, err := json.Marshal(data)

	if err != nil {
		RespondWithError(w, fmt.Sprintf("error while marshalling data: %v", err), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, writeError := w.Write(jsonData)

	if writeError != nil {
		RespondWithError(w, fmt.Sprintf("error while writing the data: %v", writeError), http.StatusBadRequest)
		return
	}
}

func RespondOK(w http.ResponseWriter, data interface{}) {
	Respond(w, data, http.StatusOK)
}

func RespondCreated(w http.ResponseWriter, data interface{}) {
	Respond(w, data, http.StatusCreated)
}

func RespondNoContent(w http.ResponseWriter, data interface{}) {
	Respond(w, data, http.StatusNoContent)
}

func RespondError(w http.ResponseWriter, msg string, errCode int, fields []string) {
	if msg == "" {
		msg = "An error has occurred while processing"
	}
	if errCode == 0 {
		errCode = http.StatusBadRequest
	}
	w.WriteHeader(errCode)

	resp := make(map[string]interface{})
	resp["message"] = msg
	if fields != nil {
		resp["fields"] = fields
	}

	encoder := json.NewEncoder(w)
	err := encoder.Encode(resp)

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"error": %v}`, err)
	}
}

func RespondWithError(w http.ResponseWriter, msg string, errCode int) {
	RespondError(w, msg, errCode, nil)
}

func RespondWithErrorFields(w http.ResponseWriter, msg string, errCode int, fields []string) {
	RespondError(w, msg, errCode, fields)
}
