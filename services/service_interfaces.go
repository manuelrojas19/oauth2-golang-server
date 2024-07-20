package services

import (
	"github.com/manuelrojas19/go-oauth2-server/models/oauth"
	"github.com/manuelrojas19/go-oauth2-server/services/commands"
	"github.com/manuelrojas19/go-oauth2-server/store/entities"
)

type AccessTokenService interface {
	GrantAccessToken(clientKey string) (string, error)
}

type OauthClientService interface {
	CreateOauthClient(command *commands.RegisterOauthClientCommand) (*oauth.Client, error)
	FindOauthClient(clientKey string) (*entities.OauthClient, error)
}
