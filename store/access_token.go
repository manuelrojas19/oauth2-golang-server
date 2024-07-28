package store

import (
	"github.com/google/uuid"
	"time"
)

type AccessToken struct {
	Id            string    `gorm:"primaryKey;type:varchar(255);unique;not null"`
	Token         string    `gorm:"type:text;unique;not null"`
	TokenType     string    `gorm:"type:varchar(255);not null"`
	Scope         string    `gorm:"type:varchar(255);not null"`
	ExpiresAt     time.Time `gorm:"not null"`
	CreatedAt     time.Time `gorm:"default:now()"`
	UpdatedAt     time.Time `gorm:"default:now()"`
	Client        *OauthClient
	ClientId      string         `gorm:"index;not null"`
	RefreshTokens []RefreshToken `gorm:"foreignKey:AccessTokenId;constraint:OnDelete:CASCADE"`
}

type AccessTokenBuilder struct {
	token     string
	tokenType string
	scope     string
	expiresAt time.Time
	clientId  string
	client    *OauthClient
}

// NewAccessTokenBuilder initializes a new builder instance.
func NewAccessTokenBuilder() *AccessTokenBuilder {
	return &AccessTokenBuilder{}
}

// WithToken sets the token value.
func (b *AccessTokenBuilder) WithToken(token string) *AccessTokenBuilder {
	b.token = token
	return b
}

// WithTokenType sets the token type value.
func (b *AccessTokenBuilder) WithTokenType(tokenType string) *AccessTokenBuilder {
	b.tokenType = tokenType
	return b
}

// WithScope sets the scope value.
func (b *AccessTokenBuilder) WithScope(scope string) *AccessTokenBuilder {
	b.scope = scope
	return b
}

// WithExpiresAt sets the expiration time.
func (b *AccessTokenBuilder) WithExpiresAt(expiresAt time.Time) *AccessTokenBuilder {
	b.expiresAt = expiresAt
	return b
}

// WithClientId sets the client Id value.
func (b *AccessTokenBuilder) WithClientId(clientId string) *AccessTokenBuilder {
	b.clientId = clientId
	return b
}

// WithClient sets the client reference.
func (b *AccessTokenBuilder) WithClient(client *OauthClient) *AccessTokenBuilder {
	b.client = client
	return b
}

// Build constructs an AccessToken instance.
func (b *AccessTokenBuilder) Build() *AccessToken {
	return &AccessToken{
		Id:        uuid.New().String(),
		Token:     b.token,
		TokenType: b.tokenType,
		Scope:     b.scope,
		ExpiresAt: b.expiresAt,
		ClientId:  b.clientId,
		Client:    b.client,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
