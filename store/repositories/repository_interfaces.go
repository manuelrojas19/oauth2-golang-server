package repositories

import "github.com/manuelrojas19/go-oauth2-server/store"

type OauthClientRepository interface {
	Save(client *store.OauthClient) (*store.OauthClient, error)
	FindByClientId(clientKey string) (*store.OauthClient, error)
}

type AccessTokenRepository interface {
	Save(token *store.AccessToken) (*store.AccessToken, error)
}

type ScopeRepository interface {
	FindByIdList(ids []string) ([]*store.Scope, error)
	Create(name, description string) (*store.Scope, error)
	FindById(id string) (*store.Scope, error)
	Exists(id string) (bool, error)
}

type RefreshTokenRepository interface {
	Save(token *store.RefreshToken) (*store.RefreshToken, error)
	FindByToken(token string) (*store.RefreshToken, error)
	InvalidateRefreshTokensByAccessTokenId(tokenId string) error
}

type AccessConsentRepository interface {
	HasUserConsented(userID, clientID, scope string) (bool, error)
	Save(userId, clientId, scope string) (*store.AccessConsent, error)
}

type AuthorizationRepository interface {
	Save(authCode *store.AuthCode) (*store.AuthCode, error)
	FindByCode(code string) (*store.AuthCode, error)
}

type UserRepository interface {
	Save(authCode *store.User) (*store.User, error)
	FindByUserId(id string) (*store.User, error)
}
