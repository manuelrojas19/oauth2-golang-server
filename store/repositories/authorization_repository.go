package repositories

import (
	"github.com/manuelrojas19/go-oauth2-server/store/entities"
	"gorm.io/gorm"
)

type authorizationRepository struct {
	Db *gorm.DB
}

func NewAuthorizationRepository(db *gorm.DB) AuthorizationRepository {
	return &authorizationRepository{
		Db: db,
	}
}

func (r *authorizationRepository) Save(authCode *entities.AuthorizationCode) (*entities.AuthorizationCode, error) {
	if err := r.Db.Create(authCode).Error; err != nil {
		return nil, err
	}
	return authCode, nil
}
