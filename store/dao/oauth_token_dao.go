package dao

import (
	"log"
	"time"

	"github.com/manuelrojas19/go-oauth2-server/store/entities"
	"gorm.io/gorm"
)

type OauthTokenDao struct {
	Db *gorm.DB
}

func (ot *OauthTokenDao) SaveToken(token *entities.OauthTokenEntity) (*entities.OauthTokenEntity, error) {
	tx := ot.Db.Begin()

	expiredTokensQuery := tx.Unscoped().Where("client_key = ?", token.ClientKey).Where("expires_at <= ?", time.Now())

	if err := expiredTokensQuery.Delete(new(entities.OauthTokenEntity)).Error; err != nil {
		log.Println("Error deleting")
		tx.Rollback()
		return nil, err
	}

	if err := tx.Create(token).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	return token, nil
}
