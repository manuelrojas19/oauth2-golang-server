package entities

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/manuelrojas19/go-oauth2-server/models/oauth/authmethodtype"
	"github.com/manuelrojas19/go-oauth2-server/models/oauth/granttype"
	"github.com/manuelrojas19/go-oauth2-server/models/oauth/responsetype"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
)

type OauthClient struct {
	ClientId                string         `gorm:"primaryKey;type:varchar(255);unique;not null"`
	ClientSecret            string         `gorm:"type:text;not null"`
	ClientName              string         `gorm:"type:varchar(255);unique;not null"`
	ResponseTypes           pq.StringArray `gorm:"type:text[];not null"`
	GrantTypes              pq.StringArray `gorm:"type:text[];not null"`
	TokenEndpointAuthMethod string         `gorm:"type:varchar(255);not null"`
	RedirectURIs            pq.StringArray `gorm:"type:text[]"`
	CreatedAt               time.Time      `gorm:"default:now()"`
	UpdatedAt               time.Time      `gorm:"default:now()"`
	IsConfidential          bool
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
}

// NewOauthClientBuilder initializes a new OauthClientBuilder.
func NewOauthClientBuilder() *OauthClientBuilder {
	return &OauthClientBuilder{}
}

// WithClientID sets the client Id.
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

// WithResponseTypes sets the responsetype types.
func (b *OauthClientBuilder) WithResponseTypes(responseTypes []responsetype.ResponseType) *OauthClientBuilder {
	b.responseTypes = responseTypes
	return b
}

// WithGrantTypes sets the granttype types.
func (b *OauthClientBuilder) WithGrantTypes(grantTypes []granttype.GrantType) *OauthClientBuilder {
	b.grantTypes = grantTypes
	return b
}

// WithTokenEndpointAuthMethod sets the token endpoint authmethodtype method.
func (b *OauthClientBuilder) WithTokenEndpointAuthMethod(authMethod authmethodtype.TokenEndpointAuthMethod) *OauthClientBuilder {
	b.tokenEndpointAuthMethod = authMethod
	return b
}

// WithRedirectURI sets the redirect URI.
func (b *OauthClientBuilder) WithRedirectURI(redirectURI []string) *OauthClientBuilder {
	b.redirectURI = redirectURI
	return b
}

// Build constructs the OauthClient object.
func (b *OauthClientBuilder) Build() *OauthClient {
	if b.clientID == "" {
		// Generate client Id if not set
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
	}
}
