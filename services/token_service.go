package services

import (
	"fmt"
	"github.com/manuelrojas19/go-oauth2-server/models/oauth"
	"github.com/manuelrojas19/go-oauth2-server/models/oauth/granttype"
	"github.com/manuelrojas19/go-oauth2-server/services/commands"
	"github.com/manuelrojas19/go-oauth2-server/store/entities"
	"github.com/manuelrojas19/go-oauth2-server/store/repositories"
	"github.com/manuelrojas19/go-oauth2-server/utils"
	"log"
	"time"
)

const AccessTokenDuration = 1 * time.Hour
const RefreshTokenDuration = 30 * 24 * time.Hour

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
	switch command.GrantType {
	case granttype.ClientCredentials:
		// Handle Client Credentials Grant
		return t.getTokenByClientCredentialsFlow(command.ClientId, command.ClientSecret)
	case granttype.RefreshToken:
		// Handle Refresh Token Grant
		return t.getTokenByRefreshTokenFlow(command.RefreshToken, command.ClientId, command.ClientSecret)
	default:
		return nil, fmt.Errorf("unsupported grant type: %s", command.GrantType)
	}
}

func (t *tokenService) getTokenByClientCredentialsFlow(clientId, clientSecret string) (*oauth.Token, error) {
	// Step 1: Retrieve and validate client
	client, err := t.client.FindOauthClient(clientId)
	if err != nil {
		return nil, err
	}

	if err := client.ValidateSecret(clientSecret); err != nil {
		return nil, err
	}

	// Step 2: Generate Access Token
	accessTokenJwt, err := utils.GenerateJWT(clientId, "user", []byte("secret"), "access")
	if err != nil {
		return nil, err
	}

	// Step 2: Generate Access Token
	accessTokenJwe, err := utils.GenerateJWE(accessTokenJwt)
	if err != nil {
		return nil, err
	}

	// Create and save Access Token
	accessToken := entities.NewAccessTokenBuilder().
		WithClient(client).
		WithClientId(clientId).
		WithToken(accessTokenJwe).
		WithTokenType("JWE").
		WithExpiresAt(time.Now().Add(AccessTokenDuration)). // Example expiration
		Build()

	accessToken, err = t.accessTokenRepository.Save(accessToken)
	if err != nil {
		return nil, err
	}

	// Step 3: Generate Refresh Token
	refreshTokenJwt, err := utils.GenerateJWT(clientId, "user", []byte("secret"), "refresh")
	if err != nil {
		return nil, err
	}

	// Create and save Refresh Token

	refreshToken := entities.NewRefreshTokenBuilder().
		WithAccessToken(accessToken).
		WithAccessTokenId(accessToken.Id).
		WithClient(client).
		WithClientId(clientId).
		WithToken(refreshTokenJwt).
		WithTokenType("JWT").
		WithExpiresAt(time.Now().Add(RefreshTokenDuration)). // Example expiration
		Build()

	refreshToken, err = t.refreshTokenRepository.Save(refreshToken)
	if err != nil {
		return nil, err
	}

	// Step 4: Build and return the Token response
	token := oauth.NewTokenBuilder().
		WithClientId(client.ClientId).
		WithUserId("user"). // Assuming a static user ID or replace with dynamic value
		WithAccessToken(accessToken.Token).
		WithAccessTokenCreatedAt(accessToken.CreatedAt).
		WithAccessTokenExpiresAt(accessToken.ExpiresAt.Sub(time.Now())).
		WithRefreshToken(refreshToken.Token).
		WithRefreshTokenCreatedAt(refreshToken.CreatedAt).
		WithRefreshTokenExpiresAt(refreshToken.ExpiresAt.Sub(time.Now())).
		WithExtension(nil). // If you have any extensions, set them here
		Build()

	return token, nil
}

func (t *tokenService) getTokenByRefreshTokenFlow(token, clientId, clientSecret string) (*oauth.Token, error) {
	log.Println("Received refresh token request for client:", clientId)

	// Step 0: Retrieve and validate client
	client, err := t.client.FindOauthClient(clientId)
	if err != nil {
		return nil, err
	}

	if err := client.ValidateSecret(clientSecret); err != nil {
		return nil, err
	}

	// Step 1: Retrieve and validate the refresh token
	refreshToken, err := t.refreshTokenRepository.FindByToken(token)
	if err != nil {
		log.Printf("Failed to find refresh token: %v", err)
		return nil, fmt.Errorf("failed to find refresh token: %w", err)
	}

	if refreshToken == nil || refreshToken.IsExpired() {
		log.Println("Invalid or expired refresh token")
		return nil, fmt.Errorf("invalid or expired refresh token")
	}

	log.Println("Refresh token validated successfully for client:", clientId)

	// Step 2: Generate a new access token
	accessTokenJwt, err := utils.GenerateJWT(refreshToken.ClientId, "user", []byte("secret"), "access")
	if err != nil {
		log.Printf("Failed to generate new access token: %v", err)
		return nil, fmt.Errorf("failed to generate new access token: %w", err)
	}

	// Create and save the new access token
	accessToken := entities.NewAccessTokenBuilder().
		WithClientId(client.ClientId).
		WithToken(accessTokenJwt).
		WithTokenType("JWT").
		WithExpiresAt(time.Now().Add(AccessTokenDuration)). // Example expiration
		Build()

	accessToken, err = t.accessTokenRepository.Save(accessToken)
	if err != nil {
		log.Printf("Failed to save new access token: %v", err)
		return nil, fmt.Errorf("failed to save new access token: %w", err)
	}

	log.Println("New access token created and saved successfully for client:", refreshToken.ClientId)

	// Step 3: Build and return the token response
	newToken := oauth.NewTokenBuilder().
		WithClientId(client.ClientId).
		WithUserId("user"). // Adjust according to your use case
		WithAccessToken(accessToken.Token).
		WithAccessTokenCreatedAt(time.Now()).
		WithAccessTokenExpiresAt(accessToken.ExpiresAt.Sub(time.Now())).
		WithRefreshToken(refreshToken.Token). // Retain the same refresh token
		WithRefreshTokenCreatedAt(refreshToken.CreatedAt).
		WithRefreshTokenExpiresAt(refreshToken.ExpiresAt.Sub(time.Now())).
		WithExtension(nil). // If you have any extensions, set them here
		Build()

	log.Println("Token response successfully built for client:", refreshToken.ClientId)

	return newToken, nil
}
