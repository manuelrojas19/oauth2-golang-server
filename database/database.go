package database

import (
	"github.com/manuelrojas19/go-oauth2-server/store/entities"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDatabaseConnection() (*gorm.DB, error) {
	dbUrl := "postgres://postgres:postgres@localhost:5432/oauthDB"

	datasource, error := gorm.Open(postgres.Open(dbUrl), &gorm.Config{})

	datasource.AutoMigrate(&entities.OauthClientEntity{}, &entities.OauthTokenEntity{})

	if error != nil {
		return nil, error
	}

	return datasource, nil
}
