package services

import (
	"github.com/google/uuid"
	"github.com/manuelrojas19/go-oauth2-server/models/oauth"
	"github.com/manuelrojas19/go-oauth2-server/models/oauth/authmethodtype"
	"github.com/manuelrojas19/go-oauth2-server/models/oauth/granttype"
	"github.com/manuelrojas19/go-oauth2-server/models/oauth/responsetype"
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
	clientSecret := uuid.New().String()

	encryptedClientSecret, err := utils.EncryptText(clientSecret)
	if err != nil {
		log.Printf("Error encrypting client secret: %v", err)
		return nil, err
	}

	// Build the client entity
	clientEntity := entities.NewOauthClientBuilder().
		WithClientName(command.ClientName).
		WithClientSecret(encryptedClientSecret).
		WithResponseTypes(command.ResponseTypes).
		WithGrantTypes(command.GrantTypes).
		WithTokenEndpointAuthMethod(command.TokenEndpointAuthMethod).
		WithRedirectURI(command.RedirectUris).
		Build()

	// Save the client entity
	savedClient, err := s.oauthClientRepository.Save(clientEntity)
	if err != nil {
		log.Printf("Error saving OAuth client: %v", err)
		return nil, err
	}

	// Map to Client model
	clientModel := oauth.NewClientBuilder().
		WithClientId(savedClient.ClientId).
		WithClientSecret(clientSecret).
		WithClientName(savedClient.ClientName).
		WithResponseTypes(responsetype.StringListToEnumList(savedClient.ResponseTypes)).
		WithGrantTypes(granttype.StringListToEnumList(savedClient.GrantTypes)).
		WithTokenEndpointAuthMethod(authmethodtype.TokenEndpointAuthMethod(savedClient.TokenEndpointAuthMethod)).
		WithRedirectUris(savedClient.RedirectURIs).
		Build()
	return clientModel, nil
}

// FindOauthClient retrieves an OAuth client by its client Id.
func (s *oauthClientService) FindOauthClient(clientId string) (*entities.OauthClient, error) {
	client, err := s.oauthClientRepository.FindByClientId(clientId)
	if err != nil {
		log.Printf("Error finding OAuth client by Id: %v", err)
		return nil, err
	}
	return client, nil
}
