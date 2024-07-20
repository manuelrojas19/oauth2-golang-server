package services

import (
	"github.com/google/uuid"
	"github.com/manuelrojas19/go-oauth2-server/mappers"
	"github.com/manuelrojas19/go-oauth2-server/models/oauth"
	"github.com/manuelrojas19/go-oauth2-server/services/commands"
	"github.com/manuelrojas19/go-oauth2-server/store/entities"
	"github.com/manuelrojas19/go-oauth2-server/store/repositories"
	"github.com/manuelrojas19/go-oauth2-server/utils"
	"log"
)

type oauthClientService struct {
	oauthClientRepository repositories.OauthClientRepository
}

// NewOauthClientService initializes a new OauthClientService.
func NewOauthClientService(oauthClientRepository repositories.OauthClientRepository) OauthClientService {
	return &oauthClientService{oauthClientRepository: oauthClientRepository}
}

// CreateOauthClient creates a new OAuth client and returns it.
func (s *oauthClientService) CreateOauthClient(command *commands.RegisterOauthClientCommand) (*oauth.Client, error) {
	// Encrypt the client secret
	clientSecret, err := utils.EncryptText(uuid.New().String())
	if err != nil {
		log.Printf("Error encrypting client secret: %v", err)
		return nil, err
	}

	// Build the client entity
	clientEntity := entities.NewOauthClientBuilder().
		SetClientName(command.ClientName).
		SetClientSecret(clientSecret).
		SetResponseTypes(command.ResponseTypes).
		SetGrantTypes(command.GrantTypes).
		SetTokenEndpointAuthMethod(command.TokenEndpointAuthMethod).
		SetRedirectURI(command.RedirectUris).
		Build()

	// Save the client entity
	savedClient, err := s.oauthClientRepository.Save(clientEntity)
	if err != nil {
		log.Printf("Error saving OAuth client: %v", err)
		return nil, err
	}

	// Map to responsetype model
	clientModel := mappers.NewClientModelFromClientEntity(savedClient)
	return clientModel, nil
}

// FindOauthClient retrieves an OAuth client by its client ID.
func (s *oauthClientService) FindOauthClient(clientID string) (*entities.OauthClient, error) {
	client, err := s.oauthClientRepository.FindByClientId(clientID)
	if err != nil {
		log.Printf("Error finding OAuth client by ID: %v", err)
		return nil, err
	}
	return client, nil
}
