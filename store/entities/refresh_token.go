package entities

import (
	"github.com/google/uuid"
	"time"
)

type RefreshToken struct {
	Id            string    `gorm:"primaryKey;type:varchar(255);unique;not null"`
	Token         string    `gorm:"type:varchar(255);unique;not null"`
	TokenType     string    `gorm:"type:varchar(255);not null"`
	ExpiresAt     time.Time `gorm:"not null"`
	CreatedAt     time.Time `gorm:"default:now()"`
	Client        *OauthClient
	ClientId      string `gorm:"index;not null"`
	AccessToken   *AccessToken
	AccessTokenId string `gorm:"index;not null"`
}

// IsExpired checks if the refresh token has expired
func (r *RefreshToken) IsExpired() bool {
	return time.Now().After(r.ExpiresAt)
}

type RefreshTokenBuilder struct {
	id            string
	token         string
	tokenType     string
	expiresAt     time.Time
	createdAt     time.Time
	client        *OauthClient
	clientId      string
	accessToken   *AccessToken
	accessTokenId string
}

func NewRefreshTokenBuilder() *RefreshTokenBuilder {
	return &RefreshTokenBuilder{}
}

func (b *RefreshTokenBuilder) WithToken(token string) *RefreshTokenBuilder {
	b.token = token
	return b
}

func (b *RefreshTokenBuilder) WithTokenType(tokenType string) *RefreshTokenBuilder {
	b.tokenType = tokenType
	return b
}

func (b *RefreshTokenBuilder) WithExpiresAt(expiresAt time.Time) *RefreshTokenBuilder {
	b.expiresAt = expiresAt
	return b
}

func (b *RefreshTokenBuilder) WithCreatedAt(createdAt time.Time) *RefreshTokenBuilder {
	b.createdAt = createdAt
	return b
}

func (b *RefreshTokenBuilder) WithClient(client *OauthClient) *RefreshTokenBuilder {
	b.client = client
	if client != nil {
		b.clientId = client.ClientId
	}
	return b
}

func (b *RefreshTokenBuilder) WithClientId(clientId string) *RefreshTokenBuilder {
	b.clientId = clientId
	return b
}

func (b *RefreshTokenBuilder) WithAccessToken(accessToken *AccessToken) *RefreshTokenBuilder {
	b.accessToken = accessToken
	if accessToken != nil {
		b.accessTokenId = accessToken.Id
	}
	return b
}

func (b *RefreshTokenBuilder) WithAccessTokenId(accessTokenId string) *RefreshTokenBuilder {
	b.accessTokenId = accessTokenId
	return b
}

func (b *RefreshTokenBuilder) Build() *RefreshToken {
	return &RefreshToken{
		Id:            uuid.New().String(),
		Token:         b.token,
		TokenType:     b.tokenType,
		ExpiresAt:     b.expiresAt,
		CreatedAt:     b.createdAt,
		Client:        b.client,
		ClientId:      b.clientId,
		AccessToken:   b.accessToken,
		AccessTokenId: b.accessTokenId,
	}
}
