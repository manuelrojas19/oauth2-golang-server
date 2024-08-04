package api

import (
	"errors"
	"fmt"
	"github.com/manuelrojas19/go-oauth2-server/oauth"
	"github.com/manuelrojas19/go-oauth2-server/oauth/authmethodtype"
	"github.com/manuelrojas19/go-oauth2-server/oauth/granttype"
	"github.com/manuelrojas19/go-oauth2-server/oauth/responsetype"
	"github.com/manuelrojas19/go-oauth2-server/utils"
	"strings"
)

type RegisterClientRequest struct {
	ClientName              string                                 `json:"client_name"`
	GrantTypes              []granttype.GrantType                  `json:"grant_types"`
	ResponseTypes           []responsetype.ResponseType            `json:"response_types"`
	TokenEndpointAuthMethod authmethodtype.TokenEndpointAuthMethod `json:"token_endpoint_auth_method"`
	RedirectUris            []string                               `json:"redirect_uris"`
	Scopes                  []oauth.Scope                          `json:"scopes"`
}

// Validate checks if the values of enum fields are valid
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
		if !utils.IsValidGrantType(grantType) {
			return fmt.Errorf("invalid grant_type: %s", grantType)
		}
	}

	// Validate ResponseTypes
	if len(r.ResponseTypes) == 0 {
		return errors.New("at least one response_type is required")
	}
	for _, responseType := range r.ResponseTypes {
		if !utils.IsValidResponseType(responseType) {
			return fmt.Errorf("invalid response_type: %s", responseType)
		}
	}

	// Validate TokenEndpointAuthMethod
	if !utils.IsValidAuthMethod(r.TokenEndpointAuthMethod) {
		return fmt.Errorf("invalid token_endpoint_auth_method: %s", r.TokenEndpointAuthMethod)
	}

	// Validate RedirectUris (if specified)
	if r.RedirectUris != nil && len(r.RedirectUris) == 0 {
		return errors.New("redirect_uris must be a non-empty array if specified")
	}

	return nil
}
