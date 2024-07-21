package commands

import "github.com/manuelrojas19/go-oauth2-server/models/oauth/granttype"

type GrantAccessTokenCommand struct {
	ClientId     string
	ClientSecret string
	GrantType    granttype.GrantType
}

func NewGrantAccessTokenCommand(clientId string, clientSecret string, grantType granttype.GrantType) *GrantAccessTokenCommand {
	return &GrantAccessTokenCommand{ClientId: clientId, ClientSecret: clientSecret, GrantType: grantType}
}
