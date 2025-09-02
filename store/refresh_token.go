package store

import (
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	Id            string    `gorm:"primaryKey;type:varchar(255);unique;not null"`
	Token         string    `gorm:"type:text;unique;not null"`
	TokenType     string    `gorm:"type:varchar(255);not null"`
	ExpiresAt     time.Time `gorm:"not null"`
	CreatedAt     time.Time `gorm:"default:now()"`
	AccessTokenId string    `gorm:"index;not null;constraint:OnDelete:CASCADE"`
	ClientId      *string   `gorm:"index"`
	UserId        *string   `gorm:"index"`
	Scope         string    `gorm:"type:varchar(255);not null"`

	AccessToken *AccessToken
	Client      *OauthClient
	User        *User
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
	clientId      *string
	accessToken   *AccessToken
	accessTokenId string
	user          *User
	userId        *string
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
		b.clientId = &client.ClientId
	}
	return b
}

func (b *RefreshTokenBuilder) WithClientId(clientId *string) *RefreshTokenBuilder {
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

// WithUserId sets the userId reference.
func (b *RefreshTokenBuilder) WithUserId(userId *string) *RefreshTokenBuilder {
	b.userId = userId
	return b
}

// WithUser sets the user reference.
func (b *RefreshTokenBuilder) WithUser(user *User) *RefreshTokenBuilder {
	b.user = user
	if user != nil {
		b.userId = &user.Id
	}
	return b
}

func (b *RefreshTokenBuilder) Build() *RefreshToken {
	return &RefreshToken{
		Id:            uuid.New().String(),
		Token:         b.token,
		TokenType:     b.tokenType,
		ExpiresAt:     b.expiresAt,
		Client:        b.client,
		ClientId:      b.clientId,
		AccessToken:   b.accessToken,
		AccessTokenId: b.accessTokenId,
		User:          b.user,
		UserId:        b.userId,
		CreatedAt:     time.Now(),
	}
}
