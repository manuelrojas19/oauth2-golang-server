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
		return t.handleClientCredentialsFlow(command.ClientId, command.ClientSecret)
	case granttype.RefreshToken:
		return t.handleRefreshTokenFlow(command.ClientId, command.ClientSecret, command.RefreshToken)
	default:
		return nil, fmt.Errorf("unsupported grant type: %s", command.GrantType)
	}
}

// handleClientCredentialsFlow processes the client credentials grant type by validating the client credentials,
// generating an access token, and issuing a refresh token.
func (t *tokenService) handleClientCredentialsFlow(clientId, clientSecret string) (*oauth.Token, error) {
	// Step 1: Retrieve and validate the client
	client, err := t.client.FindOauthClient(clientId)
	if err != nil {
		log.Printf("Error retrieving client with ID '%s': %v", clientId, err)
		return nil, fmt.Errorf("failed to find client: %w", err)
	}

	if err := client.ValidateSecret(clientSecret); err != nil {
		log.Printf("Client authentication failed for ID '%s': %v", clientId, err)
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	// Step 2: Generate a new access token
	accessTokenJwt, err := utils.GenerateJWT(clientId, "user", []byte("secret"), "access")
	if err != nil {
		log.Printf("Error generating JWT for access token: %v", err)
		return nil, fmt.Errorf("failed to generate access token JWT: %w", err)
	}

	// Create and save the new access token
	accessToken := entities.NewAccessTokenBuilder().
		WithClient(client).
		WithClientId(clientId).
		WithToken(accessTokenJwt).
		WithTokenType("JWT").
		WithExpiresAt(time.Now().Add(AccessTokenDuration)).
		Build()

	savedAccessToken, err := t.accessTokenRepository.Save(accessToken)
	if err != nil {
		log.Printf("Error saving new access token for client ID '%s': %v", clientId, err)
		return nil, fmt.Errorf("failed to save access token: %w", err)
	}

	// Step 3: Generate a new refresh token
	refreshTokenJwt, err := utils.GenerateJWT(clientId, "user", []byte("secret"), "refresh")
	if err != nil {
		log.Printf("Error generating JWT for refresh token: %v", err)
		return nil, fmt.Errorf("failed to generate refresh token JWT: %w", err)
	}

	refreshToken := entities.NewRefreshTokenBuilder().
		WithAccessToken(savedAccessToken).
		WithAccessTokenId(savedAccessToken.Id).
		WithClient(client).
		WithClientId(clientId).
		WithToken(refreshTokenJwt).
		WithTokenType("JWT").
		WithExpiresAt(time.Now().Add(RefreshTokenDuration)).
		Build()

	savedRefreshToken, err := t.refreshTokenRepository.Save(refreshToken)
	if err != nil {
		log.Printf("Error saving new refresh token for access token ID '%s': %v", savedAccessToken.Id, err)
		return nil, fmt.Errorf("failed to save refresh token: %w", err)
	}

	// Step 4: Build and return the token response
	token := oauth.NewTokenBuilder().
		WithClientId(savedAccessToken.ClientId).
		WithUserId("user").
		WithAccessToken(savedAccessToken.Token).
		WithAccessTokenCreatedAt(savedAccessToken.CreatedAt).
		WithAccessTokenExpiresAt(savedAccessToken.ExpiresAt.Sub(time.Now())).
		WithRefreshToken(savedRefreshToken.Token).
		WithRefreshTokenCreatedAt(savedRefreshToken.CreatedAt).
		WithRefreshTokenExpiresAt(savedRefreshToken.ExpiresAt.Sub(time.Now())).
		WithExtension(nil).
		Build()

	return token, nil
}

// handleRefreshTokenFlow processes the refresh token grant type by validating the refresh token,
// authenticating the client (if confidential), generating a new access token, and issuing a new refresh token.
func (t *tokenService) handleRefreshTokenFlow(clientId, clientSecret, token string) (*oauth.Token, error) {
	log.Println("Processing refresh token request")

	// Step 1: Retrieve and validate the refresh token
	refreshToken, err := t.refreshTokenRepository.FindByToken(token)
	if err != nil {
		log.Printf("Error finding refresh token with token '%s': %v", token, err)
		return nil, fmt.Errorf("failed to find refresh token: %w", err)
	}

	if clientId == "" {
		log.Println("Client ID not provided; using Client ID from refresh token")
		clientId = refreshToken.ClientId
	}

	// Step 2: Retrieve and validate the client
	client, err := t.client.FindOauthClient(clientId)
	if err != nil {
		log.Printf("Error retrieving client with ID '%s': %v", clientId, err)
		return nil, fmt.Errorf("failed to find client: %w", err)
	}

	if client.IsConfidential {
		if err := authenticateClient(clientId, clientSecret, client); err != nil {
			log.Printf("Client authentication failed for ID '%s': %v", clientId, err)
			return nil, fmt.Errorf("client authentication failed: %w", err)
		}
	}

	// Step 3: Validate the refresh token
	claims, err := utils.ValidateRefreshToken(token, []byte("secret"))
	if err != nil {
		log.Printf("Error validating refresh token '%s': %v", token, err)
		return nil, fmt.Errorf("failed to validate refresh token: %w", err)
	}
	log.Printf("Successfully validated refresh token with claims: %+v", claims)

	// Step 4: Generate a new access token
	accessTokenJwt, err := utils.GenerateJWT(refreshToken.ClientId, "user", []byte("secret"), "access")
	if err != nil {
		log.Printf("Error generating JWT for new access token: %v", err)
		return nil, fmt.Errorf("failed to generate new access token JWT: %w", err)
	}

	newAccessToken := entities.NewAccessTokenBuilder().
		WithClientId(refreshToken.ClientId).
		WithToken(accessTokenJwt).
		WithTokenType("JWT").
		WithExpiresAt(time.Now().Add(AccessTokenDuration)).
		Build()

	savedAccessToken, err := t.accessTokenRepository.Save(newAccessToken)
	if err != nil {
		log.Printf("Error saving new access token for Client ID '%s': %v", refreshToken.ClientId, err)
		return nil, fmt.Errorf("failed to save new access token: %w", err)
	}
	log.Printf("New access token successfully created and saved for Client ID '%s'", refreshToken.ClientId)

	// Step 5: Invalidate used refresh token
	err = t.refreshTokenRepository.InvalidateRefreshTokensByAccessTokenId(refreshToken.AccessTokenId)
	if err != nil {
		log.Printf("Error invalidating used refresh: %v", err)
		return nil, err
	}

	// Step 6: Generate a new refresh token
	refreshTokenJwt, err := utils.GenerateJWT(savedAccessToken.Id, "user", []byte("secret"), "refresh")
	if err != nil {
		log.Printf("Error generating JWT for new refresh token: %v", err)
		return nil, fmt.Errorf("failed to generate new refresh token JWT: %w", err)
	}

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
		log.Printf("Error saving new refresh token for Access Token ID '%s': %v", savedAccessToken.Id, err)
		return nil, fmt.Errorf("failed to save new refresh token: %w", err)
	}

	// Step 6: Build and return the token response
	newToken := oauth.NewTokenBuilder().
		WithClientId(savedAccessToken.ClientId).
		WithUserId("user").
		WithAccessToken(savedAccessToken.Token).
		WithAccessTokenCreatedAt(savedAccessToken.CreatedAt).
		WithAccessTokenExpiresAt(savedAccessToken.ExpiresAt.Sub(time.Now())).
		WithRefreshToken(savedRefreshToken.Token).
		WithRefreshTokenCreatedAt(savedRefreshToken.CreatedAt).
		WithRefreshTokenExpiresAt(savedRefreshToken.ExpiresAt.Sub(time.Now())).
		WithExtension(nil).
		Build()

	log.Printf("Token response successfully built for Client ID '%s'", savedAccessToken.ClientId)

	return newToken, nil
}

// authenticateClient checks if the client is confidential and validates the provided client secret.
func authenticateClient(clientId, clientSecret string, client *entities.OauthClient) error {
	if clientSecret == "" {
		log.Printf("Client '%s' is confidential but no client secret provided", clientId)
		return fmt.Errorf("client secret is required for confidential clients")
	}
	if err := client.ValidateSecret(clientSecret); err != nil {
		log.Printf("Authentication failed for confidential client '%s': %v", clientId, err)
		return fmt.Errorf("authentication failed: %w", err)
	}
	return nil
}
