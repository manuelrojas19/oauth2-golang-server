package repositories

import (
	"fmt"

	"github.com/manuelrojas19/go-oauth2-server/store"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type scopeRepository struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewScopeRepository(db *gorm.DB, logger *zap.Logger) ScopeRepository {
	return &scopeRepository{db: db, logger: logger}
}

func (s *scopeRepository) Create(name, description string) (*store.Scope, error) {
	s.logger.Info("Creating new scope", zap.String("name", name), zap.String("description", description))
	scope := store.NewScopeBuilder().WithName(name).WithDescription(description).Build()
	s.logger.Debug("Scope entity built for creation", zap.Any("scope", scope))

	if err := s.db.Create(scope).Error; err != nil {
		s.logger.Error("Failed to create scope in database", zap.String("name", name), zap.Error(err))
		return nil, fmt.Errorf("failed to create scope: %w", err)
	}
	s.logger.Info("Scope created successfully", zap.String("id", scope.Id), zap.String("name", scope.Name))
	return scope, nil
}

func (s *scopeRepository) FindById(id string) (*store.Scope, error) {
	s.logger.Info("Finding scope by ID", zap.String("id", id))
	var scope *store.Scope

	if err := s.db.Where("id = ?", id).First(&scope).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			s.logger.Debug("Scope not found for ID", zap.String("id", id))
			return nil, fmt.Errorf("scope with ID '%s' not found", id)
		}
		s.logger.Error("Failed to find scope by ID in database", zap.String("id", id), zap.Error(err))
		return nil, fmt.Errorf("failed to find scope: %w", err)
	}
	s.logger.Info("Scope found successfully", zap.String("id", scope.Id), zap.String("name", scope.Name))
	s.logger.Debug("Found scope details", zap.Any("scope", scope))
	return scope, nil
}

func (s *scopeRepository) FindByIdList(ids []string) ([]*store.Scope, error) {
	s.logger.Info("Finding scopes by ID list", zap.Strings("ids", ids))
	var scopes []*store.Scope

	if err := s.db.Where("id IN ?", ids).Find(&scopes).Error; err != nil {
		s.logger.Error("Failed to find scopes by ID list in database", zap.Strings("ids", ids), zap.Error(err))
		return nil, fmt.Errorf("failed to find scopes by ID list: %w", err)
	}

	if len(scopes) == 0 {
		s.logger.Info("No scopes found for the given ID list", zap.Strings("ids", ids))
		return nil, fmt.Errorf("no scopes found for the given ID list")
	}
	s.logger.Info("Scopes found successfully", zap.Int("count", len(scopes)))
	s.logger.Debug("Found scopes details", zap.Any("scopes", scopes))
	return scopes, nil
}

func (s *scopeRepository) Exists(id string) (bool, error) {
	s.logger.Info("Checking if scope exists", zap.String("id", id))
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM scopes WHERE id = ?)"

	if err := s.db.Raw(query, id).Scan(&exists).Error; err != nil {
		s.logger.Error("Failed to check scope existence in database", zap.String("id", id), zap.Error(err))
		return false, fmt.Errorf("failed to check scope existence: %w", err)
	}
	s.logger.Info("Scope existence check completed", zap.String("id", id), zap.Bool("exists", exists))
	return exists, nil
}
