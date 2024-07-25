package entities

import "time"

// UserConsent represents the user's consent for a client and scope
type UserConsent struct {
	UserId    string    `gorm:"index;not null"`
	ClientId  string    `gorm:"index;not null"`
	ScopeId   string    `gorm:"index;not null"`
	Consented bool      `gorm:"not null"`
	CreatedAt time.Time `gorm:"default:now()"`
	UpdatedAt time.Time `gorm:"default:now()"`
}
