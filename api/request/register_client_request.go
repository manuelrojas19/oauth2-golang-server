package request

import "github.com/manuelrojas19/go-oauth2-server/models/oauth"

type RegisterClientRequest struct {
	ClientName              string                        `json:"client_name"`
	GrantTypes              []oauth.GrantType             `json:"grant_types"`
	ResponseTypes           []oauth.ResponseType          `json:"response_types"`
	TokenEndpointAuthMethod oauth.TokenEndpointAuthMethod `json:"token_endpoint_auth_method"`
	RedirectUris            []string                      `json:"redirect_uris"`
}
