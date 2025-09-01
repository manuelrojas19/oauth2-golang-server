package api

import (
	"errors"
	"fmt"
	"strings"

	"github.com/manuelrojas19/go-oauth2-server/oauth/authmethodtype"
	"github.com/manuelrojas19/go-oauth2-server/oauth/granttype"
	"github.com/manuelrojas19/go-oauth2-server/oauth/responsetype"
)

type RegisterClientRequest struct {
	ClientName              string                                 `json:"client_name"`
	GrantTypes              []granttype.GrantType                  `json:"grant_types"`
	ResponseTypes           []responsetype.ResponseType            `json:"response_types"`
	TokenEndpointAuthMethod authmethodtype.TokenEndpointAuthMethod `json:"token_endpoint_auth_method"`
	RedirectUris            []string                               `json:"redirect_uris"`
	Scopes                  string                                 `json:"scope"`
}

func (r *RegisterClientRequest) Sanitize() {
	r.ClientName = strings.TrimSpace(r.ClientName)
	r.Scopes = strings.TrimSpace(r.Scopes)
	for i, uri := range r.RedirectUris {
		r.RedirectUris[i] = strings.TrimSpace(uri)
	}
}

// Validate checks if the RegisterClientRequest is valid.
func (r *RegisterClientRequest) Validate() error {
	// Check ClientName
	if strings.TrimSpace(r.ClientName) == "" {
		return errors.New("client_name is required and cannot be empty")
	}

	// Validate GrantTypes
	if len(r.GrantTypes) == 0 {
		return errors.New("at least one grant_type is required")
	}
	for _, grantType := range r.GrantTypes {
		if !IsValidGrantType(grantType) {
			return fmt.Errorf("invalid grant_type: %s", grantType)
		}
	}

	// Validate ResponseTypes
	if len(r.ResponseTypes) == 0 {
		return errors.New("at least one response_type is required")
	}
	for _, responseType := range r.ResponseTypes {
		if !IsValidResponseType(responseType) {
			return fmt.Errorf("invalid response_type: %s", responseType)
		}
	}

	// Validate TokenEndpointAuthMethod
	if !IsValidAuthMethod(r.TokenEndpointAuthMethod) {
		return fmt.Errorf("invalid token_endpoint_auth_method: %s", r.TokenEndpointAuthMethod)
	}

	// Validate RedirectUris (if specified)
	if r.RedirectUris != nil {
		if len(r.RedirectUris) == 0 {
			return errors.New("redirect_uris must be a non-empty array if specified")
		}
		for _, uri := range r.RedirectUris {
			if uri == "" {
				return errors.New("redirect_uri cannot be empty")
			}
			if !IsValidRedirectURI(uri) {
				return fmt.Errorf("malformed redirect_uri: %s", uri)
			}
		}
	}

	return nil
}
