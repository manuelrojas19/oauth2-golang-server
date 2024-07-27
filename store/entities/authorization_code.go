package entities

import (
	"github.com/google/uuid"
	"time"
)

type AuthorizationCode struct {
	Id          string `gorm:"primaryKey;type:varchar(255);unique;not null"`
	Code        string `gorm:"type:text;unique;not null"`
	Client      *OauthClient
	ClientId    string    `gorm:"type:varchar(255);not null"`
	RedirectURI string    `gorm:"type:varchar(255);not null"`
	ExpiresAt   time.Time `gorm:"default:now()"`
	CreatedAt   time.Time `gorm:"default:now()"`
}

// AuthorizationCodeBuilder helps in constructing AuthorizationCode instances
type AuthorizationCodeBuilder struct {
	authorizationCode AuthorizationCode
}

func NewAuthorizationCodeBuilder() *AuthorizationCodeBuilder {
	return &AuthorizationCodeBuilder{
		authorizationCode: AuthorizationCode{
			CreatedAt: time.Now(),
		},
	}
}

func (b *AuthorizationCodeBuilder) WithCode(code string) *AuthorizationCodeBuilder {
	b.authorizationCode.Code = code
	return b
}

func (b *AuthorizationCodeBuilder) WithClient(client *OauthClient) *AuthorizationCodeBuilder {
	b.authorizationCode.Client = client
	return b
}

func (b *AuthorizationCodeBuilder) WithClientID(clientID string) *AuthorizationCodeBuilder {
	b.authorizationCode.ClientId = clientID
	return b
}

func (b *AuthorizationCodeBuilder) WithRedirectURI(redirectURI string) *AuthorizationCodeBuilder {
	b.authorizationCode.RedirectURI = redirectURI
	return b
}

func (b *AuthorizationCodeBuilder) WithExpiresAt(expiresAt time.Time) *AuthorizationCodeBuilder {
	b.authorizationCode.ExpiresAt = expiresAt
	return b
}

func (b *AuthorizationCodeBuilder) WithCreatedAt(createdAt time.Time) *AuthorizationCodeBuilder {
	b.authorizationCode.CreatedAt = createdAt
	return b
}

func (b *AuthorizationCodeBuilder) Build() *AuthorizationCode {
	b.authorizationCode.Id = uuid.New().String()
	return &b.authorizationCode
}
