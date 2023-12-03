package respond

import (
	"encoding/json"
	"net/http"
)

// Respond responds as JSON
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
