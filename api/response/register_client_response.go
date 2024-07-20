package response

import "github.com/manuelrojas19/go-oauth2-server/models/oauth"

type RegisterClientResponse struct {
	ClientId                string                        `json:"client_id"`
	ClientName              string                        `json:"client_name"`
	GrantTypes              []oauth.GrantType             `json:"grant_types"`
	ResponseTypes           []oauth.ResponseType          `json:"response_types"`
	TokenEndpointAuthMethod oauth.TokenEndpointAuthMethod `json:"token_endpoint_auth_method"`
	RedirectUris            []string                      `json:"redirect_uris"`
}
