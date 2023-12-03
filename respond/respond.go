package respond

import (
	"encoding/json"
	"net/http"
)

// Respond responds as JSON
func Respond(w http.ResponseWriter, data interface{}) {
	encoder := json.NewEncoder(w)
	err := encoder.Encode(data)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
}
