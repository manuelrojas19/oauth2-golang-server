package request

import (
	"errors"
	"fmt"
	"github.com/manuelrojas19/go-oauth2-server/models/oauth/authmethodtype"
	"github.com/manuelrojas19/go-oauth2-server/models/oauth/granttype"
	"github.com/manuelrojas19/go-oauth2-server/models/oauth/responsetype"
	"strings"
)

type RegisterClientRequest struct {
	ClientName              string                                 `json:"client_name"`
	GrantTypes              []granttype.GrantType                  `json:"grant_type"`
	ResponseTypes           []responsetype.ResponseType            `json:"response_type"`
	TokenEndpointAuthMethod authmethodtype.TokenEndpointAuthMethod `json:"token_endpoint_auth_method"`
	RedirectUris            []string                               `json:"redirect_uris"`
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
		if !isValidGrantType(grantType) {
			return fmt.Errorf("invalid grant_type: %s", grantType)
		}
	}

	// Validate ResponseTypes
	if len(r.ResponseTypes) == 0 {
		return errors.New("at least one response_type is required")
	}
	for _, responseType := range r.ResponseTypes {
		if !isValidResponseType(responseType) {
			return fmt.Errorf("invalid response_type: %s", responseType)
		}
	}

	// Validate TokenEndpointAuthMethod
	if !isValidAuthMethod(r.TokenEndpointAuthMethod) {
		return fmt.Errorf("invalid token_endpoint_auth_method: %s", r.TokenEndpointAuthMethod)
	}

	// Validate RedirectUris (if specified)
	if r.RedirectUris != nil && len(r.RedirectUris) == 0 {
		return errors.New("redirect_uris must be a non-empty array if specified")
	}

	return nil
}

// Check if the GrantType is valid
func isValidGrantType(gt granttype.GrantType) bool {
	switch gt {
	case granttype.AuthorizationCode,
		granttype.Implicit,
		granttype.ClientCredentials,
		granttype.RefreshToken,
		granttype.Password:
		return true
	}
	return false
}

// Check if the ResponseType is valid
func isValidResponseType(rt responsetype.ResponseType) bool {
	switch rt {
	case responsetype.Code,
		responsetype.Token,
		responsetype.IDToken:
		return true
	}
	return false
}

// Check if the TokenEndpointAuthMethod is valid
func isValidAuthMethod(authMethod authmethodtype.TokenEndpointAuthMethod) bool {
	switch authMethod {
	case authmethodtype.ClientSecretBasic,
		authmethodtype.ClientSecretPost,
		authmethodtype.None:
		return true
	}
	return false
}
