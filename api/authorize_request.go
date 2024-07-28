package api

import (
	"errors"
	"fmt"
	"github.com/manuelrojas19/go-oauth2-server/oauth/responsetype"
	"github.com/manuelrojas19/go-oauth2-server/utils"
	"net/http"
	"strings"
)

// AuthorizeRequest represents the request to authorize a client
type AuthorizeRequest struct {
	ResponseType responsetype.ResponseType
	ClientId     string
	RedirectUri  string
	Scope        string
	State        string
}

// DecodeAuthorizeRequest function to handle URL encoded data
func DecodeAuthorizeRequest(r *http.Request) (*AuthorizeRequest, error) {
	// Parse URL encoded form data
	if err := r.ParseForm(); err != nil {
		return nil, fmt.Errorf("failed to parse form data: %w", err)
	}

	// Convert response_type from string to responsetype.ResponseType
	responseType := responsetype.ResponseType(r.FormValue("response_type"))

	// Extract form data into AuthorizeRequest struct
	request := &AuthorizeRequest{
		ResponseType: responseType,
		ClientId:     r.FormValue("client_id"),
		RedirectUri:  r.FormValue("redirect_uri"),
		Scope:        r.FormValue("scope"),
		State:        r.FormValue("state"),
	}

	return request, nil
}

// Validate validates the fields of the AuthorizeRequest
func (r *AuthorizeRequest) Validate() error {
	// Validate ResponseType
	if strings.TrimSpace(string(r.ResponseType)) == "" {
		return errors.New("response_type is required")
	}
	if !utils.IsValidResponseType(r.ResponseType) {
		return fmt.Errorf("invalid response_type: %s", r.ResponseType)
	}

	// Validate ClientId
	if strings.TrimSpace(r.ClientId) == "" {
		return errors.New("client_id is required")
	}

	// Validate RedirectUri
	if strings.TrimSpace(r.RedirectUri) == "" {
		return errors.New("redirect_uri is required")
	}
	if !utils.IsValidRedirectUri(r.RedirectUri) {
		return fmt.Errorf("invalid redirect_uri: %s", r.RedirectUri)
	}

	// Optionally validate Scope (depending on your application's requirements)
	if strings.TrimSpace(r.Scope) != "" && !utils.IsValidScope(r.Scope) {
		return fmt.Errorf("invalid scope: %s", r.Scope)
	}

	// State is optional but can be validated if needed
	if r.State != "" && !utils.IsValidState(r.State) {
		return fmt.Errorf("invalid state: %s", r.State)
	}

	return nil
}
