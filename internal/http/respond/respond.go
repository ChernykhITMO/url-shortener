package respond

import (
	"encoding/json"
	"net/http"
)

type ErrorBody struct {
	Error string `json:"error"`
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func WriteJSONError(w http.ResponseWriter, status int, message string) error {
	return WriteJSON(w, status, ErrorBody{Error: message})
}
