package commands

import "github.com/manuelrojas19/go-oauth2-server/models/oauth/granttype"

type GrantAccessTokenCommand struct {
	ClientId     string
	ClientSecret string
	RefreshToken string
	GrantType    granttype.GrantType
	Code         string
	RedirectUri  string
}

func NewGrantAccessTokenCommand(clientId string, clientSecret string, grantType granttype.GrantType, refreshToken string, code string, redirectUri string) *GrantAccessTokenCommand {
	return &GrantAccessTokenCommand{
		ClientId:     clientId,
		ClientSecret: clientSecret,
		GrantType:    grantType,
		RefreshToken: refreshToken,
		Code:         code,
		RedirectUri:  redirectUri,
	}
}
