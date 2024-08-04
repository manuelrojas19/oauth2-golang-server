package store

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/manuelrojas19/go-oauth2-server/oauth/authmethodtype"
	"github.com/manuelrojas19/go-oauth2-server/oauth/granttype"
	"github.com/manuelrojas19/go-oauth2-server/oauth/responsetype"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
)

type OauthClient struct {
	ClientId                string         `gorm:"primaryKey;type:varchar(255);unique;not null"`
	ClientName              string         `gorm:"type:varchar(255);unique;not null"`
	ClientSecret            string         `gorm:"type:text;not null"`
	ResponseTypes           pq.StringArray `gorm:"type:text[];not null"`
	GrantTypes              pq.StringArray `gorm:"type:text[];not null"`
	TokenEndpointAuthMethod string         `gorm:"type:varchar(255);not null"`
	RedirectURIs            pq.StringArray `gorm:"type:text[]"`
	CreatedAt               time.Time      `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt               time.Time      `gorm:"default:CURRENT_TIMESTAMP"`
	Confidential            bool
	Scopes                  []Scope `gorm:"many2many:oauth_client_scopes;foreignKey:ClientId;joinForeignKey:ClientId;References:Id;JoinReferences:ScopeId"`
}

// ValidateSecret compares a plaintext secret with a bcrypt hash and returns a boolean indicating whether they match.
func (c *OauthClient) ValidateSecret(providedPassword string) error {
	log.Println("Validating secret")
	err := bcrypt.CompareHashAndPassword([]byte(c.ClientSecret), []byte(providedPassword))
	if err != nil {
		log.Printf("Secret validation failed: %v", err)
		return fmt.Errorf("secret validation failed: %w", err)
	}
	log.Println("Secret validation succeeded")
	return nil
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
	confidential            bool
	scopes                  []Scope
}

// NewOauthClientBuilder initializes a new OauthClientBuilder.
func NewOauthClientBuilder() *OauthClientBuilder {
	return &OauthClientBuilder{}
}

// WithClientID sets the client ID.
func (b *OauthClientBuilder) WithClientID(clientID string) *OauthClientBuilder {
	b.clientID = clientID
	return b
}

// WithClientSecret sets the client secret.
func (b *OauthClientBuilder) WithClientSecret(clientSecret string) *OauthClientBuilder {
	b.clientSecret = clientSecret
	return b
}

// WithClientName sets the client name.
func (b *OauthClientBuilder) WithClientName(clientName string) *OauthClientBuilder {
	b.clientName = clientName
	return b
}

// WithResponseTypes sets the response types.
func (b *OauthClientBuilder) WithResponseTypes(responseTypes []responsetype.ResponseType) *OauthClientBuilder {
	b.responseTypes = responseTypes
	return b
}

// WithGrantTypes sets the grant types.
func (b *OauthClientBuilder) WithGrantTypes(grantTypes []granttype.GrantType) *OauthClientBuilder {
	b.grantTypes = grantTypes
	return b
}

// WithTokenEndpointAuthMethod sets the token endpoint authentication method.
func (b *OauthClientBuilder) WithTokenEndpointAuthMethod(authMethod authmethodtype.TokenEndpointAuthMethod) *OauthClientBuilder {
	b.tokenEndpointAuthMethod = authMethod
	return b
}

// WithRedirectURIs sets the redirect URI.
func (b *OauthClientBuilder) WithRedirectURIs(redirectURI []string) *OauthClientBuilder {
	b.redirectURI = redirectURI
	return b
}

// WithConfidential sets the confidentiality of the client.
func (b *OauthClientBuilder) WithConfidential(confidential bool) *OauthClientBuilder {
	b.confidential = confidential
	return b
}

// WithScopes sets the client scopes.
func (b *OauthClientBuilder) WithScopes(scopes []Scope) *OauthClientBuilder {
	b.scopes = scopes
	return b
}

// Build constructs the OauthClient object.
func (b *OauthClientBuilder) Build() *OauthClient {
	if b.clientID == "" {
		b.clientID = uuid.New().String()
	}

	return &OauthClient{
		ClientId:                b.clientID,
		ClientSecret:            b.clientSecret,
		ClientName:              b.clientName,
		ResponseTypes:           responsetype.EnumListToStringList(b.responseTypes),
		GrantTypes:              granttype.EnumListToStringList(b.grantTypes),
		TokenEndpointAuthMethod: string(b.tokenEndpointAuthMethod),
		RedirectURIs:            b.redirectURI,
		CreatedAt:               time.Now().UTC(),
		UpdatedAt:               time.Now().UTC(),
		Confidential:            b.confidential,
		Scopes:                  b.scopes,
	}
}
