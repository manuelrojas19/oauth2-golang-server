package repositories

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/manuelrojas19/go-oauth2-server/store"

	"gorm.io/gorm"
)

type accessTokenRepository struct {
	Db *gorm.DB
}

func NewAccessTokenRepository(db *gorm.DB) AccessTokenRepository {
	return &accessTokenRepository{Db: db}
}

func (ot *accessTokenRepository) Save(token *store.AccessToken) (*store.AccessToken, error) {
	log.Printf("Starting transaction to save token for client_key %s", token.ClientId)

	tx := ot.Db.Begin()

	expiredTokensQuery := tx.Unscoped().
		Where("client_id = ?", token.ClientId).
		Where("expires_at <= ?", time.Now())

	if err := expiredTokensQuery.Delete(new(store.AccessToken)).Error; err != nil {
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

func (ot *accessTokenRepository) FindByAccessToken(accessToken string) (*store.AccessToken, error) {
	var token store.AccessToken

	if err := ot.Db.Where("access_token = ?", accessToken).First(&token).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("access token not found")
		}
		return nil, fmt.Errorf("failed to find access token: %w", err)
	}

	if token.ExpiresAt.Before(time.Now()) {
		return nil, fmt.Errorf("access token expired")
	}

	return &token, nil
}

func (ot *accessTokenRepository) DeleteByAccessToken(accessToken string) error {
	result := ot.Db.Where("access_token = ?", accessToken).Delete(&store.AccessToken{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete access token: %w", result.Error)
	}
	// If no rows were affected, it means the token was not found, but we don't return an error as per RFC 7009
	return nil
}
