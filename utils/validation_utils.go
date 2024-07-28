package utils

import (
	"github.com/manuelrojas19/go-oauth2-server/oauth/authmethodtype"
	"github.com/manuelrojas19/go-oauth2-server/oauth/granttype"
	"github.com/manuelrojas19/go-oauth2-server/oauth/responsetype"
	"log"
	"net/url"
	"regexp"
)

// IsValidGrantType checks if the GrantType is valid
func IsValidGrantType(gt granttype.GrantType) bool {
	switch gt {
	case granttype.AuthorizationCode,
		granttype.Implicit,
		granttype.ClientCredentials,
		granttype.RefreshToken,
		granttype.Password:
		log.Printf("Valid grant type: %s", gt)
		return true
	}
	log.Printf("Invalid grant type: %s", gt)
	return false
}

// IsValidResponseType checks if the ResponseType is valid
func IsValidResponseType(rt responsetype.ResponseType) bool {
	switch rt {
	case responsetype.Code,
		responsetype.Token,
		responsetype.IDToken:
		log.Printf("Valid response type: %s", rt)
		return true
	}
	log.Printf("Invalid response type: %s", rt)
	return false
}

// IsValidAuthMethod checks if the TokenEndpointAuthMethod is valid
func IsValidAuthMethod(authMethod authmethodtype.TokenEndpointAuthMethod) bool {
	switch authMethod {
	case authmethodtype.ClientSecretBasic,
		authmethodtype.ClientSecretPost,
		authmethodtype.None:
		log.Printf("Valid authentication method: %s", authMethod)
		return true
	}
	log.Printf("Invalid authentication method: %s", authMethod)
	return false
}

// IsValidRedirectUri checks if the redirect_uri is a valid URL
func IsValidRedirectUri(redirectUri string) bool {
	_, err := url.ParseRequestURI(redirectUri)
	if err != nil {
		log.Printf("Invalid redirect URI: %s, error: %v", redirectUri, err)
		return false
	}
	log.Printf("Valid redirect URI: %s", redirectUri)
	return true
}

// IsValidScope checks if the scope is valid (example validation, adapt as needed)
func IsValidScope(scope string) bool {
	// Example: ensure scope contains only alphanumeric characters and spaces
	matched, err := regexp.MatchString(`^[a-zA-Z0-9 ]*$`, scope)
	if err != nil {
		log.Printf("Error validating scope: %s, error: %v", scope, err)
		return false
	}
	if matched {
		log.Printf("Valid scope: %s", scope)
		return true
	}
	log.Printf("Invalid scope: %s", scope)
	return false
}

// IsValidState checks if the state is valid (example validation, adapt as needed)
func IsValidState(state string) bool {
	// Example: ensure state is non-empty and contains only alphanumeric characters
	matched, err := regexp.MatchString(`^[a-zA-Z0-9]*$`, state)
	if err != nil {
		log.Printf("Error validating state: %s, error: %v", state, err)
		return false
	}
	if matched {
		log.Printf("Valid state: %s", state)
		return true
	}
	log.Printf("Invalid state: %s", state)
	return false
}
