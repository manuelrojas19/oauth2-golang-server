package database

import (
	"github.com/manuelrojas19/go-oauth2-server/store/entities"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDatabaseConnection() (*gorm.DB, error) {
	dbUrl := "postgres://postgres:postgres@localhost:5432/oauthDB"

	datasource, err := gorm.Open(postgres.Open(dbUrl), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	err = datasource.AutoMigrate(&entities.OauthClient{}, &entities.AccessToken{}, &entities.RefreshToken{})
	if err != nil {
		return nil, err
	}

	return datasource, nil
}
