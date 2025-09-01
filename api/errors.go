package api

import "errors"

// Pre-defined error messages for API responses
var (
	ErrInvalidRequest          = errors.New("invalid_request")
	ErrUnauthorizedClient      = errors.New("unauthorized_client")
	ErrAccessDenied            = errors.New("access_denied")
	ErrUnsupportedResponseType = errors.New("unsupported_response_type")
	ErrInvalidScope            = errors.New("invalid_scope")
	ErrServerError             = errors.New("server_error")
	ErrTemporarilyUnavailable  = errors.New("temporarily_unavailable")
	ErrUnsupportedGrantType    = errors.New("unsupported_grant_type")
	ErrUnsupportedTokenType    = errors.New("unsupported_token_type")
	ErrInvalidClient           = errors.New("invalid_client")
	ErrInvalidGrant            = errors.New("invalid_grant")
	ErrInvalidToken            = errors.New("invalid_token")
	ErrClientAlreadyExists     = errors.New("client_already_exists")
)

// ErrorResponse represents a standard OAuth2 error response.
type ErrorResponse struct {
	Error            string `json:"error"`                       // A single ASCII error code.
	ErrorDescription string `json:"error_description,omitempty"` // A human-readable ASCII text providing additional information.
	ErrorURI         string `json:"error_uri,omitempty"`         // A URI identifying a human-readable web page with information about the error.
}
