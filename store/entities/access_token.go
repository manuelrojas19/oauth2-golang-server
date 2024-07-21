package entities

import (
	"github.com/google/uuid"
	"time"
)

type AccessToken struct {
	Id        string    `gorm:"primaryKey;type:varchar(255);unique;not null"`
	Token     string    `gorm:"type:varchar(255);unique;not null"`
	Scope     string    `gorm:"type:varchar(255);not null"`
	ExpiresAt time.Time `gorm:"not null"`
	CreatedAt time.Time `gorm:"default:now()"`
	UpdatedAt time.Time `gorm:"default:now()"`
	Client    *OauthClient
	ClientId  string `gorm:"index;not null"`
}

func NewAccessToken(client *OauthClient, token string, scope string, expiresAt time.Time) *AccessToken {
	return &AccessToken{
		Id:        uuid.New().String(),
		ClientId:  client.ClientId,
		Client:    client,
		ExpiresAt: expiresAt,
		Scope:     scope,
		Token:     token,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

type AccessTokenBuilder struct {
	token     string
	scope     string
	expiresAt time.Time
	clientId  string
	client    *OauthClient
}

// NewAccessTokenBuilder initializes a new builder instance.
func NewAccessTokenBuilder() *AccessTokenBuilder {
	return &AccessTokenBuilder{}
}

// SetToken sets the token value.
func (b *AccessTokenBuilder) SetToken(token string) *AccessTokenBuilder {
	b.token = token
	return b
}

// SetScope sets the scope value.
func (b *AccessTokenBuilder) SetScope(scope string) *AccessTokenBuilder {
	b.scope = scope
	return b
}

// SetExpiresAt sets the expiration time.
func (b *AccessTokenBuilder) SetExpiresAt(expiresAt time.Time) *AccessTokenBuilder {
	b.expiresAt = expiresAt
	return b
}

// SetClientId sets the client ID value.
func (b *AccessTokenBuilder) SetClientId(clientId string) *AccessTokenBuilder {
	b.clientId = clientId
	return b
}

// SetClient sets the client reference.
func (b *AccessTokenBuilder) SetClient(client *OauthClient) *AccessTokenBuilder {
	b.client = client
	return b
}

// Build constructs an AccessToken instance.
func (b *AccessTokenBuilder) Build() *AccessToken {
	return &AccessToken{
		Token:     b.token,
		Scope:     b.scope,
		ExpiresAt: b.expiresAt,
		ClientId:  b.clientId,
		Client:    b.client,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
