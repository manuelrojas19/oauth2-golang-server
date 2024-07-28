package store

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"gorm.io/gorm"
)

type oauthClientRepository struct {
	Db *gorm.DB
}

func NewOauthClientRepository(db *gorm.DB) OauthClientRepository {
	return &oauthClientRepository{Db: db}
}

func (ocd *oauthClientRepository) Save(client *OauthClient) (*OauthClient, error) {
	// Begin a new transaction
	tx := ocd.Db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Printf("PANIC: Rolled back transaction for client Id '%s' due to: %v", client.ClientId, r)
		}
	}()

	if err := tx.Error; err != nil {
		log.Printf("ERROR: Failed to start transaction for client Id '%s': %v", client.ClientId, err)
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}

	// Check if client already exists
	if ocd.clientExists(client.ClientId) {
		log.Printf("ERROR: Client with Id '%s' already exists", client.ClientId)
		tx.Rollback()
		return nil, errors.New("client already exists")
	}

	// Create new client
	if err := tx.Create(client).Error; err != nil {
		log.Printf("ERROR: Failed to create client with Id '%s': %v", client.ClientId, err)
		tx.Rollback()
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		log.Printf("ERROR: Error committing transaction for client Id '%s': %v", client.ClientId, err)
		tx.Rollback()
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	log.Printf("INFO: Successfully saved new client with Id '%s'", client.ClientId)
	return client, nil
}

func (ocd *oauthClientRepository) FindByClientId(clientId string) (*OauthClient, error) {
	oauthClient := new(OauthClient)
	result := ocd.Db.Where("LOWER(client_id) = ?", strings.ToLower(clientId)).First(oauthClient)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("OAuth client with Id '%s' not found", clientId)
		}
		return nil, fmt.Errorf("error finding OAuth client with Id '%s': %w", clientId, result.Error)
	}

	return oauthClient, nil
}

func (ocd *oauthClientRepository) clientExists(clientKey string) bool {
	var exists bool
	result := ocd.Db.Raw("SELECT EXISTS (SELECT 1 FROM oauth_clients WHERE LOWER(client_id) = ?)", strings.ToLower(clientKey)).Scan(&exists)
	if result.Error != nil {
		log.Printf("Error checking existence of client with Id '%s': %v", clientKey, result.Error)
		return false
	}
	return exists
}