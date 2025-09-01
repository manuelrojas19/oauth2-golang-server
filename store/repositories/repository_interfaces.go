package repositories

import "github.com/manuelrojas19/go-oauth2-server/store"

type OauthClientRepository interface {
	Save(client *store.OauthClient) (*store.OauthClient, error)
	FindByClientId(clientKey string) (*store.OauthClient, error)
	ExistsByName(clientName string) bool
}

type AccessTokenRepository interface {
	Save(token *store.AccessToken) (*store.AccessToken, error)
	FindByAccessToken(accessToken string) (*store.AccessToken, error)
	DeleteByAccessToken(accessToken string) error
}

type ScopeRepository interface {
	FindByIdList(ids []string) ([]*store.Scope, error)
	Create(name, description string) (*store.Scope, error)
	FindById(id string) (*store.Scope, error)
	Exists(id string) (bool, error)
	FindByName(name string) (*store.Scope, error)
}

type RefreshTokenRepository interface {
	Save(token *store.RefreshToken) (*store.RefreshToken, error)
	FindByRefreshToken(token string) (*store.RefreshToken, error)
	InvalidateRefreshTokensByAccessTokenId(tokenId string) error
	DeleteByRefreshToken(refreshToken string) error
}

type AccessConsentRepository interface {
	HasUserConsented(userID, clientID, scope string) (bool, error)
	Save(userId, clientId, scope string) (*store.AccessConsent, error)
}

type AuthorizationRepository interface {
	Save(authCode *store.AuthCode) (*store.AuthCode, error)
	FindByCode(code string) (*store.AuthCode, error)
	Delete(code string) error
}

type UserRepository interface {
	Save(authCode *store.User) (*store.User, error)
	FindByUserId(id string) (*store.User, error)
	FindById(id string) (*store.User, error)
}
