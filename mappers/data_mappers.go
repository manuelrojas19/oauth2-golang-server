package mappers

import (
	"github.com/manuelrojas19/go-oauth2-server/api/response"
	"github.com/manuelrojas19/go-oauth2-server/models/oauth"
	"github.com/manuelrojas19/go-oauth2-server/store/entities"
)

func NewClientModelFromClientEntity(client *entities.OauthClient) *oauth.Client {
	return &oauth.Client{
		ClientId:                client.ClientId,
		ClientSecret:            client.ClientSecret,
		ClientName:              client.ClientName,
		GrantTypes:              client.GrantTypes,
		ResponseTypes:           client.ResponseTypes,
		TokenEndpointAuthMethod: client.TokenEndpointAuthMethod,
		RedirectUris:            *client.RedirectURI,
	}
}

func NewRegisterClientResponseFromClientModel(client *oauth.Client) *response.RegisterClientResponse {
	return &response.RegisterClientResponse{
		ClientId:                client.ClientId,
		ClientName:              client.ClientName,
		GrantTypes:              client.GrantTypes,
		ResponseTypes:           client.ResponseTypes,
		TokenEndpointAuthMethod: client.TokenEndpointAuthMethod,
		RedirectUris:            client.RedirectUris,
	}
}
