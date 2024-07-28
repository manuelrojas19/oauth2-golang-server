package services

import (
	"fmt"
	"github.com/manuelrojas19/go-oauth2-server/oauth"
	"github.com/manuelrojas19/go-oauth2-server/store"
	"github.com/manuelrojas19/go-oauth2-server/utils"
	"log"
	"time"
)

const (
	ErrUserNotAuthenticated = "user not authenticated"
	ErrConsentRequired      = "user consent required"
)

type AuthorizeCommand struct {
	ClientId     string
	Scope        string
	RedirectUri  string
	ResponseType string
	SessionId    string
	State        string
}

type authorizationService struct {
	oauthClientService OauthClientService
	consentService     UserConsentService
	authRepository     store.AuthorizationRepository
	sessionService     SessionService
}

// NewAuthorizationService initializes a new AuthorizationService
func NewAuthorizationService(oauthClientService OauthClientService,
	consentService UserConsentService, authRepository store.AuthorizationRepository,
	userSessionService SessionService) AuthorizationService {
	return &authorizationService{oauthClientService: oauthClientService,
		consentService: consentService,
		authRepository: authRepository,
		sessionService: userSessionService}
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

	// Check if user is authenticated
	if !a.sessionService.SessionExists(command.SessionId) {
		log.Printf("Session ID does not exist, user not authenticated")
		// If not, return an error indicating the user is not authenticated
		return nil, fmt.Errorf(ErrUserNotAuthenticated)
	}

	log.Printf("Session ID '%s' exists", command.SessionId)

	// Retrieve user ID from session
	userId, err := a.sessionService.GetUserIdFromSession(command.SessionId)
	if err != nil {
		log.Printf("Error retrieving user from session ID '%s': %v", command.SessionId, err)
		return nil, fmt.Errorf("failed to retrieve user from session: %w", err)
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
		WithClientID(client.ClientId).
		WithClient(client).
		WithRedirectURI(command.RedirectUri).
		WithExpiresAt(time.Now().Add(10 * time.Minute)). // Set an expiration time
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
		WithRedirectURI(authCodeEntity.RedirectURI).
		WithCreatedAt(authCodeEntity.CreatedAt).
		WithExpiresAt(authCodeEntity.ExpiresAt).
		Build()

	log.Printf("Successfully generated authorization code for client ID '%s' and user ID '%s'", client.ClientId, userId)

	return oauthCode, nil
}
