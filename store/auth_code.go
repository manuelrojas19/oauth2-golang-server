package store

import (
	"github.com/google/uuid"
	"time"
)

type AuthCode struct {
	Id          string    `gorm:"primaryKey;type:varchar(255);unique;not null"`
	Code        string    `gorm:"type:text;unique;not null"`
	RedirectURI string    `gorm:"type:varchar(255);not null"`
	Scope       string    `gorm:"size:255;not null"`
	Used        bool      `gorm:"not null;default:false"`
	UserId      string    `gorm:"index;not null"`
	ClientId    string    `gorm:"index;not null"`
	ExpiresAt   time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	User        *User
	Client      *OauthClient
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

func (b *AuthCodeBuilder) WithClientId(clientId string) *AuthCodeBuilder {
	b.authorizationCode.ClientId = clientId
	return b
}

func (b *AuthCodeBuilder) WithRedirectURI(redirectURI string) *AuthCodeBuilder {
	b.authorizationCode.RedirectURI = redirectURI
	return b
}

func (b *AuthCodeBuilder) WithScope(scope string) *AuthCodeBuilder {
	b.authorizationCode.Scope = scope
	return b
}

func (b *AuthCodeBuilder) WithUser(user *User) *AuthCodeBuilder {
	b.authorizationCode.User = user
	b.authorizationCode.UserId = user.Id // Automatically set UserId based on the provided user
	return b
}

func (b *AuthCodeBuilder) WithUserId(userId string) *AuthCodeBuilder {
	b.authorizationCode.UserId = userId
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
