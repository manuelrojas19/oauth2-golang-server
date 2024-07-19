package services

import (
	"github.com/manuelrojas19/go-oauth2-server/store/entities"
	"github.com/manuelrojas19/go-oauth2-server/store/repositories"
	"github.com/manuelrojas19/go-oauth2-server/utils"
	"time"
)

type accessTokenService struct {
	tokenRepository repositories.AccessTokenRepository
	client          OauthClientService
}

func NewAccessTokenService(tokenRepository repositories.AccessTokenRepository, client OauthClientService) AccessTokenService {
	return &accessTokenService{tokenRepository: tokenRepository, client: client}
}

func (t *accessTokenService) GrantAccessToken(clientKey string) (string, error) {
	client, err := t.client.FindOauthClient(clientKey)
	if err != nil {
		return "", err
	}

	generatedToken, err := utils.GenerateToken(clientKey, client.ClientId, time.Now())
	if err != nil {
		return "", err
	}

	token := entities.NewAccessToken(client, generatedToken, "GRANT", time.Now())

	if token, err := t.tokenRepository.Save(token); err != nil {
		return "", err
	} else {
		return token.Token, nil
	}
}
