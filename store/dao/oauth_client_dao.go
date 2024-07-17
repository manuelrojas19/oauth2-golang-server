package dao

import (
	"errors"

	"github.com/manuelrojas19/go-oauth2-server/store/entities"
	"gorm.io/gorm"
)

type OauthClientDao struct {
	Db *gorm.DB
}

func (ocd *OauthClientDao) SaveClient(client *entities.OauthClientEntity) (*entities.OauthClientEntity, error) {
	if ocd.clientExists(client.Key) {
		return nil, errors.New("Client Already Exists")
	}

	if error := ocd.Db.Create(client).Error; error != nil {
		return nil, error
	}

	return client, nil
}

func (ocd *OauthClientDao) FindClientByClientKey(clientKey string) (*entities.OauthClientEntity, error) {
	oauthClient := new(entities.OauthClientEntity)

	result := ocd.Db.Where("key = LOWER(?)", clientKey).First(oauthClient)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, errors.New("OAuth Client not Found")
	}

	return oauthClient, nil
}

func (ocd *OauthClientDao) clientExists(clientKey string) bool {
	_, error := ocd.FindClientByClientKey(clientKey)
	return error == nil
}
