package repositories

import (
	"errors"
	"github.com/manuelrojas19/go-oauth2-server/store/entities"
	"gorm.io/gorm"
	"log"
	"time"
)

type refreshTokenRepository struct {
	Db *gorm.DB
}

func NewRefreshTokenRepository(db *gorm.DB) RefreshTokenRepository {
	return &refreshTokenRepository{Db: db}
}

func (ot *refreshTokenRepository) Save(token *entities.RefreshToken) (*entities.RefreshToken, error) {
	log.Printf("Starting transaction to save refresh token for access_token_id %s", token.AccessTokenId)

	tx := ot.Db.Begin()

	expiredTokensQuery := tx.Unscoped().Where("access_token_id = ?", token.AccessTokenId).Where("expires_at <= ?", time.Now())

	if err := expiredTokensQuery.Delete(new(entities.RefreshToken)).Error; err != nil {
		log.Printf("Error deleting expired refresh tokens for access_token_id %s: %v", token.AccessTokenId, err)
		tx.Rollback()
		return nil, err
	}

	if err := tx.Create(token).Error; err != nil {
		log.Printf("Error creating refresh token for access_token_id %s: %v", token.AccessTokenId, err)
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("Error committing transaction for access_token_id %s: %v", token.AccessTokenId, err)
		tx.Rollback()
		return nil, err
	}

	log.Printf("Successfully saved refresh token for access_token_id %s", token.AccessTokenId)
	return token, nil
}

func (ot *refreshTokenRepository) FindByToken(token string) (*entities.RefreshToken, error) {
	refreshToken := new(entities.RefreshToken)
	result := ot.Db.Where("token = ?", token).First(refreshToken)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, errors.New("OAuth Client not Found")
	}

	return refreshToken, nil
}
