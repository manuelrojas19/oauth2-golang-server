package response

import (
	"github.com/manuelrojas19/go-oauth2-server/models/oauth/authmethodtype"
	"github.com/manuelrojas19/go-oauth2-server/models/oauth/granttype"
	"github.com/manuelrojas19/go-oauth2-server/models/oauth/responsetype"
)

type RegisterClientResponse struct {
	ClientId                string                                 `json:"client_id"`
	ClientSecret            string                                 `json:"client_secret,omitempty"`
	ClientName              string                                 `json:"client_name"`
	GrantTypes              []granttype.GrantType                  `json:"granttype"`
	ResponseTypes           []responsetype.ResponseType            `json:"responsetype"`
	TokenEndpointAuthMethod authmethodtype.TokenEndpointAuthMethod `json:"token_endpoint_auth_method"`
	RedirectUris            []string                               `json:"redirect_uris"`
}
