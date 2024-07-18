package oauth

import (
	"time"

	"github.com/google/uuid"
	"github.com/manuelrojas19/go-oauth2-server/store/dao"
	"github.com/manuelrojas19/go-oauth2-server/store/entities"
	"github.com/manuelrojas19/go-oauth2-server/utils"
)

type OauthToken struct {
	OauthTokenDao *dao.OauthTokenDao
	OauthClient   *OauthClient
}

func (ot *OauthToken) GrantAccessToken(clientKey string) (string, error) {
	client, err := ot.OauthClient.FindOauthClient(clientKey)
	if err != nil {
		return "", err
	}

	generatedToken, err := utils.Token(clientKey, client.Id, time.Now())
	if err != nil {
		return "", err
	}

	token := &entities.OauthTokenEntity{
		BaseGormEntity: entities.BaseGormEntity{
			Id:        uuid.New().String(),
			CreatedAt: time.Now(),
		},
		ClientId:  client.Id,
		ClientKey: client.Key,
		Client:    client,
		ExpiresAt: time.Now(),
		Scope:     "Grant",
		Token:     generatedToken,
	}

	if token, err := ot.OauthTokenDao.SaveToken(token); err != nil {
		return "", err
	} else {
		return token.Token, nil
	}
}
