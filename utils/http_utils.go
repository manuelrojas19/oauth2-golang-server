package utils

import (
	"encoding/json"
	"net/http"
)

// DecodeJSON parses the request body into the provided requestBody object.
func DecodeJSON[T any](r *http.Request, requestBody *T) error {
	return json.NewDecoder(r.Body).Decode(requestBody)
}

// RespondWithJSON sends a JSON response type with the specified status code and body.
func RespondWithJSON(w http.ResponseWriter, statusCode int, responseBody interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	responseError := json.NewEncoder(w).Encode(responseBody)
	if responseError != nil {
		http.Error(w, "Failed to encode response type", http.StatusInternalServerError)
	}
}

func ErrorResponseBody(err error) map[string]string {
	return map[string]string{"error": err.Error()}
}
