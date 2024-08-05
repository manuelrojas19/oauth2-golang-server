package services

import (
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

	scope, err := s.repo.Create(scopeName, scopeDescription)
	if err != nil {
		s.logger.Error("Failed to save scope", zap.Error(err))
		return nil, err
	}

	s.logger.Info("Scope saved successfully", zap.String("id", scope.Id))

	createdScope := &oauth.Scope{
		Id:          scope.Id,
		Name:        scope.Name,
		Description: scope.Description,
	}

	return createdScope, nil
}

func (s *scopeService) FindById(scopeId string) (*oauth.Scope, bool) {
	s.logger.Info("Finding scope by Id", zap.String("id", scopeId))

	scopeEntity, err := s.repo.FindById(scopeId)
	if err != nil {
		s.logger.Error("Failed to find scope", zap.Error(err))
		return nil, false
	}

	s.logger.Info("Scope found", zap.String("id", scopeEntity.Id))

	scope := &oauth.Scope{
		Id:          scopeEntity.Id,
		Name:        scopeEntity.Name,
		Description: scopeEntity.Description,
	}

	return scope, true
}

func (s *scopeService) FindByIdList(scopeIds []string) (*oauth.Scope, error) {
	s.logger.Info("Finding scopes by Id list", zap.Strings("ids", scopeIds))

	scopes, err := s.repo.FindByIdList(scopeIds)
	if err != nil {
		s.logger.Error("Failed to find scopes", zap.Error(err))
		return nil, err
	}

	s.logger.Info("Scopes found", zap.Int("count", len(scopes)))
	return nil, nil
}
