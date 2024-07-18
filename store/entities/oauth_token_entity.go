package entities

import (
	"time"
)

type OauthTokenEntity struct {
	BaseGormEntity
	Client    *OauthClientEntity
	ClientId  string    `sql:"index;not null"`
	ClientKey string    `sql:"not null"`
	Token     string    `sql:"type:varchar(40);unique;not null"`
	ExpiresAt time.Time `sql:"not null"`
	Scope     string    `sql:"type:varchar(200);not null"`
}
