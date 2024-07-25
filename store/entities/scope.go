package entities

import (
	"github.com/google/uuid"
	"time"
)

type Scope struct {
	Id          string    `gorm:"primaryKey;type:varchar(255);unique;not null"`
	Name        string    `gorm:"type:varchar(255);not null"`
	Description string    `gorm:"type:varchar(800);not null"`
	CreatedAt   time.Time `gorm:"default:now()"`
	UpdatedAt   time.Time `gorm:"default:now()"`
}

type ScopeBuilder struct {
	name        string
	description string
}

// NewScopeBuilder creates a new instance of ScopeBuilder.
func NewScopeBuilder() *ScopeBuilder {
	return &ScopeBuilder{}
}

// WithName sets the Name field in the builder.
func (s *ScopeBuilder) WithName(name string) *ScopeBuilder {
	s.name = name
	return s
}

// WithDescription sets the description field in the builder.
func (s *ScopeBuilder) WithDescription(description string) *ScopeBuilder {
	s.description = description
	return s
}

// Build creates a new ClientScope instance using the builder's settings.
func (b *ScopeBuilder) Build() *Scope {
	return &Scope{
		Id:          uuid.New().String(),
		Name:        b.name,
		Description: b.description,
	}
}
