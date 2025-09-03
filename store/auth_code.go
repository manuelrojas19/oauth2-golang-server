package store

import (
	"time"

	"github.com/google/uuid"
)

type AuthCode struct {
	Id                  string    `gorm:"primaryKey;type:varchar(255);unique;not null"`
	Code                string    `gorm:"type:text;unique;not null"`
	RedirectURI         string    `gorm:"type:varchar(255);not null"`
	Used                bool      `gorm:"not null;default:false"`
	UserId              *string   `gorm:"index"`
	ClientId            *string   `gorm:"index"`
	ExpiresAt           time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	CreatedAt           time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	CodeChallenge       string    `gorm:"type:varchar(255)"`
	CodeChallengeMethod string    `gorm:"type:varchar(255)"`
	User                *User
	Client              *OauthClient
	Scopes              []Scope `gorm:"many2many:auth_code_scopes;"`
}

// AuthCodeBuilder helps in constructing AuthCode instances
type AuthCodeBuilder struct {
	authorizationCode AuthCode
	scopes            []Scope
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

func (b *AuthCodeBuilder) WithClientId(clientId *string) *AuthCodeBuilder {
	b.authorizationCode.ClientId = clientId
	return b
}

func (b *AuthCodeBuilder) WithRedirectURI(redirectURI string) *AuthCodeBuilder {
	b.authorizationCode.RedirectURI = redirectURI
	return b
}

func (b *AuthCodeBuilder) WithScopes(scopes []Scope) *AuthCodeBuilder {
	b.authorizationCode.Scopes = scopes
	return b
}

func (b *AuthCodeBuilder) WithUser(user *User) *AuthCodeBuilder {
	b.authorizationCode.User = user
	b.authorizationCode.UserId = &user.Id // Automatically set UserId based on the provided user
	return b
}

func (b *AuthCodeBuilder) WithUserId(userId *string) *AuthCodeBuilder {
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

func (b *AuthCodeBuilder) WithCodeChallenge(codeChallenge string) *AuthCodeBuilder {
	b.authorizationCode.CodeChallenge = codeChallenge
	return b
}

func (b *AuthCodeBuilder) WithCodeChallengeMethod(codeChallengeMethod string) *AuthCodeBuilder {
	b.authorizationCode.CodeChallengeMethod = codeChallengeMethod
	return b
}

func (b *AuthCodeBuilder) Build() *AuthCode {
	b.authorizationCode.Id = uuid.New().String()
	b.authorizationCode.Scopes = b.scopes
	return &b.authorizationCode
}
