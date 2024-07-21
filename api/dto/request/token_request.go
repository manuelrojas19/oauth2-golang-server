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
	GrantType    granttype.GrantType `json:"grant_type"`
}

func (r *TokenRequest) Validate() error {
	// Check ClientName
	if strings.TrimSpace(r.ClientId) == "" {
		return errors.New("client_id is required and cannot be empty")
	}

	// Check ClientSecret
	if strings.TrimSpace(r.ClientSecret) == "" {
		return errors.New("client_secret is required and cannot be empty")
	}

	// Validate GrantType
	if !isValidGrantType(r.GrantType) {
		return fmt.Errorf("invalid grant_type: %s", r.GrantType)
	}

	return nil
}
