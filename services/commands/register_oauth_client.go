package commands

import "github.com/manuelrojas19/go-oauth2-server/models/oauth"

type RegisterOauthClientCommand struct {
	ClientName              string
	GrantTypes              []oauth.GrantType
	ResponseTypes           []oauth.ResponseType
	TokenEndpointAuthMethod oauth.TokenEndpointAuthMethod
	RedirectUris            []string
}
