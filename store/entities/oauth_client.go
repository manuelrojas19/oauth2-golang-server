package entities

import (
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/manuelrojas19/go-oauth2-server/models/oauth/authmethodtype"
	"github.com/manuelrojas19/go-oauth2-server/models/oauth/granttype"
	"github.com/manuelrojas19/go-oauth2-server/models/oauth/responsetype"
	"time"
)

type OauthClient struct {
	ClientId                string         `gorm:"primaryKey;type:varchar(255);unique;not null"`
	ClientSecret            string         `gorm:"type:varchar(255);unique;not null"`
	ClientName              string         `gorm:"type:varchar(255);unique;not null"`
	ResponseTypes           pq.StringArray `gorm:"type:text[];not null"`
	GrantTypes              pq.StringArray `gorm:"type:text[];not null"`
	TokenEndpointAuthMethod string         `gorm:"type:varchar(255);not null"`
	RedirectURI             pq.StringArray `gorm:"type:text[]"`
	CreatedAt               time.Time      `gorm:"default:now()"`
	UpdatedAt               time.Time      `gorm:"default:now()"`
}

// OauthClientBuilder helps build an OauthClient with optional configurations.
type OauthClientBuilder struct {
	clientID                string
	clientSecret            string
	clientName              string
	responseTypes           []responsetype.ResponseType
	grantTypes              []granttype.GrantType
	tokenEndpointAuthMethod authmethodtype.TokenEndpointAuthMethod
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

// SetResponseTypes sets the responsetype types.
func (b *OauthClientBuilder) SetResponseTypes(responseTypes []responsetype.ResponseType) *OauthClientBuilder {
	b.responseTypes = responseTypes
	return b
}

// SetGrantTypes sets the granttype types.
func (b *OauthClientBuilder) SetGrantTypes(grantTypes []granttype.GrantType) *OauthClientBuilder {
	b.grantTypes = grantTypes
	return b
}

// SetTokenEndpointAuthMethod sets the token endpoint authmethodtype method.
func (b *OauthClientBuilder) SetTokenEndpointAuthMethod(authMethod authmethodtype.TokenEndpointAuthMethod) *OauthClientBuilder {
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
		ResponseTypes:           responsetype.EnumListToStringList(b.responseTypes),
		GrantTypes:              granttype.EnumListToStringList(b.grantTypes),
		TokenEndpointAuthMethod: string(b.tokenEndpointAuthMethod),
		RedirectURI:             b.redirectURI,
		CreatedAt:               time.Now().UTC(),
		UpdatedAt:               time.Now().UTC(),
	}
}
