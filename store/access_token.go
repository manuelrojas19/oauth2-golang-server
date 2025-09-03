package store

import (
	"time"

	"github.com/google/uuid"
)

type AccessToken struct {
	Id            string    `gorm:"primaryKey;type:varchar(255);unique;not null"`
	Token         string    `gorm:"type:text;unique;not null"`
	TokenType     string    `gorm:"type:varchar(255);not null"`
	ExpiresAt     time.Time `gorm:"not null"`
	CreatedAt     time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	Code          string    `gorm:"type:text"` // Reference to authorization code
	UserId        *string   `gorm:"index"`
	ClientId      *string   `gorm:"index"`
	User          *User
	Client        *OauthClient
	RefreshTokens []RefreshToken `gorm:"foreignKey:AccessTokenId;constraint:OnDelete:CASCADE"`
	Scopes        []Scope        `gorm:"many2many:access_token_scopes;constraint:OnDelete:CASCADE"`
}

type AccessTokenBuilder struct {
	token     string
	tokenType string
	expiresAt time.Time
	clientId  *string
	client    *OauthClient
	userId    *string
	user      *User
	code      string
	scopes    []Scope
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
func (b *AccessTokenBuilder) WithScopes(scopes []Scope) *AccessTokenBuilder {
	b.scopes = scopes
	return b
}

// WithExpiresAt sets the expiration time.
func (b *AccessTokenBuilder) WithExpiresAt(expiresAt time.Time) *AccessTokenBuilder {
	b.expiresAt = expiresAt
	return b
}

// WithClientId sets the client ScopeId value.
func (b *AccessTokenBuilder) WithClientId(clientId *string) *AccessTokenBuilder {
	b.clientId = clientId
	return b
}

// WithClient sets the client reference.
func (b *AccessTokenBuilder) WithClient(client *OauthClient) *AccessTokenBuilder {
	b.client = client
	return b
}

// WithUserId sets the userId reference.
func (b *AccessTokenBuilder) WithUserId(userId *string) *AccessTokenBuilder {
	b.userId = userId
	return b
}

// WithUser sets the user reference.
func (b *AccessTokenBuilder) WithUser(user *User) *AccessTokenBuilder {
	b.user = user
	if user != nil {
		b.userId = &user.Id
	}
	return b
}

func (b *AccessTokenBuilder) WithCode(code string) *AccessTokenBuilder {
	b.code = code
	return b
}

// Build constructs an AccessToken instance.
func (b *AccessTokenBuilder) Build() *AccessToken {
	return &AccessToken{
		Id:        uuid.New().String(),
		Token:     b.token,
		TokenType: b.tokenType,
		ExpiresAt: b.expiresAt,
		ClientId:  b.clientId,
		Client:    b.client,
		UserId:    b.userId,
		User:      b.user,
		Code:      b.code,
		CreatedAt: time.Now(),
		Scopes:    b.scopes,
	}
}
