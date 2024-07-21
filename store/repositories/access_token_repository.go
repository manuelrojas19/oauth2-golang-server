package repositories

import (
	"log"
	"time"

	"github.com/manuelrojas19/go-oauth2-server/store/entities"
	"gorm.io/gorm"
)

type accessTokenRepository struct {
	Db *gorm.DB
}

func NewAccessTokenRepository(db *gorm.DB) AccessTokenRepository {
	return &accessTokenRepository{Db: db}
}

func (ot *accessTokenRepository) Save(token *entities.AccessToken) (*entities.AccessToken, error) {
	log.Printf("Starting transaction to save token for client_key %s", token.ClientId)

	tx := ot.Db.Begin()

	expiredTokensQuery := tx.Unscoped().Where("client_id = ?", token.ClientId).Where("expires_at <= ?", time.Now())

	if err := expiredTokensQuery.Delete(new(entities.AccessToken)).Error; err != nil {
		log.Printf("Error deleting expired tokens for client_key %s: %v", token.ClientId, err)
		tx.Rollback()
		return nil, err
	}

	if err := tx.Create(token).Error; err != nil {
		log.Printf("Error creating token for client_key %s: %v", token.ClientId, err)
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("Error committing transaction for client_key %s: %v", token.ClientId, err)
		tx.Rollback()
		return nil, err
	}

	log.Printf("Successfully saved token for client_key %s", token.ClientId)
	return token, nil
}
