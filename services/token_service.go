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

const AccessTokenDuration = 1*time.Hour + 1*time.Second
const RefreshTokenDuration = 30*24*time.Hour + 1*time.Second

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
		return t.getTokenByRefreshTokenFlow(command.ClientId, command.ClientSecret, command.RefreshToken)
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

// getTokenByRefreshTokenFlow handles the refresh token flow by validating the refresh token, authenticating the client
// (if it's confidential), generating a new access token, and issuing a new refresh token.
func (t *tokenService) getTokenByRefreshTokenFlow(clientId, clientSecret, token string) (*oauth.Token, error) {
	log.Println("Starting refresh token request processing")

	// Step 1: Retrieve and validate the refresh token
	refreshToken, err := t.refreshTokenRepository.FindByToken(token)
	if err != nil {
		log.Printf("Error finding refresh token with token '%s': %v", token, err)
		return nil, fmt.Errorf("failed to find refresh token: %w", err)
	}

	// Use Client ID from refresh token if not provided in the request
	if clientId == "" {
		log.Println("Client ID from the request is empty; using Client ID from the refresh token")
		clientId = refreshToken.ClientId
	}

	// Step 2: Retrieve and validate the client
	client, err := t.client.FindOauthClient(clientId)
	if err != nil {
		log.Printf("Error finding client with Client ID '%s': %v", clientId, err)
		return nil, fmt.Errorf("failed to find client: %w", err)
	}

	// Authenticate the client if it's confidential
	if client.IsConfidential {
		if err := authenticateClient(clientId, clientSecret, client); err != nil {
			return nil, fmt.Errorf("client authentication failed: %w", err)
		}
	}

	// Step 3: Validate the refresh token
	claims, err := utils.ValidateRefreshToken(token, []byte("secret"))
	if err != nil {
		log.Printf("Failed to validate refresh token '%s': %v", token, err)
		return nil, fmt.Errorf("failed to validate refresh token: %w", err)
	}
	log.Printf("Successfully validated refresh token with claims: %+v", claims)

	// Step 4: Generate a new access token
	accessTokenJwt, err := utils.GenerateJWT(refreshToken.ClientId, "user", []byte("secret"), "access")
	if err != nil {
		log.Printf("Failed to generate new access token for Client ID '%s': %v", refreshToken.ClientId, err)
		return nil, fmt.Errorf("failed to generate new access token: %w", err)
	}

	accessTokenJwe, err := utils.GenerateJWE(accessTokenJwt)
	if err != nil {
		log.Printf("Failed to generate JWE for access token: %v", err)
		return nil, fmt.Errorf("failed to generate JWE: %w", err)
	}

	// Create and save the new access token
	newAccessToken := entities.NewAccessTokenBuilder().
		WithClientId(refreshToken.ClientId).
		WithToken(accessTokenJwe).
		WithTokenType("JWE").
		WithExpiresAt(time.Now().Add(AccessTokenDuration)).
		Build()

	savedAccessToken, err := t.accessTokenRepository.Save(newAccessToken)
	if err != nil {
		log.Printf("Failed to save new access token for Client ID '%s': %v", refreshToken.ClientId, err)
		return nil, fmt.Errorf("failed to save new access token: %w", err)
	}
	log.Printf("New access token successfully created and saved for Client ID '%s'", refreshToken.ClientId)

	// Step 5: Generate a new refresh token
	refreshTokenJwt, err := utils.GenerateJWT(savedAccessToken.Id, "user", []byte("secret"), "refresh")
	if err != nil {
		log.Printf("Failed to generate new refresh token for Access Token ID '%s': %v", savedAccessToken.Id, err)
		return nil, fmt.Errorf("failed to generate new refresh token: %w", err)
	}

	// Create and save the new refresh token
	newRefreshToken := entities.NewRefreshTokenBuilder().
		WithAccessToken(savedAccessToken).
		WithAccessTokenId(savedAccessToken.Id).
		WithClient(savedAccessToken.Client).
		WithClientId(savedAccessToken.ClientId).
		WithToken(refreshTokenJwt).
		WithTokenType("JWT").
		WithExpiresAt(time.Now().Add(RefreshTokenDuration)).
		Build()

	savedRefreshToken, err := t.refreshTokenRepository.Save(newRefreshToken)
	if err != nil {
		log.Printf("Failed to save new refresh token for Access Token ID '%s': %v", savedAccessToken.Id, err)
		return nil, fmt.Errorf("failed to save new refresh token: %w", err)
	}

	// Step 6: Build and return the token response
	newToken := oauth.NewTokenBuilder().
		WithClientId(savedAccessToken.ClientId).
		WithUserId("user"). // Adjust according to your use case
		WithAccessToken(savedAccessToken.Token).
		WithAccessTokenCreatedAt(time.Now()).
		WithAccessTokenExpiresAt(savedAccessToken.ExpiresAt.Sub(time.Now())).
		WithRefreshToken(savedRefreshToken.Token).
		WithRefreshTokenCreatedAt(savedRefreshToken.CreatedAt).
		WithRefreshTokenExpiresAt(savedRefreshToken.ExpiresAt.Sub(time.Now())).
		WithExtension(nil). // If you have any extensions, set them here
		Build()

	log.Printf("Token response successfully built for Client ID '%s'", savedAccessToken.ClientId)

	return newToken, nil
}

// authenticateClient checks if the client is confidential and validates the provided client secret.
func authenticateClient(clientId, clientSecret string, client *entities.OauthClient) error {
	if clientSecret == "" {
		log.Printf("Client '%s' is confidential, but no client secret was provided", clientId)
		return fmt.Errorf("client secret is required for confidential clients")
	}
	if err := client.ValidateSecret(clientSecret); err != nil {
		log.Printf("Client '%s' is confidential and authentication failed: %v", clientId, err)
		return fmt.Errorf("authentication failed for confidential client: %w", err)
	}
	return nil
}
