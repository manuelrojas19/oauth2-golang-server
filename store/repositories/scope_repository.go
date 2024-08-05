package repositories

import (
	"github.com/manuelrojas19/go-oauth2-server/store"
	"gorm.io/gorm"
)

type scopeRepository struct {
	db *gorm.DB
}

func NewScopeRepository(db *gorm.DB) ScopeRepository {
	return &scopeRepository{db: db}
}

func (s *scopeRepository) Create(name, description string) (*store.Scope, error) {
	scope := store.NewScopeBuilder().WithName(name).WithDescription(description).Build()
	if err := s.db.Create(scope).Error; err != nil {
		return nil, err
	}
	return scope, nil
}

func (s *scopeRepository) FindById(id string) (*store.Scope, error) {
	var scope *store.Scope
	if err := s.db.Where("id = ?", id).Find(&scope).Error; err != nil {
		return nil, err
	}
	return scope, nil
}

func (s *scopeRepository) FindByIdList(ids []string) ([]*store.Scope, error) {
	var scopes []*store.Scope
	if err := s.db.Where("id IN ?", ids).Find(&scopes).Error; err != nil {
		return nil, err
	}
	return scopes, nil
}

func (s *scopeRepository) Exists(id string) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM scopes WHERE id = ?)"
	if err := s.db.Raw(query, id).Scan(&exists).Error; err != nil {
		return false, err
	}
	return exists, nil
}
