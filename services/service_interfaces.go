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

type AuthorizationService interface {
	Authorize(command *commands.Authorize) (*oauth.AuthCode, error)
}

type UserConsentService interface {
	Save(userID, clientID, scope string) error
	HasUserConsented(userID, clientID, scope string) bool
}

type SessionService interface {
	CreateSession(userId, email string) (string, error)
	SessionExists(sessionID string) bool
	GetUserIdFromSession(sessionID string) (string, error)
}
