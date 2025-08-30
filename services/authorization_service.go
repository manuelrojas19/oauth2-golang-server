package services

import (
	"fmt"
	"slices"
	"time"

	"github.com/manuelrojas19/go-oauth2-server/configuration"
	"github.com/manuelrojas19/go-oauth2-server/errors"
	"github.com/manuelrojas19/go-oauth2-server/oauth"
	"github.com/manuelrojas19/go-oauth2-server/oauth/responsetype"
	"github.com/manuelrojas19/go-oauth2-server/store"
	"github.com/manuelrojas19/go-oauth2-server/store/repositories"
	"github.com/manuelrojas19/go-oauth2-server/utils"
	"go.uber.org/zap"
)

type AuthorizeCommand struct {
	ClientId            string
	Scope               string
	RedirectUri         string
	ResponseType        responsetype.ResponseType
	SessionId           string
	State               string
	CodeChallenge       string
	CodeChallengeMethod string
}

type authorizationService struct {
	oauthClientService OauthClientService
	consentService     UserConsentService
	authRepository     repositories.AuthorizationRepository
	sessionService     SessionService
	userRepository     repositories.UserRepository
	logger             *zap.Logger
}

// NewAuthorizationService initializes a new AuthorizationService
func NewAuthorizationService(oauthClientService OauthClientService,
	consentService UserConsentService,
	authRepository repositories.AuthorizationRepository,
	userSessionService SessionService,
	userRepository repositories.UserRepository,
	logger *zap.Logger,
) AuthorizationService {
	return &authorizationService{
		oauthClientService: oauthClientService,
		consentService:     consentService,
		authRepository:     authRepository,
		sessionService:     userSessionService,
		userRepository:     userRepository,
		logger:             logger,
	}
}

