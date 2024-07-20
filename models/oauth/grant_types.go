package oauth

import "fmt"

type GrantType string

const (
	AuthorizationCode GrantType = "authorization_code"
	Implicit          GrantType = "implicit"
	Password          GrantType = "password"
	ClientCredentials GrantType = "client_credentials"
	RefreshToken      GrantType = "refresh_token"
)

type ResponseType string

const (
	Code    ResponseType = "code"
	Token   ResponseType = "token"
	IDToken ResponseType = "id_token"
)

type TokenEndpointAuthMethod string

const (
	ClientSecretBasic TokenEndpointAuthMethod = "client_secret_basic"
	ClientSecretPost  TokenEndpointAuthMethod = "client_secret_post"
	None              TokenEndpointAuthMethod = "none"
)

// Convert a list of GrantType to a list of strings
func GrantTypeListToStringList(grantTypes []GrantType) []string {
	var strings []string
	for _, gt := range grantTypes {
		strings = append(strings, string(gt))
	}
	return strings
}

// Convert a list of ResponseType to a list of strings
func ResponseTypeListToStringList(responseTypes []ResponseType) []string {
	var strings []string
	for _, rt := range responseTypes {
		strings = append(strings, string(rt))
	}
	return strings
}

// Convert a list of TokenEndpointAuthMethod to a list of strings
func TokenEndpointAuthMethodListToStringList(authMethods []TokenEndpointAuthMethod) []string {
	var strings []string
	for _, am := range authMethods {
		strings = append(strings, string(am))
	}
	return strings
}

// Convert a list of strings to a list of GrantType
func StringListToGrantTypeList(strings []string) []GrantType {
	var grantTypes []GrantType
	for _, s := range strings {
		switch GrantType(s) {
		case AuthorizationCode, Implicit, Password, ClientCredentials, RefreshToken:
			grantTypes = append(grantTypes, GrantType(s))
		default:
			_ = fmt.Errorf("invalid GrantType: %s", s)
		}
	}
	return grantTypes
}

// Convert a list of strings to a list of ResponseType
func StringListToResponseTypeList(strings []string) []ResponseType {
	var responseTypes []ResponseType
	for _, s := range strings {
		switch ResponseType(s) {
		case Code, Token, IDToken:
			responseTypes = append(responseTypes, ResponseType(s))
		default:
			_ = fmt.Errorf("invalid GrantType: %s", s)
		}
	}
	return responseTypes
}

// Convert a list of strings to a list of TokenEndpointAuthMethod
func StringListToTokenEndpointAuthMethodList(strings []string) []TokenEndpointAuthMethod {
	var authMethods []TokenEndpointAuthMethod
	for _, s := range strings {
		switch TokenEndpointAuthMethod(s) {
		case ClientSecretBasic, ClientSecretPost, None:
			authMethods = append(authMethods, TokenEndpointAuthMethod(s))
		default:
			_ = fmt.Errorf("invalid GrantType: %s", s)
		}
	}
	return authMethods
}
