package entities

import (
	"github.com/google/uuid"
	"time"
)

type OauthClient struct {
	ClientId                string  `gorm:"primaryKey;type:varchar(255)"`
	ClientSecret            string  `gorm:"type:varchar(255);not null"`
	ClientName              string  `gorm:"type:varchar(255);"`
	ResponseTypes           string  `gorm:"type:varchar(255);"`
	GrantTypes              string  `gorm:"type:varchar(255);"`
	TokenEndpointAuthMethod string  `gorm:"type:varchar(255);"`
	RedirectURI             *string `gorm:"type:varchar(255)"`
	CreatedAt               time.Time
	UpdatedAt               time.Time
}

// OauthClientBuilder helps build an OauthClient with optional configurations.
type OauthClientBuilder struct {
	clientID                string
	clientSecret            string
	clientName              string
	responseTypes           string
	grantTypes              string
	tokenEndpointAuthMethod string
	redirectURI             *string
}

// NewOauthClientBuilder initializes a new OauthClientBuilder.
func NewOauthClientBuilder() *OauthClientBuilder {
	return &OauthClientBuilder{}
}

// SetClientID sets the client ID.
func (b *OauthClientBuilder) SetClientID(clientID string) *OauthClientBuilder {
	b.clientID = clientID
	return b
}

// SetClientSecret sets the client secret.
func (b *OauthClientBuilder) SetClientSecret(clientSecret string) *OauthClientBuilder {
	b.clientSecret = clientSecret
	return b
}

// SetClientName sets the client name.
func (b *OauthClientBuilder) SetClientName(clientName string) *OauthClientBuilder {
	b.clientName = clientName
	return b
}

// SetResponseTypes sets the response types.
func (b *OauthClientBuilder) SetResponseTypes(responseTypes string) *OauthClientBuilder {
	b.responseTypes = responseTypes
	return b
}

// SetGrantTypes sets the grant types.
func (b *OauthClientBuilder) SetGrantTypes(grantTypes string) *OauthClientBuilder {
	b.grantTypes = grantTypes
	return b
}

// SetTokenEndpointAuthMethod sets the token endpoint auth method.
func (b *OauthClientBuilder) SetTokenEndpointAuthMethod(authMethod string) *OauthClientBuilder {
	b.tokenEndpointAuthMethod = authMethod
	return b
}

// SetRedirectURI sets the redirect URI.
func (b *OauthClientBuilder) SetRedirectURI(redirectURI *string) *OauthClientBuilder {
	b.redirectURI = redirectURI
	return b
}

// Build constructs the OauthClient object.
func (b *OauthClientBuilder) Build() *OauthClient {
	if b.clientID == "" {
		// Generate client ID if not set
		b.clientID = uuid.New().String()
	}

	return &OauthClient{
		ClientId:                b.clientID,
		ClientSecret:            b.clientSecret,
		ClientName:              b.clientName,
		ResponseTypes:           b.responseTypes,
		GrantTypes:              b.grantTypes,
		TokenEndpointAuthMethod: b.tokenEndpointAuthMethod,
		RedirectURI:             b.redirectURI,
		CreatedAt:               time.Now().UTC(),
		UpdatedAt:               time.Now().UTC(),
	}
}
