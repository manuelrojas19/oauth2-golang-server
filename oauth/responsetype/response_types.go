package responsetype

import "fmt"

type ResponseType string

const (
	Code    ResponseType = "code"
	Token   ResponseType = "token"
	IDToken ResponseType = "id_token"
)

// EnumListToStringList Convert a list of ResponseType to a list of strings
func EnumListToStringList(responseTypes []ResponseType) []string {
	var strings []string
	for _, rt := range responseTypes {
		strings = append(strings, string(rt))
	}
	return strings
}

// StringListToEnumList Convert a list of strings to a list of ResponseType
func StringListToEnumList(strings []string) []ResponseType {
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
