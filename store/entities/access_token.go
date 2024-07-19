package entities

import (
	"github.com/google/uuid"
	"time"
)

type AccessToken struct {
	Id        string    `gorm:"primary_key"`
	Token     string    `sql:"type:varchar(40);unique;not null"`
	Scope     string    `sql:"type:varchar(200);not null"`
	ExpiresAt time.Time `sql:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Client    *OauthClient
	ClientId  string `sql:"index;not null"`
}

func NewAccessToken(client *OauthClient, token string, scope string, expiresAt time.Time) *AccessToken {
	return &AccessToken{
		Id:        uuid.New().String(),
		ClientId:  client.ClientId,
		Client:    client,
		ExpiresAt: expiresAt,
		Scope:     scope,
		Token:     token,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
