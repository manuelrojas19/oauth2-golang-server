package configuration

import (
	"github.com/manuelrojas19/go-oauth2-server/store"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

var (
	DatabaseUrl string
)

func LoadDbSecrets() {
	DatabaseUrl = os.Getenv("DATABASE_URL")
}

func InitDatabaseConnection() (*gorm.DB, error) {
	datasource, err := gorm.Open(postgres.Open(DatabaseUrl), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	err = datasource.AutoMigrate(
		&store.OauthClient{},
		&store.AccessToken{},
		&store.RefreshToken{},
		&store.User{},
		&store.AuthCode{},
	)

	if err != nil {
		return nil, err
	}

	return datasource, nil
}
