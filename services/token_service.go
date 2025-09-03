package services

import (
	"fmt"
	"time"

	"github.com/manuelrojas19/go-oauth2-server/oauth"
	"github.com/manuelrojas19/go-oauth2-server/oauth/granttype"
	"github.com/manuelrojas19/go-oauth2-server/store"
	"github.com/manuelrojas19/go-oauth2-server/store/repositories"
	"github.com/manuelrojas19/go-oauth2-server/utils"
	"go.uber.org/zap"
)

const AccessTokenDuration = 1*time.Hour
const RefreshTokenDuration = 30*24*time.Hour

type GrantAccessTokenCommand struct {
	ClientId     string
	ClientSecret string
	RefreshToken string
	GrantType    granttype.GrantType
	Code         string
	RedirectUri  string
	CodeVerifier string
}

func NewGrantAccessTokenCommand(clientId string, clientSecret string, grantType granttype.GrantType, refreshToken string, code string, redirectUri string, codeVerifier string) *GrantAccessTokenCommand {
	return &GrantAccessTokenCommand{
		ClientId:     clientId,
		ClientSecret: clientSecret,
		GrantType:    grantType,
		RefreshToken: refreshToken,
		Code:         code,
		RedirectUri:  redirectUri,
		CodeVerifier: codeVerifier,
	}
}

type tokenService struct {
	accessTokenRepository  repositories.AccessTokenRepository
	refreshTokenRepository repositories.RefreshTokenRepository
	authRepository         repositories.AuthorizationRepository
	client                 OauthClientService
	logger                 *zap.Logger
}

func NewTokenService(
	accessTokenRepository repositories.AccessTokenRepository,
	refreshTokenRepository repositories.RefreshTokenRepository,
	authRepository repositories.AuthorizationRepository,
	client OauthClientService,
	logger *zap.Logger) TokenService {
	return &tokenService{
		accessTokenRepository:  accessTokenRepository,
		refreshTokenRepository: refreshTokenRepository,
		authRepository:         authRepository,
		client:                 client,
		logger:                 logger,
	}
}

func (t *tokenService) GrantAccessToken(command *GrantAccessTokenCommand) (*oauth.Token, error) {
	t.logger.Info("Granting access token", zap.String("grantType", string(command.GrantType)), zap.String("clientId", command.ClientId))
	switch command.GrantType {
	case granttype.ClientCredentials:
		return t.handleClientCredentialsFlow(command.ClientId, command.ClientSecret)
	case granttype.RefreshToken:
		return t.handleRefreshTokenFlow(command.ClientId, command.ClientSecret, command.RefreshToken)
	case granttype.AuthorizationCode:
		return t.handleAuthorizationCodeFlow(command.ClientId, command.ClientSecret, command.Code, command.RedirectUri, command.CodeVerifier)
	default:
		t.logger.Warn("Unsupported grant type", zap.String("grantType", string(command.GrantType)))
		return nil, fmt.Errorf("unsupported grant type: %s", command.GrantType)
	}
}

