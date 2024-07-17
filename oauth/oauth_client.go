package oauth

import (
	"time"

	"github.com/google/uuid"
	"github.com/manuelrojas19/go-oauth2-server/store/dao"
	"github.com/manuelrojas19/go-oauth2-server/store/entities"
	"golang.org/x/crypto/bcrypt"
)

type OauthClient struct {
	OauthClientDao *dao.OauthClientDao
}

func (oc *OauthClient) CreateOauthClient(clientKey string, clientSecret string, redirectUri string) (*entities.OauthClientEntity, error) {
	clientSecretHash, err := bcrypt.GenerateFromPassword([]byte(clientSecret), 3)
	if err != nil {
		return nil, err
	}

	client := &entities.OauthClientEntity{
		BaseGormEntity: entities.BaseGormEntity{
			Id:        uuid.New().String(),
			CreatedAt: time.Now().UTC(),
		},
		Key:         clientKey,
		Secret:      string(clientSecretHash),
		RedirectURI: redirectUri,
	}

	return oc.OauthClientDao.SaveClient(client)
}

func (oc *OauthClient) FindOauthClient(clientKey string) (*entities.OauthClientEntity, error) {
	return oc.OauthClientDao.FindClientByClientKey(clientKey)
}
