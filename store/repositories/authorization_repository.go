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

func (authorizationRepository) Save(authCode *entities.AuthorizationCode) (*entities.AuthorizationCode, error) {
	//TODO implement me
	panic("implement me")
}