func (a *authorizationService) Authorize(command *AuthorizeCommand) (*oauth.AuthCode, error) {
	clientId := command.ClientId

	a.logger.Info("Authorize request will be processed",
		zap.String("clientId", clientId),
		zap.String("responseType", string(command.ResponseType)),
		zap.String("redirectUri", command.RedirectUri),
		zap.String("sessionId", command.SessionId),
		zap.String("scope", command.Scope),
	)
	a.logger.Debug("Authorization command details", zap.Any("command", command))

	start := time.Now()

	// Validate response type
	if command.ResponseType != responsetype.Code {
		a.logger.Error("Unsupported response type",
			zap.String("responseType", string(command.ResponseType)),
			zap.String("clientId", clientId),
			zap.Duration("duration", time.Since(start)),
			zap.Stack("stacktrace"),
		)
		return nil, fmt.Errorf(errors.ErrUnsupportedResponseType)
	}
	a.logger.Debug("Response type validated", zap.String("responseType", string(command.ResponseType)))

	// Retrieve the OAuth client
	client, err := a.oauthClientService.FindOauthClient(clientId)
	if err != nil {
		a.logger.Error("Error retrieving client",
			zap.String("clientId", clientId),
			zap.Error(err),
			zap.Duration("duration", time.Since(start)),
			zap.Stack("stacktrace"),
		)
		return nil, fmt.Errorf("failed to retrieve client: %w", err)
	}
	a.logger.Info("Successfully retrieved Oauth client",
		zap.String("clientId", clientId),
		zap.Duration("duration", time.Since(start)),
	)
	a.logger.Debug("Retrieved OAuth client details", zap.Any("client", client))

	// Validate redirect URI
	if !isRegisteredRedirectUri(command, client) {
		a.logger.Error("Invalid redirect URI",
			zap.String("redirectUri", command.RedirectUri),
			zap.String("clientId", clientId),
			zap.Duration("duration", time.Since(start)),
			zap.Stack("stacktrace"),
		)
		return nil, fmt.Errorf("%s: %s", errors.ErrInvalidRedirectUri, command.RedirectUri)
	}
	a.logger.Debug("Redirect URI validated", zap.String("redirectUri", command.RedirectUri))

	// Check if user is authenticated
	if !a.sessionService.SessionExists(command.SessionId) {
		a.logger.Warn("User not authenticated",
			zap.String("sessionId", command.SessionId),
			zap.Duration("duration", time.Since(start)),
			zap.Stack("stacktrace"),
		)
		return nil, fmt.Errorf(errors.ErrUserNotAuthenticated)
	}
	a.logger.Info("Session exists",
		zap.String("sessionId", command.SessionId),
		zap.Duration("duration", time.Since(start)),
	)
	a.logger.Debug("User session exists and is valid", zap.String("sessionId", command.SessionId))

	// Retrieve user Id from session
	userId, err := a.sessionService.GetUserIdFromSession(command.SessionId)
	if err != nil {
		a.logger.Error("Error retrieving user from session",
			zap.String("sessionId", command.SessionId),
			zap.Error(err),
			zap.Duration("duration", time.Since(start)),
			zap.Stack("stacktrace"),
		)
		return nil, fmt.Errorf("failed to retrieve user from session: %w", err)
	}
	a.logger.Debug("Retrieved userId from session", zap.String("userId", userId))

	// Validate user in the database
	user, err := a.userRepository.FindByUserId(userId)
	if err != nil {
		a.logger.Error("Error retrieving user",
			zap.String("userId", userId),
			zap.Error(err),
			zap.Duration("duration", time.Since(start)),
			zap.Stack("stacktrace"),
		)
		return nil, fmt.Errorf("failed to retrieve user from user Id: %w", err)
	}
	a.logger.Debug("User found in database", zap.String("userId", user.Id))

	// Validate access consent
	if !a.consentService.HasUserConsented(user.Id, client.ClientId, command.Scope) {
		a.logger.Warn("User consent required",
			zap.String("userId", user.Id),
			zap.String("clientId", client.ClientId),
			zap.String("scope", command.Scope),
			zap.Duration("duration", time.Since(start)),
			zap.Stack("stacktrace"),
		)
		return nil, fmt.Errorf(errors.ErrConsentRequired)
	}
	a.logger.Debug("User consent confirmed")

	// Generate authorization code
	code, err := utils.GenerateAuthCode(client.ClientId, userId)
	if err != nil {
		a.logger.Error("Error generating authorization code",
			zap.String("clientId", client.ClientId),
			zap.String("userId", userId),
			zap.Error(err),
			zap.Duration("duration", time.Since(start)),
			zap.Stack("stacktrace"),
		)
		return nil, fmt.Errorf(errors.ErrConsentRequired)
	}
	a.logger.Debug("Authorization code generated", zap.String("code", code))

	// Build authorization code entity
	authCodeEntity := store.NewAuthorizationCodeBuilder().
		WithCode(code).
		WithClientId(client.ClientId).
		WithClient(client).
		WithUserId(user.Id).
		WithScope(command.Scope).
		WithRedirectURI(command.RedirectUri).
		WithCodeChallenge(command.CodeChallenge).
		WithCodeChallengeMethod(command.CodeChallengeMethod).
		WithExpiresAt(time.Now().Add(configuration.AuthCodeExpireTime)).
		Build()
	a.logger.Debug("Authorization code entity built", zap.Any("authCodeEntity", authCodeEntity))

	// Save authorization code entity to repository
	authCodeEntity, err = a.authRepository.Save(authCodeEntity)
	if err != nil {
		a.logger.Error("Error saving authorization code entity",
			zap.Error(err),
			zap.Duration("duration", time.Since(start)),
			zap.Stack("stacktrace"),
			zap.String("authCode", authCodeEntity.Code),
		)
		return nil, fmt.Errorf("failed to save authorization code entity: %w", err)
	}
	a.logger.Info("Authorization code entity saved successfully", zap.String("authCode", authCodeEntity.Code))

	// Build the OAuth authorization code response
	oauthCode := oauth.NewAuthCodeBuilder().
		WithCode(authCodeEntity.Code).
		WithClientId(authCodeEntity.ClientId).
		WithScope(authCodeEntity.Scope).
		WithRedirectURI(authCodeEntity.RedirectURI).
		WithCreatedAt(authCodeEntity.CreatedAt).
		WithExpiresAt(authCodeEntity.ExpiresAt).
		Build()

	a.logger.Info("Successfully generated authorization code",
		zap.String("clientId", client.ClientId),
		zap.String("userId", userId),
		zap.Duration("duration", time.Since(start)),
	)

	return oauthCode, nil
}

func isRegisteredRedirectUri(command *AuthorizeCommand, client *store.OauthClient) bool {
	return slices.Contains(client.RedirectURIs, command.RedirectUri)
}
