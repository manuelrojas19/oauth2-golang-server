package oauth

import "time"

type AuthCode struct {
	Code        string
	ClientId    string
	RedirectURI string
	Scope       string
	ExpiresAt   time.Time
	CreatedAt   time.Time
}

// AuthCodeBuilder helps in constructing AuthCode instances
type AuthCodeBuilder struct {
	authCode AuthCode
}

func NewAuthCodeBuilder() *AuthCodeBuilder {
	return &AuthCodeBuilder{
		authCode: AuthCode{
			CreatedAt: time.Now(),
		},
	}
}

func (b *AuthCodeBuilder) WithCode(code string) *AuthCodeBuilder {
	b.authCode.Code = code
	return b
}

func (b *AuthCodeBuilder) WithClientId(clientId string) *AuthCodeBuilder {
	b.authCode.ClientId = clientId
	return b
}

func (b *AuthCodeBuilder) WithRedirectURI(redirectURI string) *AuthCodeBuilder {
	b.authCode.RedirectURI = redirectURI
	return b
}

func (b *AuthCodeBuilder) WithExpiresAt(expiresAt time.Time) *AuthCodeBuilder {
	b.authCode.ExpiresAt = expiresAt
	return b
}

func (b *AuthCodeBuilder) WithCreatedAt(createdAt time.Time) *AuthCodeBuilder {
	b.authCode.CreatedAt = createdAt
	return b
}

func (b *AuthCodeBuilder) WithScope(scope string) *AuthCodeBuilder {
	b.authCode.Scope = scope
	return b
}

func (b *AuthCodeBuilder) Build() *AuthCode {
	return &b.authCode
}

// Equal compares two AuthCode instances for equality, ignoring ext field of time.
func (a *AuthCode) Equal(b *AuthCode) bool {
	if a == b {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return a.Code == b.Code &&
		a.ClientId == b.ClientId &&
		a.RedirectURI == b.RedirectURI &&
		a.Scope == b.Scope
}
