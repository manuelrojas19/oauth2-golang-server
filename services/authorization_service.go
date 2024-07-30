package services

import (
	"fmt"
	"github.com/manuelrojas19/go-oauth2-server/configuration"
	"github.com/manuelrojas19/go-oauth2-server/errors"
	"github.com/manuelrojas19/go-oauth2-server/oauth"
	"github.com/manuelrojas19/go-oauth2-server/oauth/responsetype"
	"github.com/manuelrojas19/go-oauth2-server/store"
	"github.com/manuelrojas19/go-oauth2-server/store/repositories"
	"github.com/manuelrojas19/go-oauth2-server/utils"
	"log"
	"time"
)

type AuthorizeCommand struct {
	ClientId     string
	Scope        string
	RedirectUri  string
	ResponseType responsetype.ResponseType
	SessionId    string
	State        string
}

type authorizationService struct {
	oauthClientService OauthClientService
	consentService     UserConsentService
	authRepository     repositories.AuthorizationRepository
	sessionService     SessionService
	userRepository     repositories.UserRepository
}

// NewAuthorizationService initializes a new AuthorizationService
func NewAuthorizationService(oauthClientService OauthClientService,
	consentService UserConsentService,
	authRepository repositories.AuthorizationRepository,
	userSessionService SessionService,
	userRepository repositories.UserRepository,
) AuthorizationService {
	return &authorizationService{oauthClientService: oauthClientService,
		consentService: consentService,
		authRepository: authRepository,
		sessionService: userSessionService,
		userRepository: userRepository,
	}
}

// Authorize authorizes and generate an Auth Code
func (a authorizationService) Authorize(command *AuthorizeCommand) (*oauth.AuthCode, error) {
	clientId := command.ClientId

	// Retrieve the OAuth client
	client, err := a.oauthClientService.FindOauthClient(clientId)
	if err != nil {
		log.Printf("Error retrieving client with ID '%s': %v", clientId, err)
		return nil, fmt.Errorf("failed to retrieve client: %w", err)
	}

	// Validate if response type request is supported, on authorize only supported Code response type
	if !(command.ResponseType == responsetype.Code) {
		log.Printf("Response type '%s' is not supported by the client '%s'", command.ResponseType, clientId)
		// If not, return an error indicating the response type is not supported
		return nil, fmt.Errorf(errors.ErrUnsupportedResponseType)
	}

	// Validate if redirect URI is in the client's list of valid redirect URIs
	if !isRegisteredRedirectUri(command, client) {
		log.Printf("Redirect URI '%s' is not registered for client '%s'", command.RedirectUri, clientId)
		return nil, fmt.Errorf(errors.ErrInvalidRedirectUri)
	}

	// Check if user is authenticated
	if !a.sessionService.SessionExists(command.SessionId) {
		log.Printf("Session ID does not exist, user not authenticated")
		// If not, return an error indicating the user is not authenticated
		return nil, fmt.Errorf(errors.ErrUserNotAuthenticated)
	}

	log.Printf("Session ID '%s' exists", command.SessionId)

	// Retrieve user ID from session
	userId, err := a.sessionService.GetUserIdFromSession(command.SessionId)
	if err != nil {
		log.Printf("Error retrieving user from session ID '%s': %v", command.SessionId, err)
		return nil, fmt.Errorf("failed to retrieve user from session: %w", err)
	}

	// Validate if user is on database in order to add user db reference to auth code
	user, err := a.userRepository.FindByUserId(userId)

	if err != nil {
		log.Printf("Error retrieving user from user ID '%s': %v", userId, err)
		return nil, fmt.Errorf("failed to retrieve user from user ID: %w", err)
	}

	// Generate authorization code
	code, err := utils.GenerateAuthCode(client.ClientId, userId)
	if err != nil {
		log.Printf("Error generating authorization code for client ID '%s' and user ID '%s': %v", client.ClientId, userId, err)
		return nil, fmt.Errorf("failed to generate authorization code: %w", err)
	}

	// Build authorization code entity
	authCodeEntity := store.NewAuthorizationCodeBuilder().
		WithCode(code).
		WithClientId(client.ClientId).
		WithClient(client).
		WithUserId(user.Id).
		WithScope(command.Scope).
		WithRedirectURI(command.RedirectUri).
		WithExpiresAt(time.Now().Add(configuration.AuthCodeExpireTime)). // Set an expiration time
		Build()

	// Save authorization code entity to repository
	authCodeEntity, err = a.authRepository.Save(authCodeEntity)
	if err != nil {
		log.Printf("Error saving authorization code entity: %v", err)
		return nil, fmt.Errorf("failed to save authorization code entity: %w", err)
	}

	// Build the OAuth authorization code response
	oauthCode := oauth.NewAuthCodeBuilder().
		WithCode(authCodeEntity.Code).
		WithClientId(authCodeEntity.ClientId).
		WithScope(authCodeEntity.Scope).
		WithRedirectURI(authCodeEntity.RedirectURI).
		WithCreatedAt(authCodeEntity.CreatedAt).
		WithExpiresAt(authCodeEntity.ExpiresAt).
		Build()

	log.Printf("Successfully generated authorization code for client ID '%s' and user ID '%s'", client.ClientId, userId)

	return oauthCode, nil
}

func isRegisteredRedirectUri(command *AuthorizeCommand, client *store.OauthClient) bool {
	validRedirectUri := false
	for _, uri := range client.RedirectURIs {
		if command.RedirectUri == uri {
			validRedirectUri = true
			break
		}
	}
	return validRedirectUri
}
