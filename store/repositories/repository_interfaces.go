package repositories

import "github.com/manuelrojas19/go-oauth2-server/store/entities"

type OauthClientRepository interface {
	Save(client *entities.OauthClient) (*entities.OauthClient, error)
	FindByClientId(clientKey string) (*entities.OauthClient, error)
}

type AccessTokenRepository interface {
	Save(token *entities.AccessToken) (*entities.AccessToken, error)
}

type RefreshTokenRepository interface {
	Save(token *entities.RefreshToken) (*entities.RefreshToken, error)
	FindByToken(token string) (*entities.RefreshToken, error)
}
