package oauth

import (
	"github.com/manuelrojas19/go-oauth2-server/models/oauth/authmethodtype"
	"github.com/manuelrojas19/go-oauth2-server/models/oauth/granttype"
	"github.com/manuelrojas19/go-oauth2-server/models/oauth/responsetype"
)

type Client struct {
	ClientId                string
	ClientSecret            string
	ClientName              string
	GrantTypes              []granttype.GrantType
	ResponseTypes           []responsetype.ResponseType
	TokenEndpointAuthMethod authmethodtype.TokenEndpointAuthMethod
	RedirectUris            []string
}
