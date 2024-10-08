package api

import (
	"errors"
	"fmt"
	"github.com/manuelrojas19/go-oauth2-server/oauth/responsetype"
	"github.com/manuelrojas19/go-oauth2-server/utils"
	"net/http"
	"net/url"
	"regexp"
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

	err := sanitizeAuthorizeRequest(request)
	if err != nil {
		return nil, err
	}

	return request, nil
}

// sanitizeAuthorizeRequest sanitizes and validates the input values of an AuthorizeRequest
func sanitizeAuthorizeRequest(request *AuthorizeRequest) error {
	// Trim whitespace from all fields
	request.ResponseType = responsetype.ResponseType(strings.TrimSpace(string(request.ResponseType)))
	request.ClientId = strings.TrimSpace(request.ClientId)
	request.RedirectUri = strings.TrimSpace(request.RedirectUri)
	request.Scope = strings.TrimSpace(request.Scope)
	request.State = strings.TrimSpace(request.State)

	// Validate ClientId length
	if len(request.ClientId) < 1 || len(request.ClientId) > 256 {
		return errors.New("client_id length is invalid")
	}

	// Optionally validate State length
	if len(request.State) > 256 {
		return errors.New("state length is invalid")
	}

	// Validate RedirectUri
	if _, err := url.ParseRequestURI(request.RedirectUri); err != nil {
		return errors.New("redirect_uri is invalid")
	}

	// Additional checks for potential injection attacks
	if containsInjectionPatterns(request.ClientId) || containsInjectionPatterns(request.State) {
		return errors.New("client_id or state contains invalid characters")
	}

	return nil
}

// containsInjectionPatterns checks for common injection patterns
func containsInjectionPatterns(s string) bool {
	// Define common patterns for injection attacks
	injectionPatterns := []string{
		`(<[^>]+>)`, // HTML tags
		`(\bselect\b|\bunion\b|\bupdate\b|\bdelete\b|\binsert\b)`, // SQL injection
		`(\bscript\b)`, // JavaScript
	}

	for _, pattern := range injectionPatterns {
		matched, _ := regexp.MatchString(pattern, s)
		if matched {
			return true
		}
	}
	return false
}

// Validate validates the fields of the AuthorizeRequest
func (r *AuthorizeRequest) Validate() error {
	// Validate ResponseType
	if strings.TrimSpace(string(r.ResponseType)) == "" {
		return errors.New("response_type is required")
	}
	if !utils.IsValidResponseType(r.ResponseType) {
		return fmt.Errorf("the authorization server does not support obtaining an authorization code using this method")
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
		return fmt.Errorf("the requested scope is invalid, unknown, or malformed")
	}

	// State is optional but can be validated if needed
	if r.State != "" && !utils.IsValidState(r.State) {
		return fmt.Errorf("invalid state: %s", r.State)
	}

	return nil
}
