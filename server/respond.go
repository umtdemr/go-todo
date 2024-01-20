package server

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Respond is a helper function to respond with the given data and status code
func Respond(w http.ResponseWriter, data interface{}, statusCode int) {
	jsonData, err := json.Marshal(data)

	// if there is an error while marshalling the data, respond with an error
	if err != nil {
		RespondWithError(w, fmt.Sprintf("error while marshalling data: %v", err), http.StatusBadRequest)
		return
	}

	// set the content type header
	w.Header().Set("Content-Type", "application/json")

	// set the status code
	w.WriteHeader(statusCode)

	// write the data
	_, writeError := w.Write(jsonData)

	// if there is an error while writing the data, respond with an error
	if writeError != nil {
		RespondWithError(w, fmt.Sprintf("error while writing the data: %v", writeError), http.StatusBadRequest)
		return
	}
}

// RespondOK is a helper function to respond with 200 OK
func RespondOK(w http.ResponseWriter, data interface{}) {
	Respond(w, data, http.StatusOK)
}

// RespondCreated is a helper function to respond with 201 Created
func RespondCreated(w http.ResponseWriter, data interface{}) {
	Respond(w, data, http.StatusCreated)
}

// RespondNoContent is a helper function to respond with 204 No Content
func RespondNoContent(w http.ResponseWriter, data interface{}) {
	Respond(w, data, http.StatusNoContent)
}

// RespondError is a helper function to respond with an error
// fields is an optional parameter to specify the fields that caused the error
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

// RespondWithError is an alias for RespondError with fields set to nil
func RespondWithError(w http.ResponseWriter, msg string, errCode int) {
	RespondError(w, msg, errCode, nil)
}

// RespondWithErrorFields is an alias function for RespondError with the given fields
func RespondWithErrorFields(w http.ResponseWriter, msg string, errCode int, fields []string) {
	RespondError(w, msg, errCode, fields)
}
