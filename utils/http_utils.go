package utils

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/manuelrojas19/go-oauth2-server/api"
	"go.uber.org/zap"
)

// DecodeJSON parses the request body into the provided requestBody object.
func DecodeJSON[T any](r *http.Request, requestBody *T) error {
	return json.NewDecoder(r.Body).Decode(requestBody)
}

// RespondWithJSON sends a JSON response with the given status code and payload.
func RespondWithJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(response)
}

// HandleErrorResponse centralizes error handling for HTTP responses.
// It maps specific API errors to appropriate HTTP status codes and responds with an ErrorResponseBody.
func HandleErrorResponse(w http.ResponseWriter, logger *zap.Logger, err error) {
	
	logger.Error("Error processing request", zap.Error(err))

	var status int
	var apiError api.ErrorResponse

	switch {
	case errors.Is(err, api.ErrInvalidClient):
		status = http.StatusUnauthorized
		apiError = api.ErrorResponseBody(api.ErrInvalidClient)
	case errors.Is(err, api.ErrInvalidGrant):
		status = http.StatusBadRequest
		apiError = api.ErrorResponseBody(api.ErrInvalidGrant)
	case errors.Is(err, api.ErrUnsupportedGrantType):
		status = http.StatusBadRequest
		apiError = api.ErrorResponseBody(api.ErrUnsupportedGrantType)
	case errors.Is(err, api.ErrInvalidScope):
		status = http.StatusBadRequest
		apiError = api.ErrorResponseBody(api.ErrInvalidScope)
	case errors.Is(err, api.ErrInvalidRequest):
		status = http.StatusBadRequest
		apiError = api.ErrorResponseBody(api.ErrInvalidRequest)
	case errors.Is(err, api.ErrInvalidRedirectUri):
		status = http.StatusBadRequest
		apiError = api.ErrorResponseBody(api.ErrInvalidRedirectUri)
	case errors.Is(err, api.ErrClientAlreadyExists):
		status = http.StatusConflict
		apiError = api.ErrorResponseBody(api.ErrClientAlreadyExists, err.Error())
	default:
		status = http.StatusInternalServerError
		apiError = api.ErrorResponseBody(api.ErrServerError)
	}
	RespondWithJSON(w, status, apiError)
}
