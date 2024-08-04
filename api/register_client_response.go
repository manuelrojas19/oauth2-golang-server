package api

import (
	"github.com/manuelrojas19/go-oauth2-server/oauth"
	"github.com/manuelrojas19/go-oauth2-server/oauth/authmethodtype"
	"github.com/manuelrojas19/go-oauth2-server/oauth/granttype"
	"github.com/manuelrojas19/go-oauth2-server/oauth/responsetype"
)

type RegisterClientResponse struct {
	ClientId                string                                 `json:"client_id"`
	ClientSecret            string                                 `json:"client_secret,omitempty"`
	ClientName              string                                 `json:"client_name"`
	GrantTypes              []granttype.GrantType                  `json:"grant_type"`
	ResponseTypes           []responsetype.ResponseType            `json:"response_type"`
	TokenEndpointAuthMethod authmethodtype.TokenEndpointAuthMethod `json:"token_endpoint_auth_method"`
	RedirectUris            []string                               `json:"redirect_uris"`
	Scopes                  []oauth.Scope                          `json:"scopes"`
}
