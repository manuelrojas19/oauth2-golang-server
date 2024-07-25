package services

import (
	"fmt"
	"github.com/manuelrojas19/go-oauth2-server/models/oauth"
	"github.com/manuelrojas19/go-oauth2-server/services/commands"
	"github.com/manuelrojas19/go-oauth2-server/store/entities"
	"github.com/manuelrojas19/go-oauth2-server/store/repositories"
	"log"
)

const (
	ErrNotAuthenticated = "user not authenticated"
	ErrConsentRequired  = "user consent required"
)

type authorizationService struct {
	oauthClientService OauthClientService
	consentService     UserConsentService
	authRepository     repositories.AuthorizationRepository
}

// NewAuthorizationService initializes a new AuthorizationService
func NewAuthorizationService(oauthClientService OauthClientService,
	consentService UserConsentService,
	authRepository repositories.AuthorizationRepository) AuthorizationService {
	return &authorizationService{
		oauthClientService: oauthClientService,
		consentService:     consentService,
		authRepository:     authRepository,
	}
}

func (a authorizationService) Authorize(command *commands.Authorization) (*oauth.AuthCode, error) {
	clientId := command.ClientId
	client, err := a.oauthClientService.FindOauthClient(clientId)
	if err != nil {
		log.Printf("Error retrieving client with ID '%s': %v", clientId, err)
		return nil, err
	}

	user := entities.User{
		ID: "id",
	}

	// Get user consent
	if !userIsAuthenticated(user.ID, command.ClientId, command.Scope) {
		// If not, present the consent screen (this could be a redirect to a consent page)
		return nil, fmt.Errorf(ErrNotAuthenticated)
	}

	// Get user consent
	if !a.consentService.HasUserConsented(user.ID, command.ClientId, command.Scope) {
		// If not, present the consent screen (this could be a redirect to a consent page)
		return nil, fmt.Errorf(ErrConsentRequired)
	}

	// Generate authorization code
	code, err := generateAuthCode(client.ClientId, user.ID, command.Scope)
	if err != nil {
		return nil, fmt.Errorf("failed to generate authorization code: %w", err)
	}

	return code, nil
}

func generateAuthCode(clientId string, userId string, scope string) (*oauth.AuthCode, error) {
	return &oauth.AuthCode{}, nil
}

func userIsAuthenticated(userId string, clientId string, scope string) bool {
	return false
}

func hasUserConsented(userId string, clientId string, scope string) bool {
	return true
}

func authenticateUser(userId string) (*entities.User, error) {
	return &entities.User{}, nil
}
