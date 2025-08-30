package services

import (
	"github.com/manuelrojas19/go-oauth2-server/store/repositories"
	"go.uber.org/zap"
)

type userConsentService struct {
	consentRepo repositories.AccessConsentRepository
	logger      *zap.Logger
}

func NewUserConsentService(consentRepo repositories.AccessConsentRepository, logger *zap.Logger) UserConsentService {
	return &userConsentService{
		consentRepo: consentRepo,
		logger:      logger,
	}
}

func (c *userConsentService) Save(userId, clientId, scopeId string) error {
	c.logger.Info("Attempting to save user consent", zap.String("userId", userId), zap.String("clientId", clientId), zap.String("scopeId", scopeId))
	c.logger.Debug("Calling consent repository to save consent", zap.String("userId", userId), zap.String("clientId", clientId), zap.String("scopeId", scopeId))
	_, err := c.consentRepo.Save(userId, clientId, scopeId)
	if err != nil {
		c.logger.Error("Failed to save user consent", zap.String("userId", userId), zap.String("clientId", clientId), zap.String("scopeId", scopeId), zap.Error(err))
		return err
	}
	c.logger.Info("Successfully saved user consent", zap.String("userId", userId), zap.String("clientId", clientId), zap.String("scopeId", scopeId))
	return nil
}

func (c *userConsentService) HasUserConsented(userId, clientId, scopeId string) bool {
	c.logger.Info("Checking if user has consented", zap.String("userId", userId), zap.String("clientId", clientId), zap.String("scopeId", scopeId))
	c.logger.Debug("Calling consent repository to check user consent", zap.String("userId", userId), zap.String("clientId", clientId), zap.String("scopeId", scopeId))

	consented, err := c.consentRepo.HasUserConsented(userId, clientId, scopeId)
	if err != nil {
		c.logger.Error("Failed to check user consent", zap.String("userId", userId), zap.String("clientId", clientId), zap.String("scopeId", scopeId), zap.Error(err))
		return false
	}

	c.logger.Info("User consent status retrieved", zap.String("userId", userId), zap.String("clientId", clientId), zap.String("scopeId", scopeId), zap.Bool("consented", consented))
	return consented
}
