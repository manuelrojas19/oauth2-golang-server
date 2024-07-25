package entities

import "time"

type ClientScope struct {
	ClientId  string    `gorm:"index;not null"`
	ScopeId   string    `gorm:"index;not null"`
	CreatedAt time.Time `gorm:"default:now()"`
	UpdatedAt time.Time `gorm:"default:now()"`
}

type ClientScopeBuilder struct {
	clientId string
	scopeId  string
}

// NewClientScopeBuilder creates a new instance of ClientScopeBuilder.
func NewClientScopeBuilder() *ClientScopeBuilder {
	return &ClientScopeBuilder{}
}

// WithClientId sets the ClientId field in the builder.
func (b *ClientScopeBuilder) WithClientId(clientId string) *ClientScopeBuilder {
	b.clientId = clientId
	return b
}

// WithScopeId sets the ScopeId field in the builder.
func (b *ClientScopeBuilder) WithScopeId(scopeId string) *ClientScopeBuilder {
	b.scopeId = scopeId
	return b
}

// Build creates a new ClientScope instance using the builder's settings.
func (b *ClientScopeBuilder) Build() *ClientScope {
	return &ClientScope{
		ClientId: b.clientId,
		ScopeId:  b.scopeId,
	}
}
