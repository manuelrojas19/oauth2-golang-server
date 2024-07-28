package mappers

import (
	"github.com/manuelrojas19/go-oauth2-server/api"
	"github.com/manuelrojas19/go-oauth2-server/oauth"
	"github.com/manuelrojas19/go-oauth2-server/services"
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

func NewCreateOauthClientCommandFromRequest(req *api.RegisterClientRequest) *services.RegisterOauthClientCommand {
	return &services.RegisterOauthClientCommand{
		ClientName:              req.ClientName,
		GrantTypes:              req.GrantTypes,
		ResponseTypes:           req.ResponseTypes,
		TokenEndpointAuthMethod: req.TokenEndpointAuthMethod,
		RedirectUris:            req.RedirectUris,
	}
}
