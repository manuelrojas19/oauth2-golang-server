package store

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	Id        string          `gorm:"primaryKey;type:varchar(255);unique;not null"`
	Name      string          `gorm:"type:varchar(255);not null"`
	Email     string          `gorm:"type:varchar(255);unique"`
	IdpName   string          `gorm:"type:varchar(255);not null"`
	CreatedAt time.Time       `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time       `gorm:"default:CURRENT_TIMESTAMP"`
	Consents  []AccessConsent `gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE"`
}

// UserBuilder helps in constructing User instances with optional configurations.
type UserBuilder struct {
	id      string
	name    string
	email   string
	idpName string
}

// NewUserBuilder initializes a new UserBuilder.
func NewUserBuilder() *UserBuilder {
	return &UserBuilder{}
}

// WithID sets the ScopeId field in the builder.
func (b *UserBuilder) WithID(id string) *UserBuilder {
	b.id = id
	return b
}

// WithName sets the Name field in the builder.
func (b *UserBuilder) WithName(name string) *UserBuilder {
	b.name = name
	return b
}

// WithEmail sets the Email field in the builder.
func (b *UserBuilder) WithEmail(email string) *UserBuilder {
	b.email = email
	return b
}

// WithIdpName sets the IdpName field in the builder.
func (b *UserBuilder) WithIdpName(idpName string) *UserBuilder {
	b.idpName = idpName
	return b
}

// Build creates a new User instance using the builder's settings.
func (b *UserBuilder) Build() *User {
	if b.id == "" {
		b.id = uuid.New().String() // Generate a new UUID if Id is not provided
	}

	return &User{
		Id:        b.id,
		Name:      b.name,
		Email:     b.email,
		IdpName:   b.idpName,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
}
