package request

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/manuelrojas19/go-oauth2-server/models/oauth/granttype"
	"net/http"
	"strings"
)

// TokenRequest represents the request to obtain an access token.
type TokenRequest struct {
	ClientId     string
	ClientSecret string
	RefreshToken string
	GrantType    granttype.GrantType
	AuthCode     string
	RedirectUri  string
}

// DecodeTokenRequest function to handle URL encoded data and Authorization header.
func DecodeTokenRequest(r *http.Request, request *TokenRequest) error {
	// Parse URL encoded form data
	if err := r.ParseForm(); err != nil {
		return fmt.Errorf("failed to parse form data: %w", err)
	}

	// Extract grant_type from form data
	grantTypeStr := r.FormValue("grant_type")
	request.GrantType = granttype.GrantType(grantTypeStr)

	// Handle different grant types
	switch request.GrantType {
	case granttype.AuthorizationCode:
		// For these grant type, client credentials, code and redirect Uri are required
		request.ClientId = r.FormValue("client_id")
		request.ClientSecret = r.FormValue("client_secret")
		request.AuthCode = r.FormValue("code")
		request.RedirectUri = r.FormValue("redirect_uri")
	case granttype.Implicit, granttype.Password, granttype.ClientCredentials:
		// For these grant types, client credentials are required
		request.ClientId = r.FormValue("client_id")
		request.ClientSecret = r.FormValue("client_secret")
	case granttype.RefreshToken:
		// For Refresh Token grant type, client credentials are optional
		request.RefreshToken = r.FormValue("refresh_token")
		request.ClientId = r.FormValue("client_id")
		request.ClientSecret = r.FormValue("client_secret")
		// Check Authorization header if client credentials are not provided on form value
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" {
			if strings.HasPrefix(authHeader, "Basic ") {
				if err := parseBasicAuth(authHeader, request); err != nil {
					return err
				}
			} else {
				return errors.New("unsupported Authorization header format")
			}
		}
	default:
		return fmt.Errorf("unsupported grant_type: %s", request.GrantType)
	}
	return nil
}

// parseBasicAuth extracts client credentials from the Basic Authentication header.
func parseBasicAuth(authHeader string, request *TokenRequest) error {
	// Extract the Base64 encoded credentials
	encodedCredentials := strings.TrimPrefix(authHeader, "Basic ")
	decodedCredentials, err := base64.StdEncoding.DecodeString(encodedCredentials)
	if err != nil {
		return fmt.Errorf("failed to decode Basic Authentication header: %w", err)
	}
	// Split credentials into client Id and client secret
	credentials := strings.SplitN(string(decodedCredentials), ":", 2)
	if len(credentials) != 2 {
		return errors.New("invalid Authorization header format")
	}
	request.ClientId = credentials[0]
	request.ClientSecret = credentials[1]
	return nil
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
