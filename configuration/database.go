package configuration

import (
	"github.com/manuelrojas19/go-oauth2-server/store"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
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
		&store.Scope{},
		&store.OauthClient{},
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

	// Create composite unique index
	err = datasource.Exec(`
        CREATE UNIQUE INDEX IF NOT EXISTS idx_oauth_client_scope_unique
        ON oauth_client_scopes (client_id, scope_id);
    `).Error
	if err != nil {
		return nil, err
	}

	return datasource, nil
}
