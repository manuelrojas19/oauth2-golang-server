package entities

import (
	"github.com/google/uuid"
	"github.com/manuelrojas19/go-oauth2-server/models/oauth"
	"time"
)

type OauthClient struct {
	ClientId                string  `gorm:"primaryKey;type:varchar(255)"`
	ClientSecret            string  `gorm:"type:varchar(255);not null"`
	ClientName              string  `gorm:"type:varchar(255);"`
	ResponseTypes           string  `gorm:"type:varchar(255);"`
	GrantTypes              string  `gorm:"type:varchar(255);"`
	TokenEndpointAuthMethod string  `gorm:"type:varchar(255);"`
	RedirectURI             *string `gorm:"type:varchar(255)"`
	CreatedAt               time.Time
	UpdatedAt               time.Time
}

func NewOauthClientEntityFromModel(client *oauth.Client) *OauthClient {
	return &OauthClient{
		ClientId:                uuid.New().String(),
		ClientSecret:            client.ClientSecret,
		ClientName:              client.ClientName,
		ResponseTypes:           client.ResponseTypes,
		GrantTypes:              client.GrantTypes,
		TokenEndpointAuthMethod: client.TokenEndpointAuthMethod,
		RedirectURI:             &client.RedirectUris,
		CreatedAt:               time.Now().UTC(),
		UpdatedAt:               time.Now().UTC(),
	}
}
