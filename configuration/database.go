package configuration

import (
	"github.com/manuelrojas19/go-oauth2-server/store"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func NewDatabaseConnection() (*gorm.DB, error) {
	datasource, err := gorm.Open(postgres.Open(DatabaseUrl), &gorm.Config{})

	if err != nil {
		return nil, err
	}
	err = datasource.AutoMigrate(
		&store.Scope{},
		&store.OauthResource{},
		&store.OauthClient{},
		&store.OauthResource{},
		&store.AccessToken{},
		&store.RefreshToken{},
		&store.User{},
		&store.AuthCode{},
		&store.AccessConsent{},
	)

	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return datasource, nil
}
