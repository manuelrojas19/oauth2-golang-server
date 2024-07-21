package services

import (
	"github.com/manuelrojas19/go-oauth2-server/models/oauth"
	"github.com/manuelrojas19/go-oauth2-server/services/commands"
	"github.com/manuelrojas19/go-oauth2-server/store/entities"
	"github.com/manuelrojas19/go-oauth2-server/store/repositories"
	"github.com/manuelrojas19/go-oauth2-server/utils"
	"time"
)

const RefreshTokenDuration = 30 * 24 * time.Hour
const AccessTokenDuration = 1 * time.Hour

type tokenService struct {
	accessTokenRepository  repositories.AccessTokenRepository
	refreshTokenRepository repositories.RefreshTokenRepository
	client                 OauthClientService
}

func NewTokenService(tokenRepository repositories.AccessTokenRepository,
	refreshTokenRepository repositories.RefreshTokenRepository,
	client OauthClientService) TokenService {
	return &tokenService{accessTokenRepository: tokenRepository,
		refreshTokenRepository: refreshTokenRepository,
		client:                 client}
}

func (t *tokenService) GrantAccessToken(command *commands.GrantAccessTokenCommand) (*oauth.Token, error) {
	// Step 1: Retrieve and validate client
	client, err := t.client.FindOauthClient(command.ClientId)
	if err != nil {
		return nil, err
	}

	if err := client.ValidateSecret(command.ClientSecret); err != nil {
		return nil, err
	}

	// Step 2: Generate Access Token
	accessTokenJwt, err := utils.GenerateJWT(command.ClientId, "user", []byte("secret"), "access")
	if err != nil {
		return nil, err
	}

	// Create and save Access Token
	accessToken := entities.NewAccessTokenBuilder().
		WithClient(client).
		WithClientId(command.ClientId).
		WithToken(accessTokenJwt).
		WithTokenType("JWT").
		WithExpiresAt(time.Now().Add(RefreshTokenDuration)). // Example expiration
		Build()

	accessToken, err = t.accessTokenRepository.Save(accessToken)
	if err != nil {
		return nil, err
	}

	// Step 3: Generate Refresh Token
	refreshTokenJwt, err := utils.GenerateJWT(command.ClientId, "user", []byte("secret"), "refresh")
	if err != nil {
		return nil, err
	}

	// Create and save Refresh Token

	refreshToken := entities.NewRefreshTokenBuilder().
		WithAccessToken(accessToken).
		WithAccessTokenId(accessToken.Id).
		WithClient(client).
		WithClientId(command.ClientId).
		WithToken(refreshTokenJwt).
		WithTokenType("JWT").
		WithExpiresAt(time.Now().Add(AccessTokenDuration)). // Example expiration
		Build()

	refreshToken, err = t.refreshTokenRepository.Save(refreshToken)
	if err != nil {
		return nil, err
	}

	// Step 4: Build and return the Token response
	token := oauth.NewTokenBuilder().
		WithClientId(client.ClientId).
		WithUserId("user"). // Assuming a static user ID or replace with dynamic value
		WithAccessToken(accessTokenJwt).
		WithAccessTokenCreatedAt(time.Now()).
		WithAccessTokenExpiresAt(AccessTokenDuration).
		WithRefreshToken(refreshTokenJwt).
		WithRefreshTokenCreatedAt(time.Now()).
		WithRefreshTokenExpiresAt(RefreshTokenDuration).
		WithExtension(nil). // If you have any extensions, set them here
		Build()

	return token, nil
}
