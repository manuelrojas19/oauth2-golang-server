package mappers

import (
	"github.com/manuelrojas19/go-oauth2-server/api/dto/request"
	"github.com/manuelrojas19/go-oauth2-server/api/dto/response"
	"github.com/manuelrojas19/go-oauth2-server/models/oauth"
	"github.com/manuelrojas19/go-oauth2-server/services/commands"
)

func NewRegisterClientResponseFromClientModel(client *oauth.Client) *response.RegisterClientResponse {
	return &response.RegisterClientResponse{
		ClientId:                client.ClientId,
		ClientSecret:            client.ClientSecret,
		ClientName:              client.ClientName,
		GrantTypes:              client.GrantTypes,
		ResponseTypes:           client.ResponseTypes,
		TokenEndpointAuthMethod: client.TokenEndpointAuthMethod,
		RedirectUris:            client.RedirectUris,
	}
}

func NewCreateOauthClientCommandFromRequest(req *request.RegisterClientRequest) *commands.RegisterOauthClientCommand {
	return &commands.RegisterOauthClientCommand{
		ClientName:              req.ClientName,
		GrantTypes:              req.GrantTypes,
		ResponseTypes:           req.ResponseTypes,
		TokenEndpointAuthMethod: req.TokenEndpointAuthMethod,
		RedirectUris:            req.RedirectUris,
	}
}
