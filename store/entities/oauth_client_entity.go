package entities

type OauthClientEntity struct {
	BaseGormEntity
	Key         string `sql:"type:varchar(254);unique;not null"`
	Secret      string `sql:"type:varchar(60);not null"`
	RedirectURI string `sql:"type:varchar(200)"`
}
