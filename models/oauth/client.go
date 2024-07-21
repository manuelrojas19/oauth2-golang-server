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

type ClientBuilder struct {
	client Client
}

// NewClientBuilder initializes and returns a new ClientBuilder.
func NewClientBuilder() *ClientBuilder {
	return &ClientBuilder{}
}

// WithClientId sets the ClientId for the builder.
func (b *ClientBuilder) WithClientId(clientId string) *ClientBuilder {
	b.client.ClientId = clientId
	return b
}

// WithClientSecret sets the ClientSecret for the builder.
func (b *ClientBuilder) WithClientSecret(clientSecret string) *ClientBuilder {
	b.client.ClientSecret = clientSecret
	return b
}

// WithClientName sets the ClientName for the builder.
func (b *ClientBuilder) WithClientName(clientName string) *ClientBuilder {
	b.client.ClientName = clientName
	return b
}

// WithGrantTypes sets the GrantTypes for the builder.
func (b *ClientBuilder) WithGrantTypes(grantTypes []granttype.GrantType) *ClientBuilder {
	b.client.GrantTypes = grantTypes
	return b
}

// WithResponseTypes sets the ResponseTypes for the builder.
func (b *ClientBuilder) WithResponseTypes(responseTypes []responsetype.ResponseType) *ClientBuilder {
	b.client.ResponseTypes = responseTypes
	return b
}

// WithTokenEndpointAuthMethod sets the TokenEndpointAuthMethod for the builder.
func (b *ClientBuilder) WithTokenEndpointAuthMethod(authMethod authmethodtype.TokenEndpointAuthMethod) *ClientBuilder {
	b.client.TokenEndpointAuthMethod = authMethod
	return b
}

// WithRedirectUris sets the RedirectUris for the builder.
func (b *ClientBuilder) WithRedirectUris(redirectUris []string) *ClientBuilder {
	b.client.RedirectUris = redirectUris
	return b
}

// Build constructs and returns the Client instance.
func (b *ClientBuilder) Build() *Client {
	return &b.client
}
