package store

import "time"

// AccessConsent represents the user's consent for a client and scope.
type AccessConsent struct {
	Id         string    `gorm:"primaryKey;type:varchar(255);unique;not null"`
	UserId     string    `gorm:"index;not null"`
	ClientId   string    `gorm:"index;not null"`
	ResourceId string    `gorm:"index;not null"`
	ScopeId    string    `gorm:"index;not null"`
	Consented  bool      `gorm:"not null;default:false"`
	CreatedAt  time.Time `gorm:"default:now()"`
	UpdatedAt  time.Time `gorm:"default:now()"`

	Client   *OauthClient
	User     *User
	Resource *OauthResource
	Scope    *Scope
}

// ConsentBuilder helps in constructing AccessConsent instances with optional configurations.
type ConsentBuilder struct {
	userId    string
	clientId  string
	scopeId   string
	consented bool
}

// NewUserConsentBuilder initializes a new ConsentBuilder.
func NewUserConsentBuilder() *ConsentBuilder {
	return &ConsentBuilder{}
}

// WithUserId sets the UserId field in the builder.
func (b *ConsentBuilder) WithUserId(userId string) *ConsentBuilder {
	b.userId = userId
	return b
}

// WithClientId sets the ClientId field in the builder.
func (b *ConsentBuilder) WithClientId(clientId string) *ConsentBuilder {
	b.clientId = clientId
	return b
}

// WithScopeId sets the ScopeId field in the builder.
func (b *ConsentBuilder) WithScopeId(scopeId string) *ConsentBuilder {
	b.scopeId = scopeId
	return b
}

// WithConsented sets the Consented field in the builder.
func (b *ConsentBuilder) WithConsented(consented bool) *ConsentBuilder {
	b.consented = consented
	return b
}

// Build creates a new AccessConsent instance using the builder's settings.
func (b *ConsentBuilder) Build() *AccessConsent {
	return &AccessConsent{
		UserId:    b.userId,
		ClientId:  b.clientId,
		ScopeId:   b.scopeId,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
}
