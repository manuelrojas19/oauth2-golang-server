package repositories

import "github.com/manuelrojas19/go-oauth2-server/store"

type OauthClientRepository interface {
	Save(client *store.OauthClient) (*store.OauthClient, error)
	FindByClientId(clientKey string) (*store.OauthClient, error)
}

type AccessTokenRepository interface {
	Save(token *store.AccessToken) (*store.AccessToken, error)
}

type RefreshTokenRepository interface {
	Save(token *store.RefreshToken) (*store.RefreshToken, error)
	FindByToken(token string) (*store.RefreshToken, error)
	InvalidateRefreshTokensByAccessTokenId(tokenId string) error
}

type UserConsentRepository interface {
	HasUserConsented(userID, clientID, scope string) (bool, error)
	Save(userID, clientID, scope string) (bool, error)
}

type AuthorizationRepository interface {
	Save(authCode *store.AuthCode) (*store.AuthCode, error)
	FindByCode(code string) (*store.AuthCode, error)
}

type UserRepository interface {
	Save(authCode *store.User) (*store.User, error)
	FindByUserId(id string) (*store.User, error)
}
