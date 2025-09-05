package api

import "errors"

// Pre-defined error messages for API responses
var (
	ErrInvalidRequest          = errors.New("invalid_request")
	ErrInvalidRedirectUri      = errors.New("invalid_redirect_uri")
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

// errorDescriptions provides default human-readable descriptions for the API errors.
var errorDescriptions = map[error]string{
	ErrInvalidRequest:          "The request is missing a required parameter, includes an invalid parameter value, or is otherwise malformed.",
	ErrInvalidRedirectUri:      "One or more redirect URIs are invalid or missing.",
	ErrUnauthorizedClient:      "The client is not authorized to request an authorization code using this method.",
	ErrAccessDenied:            "The resource owner or authorization server denied the request.",
	ErrUnsupportedResponseType: "The authorization server does not support the requested response type.",
	ErrInvalidScope:            "The requested scope is invalid, unknown, or malformed.",
	ErrServerError:             "The authorization server encountered an unexpected condition that prevented it from fulfilling the request.",
	ErrTemporarilyUnavailable:  "The authorization server is currently unable to handle the request due to a temporary overloading or maintenance of the server.",
	ErrUnsupportedGrantType:    "The authorization grant type is not supported by the authorization server.",
	ErrUnsupportedTokenType:    "The authorization server does not support the requested token type.",
	ErrInvalidClient:           "Client authentication failed (e.g., unknown client, no client authentication included, or unsupported authentication method).",
	ErrInvalidGrant:            "The provided authorization grant (e.g., authorization code, refresh token) or refresh token is invalid, expired, revoked, does not match the redirection URI used in the authorization request, or was issued to another client.",
	ErrInvalidToken:            "The access token provided is expired, revoked, malformed, or invalid for other reasons.",
	ErrClientAlreadyExists:     "A client with the provided name already exists.",
}

// ErrorResponse represents a standard OAuth2 error response.
type ErrorResponse struct {
	Error            string `json:"error"`                       // A single ASCII error code.
	ErrorDescription string `json:"error_description,omitempty"` // A human-readable ASCII text providing additional information.
	ErrorURI         string `json:"error_uri,omitempty"`         // A URI identifying a human-readable web page with information about the error.
}
