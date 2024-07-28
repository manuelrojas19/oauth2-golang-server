package services

import (
	"github.com/lestrrat-go/jwx/jwk"
	oauth2 "github.com/manuelrojas19/go-oauth2-server/oauth"
	"github.com/manuelrojas19/go-oauth2-server/store"
)

type GrantService interface {
	ResolveGrantType(command *GrantAccessTokenCommand) (string, error)
}

type TokenService interface {
	GrantAccessToken(command *GrantAccessTokenCommand) (*oauth2.Token, error)
}

type AuthorizationService interface {
	Authorize(command *AuthorizeCommand) (*oauth2.AuthCode, error)
}

type OauthClientService interface {
	CreateOauthClient(command *RegisterOauthClientCommand) (*oauth2.Client, error)
	FindOauthClient(clientId string) (*store.OauthClient, error)
}

type WellKnownService interface {
	GetJwk() (*jwk.Set, error)
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
