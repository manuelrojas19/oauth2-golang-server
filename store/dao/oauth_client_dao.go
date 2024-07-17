package dao

import (
	"github.com/manuelrojas19/go-oauth2-server/store/entities"
	"gorm.io/gorm"
)

type OauthClientDao struct {
	Db *gorm.DB
}

func (ocd *OauthClientDao) Save(client *entities.OauthClientEntity) (*entities.OauthClientEntity, error) {
	if error := ocd.Db.Create(client).Error; error != nil {
		return nil, error
	}
	return client, nil
}
