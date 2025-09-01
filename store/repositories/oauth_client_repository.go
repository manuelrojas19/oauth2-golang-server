package repositories

import (
	"errors"
	"fmt"
	"strings"

	"github.com/manuelrojas19/go-oauth2-server/store"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type oauthClientRepository struct {
	Db     *gorm.DB
	logger *zap.Logger
}

func NewOauthClientRepository(db *gorm.DB, logger *zap.Logger) OauthClientRepository {
	return &oauthClientRepository{Db: db, logger: logger}
}

func (ocd *oauthClientRepository) Save(client *store.OauthClient) (*store.OauthClient, error) {
	ocd.logger.Info("Starting transaction to save OAuth client", zap.String("clientId", client.ClientId))
	ocd.logger.Debug("OAuth client details to be saved", zap.Any("client", client))

	tx := ocd.Db.Begin()
	if tx.Error != nil {
		ocd.logger.Error("Failed to begin transaction for saving OAuth client", zap.Error(tx.Error))
		return nil, fmt.Errorf("failed to start transaction: %w", tx.Error)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			ocd.logger.Error("PANIC: Rolled back transaction for OAuth client", zap.String("clientId", client.ClientId), zap.Any("panicReason", r), zap.Stack("stacktrace"))
		}
	}()

	// Check if client already exists
	ocd.logger.Debug("Checking if OAuth client already exists", zap.String("clientName", client.ClientName))
	if ocd.ExistsByName(client.ClientName) {
		ocd.logger.Warn("OAuth client with client name already exists", zap.String("clientName", client.ClientName))
		tx.Rollback()
		return nil, fmt.Errorf("client with name '%s' already exists", client.ClientName)
	}
	ocd.logger.Debug("OAuth client does not exist, proceeding with creation")

	// Create new client
	if err := tx.Create(client).Error; err != nil {
		ocd.logger.Error("Failed to create OAuth client in database", zap.String("clientId", client.ClientId), zap.Error(err))
		tx.Rollback()
		return nil, fmt.Errorf("failed to create client: %w", err)
	}
	ocd.logger.Debug("OAuth client created in database", zap.String("clientId", client.ClientId))

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		ocd.logger.Error("Error committing transaction for saving OAuth client", zap.String("clientId", client.ClientId), zap.Error(err))
		tx.Rollback()
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	ocd.logger.Info("Successfully saved new OAuth client", zap.String("clientId", client.ClientId), zap.String("clientName", client.ClientName))
	return client, nil
}

func (ocd *oauthClientRepository) FindByClientId(clientId string) (*store.OauthClient, error) {
	ocd.logger.Info("Attempting to find OAuth client by client ID", zap.String("clientId", clientId))

	oauthClient := new(store.OauthClient)
	result := ocd.Db.Where("LOWER(client_id) = ?", strings.ToLower(clientId)).First(oauthClient)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			ocd.logger.Debug("OAuth client not found for client ID", zap.String("clientId", clientId))
			return nil, fmt.Errorf("OAuth client with clientId '%s' not found", clientId)
		}
		ocd.logger.Error("Error finding OAuth client in database", zap.String("clientId", clientId), zap.Error(result.Error))
		return nil, fmt.Errorf("error finding OAuth client with ClientId '%s': %w", clientId, result.Error)
	}
	ocd.logger.Info("OAuth client found successfully", zap.String("clientId", oauthClient.ClientId))
	ocd.logger.Debug("Found OAuth client details", zap.Any("oauthClient", oauthClient))

	return oauthClient, nil
}

func (ocd *oauthClientRepository) ExistsByName(clientName string) bool {
	ocd.logger.Debug("Checking existence of client by name", zap.String("clientName", clientName))
	var exists bool
	result := ocd.Db.Raw("SELECT EXISTS (SELECT 1 FROM oauth_clients WHERE LOWER(client_name) = ?)", strings.ToLower(clientName)).Scan(&exists)
	if result.Error != nil {
		ocd.logger.Error("Error checking existence of client by name in database", zap.String("clientName", clientName), zap.Error(result.Error))
		return false
	}
	ocd.logger.Debug("Client existence check result by name", zap.String("clientName", clientName), zap.Bool("exists", exists))
	return exists
}
