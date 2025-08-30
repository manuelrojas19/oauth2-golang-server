package services

import (
	"fmt"

	"github.com/manuelrojas19/go-oauth2-server/oauth"
	"github.com/manuelrojas19/go-oauth2-server/store/repositories"
	"go.uber.org/zap"
)

type scopeService struct {
	repo   repositories.ScopeRepository
	logger *zap.Logger
}

func NewScopeService(repo repositories.ScopeRepository, logger *zap.Logger) ScopeService {
	return &scopeService{repo: repo, logger: logger}
}

func (s *scopeService) Save(scopeName, scopeDescription string) (*oauth.Scope, error) {
	s.logger.Info("Saving scope", zap.String("name", scopeName), zap.String("description", scopeDescription))
	s.logger.Debug("Calling repository to create scope", zap.String("name", scopeName))

	scope, err := s.repo.Create(scopeName, scopeDescription)
	if err != nil {
		s.logger.Error("Failed to save scope", zap.Error(err), zap.String("scopeName", scopeName))
		return nil, err
	}

	s.logger.Info("Scope saved successfully", zap.String("id", scope.Id))
	s.logger.Debug("Created scope details", zap.Any("scope", scope))

	createdScope := &oauth.Scope{
		Id:          scope.Id,
		Name:        scope.Name,
		Description: scope.Description,
	}

	return createdScope, nil
}

func (s *scopeService) FindById(scopeId string) (*oauth.Scope, bool) {
	s.logger.Info("Finding scope by Id", zap.String("id", scopeId))
	s.logger.Debug("Calling repository to find scope by Id", zap.String("id", scopeId))

	scopeEntity, err := s.repo.FindById(scopeId)
	if err != nil {
		s.logger.Error("Failed to find scope by Id", zap.Error(err), zap.String("scopeId", scopeId))
		return nil, false
	}

	if scopeEntity == nil {
		s.logger.Info("Scope not found for Id", zap.String("id", scopeId))
		return nil, false
	}

	s.logger.Info("Scope found", zap.String("id", scopeEntity.Id))
	s.logger.Debug("Found scope details", zap.Any("scopeEntity", scopeEntity))

	scope := &oauth.Scope{
		Id:          scopeEntity.Id,
		Name:        scopeEntity.Name,
		Description: scopeEntity.Description,
	}

	return scope, true
}

func (s *scopeService) FindByIdList(scopeIds []string) ([]oauth.Scope, error) {
	s.logger.Info("Finding scopes by Id list", zap.Strings("ids", scopeIds))
	s.logger.Debug("Calling repository to find scopes by Id list", zap.Strings("ids", scopeIds))

	scopeEntities, err := s.repo.FindByIdList(scopeIds)
	if err != nil {
		s.logger.Error("Failed to find scopes by Id list", zap.Error(err), zap.Strings("scopeIds", scopeIds))
		return nil, err
	}

	if len(scopeEntities) == 0 {
		s.logger.Info("No scopes found for the given Id list", zap.Strings("ids", scopeIds))
		return nil, fmt.Errorf("no scopes found for the given Id list")
	}

	var oauthScopes []oauth.Scope
	for _, scopeEntity := range scopeEntities {
		oauthScopes = append(oauthScopes, oauth.Scope{
			Id:          scopeEntity.Id,
			Name:        scopeEntity.Name,
			Description: scopeEntity.Description,
		})
	}

	s.logger.Info("Scopes found", zap.Int("count", len(oauthScopes)))
	s.logger.Debug("Found scopes details", zap.Any("scopes", oauthScopes))
	return oauthScopes, nil
}
