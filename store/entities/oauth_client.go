package entities

import (
	"github.com/google/uuid"
	"time"
)

type OauthClient struct {
	ClientId     string  `gorm:"primaryKey;type:varchar(255)"`
	ClientSecret string  `gorm:"type:varchar(255);not null"`
	RedirectURI  *string `gorm:"type:varchar(255)"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func NewOauthClient(clientSecret string, redirectUri *string) *OauthClient {
	return &OauthClient{
		ClientId:     uuid.New().String(),
		ClientSecret: clientSecret,
		RedirectURI:  redirectUri,
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}
}
