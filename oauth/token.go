package oauth

import (
	"net/url"
	"time"
)

type Token struct {
	ClientId              string
	UserId                string
	RedirectURI           string
	Scope                 string
	AccessToken           string
	AccessTokenCreatedAt  time.Time
	AccessTokenExpiresAt  time.Duration
	RefreshToken          string
	RefreshTokenCreatedAt time.Time
	RefreshTokenExpiresAt time.Duration
	Extension             url.Values
}

type TokenBuilder struct {
	clientId              string
	userId                string
	redirectURI           string
	scope                 string
	accessToken           string
	accessTokenCreatedAt  time.Time
	accessTokenExpiresAt  time.Duration
	refreshToken          string
	refreshTokenCreatedAt time.Time
	refreshTokenExpiresAt time.Duration
	extension             url.Values
}

func NewTokenBuilder() *TokenBuilder {
	return &TokenBuilder{}
}

func (b *TokenBuilder) WithClientId(clientId string) *TokenBuilder {
	b.clientId = clientId
	return b
}

func (b *TokenBuilder) WithUserId(userId string) *TokenBuilder {
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

func (b *TokenBuilder) WithAccessToken(accessToken string) *TokenBuilder {
	b.accessToken = accessToken
	return b
}

func (b *TokenBuilder) WithAccessTokenCreatedAt(accessTokenCreatedAt time.Time) *TokenBuilder {
	b.accessTokenCreatedAt = accessTokenCreatedAt
	return b
}

func (b *TokenBuilder) WithAccessTokenExpiresAt(accessTokenExpiresAt time.Duration) *TokenBuilder {
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

func (b *TokenBuilder) WithRefreshTokenExpiresAt(refreshTokenExpiresAt time.Duration) *TokenBuilder {
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
		AccessToken:           b.accessToken,
		AccessTokenCreatedAt:  b.accessTokenCreatedAt,
		AccessTokenExpiresAt:  b.accessTokenExpiresAt,
		RefreshToken:          b.refreshToken,
		RefreshTokenCreatedAt: b.refreshTokenCreatedAt,
		RefreshTokenExpiresAt: b.refreshTokenExpiresAt,
		Extension:             b.extension,
	}
}
