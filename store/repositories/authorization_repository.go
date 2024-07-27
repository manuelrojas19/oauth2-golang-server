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

// FindByCode retrieves an authorization code by its code string.
func (r *authorizationRepository) FindByCode(code string) (*entities.AuthorizationCode, error) {
	var authCode entities.AuthorizationCode
	if err := r.Db.Where("code = ?", code).First(&authCode).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // Return nil if the record is not found
		}
		return nil, err
	}
	return &authCode, nil
}
