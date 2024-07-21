package utils

import (
	"encoding/json"
	"github.com/manuelrojas19/go-oauth2-server/api/dto/request"
	"github.com/manuelrojas19/go-oauth2-server/models/oauth/granttype"
	"net/http"
)

// Decode parses the request body into the provided requestBody object.
func Decode[T any](r *http.Request, requestBody *T) error {
	return json.NewDecoder(r.Body).Decode(requestBody)
}

// DecodeTokenRequest function to handle URL encoded data.
func DecodeTokenRequest(r *http.Request, requestBody *request.TokenRequest) error {
	// Parse URL encoded form data
	if err := r.ParseForm(); err != nil {
		return err
	}

	// Populate requestBody fields from form values
	requestBody.ClientId = r.FormValue("client_id")
	requestBody.ClientSecret = r.FormValue("client_secret")
	requestBody.GrantType = granttype.GrantType(r.FormValue("grant_type"))
	requestBody.RefreshToken = r.FormValue("refresh_token")

	return nil
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
