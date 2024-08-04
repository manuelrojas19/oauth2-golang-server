package services

import (
	"github.com/google/uuid"
	"github.com/manuelrojas19/go-oauth2-server/oauth"
	"github.com/manuelrojas19/go-oauth2-server/oauth/authmethodtype"
	"github.com/manuelrojas19/go-oauth2-server/oauth/granttype"
	"github.com/manuelrojas19/go-oauth2-server/oauth/responsetype"
	"github.com/manuelrojas19/go-oauth2-server/store"
	"github.com/manuelrojas19/go-oauth2-server/store/repositories"
	"github.com/manuelrojas19/go-oauth2-server/utils"
	"go.uber.org/zap"
	"time"
)

type RegisterOauthClientCommand struct {
	ClientName              string
	GrantTypes              []granttype.GrantType
	ResponseTypes           []responsetype.ResponseType
	TokenEndpointAuthMethod authmethodtype.TokenEndpointAuthMethod
	RedirectUris            []string
	Scopes                  []oauth.Scope
}

type oauthClientService struct {
	oauthClientRepository repositories.OauthClientRepository
	logger                *zap.Logger
}

// NewOauthClientService initializes a new OauthClientService.
func NewOauthClientService(oauthClientRepository repositories.OauthClientRepository, logger *zap.Logger) OauthClientService {
	return &oauthClientService{oauthClientRepository: oauthClientRepository, logger: logger}
}

// CreateOauthClient creates a new OAuth client and returns it.
func (s *oauthClientService) CreateOauthClient(command *RegisterOauthClientCommand) (*oauth.Client, error) {
	start := time.Now()
	clientSecret := uuid.New().String()

	// Encrypt the client secret
	encryptedClientSecret, err := utils.EncryptText(clientSecret)
	if err != nil {
		s.logger.Error("Error encrypting client secret",
			zap.String("clientName", command.ClientName),
			zap.Error(err),
			zap.Duration("duration", time.Since(start)),
			zap.Stack("stacktrace"),
		)
		return nil, err
	}

	// Adding scopes
	scopes := make([]store.Scope, len(command.Scopes))
	for i, scopeData := range command.Scopes {
		scope := *store.NewScopeBuilder().WithName(scopeData.Name).WithDescription(scopeData.Description).Build()
		scopes[i] = scope
	}

	// Build the client entity
	clientEntity := store.NewOauthClientBuilder().
		WithClientName(command.ClientName).
		WithClientSecret(encryptedClientSecret).
		WithResponseTypes(command.ResponseTypes).
		WithGrantTypes(command.GrantTypes).
		WithTokenEndpointAuthMethod(command.TokenEndpointAuthMethod).
		WithRedirectURIs(command.RedirectUris).
		WithScopes(scopes).
		Build()

	s.logger.Info("Client to be created", zap.Any("client", clientEntity))

	// Save the client entity
	savedClient, err := s.oauthClientRepository.Save(clientEntity)
	if err != nil {
		s.logger.Error("Error saving OAuth client",
			zap.String("clientName", command.ClientName),
			zap.Error(err),
			zap.Duration("duration", time.Since(start)),
			zap.Stack("stacktrace"),
		)
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
		WithScopes(command.Scopes).
		Build()

	s.logger.Info("Successfully created OAuth client",
		zap.String("clientId", savedClient.ClientId),
		zap.String("clientName", savedClient.ClientName),
		zap.Duration("duration", time.Since(start)),
	)

	return clientModel, nil
}

// FindOauthClient retrieves an OAuth client by its client ID.
func (s *oauthClientService) FindOauthClient(clientId string) (*store.OauthClient, error) {
	start := time.Now()
	client, err := s.oauthClientRepository.FindByClientId(clientId)
	if err != nil {
		s.logger.Error("Error finding OAuth client by clientId",
			zap.String("clientId", clientId),
			zap.Error(err),
			zap.Duration("duration", time.Since(start)),
			zap.Stack("stacktrace"),
		)
		return nil, err
	}

	s.logger.Info("Successfully found OAuth client",
		zap.String("clientId", clientId),
		zap.Duration("duration", time.Since(start)),
	)

	return client, nil
}
