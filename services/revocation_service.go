package services

import (
	"fmt"

	"github.com/manuelrojas19/go-oauth2-server/api"
	"github.com/manuelrojas19/go-oauth2-server/store/repositories"
	"go.uber.org/zap"
)

type RevokeCommand struct {
	Token         string
	TokenTypeHint string
}

type revocationService struct {
	accessTokenRepository  repositories.AccessTokenRepository
	refreshTokenRepository repositories.RefreshTokenRepository
	logger                 *zap.Logger
}

func NewRevocationService(accessTokenRepository repositories.AccessTokenRepository, refreshTokenRepository repositories.RefreshTokenRepository, logger *zap.Logger) RevocationService {
	return &revocationService{
		accessTokenRepository:  accessTokenRepository,
		refreshTokenRepository: refreshTokenRepository,
		logger:                 logger,
	}
}

func (s *revocationService) Revoke(command *RevokeCommand) error {
	s.logger.Info("Attempting to revoke token", zap.String("token", command.Token), zap.String("tokenTypeHint", command.TokenTypeHint))

	if command.TokenTypeHint == "access_token" || command.TokenTypeHint == "" {
		err := s.accessTokenRepository.DeleteByAccessToken(command.Token)
		if err != nil {
			s.logger.Error("Error revoking access token", zap.Error(err))
			return fmt.Errorf(api.ErrServerError.Error())
		}
		s.logger.Info("Access token revoked successfully", zap.String("token", command.Token))
		return nil
	}

	if command.TokenTypeHint == "refresh_token" || command.TokenTypeHint == "" {
		err := s.refreshTokenRepository.DeleteByRefreshToken(command.Token)
		if err != nil {
			s.logger.Error("Error revoking refresh token", zap.Error(err))
			return fmt.Errorf(api.ErrServerError.Error())
		}
		s.logger.Info("Refresh token revoked successfully", zap.String("token", command.Token))
		return nil
	}

	s.logger.Warn("Unsupported token type hint for revocation", zap.String("tokenTypeHint", command.TokenTypeHint))
	return api.ErrUnsupportedTokenType
}
