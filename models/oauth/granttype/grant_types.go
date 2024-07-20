package granttype

import "fmt"

type GrantType string

const (
	AuthorizationCode GrantType = "authorization_code"
	Implicit          GrantType = "implicit"
	Password          GrantType = "password"
	ClientCredentials GrantType = "client_credentials"
	RefreshToken      GrantType = "refresh_token"
)

// EnumListToStringList Convert a list of GrantType to a list of strings
func EnumListToStringList(grantTypes []GrantType) []string {
	var strings []string
	for _, gt := range grantTypes {
		strings = append(strings, string(gt))
	}
	return strings
}

// StringListToEnumList Convert a list of strings to a list of GrantType
func StringListToEnumList(strings []string) []GrantType {
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
