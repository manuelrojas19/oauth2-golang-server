package oauth

import (
	"net/url"
	"time"
)

type Token struct {
	ClientId              *string
	UserId                *string
	RedirectURI           string
	Scope                 string
	TokenType             string
	AccessToken           string
	AccessTokenCreatedAt  time.Time
	AccessTokenExpiresIn  int
	AccessTokenExpiresAt  time.Time
	RefreshToken          string
	RefreshTokenCreatedAt time.Time
	RefreshTokenExpiresAt time.Time
	Extension             url.Values
}

type TokenBuilder struct {
	clientId              *string
	userId                *string
	redirectURI           string
	scope                 string
	tokenType             string
	accessToken           string
	accessTokenCreatedAt  time.Time
	accessTokenExpiresIn  int
	accessTokenExpiresAt  time.Time
	refreshToken          string
	refreshTokenCreatedAt time.Time
	refreshTokenExpiresAt time.Time
	extension             url.Values
}

func NewTokenBuilder() *TokenBuilder {
	return &TokenBuilder{}
}

func (b *TokenBuilder) WithClientId(clientId *string) *TokenBuilder {
	b.clientId = clientId
	return b
}

func (b *TokenBuilder) WithUserId(userId *string) *TokenBuilder {
	b.userId = userId
	return b
}

func (b *TokenBuilder) WithRedirectURI(redirectURI string) *TokenBuilder {
	b.redirectURI = redirectURI
	return b
}

func (b *TokenBuilder) WithScope(scope string) *TokenBuilder {
	b.scope = scope
	return b
}

func (b *TokenBuilder) WithTokenType(tokenType string) *TokenBuilder {
	b.tokenType = tokenType
	return b
}

func (b *TokenBuilder) WithAccessToken(accessToken string) *TokenBuilder {
	b.accessToken = accessToken
	return b
}

func (b *TokenBuilder) WithAccessTokenCreatedAt(accessTokenCreatedAt time.Time) *TokenBuilder {
	b.accessTokenCreatedAt = accessTokenCreatedAt
	return b
}

func (b *TokenBuilder) WithAccessTokenExpiresIn(accessTokenExpiresIn int) *TokenBuilder {
	b.accessTokenExpiresIn = accessTokenExpiresIn
	return b
}

func (b *TokenBuilder) WithAccessTokenExpiresAt(accessTokenExpiresAt time.Time) *TokenBuilder {
	b.accessTokenExpiresAt = accessTokenExpiresAt
	return b
}

func (b *TokenBuilder) WithRefreshToken(refreshToken string) *TokenBuilder {
	b.refreshToken = refreshToken
	return b
}

func (b *TokenBuilder) WithRefreshTokenCreatedAt(refreshTokenCreatedAt time.Time) *TokenBuilder {
	b.refreshTokenCreatedAt = refreshTokenCreatedAt
	return b
}

func (b *TokenBuilder) WithRefreshTokenExpiresAt(refreshTokenExpiresAt time.Time) *TokenBuilder {
	b.refreshTokenExpiresAt = refreshTokenExpiresAt
	return b
}

func (b *TokenBuilder) WithExtension(extension url.Values) *TokenBuilder {
	b.extension = extension
	return b
}

func (b *TokenBuilder) Build() *Token {
	return &Token{
		ClientId:              b.clientId,
		UserId:                b.userId,
		RedirectURI:           b.redirectURI,
		Scope:                 b.scope,
		TokenType:             b.tokenType,
		AccessToken:           b.accessToken,
		AccessTokenCreatedAt:  b.accessTokenCreatedAt,
		AccessTokenExpiresIn:  b.accessTokenExpiresIn,
		AccessTokenExpiresAt:  b.accessTokenExpiresAt,
		RefreshToken:          b.refreshToken,
		RefreshTokenCreatedAt: b.refreshTokenCreatedAt,
		RefreshTokenExpiresAt: b.refreshTokenExpiresAt,
		Extension:             b.extension,
	}
}
