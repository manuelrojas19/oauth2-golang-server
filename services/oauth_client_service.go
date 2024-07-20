package services

import (
	"github.com/google/uuid"
	"github.com/manuelrojas19/go-oauth2-server/api/request"
	"github.com/manuelrojas19/go-oauth2-server/mappers"
	"github.com/manuelrojas19/go-oauth2-server/models/oauth"
	"github.com/manuelrojas19/go-oauth2-server/store/entities"
	"github.com/manuelrojas19/go-oauth2-server/store/repositories"
	"github.com/manuelrojas19/go-oauth2-server/utils"
)

type oauthClientService struct {
	oauthClientRepository repositories.OauthClientRepository
}

func NewOauthClientService(oauthClientRepository repositories.OauthClientRepository) OauthClientService {
	return &oauthClientService{oauthClientRepository: oauthClientRepository}
}

func (c *oauthClientService) CreateOauthClient(request *request.RegisterClientRequest) (*oauth.Client, error) {
	client := mappers.NewClientModelFromRegisterClientRequest(request)

	clientSecret, err := utils.EncryptText(uuid.New().String())
	if err != nil {
		return nil, err
	}
	client.ClientSecret = clientSecret

	savedClient, err := c.oauthClientRepository.Save(entities.NewOauthClientEntityFromModel(client))
	if err != nil {
		return nil, err
	}

	return mappers.NewClientModelFromClientEntity(savedClient), nil

}

func (c *oauthClientService) FindOauthClient(clientKey string) (*entities.OauthClient, error) {
	return c.oauthClientRepository.FindByClientId(clientKey)
}
