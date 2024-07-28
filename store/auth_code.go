package store

import (
	"github.com/google/uuid"
	"time"
)

type AuthCode struct {
	Id          string `gorm:"primaryKey;type:varchar(255);unique;not null"`
	Code        string `gorm:"type:text;unique;not null"`
	Client      *OauthClient
	ClientId    string    `gorm:"type:varchar(255);not null"`
	RedirectURI string    `gorm:"type:varchar(255);not null"`
	ExpiresAt   time.Time `gorm:"default:now()"`
	CreatedAt   time.Time `gorm:"default:now()"`
}

// AuthCodeBuilder helps in constructing AuthCode instances
type AuthCodeBuilder struct {
	authorizationCode AuthCode
}

func NewAuthorizationCodeBuilder() *AuthCodeBuilder {
	return &AuthCodeBuilder{
		authorizationCode: AuthCode{
			CreatedAt: time.Now(),
		},
	}
}

func (b *AuthCodeBuilder) WithCode(code string) *AuthCodeBuilder {
	b.authorizationCode.Code = code
	return b
}

func (b *AuthCodeBuilder) WithClient(client *OauthClient) *AuthCodeBuilder {
	b.authorizationCode.Client = client
	return b
}

func (b *AuthCodeBuilder) WithClientID(clientID string) *AuthCodeBuilder {
	b.authorizationCode.ClientId = clientID
	return b
}

func (b *AuthCodeBuilder) WithRedirectURI(redirectURI string) *AuthCodeBuilder {
	b.authorizationCode.RedirectURI = redirectURI
	return b
}

func (b *AuthCodeBuilder) WithExpiresAt(expiresAt time.Time) *AuthCodeBuilder {
	b.authorizationCode.ExpiresAt = expiresAt
	return b
}

func (b *AuthCodeBuilder) WithCreatedAt(createdAt time.Time) *AuthCodeBuilder {
	b.authorizationCode.CreatedAt = createdAt
	return b
}

func (b *AuthCodeBuilder) Build() *AuthCode {
	b.authorizationCode.Id = uuid.New().String()
	return &b.authorizationCode
}
