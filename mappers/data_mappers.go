package mappers

import (
	"github.com/manuelrojas19/go-oauth2-server/api"
	"github.com/manuelrojas19/go-oauth2-server/oauth"
)

func NewRegisterClientResponseFromClientModel(client *oauth.Client) *api.RegisterClientResponse {
	return &api.RegisterClientResponse{
		ClientId:                client.ClientId,
		ClientSecret:            client.ClientSecret,
		ClientName:              client.ClientName,
		GrantTypes:              client.GrantTypes,
		ResponseTypes:           client.ResponseTypes,
		TokenEndpointAuthMethod: client.TokenEndpointAuthMethod,
		RedirectUris:            client.RedirectUris,
	}
}
