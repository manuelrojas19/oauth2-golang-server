package entities

import (
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/manuelrojas19/go-oauth2-server/models/oauth"
	"time"
)

type OauthClient struct {
	ClientId                string         `gorm:"primaryKey;type:varchar(255)"`
	ClientSecret            string         `gorm:"type:varchar(255);not null"`
	ClientName              string         `gorm:"type:varchar(255);"`
	ResponseTypes           pq.StringArray `gorm:"type:text[]"`
	GrantTypes              pq.StringArray `gorm:"type:text[]"`
	TokenEndpointAuthMethod string         `gorm:"type:varchar(255);"`
	RedirectURI             pq.StringArray `gorm:"type:text[]"`
	CreatedAt               time.Time
	UpdatedAt               time.Time
}

// OauthClientBuilder helps build an OauthClient with optional configurations.
type OauthClientBuilder struct {
	clientID                string
	clientSecret            string
	clientName              string
	responseTypes           []oauth.ResponseType
	grantTypes              []oauth.GrantType
	tokenEndpointAuthMethod oauth.TokenEndpointAuthMethod
	redirectURI             []string
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
func (b *OauthClientBuilder) SetResponseTypes(responseTypes []oauth.ResponseType) *OauthClientBuilder {
	b.responseTypes = responseTypes
	return b
}

// SetGrantTypes sets the grant types.
func (b *OauthClientBuilder) SetGrantTypes(grantTypes []oauth.GrantType) *OauthClientBuilder {
	b.grantTypes = grantTypes
	return b
}

// SetTokenEndpointAuthMethod sets the token endpoint auth method.
func (b *OauthClientBuilder) SetTokenEndpointAuthMethod(authMethod oauth.TokenEndpointAuthMethod) *OauthClientBuilder {
	b.tokenEndpointAuthMethod = authMethod
	return b
}

// SetRedirectURI sets the redirect URI.
func (b *OauthClientBuilder) SetRedirectURI(redirectURI []string) *OauthClientBuilder {
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
		ResponseTypes:           oauth.ResponseTypeListToStringList(b.responseTypes),
		GrantTypes:              oauth.GrantTypeListToStringList(b.grantTypes),
		TokenEndpointAuthMethod: string(b.tokenEndpointAuthMethod),
		RedirectURI:             b.redirectURI,
		CreatedAt:               time.Now().UTC(),
		UpdatedAt:               time.Now().UTC(),
	}
}
