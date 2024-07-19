package utils

import (
	"encoding/json"
	"github.com/manuelrojas19/go-oauth2-server/models/request"
	"net/http"
)

// Decode parses the request body into the provided requestBody object.
func Decode(r *http.Request, requestBody *request.RegisterClientRequest) error {
	return json.NewDecoder(r.Body).Decode(requestBody)
}

// RespondWithJSON sends a JSON response with the specified status code and body.
func RespondWithJSON(w http.ResponseWriter, statusCode int, responseBody interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(responseBody); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
