package oauth

import (
	"net/url"
	"time"
)

type AccessToken struct {
	ClientId         string
	UserId           string
	RedirectURI      string
	Scope            string
	Access           string
	AccessCreatedAt  time.Time
	AccessExpiresAt  time.Duration
	Refresh          string
	RefreshCreatedAt time.Time
	RefreshExpiresAt time.Duration
	Extension        url.Values
}

func NewToken() *AccessToken {
	return &AccessToken{Extension: make(url.Values)}
}
