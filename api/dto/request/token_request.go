package request

import (
	"errors"
	"fmt"
	"github.com/manuelrojas19/go-oauth2-server/models/oauth/granttype"
	"strings"
)

// TokenRequest represents the request to obtain an access token.
type TokenRequest struct {
	ClientId     string              `json:"client_id"`
	ClientSecret string              `json:"client_secret"`
	RefreshToken string              `json:"refresh_token"`
	GrantType    granttype.GrantType `json:"grant_type"`
}

func (r *TokenRequest) Validate() error {
	// Validate GrantType
	if !isValidGrantType(r.GrantType) {
		return fmt.Errorf("invalid grant_type: %s", r.GrantType)
	}

	// Validate ClientId and ClientSecret based on GrantType
	switch r.GrantType {
	case granttype.AuthorizationCode, granttype.Implicit, granttype.Password, granttype.ClientCredentials:
		// Ensure ClientId and ClientSecret are not empty for these grant types
		if strings.TrimSpace(r.ClientId) == "" {
			return errors.New("client_id is required for the grant_type: " + string(r.GrantType))
		}
		if strings.TrimSpace(r.ClientSecret) == "" {
			return errors.New("client_secret is required for the grant_type: " + string(r.GrantType))
		}

	case granttype.RefreshToken:
		// ClientId and ClientSecret are optional for Refresh Token Grant Type
		if strings.TrimSpace(r.RefreshToken) == "" {
			return errors.New("refresh_token is required for refresh_token grant type")
		}

	default:
		return fmt.Errorf("unsupported grant_type: %s", r.GrantType)
	}

	return nil
}