// handleClientCredentialsFlow processes the client credentials grant type by validating the client credentials,
// generating an access token, and issuing a refresh token.
func (t *tokenService) handleClientCredentialsFlow(clientId, clientSecret string) (*oauth.Token, error) {

	t.logger.Info("Handling Client Credentials Flow", zap.String("clientId", clientId))

	// Step 1: Retrieve and validate the client
	client, err := t.client.FindOauthClient(clientId)
	if err != nil {
		t.logger.Error("Error retrieving client for Client Credentials Flow", zap.String("clientId", clientId), zap.Error(err))
		return nil, fmt.Errorf("failed to find client: %w", err)
	}

	if err := client.ValidateSecret(clientSecret); err != nil {
		t.logger.Error("Client authentication failed for Client Credentials Flow", zap.String("clientId", clientId), zap.Error(err))
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	// Preload scopes for the client
	err = t.client.PreloadOauthClientScopes(client)
	if err != nil {
		t.logger.Error("Error preloading client scopes", zap.String("clientId", clientId), zap.Error(err))
		return nil, fmt.Errorf("failed to preload client scopes: %w", err)
	}

	t.logger.Debug("Client authenticated successfully for Client Credentials Flow", zap.String("clientId", clientId))

	// Step 2: Generate a new access token
	accessTokenJwt, err := utils.GenerateJWT(&clientId, nil, []byte("secret"), "access")
	if err != nil {
		t.logger.Error("Error generating JWT for access token in Client Credentials Flow", zap.String("clientId", clientId), zap.Error(err))
		return nil, fmt.Errorf("failed to generate access token JWT: %w", err)
	}
	t.logger.Debug("Access token JWT generated for Client Credentials Flow")

	// Create and save the new access token
	accessToken := store.NewAccessTokenBuilder().
		WithClient(client).
		WithClientId(&clientId).
		WithToken(accessTokenJwt).
		WithTokenType("Bearer").
		WithExpiresAt(time.Now().Add(AccessTokenDuration)).
		WithScopes(client.Scopes).
		Build()

	savedAccessToken, err := t.accessTokenRepository.Save(accessToken)
	if err != nil {
		t.logger.Error("Error saving new access token for Client Credentials Flow", zap.String("clientId", utils.StringDeref(accessToken.ClientId)), zap.Error(err))
		return nil, fmt.Errorf("failed to save access token: %w", err)
	}
	t.logger.Info("New access token saved successfully for Client Credentials Flow", zap.String("accessTokenId", savedAccessToken.Id))

	// Step 3: Build and return the token response
	token := oauth.NewTokenBuilder().
		WithClientId(savedAccessToken.ClientId).
		WithAccessToken(savedAccessToken.Token).
		WithTokenType(savedAccessToken.TokenType).
		WithAccessTokenCreatedAt(savedAccessToken.CreatedAt).
		WithAccessTokenExpiresIn(int(AccessTokenDuration.Seconds())).
		WithAccessTokenExpiresAt(savedAccessToken.ExpiresAt).
		WithExtension(nil).
		WithScope(utils.ScopesToStringSlice(savedAccessToken.Scopes)).
		Build()
	t.logger.Info("Token response built for Client Credentials Flow", zap.String("clientId", utils.StringDeref(savedAccessToken.ClientId)))

	return token, nil
}

// handleRefreshTokenFlow processes the refresh token grant type by validating the refresh token,
// authenticating the client (if confidential), generating a new access token, and issuing a new refresh token.
func (t *tokenService) handleRefreshTokenFlow(clientId, clientSecret, token string) (*oauth.Token, error) {
	t.logger.Info("Processing refresh token request", zap.String("clientId", clientId), zap.String("refreshToken", token))

	// Step 1: Retrieve and validate the refresh token
	refreshToken, err := t.refreshTokenRepository.FindByRefreshToken(token)
	if err != nil {
		t.logger.Error("Error finding refresh token", zap.String("refreshToken", token), zap.Error(err))
		return nil, fmt.Errorf("failed to find refresh token: %w", err)
	}
	t.logger.Debug("Refresh token retrieved", zap.String("refreshTokenId", refreshToken.Id))

	if clientId == "" {
		if refreshToken.ClientId != nil {
			t.logger.Debug("Client ID not provided; using Client ID from refresh token", zap.String("refreshTokenClientId", *refreshToken.ClientId))
			clientId = *refreshToken.ClientId
		}
	} else {
		// If clientId is provided in the request, it must match the one in the refresh token (if present)
		if refreshToken.ClientId != nil && clientId != *refreshToken.ClientId {
			return nil, fmt.Errorf("client ID mismatch: provided client ID does not match refresh token's client ID")
		}
	}

	// Step 2: Retrieve and validate the client
	client, err := t.client.FindOauthClient(clientId)
	if err != nil {
		t.logger.Error("Error retrieving client for Refresh Token Flow", zap.String("clientId", clientId), zap.Error(err))
		return nil, fmt.Errorf("failed to find client: %w", err)
	}

	// Preload scopes for the client
	err = t.client.PreloadOauthClientScopes(client)
	if err != nil {
		t.logger.Error("Error preloading client scopes for Refresh Token Flow", zap.String("clientId", clientId), zap.Error(err))
		return nil, fmt.Errorf("failed to preload client scopes: %w", err)
	}

	t.logger.Debug("Client retrieved for Refresh Token Flow", zap.String("clientId", client.ClientId))

	if client.Confidential {
		if err := t.authenticateClient(clientId, clientSecret, client); err != nil {
			t.logger.Error("Client authentication failed for Refresh Token Flow", zap.String("clientId", clientId), zap.Error(err))
			return nil, fmt.Errorf("client authentication failed: %w", err)
		}
		t.logger.Debug("Confidential client authenticated for Refresh Token Flow", zap.String("clientId", clientId))
	}

	// Step 3: Validate the refresh token
	claims, err := utils.ValidateRefreshToken(token, []byte("secret"))
	if err != nil {
		t.logger.Error("Error validating refresh token", zap.String("refreshToken", token), zap.Error(err))
		return nil, fmt.Errorf("failed to validate refresh token: %w", err)
	}
	t.logger.Debug("Successfully validated refresh token", zap.Any("claims", claims))

	// Step 4: Generate a new access token
	accessTokenJwt, err := utils.GenerateJWT(refreshToken.ClientId, refreshToken.UserId, []byte("secret"), "access")
	if err != nil {
		t.logger.Error("Error generating JWT for new access token in Refresh Token Flow", zap.String("clientId", utils.StringDeref(refreshToken.ClientId)), zap.Error(err))
		return nil, fmt.Errorf("failed to generate new access token JWT: %w", err)
	}
	t.logger.Debug("New access token JWT generated for Refresh Token Flow")

	newAccessToken := store.NewAccessTokenBuilder().
		WithClientId(refreshToken.ClientId).
		WithToken(accessTokenJwt).
		WithTokenType("Bearer").
		WithExpiresAt(time.Now().Add(AccessTokenDuration)).
		WithUserId(refreshToken.UserId).
		WithScopes(client.Scopes).
		Build()

	savedAccessToken, err := t.accessTokenRepository.Save(newAccessToken)
	if err != nil {
		t.logger.Error("Error saving new access token for Refresh Token Flow", zap.String("clientId", utils.StringDeref(refreshToken.ClientId)), zap.Error(err))
		return nil, fmt.Errorf("failed to save new access token: %w", err)
	}
	t.logger.Info("New access token created and saved successfully for Refresh Token Flow", zap.String("accessTokenId", savedAccessToken.Id))

	// Step 5: Invalidate used refresh token
	err = t.refreshTokenRepository.InvalidateRefreshTokensByAccessTokenId(refreshToken.AccessTokenId)
	if err != nil {
		t.logger.Error("Error invalidating used refresh token", zap.String("accessTokenId", refreshToken.AccessTokenId), zap.Error(err))
		return nil, err
	}
	t.logger.Debug("Old refresh token invalidated successfully")

	// Step 6: Generate a new refresh token
	refreshTokenJwt, err := utils.GenerateJWT(savedAccessToken.ClientId, savedAccessToken.UserId, []byte("secret"), "refresh")
	if err != nil {
		t.logger.Error("Error generating JWT for new refresh token in Refresh Token Flow", zap.String("accessTokenId", savedAccessToken.Id), zap.Error(err))
		return nil, fmt.Errorf("failed to generate new refresh token JWT: %w", err)
	}
	t.logger.Debug("New refresh token JWT generated")

	newRefreshToken := store.NewRefreshTokenBuilder().
		WithAccessToken(savedAccessToken).
		WithAccessTokenId(savedAccessToken.Id).
		WithClient(savedAccessToken.Client).
		WithClientId(savedAccessToken.ClientId).
		WithToken(refreshTokenJwt).
		WithTokenType("Bearer").
		WithExpiresAt(time.Now().Add(RefreshTokenDuration)).
		WithUserId(savedAccessToken.UserId).
		Build()

	savedRefreshToken, err := t.refreshTokenRepository.Save(newRefreshToken)
	if err != nil {
		t.logger.Error("Error saving new refresh token for Refresh Token Flow", zap.String("accessTokenId", savedAccessToken.Id), zap.Error(err))
		return nil, fmt.Errorf("failed to save new refresh token: %w", err)
	}
	t.logger.Info("New refresh token saved successfully for Refresh Token Flow", zap.String("refreshTokenId", savedRefreshToken.Id))

	// Step 6: Build and return the token response
	newToken := oauth.NewTokenBuilder().
		WithClientId(savedAccessToken.ClientId).
		WithUserId(savedAccessToken.UserId).
		WithAccessToken(savedAccessToken.Token).
		WithTokenType(savedAccessToken.TokenType).
		WithAccessTokenCreatedAt(savedAccessToken.CreatedAt).
		WithAccessTokenExpiresIn(int(AccessTokenDuration.Seconds())).
		WithAccessTokenExpiresAt(savedAccessToken.ExpiresAt).
		WithRefreshToken(savedRefreshToken.Token).
		WithRefreshTokenCreatedAt(savedRefreshToken.CreatedAt).
		WithRefreshTokenExpiresAt(savedRefreshToken.ExpiresAt).
		WithExtension(nil).
		WithScope(utils.ScopesToStringSlice(savedAccessToken.Scopes)).
		Build()

	t.logger.Info("Token response successfully built for Refresh Token Flow", zap.String("clientId", utils.StringDeref(savedAccessToken.ClientId)))

	return newToken, nil
}

// authenticateClient checks if the client is confidential and validates the provided client secret.
func (t *tokenService) authenticateClient(clientId, clientSecret string, client *store.OauthClient) error {
	if clientSecret == "" {
		t.logger.Warn("Client is confidential but no client secret provided", zap.String("clientId", utils.StringDeref(&clientId)))
		return fmt.Errorf("client secret is required for confidential clients")
	}
	if err := client.ValidateSecret(clientSecret); err != nil {
		t.logger.Error("Authentication failed for confidential client", zap.String("clientId", utils.StringDeref(&clientId)), zap.Error(err))
		return fmt.Errorf("authentication failed: %w", err)
	}
	t.logger.Debug("Client authenticated successfully", zap.String("clientId", utils.StringDeref(&clientId)))
	return nil
}

// handleAuthorizationCodeFlow processes the authorization code grant type by validating the authorization code,
// generating an access token, and issuing a refresh token.
func (t *tokenService) handleAuthorizationCodeFlow(clientId, clientSecret, code, redirectUri, codeVerifier string) (*oauth.Token, error) {
	t.logger.Info("Handling Authorization Code Flow", zap.String("clientId", clientId), zap.String("code", code))
	// Step 1: Retrieve and validate the authorization code
	authCode, err := t.authRepository.FindByCode(code)
	if err != nil {
		t.logger.Error("Error finding authorization code", zap.String("code", code), zap.Error(err))
		return nil, fmt.Errorf("failed to find authorization code: %w", err)
	}
	t.logger.Debug("Authorization code retrieved", zap.Any("authCode", authCode))

	if authCode.ClientId != nil && clientId != *authCode.ClientId {
		t.logger.Warn("Client ID mismatch", zap.String("expectedClientId", utils.StringDeref(authCode.ClientId)), zap.String("receivedClientId", clientId), zap.String("code", code))
		return nil, fmt.Errorf("client ID mismatch")
	}

	if authCode.RedirectURI != redirectUri {
		t.logger.Warn("Redirect URI mismatch", zap.String("expectedRedirectUri", authCode.RedirectURI), zap.String("receivedRedirectUri", redirectUri), zap.String("code", code))
		return nil, fmt.Errorf("redirect URI mismatch")
	}

	if time.Now().After(authCode.ExpiresAt) {
		t.logger.Warn("Authorization code has expired", zap.String("code", code), zap.Time("expiresAt", authCode.ExpiresAt))
		return nil, fmt.Errorf("authorization code has expired")
	}

	// PKCE validation
	if authCode.CodeChallenge != "" && authCode.CodeChallengeMethod == "S256" {
		if codeVerifier == "" {
			t.logger.Warn("Code verifier is missing for PKCE enabled authorization code", zap.String("code", code))
			return nil, fmt.Errorf("code verifier required for PKCE")
		}
		// Calculate the S256 code_challenge from the code_verifier
		calculatedCodeChallenge := utils.S256Challenge(codeVerifier)
		if calculatedCodeChallenge != authCode.CodeChallenge {
			t.logger.Warn("Code challenge mismatch", zap.String("code", code), zap.String("expectedCodeChallenge", authCode.CodeChallenge), zap.String("receivedCodeChallenge", calculatedCodeChallenge))
			return nil, fmt.Errorf("code challenge mismatch")
		}
		t.logger.Debug("PKCE code challenge validated successfully", zap.String("code", code))
	} else if authCode.CodeChallenge != "" && authCode.CodeChallengeMethod == "" {
		t.logger.Warn("Code challenge method is missing for PKCE enabled authorization code", zap.String("code", code))
		return nil, fmt.Errorf("code challenge method required for PKCE")
	}

	// Step 2: Retrieve and validate the client
	client, err := t.client.FindOauthClient(clientId)
	if err != nil {
		t.logger.Error("Error retrieving client for Authorization Code Flow", zap.String("clientId", clientId), zap.Error(err))
		return nil, fmt.Errorf("failed to find client: %w", err)
	}
	t.logger.Debug("Client retrieved for Authorization Code Flow", zap.String("clientId", clientId))

	if err := t.authenticateClient(clientId, clientSecret, client); err != nil {
		t.logger.Error("Client authentication failed for Authorization Code Flow", zap.String("clientId", clientId), zap.Error(err))
		return nil, fmt.Errorf("client authentication failed: %w", err)
	}
	t.logger.Debug("Confidential client authenticated for Authorization Code Flow", zap.String("clientId", clientId))

	// Step 2.5: Invalidate the authorization code to prevent replay attacks
	t.logger.Debug("Invalidating authorization code to prevent replay attacks", zap.String("code", authCode.Code))
	err = t.authRepository.Delete(authCode.Code)
	if err != nil {
		t.logger.Error("Error deleting authorization code", zap.String("code", authCode.Code), zap.Error(err))
		return nil, fmt.Errorf("failed to invalidate authorization code: %w", err)
	}
	t.logger.Info("Authorization code invalidated successfully", zap.String("code", authCode.Code))

	// Step 3: Generate a new access token
	accessTokenJwt, err := utils.GenerateJWT(authCode.ClientId, authCode.UserId, []byte("secret"), "access")
	if err != nil {
		t.logger.Error("Error generating JWT for access token in Authorization Code Flow", zap.String("clientId", utils.StringDeref(authCode.ClientId)), zap.Error(err))
		return nil, fmt.Errorf("failed to generate access token JWT: %w", err)
	}
	t.logger.Debug("Access token JWT generated for Authorization Code Flow")

	newAccessToken := store.NewAccessTokenBuilder().
		WithClientId(authCode.ClientId).
		WithToken(accessTokenJwt).
		WithCode(code).
		WithTokenType("Bearer").
		WithExpiresAt(time.Now().Add(AccessTokenDuration)).
		WithUserId(authCode.UserId).
		Build()

	savedAccessToken, err := t.accessTokenRepository.Save(newAccessToken)
	if err != nil {
		t.logger.Error("Error saving new access token for Authorization Code Flow", zap.String("clientId", utils.StringDeref(authCode.ClientId)), zap.Error(err))
		return nil, fmt.Errorf("failed to save new access token: %w", err)
	}
	t.logger.Info("New access token created and saved successfully for Authorization Code Flow", zap.String("accessTokenId", savedAccessToken.Id))

	// Step 4: Generate a new refresh token
	refreshTokenJwt, err := utils.GenerateJWT(savedAccessToken.ClientId, savedAccessToken.UserId, []byte("secret"), "refresh")
	if err != nil {
		t.logger.Error("Error generating JWT for refresh token in Authorization Code Flow", zap.String("accessTokenId", savedAccessToken.Id), zap.Error(err))
		return nil, fmt.Errorf("failed to generate refresh token JWT: %w", err)
	}
	t.logger.Debug("New refresh token JWT generated for Authorization Code Flow")

	newRefreshToken := store.NewRefreshTokenBuilder().
		WithAccessToken(savedAccessToken).
		WithAccessTokenId(savedAccessToken.Id).
		WithClient(savedAccessToken.Client).
		WithClientId(savedAccessToken.ClientId).
		WithToken(refreshTokenJwt).
		WithTokenType("Bearer").
		WithExpiresAt(time.Now().Add(RefreshTokenDuration)).
		WithUserId(savedAccessToken.UserId).
		Build()

	savedRefreshToken, err := t.refreshTokenRepository.Save(newRefreshToken)
	if err != nil {
		t.logger.Error("Error saving new refresh token for Authorization Code Flow", zap.String("accessTokenId", savedAccessToken.Id), zap.Error(err))
		return nil, fmt.Errorf("failed to save new refresh token: %w", err)
	}
	t.logger.Info("New refresh token saved successfully for Authorization Code Flow", zap.String("refreshTokenId", savedRefreshToken.Id))

	// Step 5: Build and return the token response
	token := oauth.NewTokenBuilder().
		WithClientId(savedAccessToken.ClientId).
		WithUserId(savedAccessToken.UserId).
		WithAccessToken(savedAccessToken.Token).
		WithTokenType(savedAccessToken.TokenType).
		WithAccessTokenCreatedAt(savedAccessToken.CreatedAt).
		WithAccessTokenExpiresIn(int(AccessTokenDuration.Seconds())).
		WithAccessTokenExpiresAt(savedAccessToken.ExpiresAt).
		WithRefreshToken(savedRefreshToken.Token).
		WithRefreshTokenCreatedAt(savedRefreshToken.CreatedAt).
		WithRefreshTokenExpiresAt(savedRefreshToken.ExpiresAt).
		WithExtension(nil).
		Build()

	t.logger.Info("Token response successfully built for Authorization Code Flow", zap.String("clientId", utils.StringDeref(savedAccessToken.ClientId)))

	return token, nil
}
