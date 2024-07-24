package services

import (
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/manuelrojas19/go-oauth2-server/models/oauth"
	"github.com/manuelrojas19/go-oauth2-server/services/commands"
	"github.com/manuelrojas19/go-oauth2-server/store/entities"
)

type GrantService interface {
	ResolveGrantType(command *commands.GrantAccessTokenCommand) (string, error)
}

type TokenService interface {
	GrantAccessToken(command *commands.GrantAccessTokenCommand) (*oauth.Token, error)
}

type OauthClientService interface {
	CreateOauthClient(command *commands.RegisterOauthClientCommand) (*oauth.Client, error)
	FindOauthClient(clientId string) (*entities.OauthClient, error)
}

type WellKnownService interface {
	GetJwk() (*jwk.Set, error)
}
