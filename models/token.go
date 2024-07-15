package models

import (
	"net/url"
	"time"
)

type Token struct {
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

func NewToken() *Token {
	return &Token{Extension: make(url.Values)}
}
