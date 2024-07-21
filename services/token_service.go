package services

import (
	"github.com/manuelrojas19/go-oauth2-server/services/commands"
	"github.com/manuelrojas19/go-oauth2-server/store/entities"
	"github.com/manuelrojas19/go-oauth2-server/store/repositories"
	"github.com/manuelrojas19/go-oauth2-server/utils"
	"time"
)

type tokenService struct {
	tokenRepository repositories.AccessTokenRepository
	client          OauthClientService
}

func NewTokenService(tokenRepository repositories.AccessTokenRepository, client OauthClientService) TokenService {
	return &tokenService{tokenRepository: tokenRepository, client: client}
}

func (t *tokenService) GrantAccessToken(command *commands.GrantAccessTokenCommand) (string, error) {
	client, err := t.client.FindOauthClient(command.ClientId)
	if err != nil {
		return "", err
	}

	if err := client.ValidateSecret(command.ClientSecret); err != nil {
		return "", err
	}

	generatedToken, err := utils.GenerateJWT(command.ClientId, "user", []byte("secret"))

	if err != nil {
		return "", err
	}

	token := entities.NewAccessToken(client, generatedToken, "", time.Now().Add(24*time.Hour))

	if token, err := t.tokenRepository.Save(token); err != nil {
		return "", err
	} else {
		return token.Token, nil
	}
}
