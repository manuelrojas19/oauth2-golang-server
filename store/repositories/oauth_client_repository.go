package repositories

import (
	"errors"

	"github.com/manuelrojas19/go-oauth2-server/store/entities"
	"gorm.io/gorm"
)

type oauthClientRepository struct {
	Db *gorm.DB
}

func NewOauthClientRepository(db *gorm.DB) OauthClientRepository {
	return &oauthClientRepository{Db: db}
}

func (ocd *oauthClientRepository) Save(client *entities.OauthClient) (*entities.OauthClient, error) {
	if ocd.clientExists(client.ClientId) {
		return nil, errors.New("client Already Exists")
	}

	if err := ocd.Db.Create(client).Error; err != nil {
		return nil, err
	}

	return client, nil
}

func (ocd *oauthClientRepository) FindByClientId(clientId string) (*entities.OauthClient, error) {
	oauthClient := new(entities.OauthClient)

	result := ocd.Db.Where("client_id = LOWER(?)", clientId).First(oauthClient)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, errors.New("OAuth Client not Found")
	}

	return oauthClient, nil
}

func (ocd *oauthClientRepository) clientExists(clientKey string) bool {
	_, err := ocd.FindByClientId(clientKey)
	return err == nil
}
