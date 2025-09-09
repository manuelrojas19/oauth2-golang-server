package services

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/manuelrojas19/go-oauth2-server/api"
	"github.com/manuelrojas19/go-oauth2-server/oauth"
	"github.com/manuelrojas19/go-oauth2-server/oauth/authmethodtype"
	"github.com/manuelrojas19/go-oauth2-server/oauth/granttype"
	"github.com/manuelrojas19/go-oauth2-server/oauth/responsetype"
	"github.com/manuelrojas19/go-oauth2-server/store"
	"github.com/manuelrojas19/go-oauth2-server/store/repositories"
	"github.com/manuelrojas19/go-oauth2-server/utils"
	"go.uber.org/zap"
)

type RegisterOauthClientCommand struct {
	ClientName              string
	GrantTypes              []granttype.GrantType
	ResponseTypes           []responsetype.ResponseType
	TokenEndpointAuthMethod authmethodtype.TokenEndpointAuthMethod
	RedirectUris            []string
	Scopes                  string
}

type oauthClientService struct {
	oauthClientRepository repositories.OauthClientRepository
	logger                *zap.Logger
	scopeRepository       repositories.ScopeRepository
}

// NewOauthClientService initializes a new OauthClientService.
func NewOauthClientService(oauthClientRepository repositories.OauthClientRepository, scopeRepository repositories.ScopeRepository, logger *zap.Logger) OauthClientService {
	return &oauthClientService{oauthClientRepository: oauthClientRepository, scopeRepository: scopeRepository, logger: logger}
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
		return nil, api.ErrServerError
	}

	// Validate and fetch scopes
	var clientScopes []store.Scope
	scopeNames := splitAndTrim(command.Scopes)

	for _, scopeName := range scopeNames {
		scope, err := s.scopeRepository.FindByName(scopeName)
		if err != nil {
			s.logger.Warn("Scope validation failed",
				zap.String("scope", scopeName),
				zap.Error(err),
			)
			return nil, api.ErrInvalidScope
		}
		clientScopes = append(clientScopes, *scope)
	}

	s.logger.Debug("Scopes successfully validated",
		zap.Int("count", len(clientScopes)),
		zap.Any("scopes", clientScopes),
	)

	if s.oauthClientRepository.ExistsByName(command.ClientName) {
		s.logger.Warn("Client name already in use", zap.String("clientName", command.ClientName))
		return nil, api.ErrClientAlreadyExists
	}

	// Build the client entity
	clientEntity := store.NewOauthClientBuilder().
		WithClientName(command.ClientName).
		WithClientSecret(encryptedClientSecret).
		WithResponseTypes(command.ResponseTypes).
		WithGrantTypes(command.GrantTypes).
		WithTokenEndpointAuthMethod(command.TokenEndpointAuthMethod).
		WithRedirectURIs(command.RedirectUris).
		WithScopes(clientScopes).
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
		return nil, api.ErrServerError
	}

	// Map to Client model
	clientModel := oauth.NewClientBuilder().
		WithClientId(savedClient.ClientId).
		WithClientSecret(clientSecret).
		WithClientName(savedClient.ClientName).
		WithClientIdIssuedAt(savedClient.CreatedAt.Unix()).
		WithClientSecretExpiresAt(savedClient.ClientSecretExpiresAt).
		WithResponseTypes(responsetype.StringListToEnumList(savedClient.ResponseTypes)).
		WithGrantTypes(granttype.StringListToEnumList(savedClient.GrantTypes)).
		WithTokenEndpointAuthMethod(authmethodtype.TokenEndpointAuthMethod(savedClient.TokenEndpointAuthMethod)).
		WithRedirectUris(savedClient.RedirectURIs).
		WithScopes(oauthScopesFromStoreScopes(savedClient.Scopes)).
		Build()

	s.logger.Info("Successfully created OAuth client",
		zap.String("clientId", savedClient.ClientId),
		zap.String("clientName", savedClient.ClientName),
		zap.Duration("duration", time.Since(start)),
	)

	return clientModel, nil
}

// FindOauthClient retrieves an OAuth client by its client Id.
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

// PreloadOauthClientScopes preloads the associated scopes for a given OauthClient.
func (s *oauthClientService) PreloadOauthClientScopes(client *store.OauthClient) error {
	start := time.Now()

	err := s.oauthClientRepository.PreloadScopes(client)
	if err != nil {
		s.logger.Error("Error preloading scopes for OAuth client",
			zap.String("clientId", client.ClientId),
			zap.Error(err),
			zap.Duration("duration", time.Since(start)),
			zap.Stack("stacktrace"),
		)
		return fmt.Errorf("failed to preload client scopes: %w", err)
	}

	s.logger.Debug("Scopes preloaded successfully for OAuth client",
		zap.String("clientId", client.ClientId),
		zap.Int("scopeCount", len(client.Scopes)),
		zap.Duration("duration", time.Since(start)),
	)
	return nil
}

func oauthScopesFromStoreScopes(storeScopes []store.Scope) []oauth.Scope {
	oauthScopes := make([]oauth.Scope, 0, len(storeScopes))
	for _, storeScope := range storeScopes {
		oauthScopes = append(oauthScopes, oauth.Scope{
			Name:        storeScope.Name,
			Description: storeScope.Description,
		})
	}
	return oauthScopes
}

func splitAndTrim(scopes string) []string {
	if scopes == "" {
		return []string{}
	}
	rawScopes := strings.Fields(scopes)
	trimmedScopes := make([]string, 0, len(rawScopes))
	for _, s := range rawScopes {
		t := strings.TrimSpace(s)
		if t != "" {
			trimmedScopes = append(trimmedScopes, t)
		}
	}
	return trimmedScopes
}
