package entities

import (
	"time"
)

type AuthorizationCode struct {
	Id          string `gorm:"primaryKey;type:varchar(255);unique;not null"`
	Code        string `gorm:"type:text;unique;not null"`
	Client      *OauthClient
	ClientID    string    `gorm:"type:varchar(255);not null"`
	RedirectURI string    `gorm:"type:varchar(255);not null"`
	ExpiresAt   time.Time `gorm:"default:now()"`
	CreatedAt   time.Time `gorm:"default:now()"`
}
