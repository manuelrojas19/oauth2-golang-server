package services

import (
	"time"

	"github.com/manuelrojas19/go-oauth2-server/store/repositories"
	"go.uber.org/zap"
)

type IntrospectCommand struct {
	Token         string
	TokenTypeHint string
}

type IntrospectionResponse struct {
	Active    bool   `json:"active"`
	Scope     string `json:"scope,omitempty"`
	ClientId  string `json:"client_id,omitempty"`
	ExpiresAt int64  `json:"exp,omitempty"`
	IssuedAt  int64  `json:"iat,omitempty"`
	Subject   string `json:"sub,omitempty"`
	Audience  string `json:"aud,omitempty"`
	Issuer    string `json:"iss,omitempty"`
	TokenType string `json:"token_type,omitempty"`
}

type introspectionService struct {
	accessTokenRepository  repositories.AccessTokenRepository
	refreshTokenRepository repositories.RefreshTokenRepository
	logger                 *zap.Logger
}

func NewIntrospectionService(accessTokenRepository repositories.AccessTokenRepository, refreshTokenRepository repositories.RefreshTokenRepository, logger *zap.Logger) IntrospectionService {
	return &introspectionService{
		accessTokenRepository:  accessTokenRepository,
		refreshTokenRepository: refreshTokenRepository,
		logger:                 logger,
	}
}

func (s *introspectionService) Introspect(command *IntrospectCommand) (*IntrospectionResponse, error) {
	s.logger.Info("Attempting to introspect token", zap.String("token", command.Token), zap.String("tokenTypeHint", command.TokenTypeHint))
	s.logger.Debug("Introspection command details", zap.Any("command", command))

	// Try to find the token as an access token
	s.logger.Debug("Attempting to find token as an access token")
	accessTokenEntity, err := s.accessTokenRepository.FindByAccessToken(command.Token)
	if err == nil && accessTokenEntity != nil {
		s.logger.Debug("Access token found", zap.String("accessTokenId", accessTokenEntity.Id))
		return s.buildIntrospectionResponse(*accessTokenEntity.UserId, *accessTokenEntity.ClientId, accessTokenEntity.Scope, accessTokenEntity.CreatedAt, time.Until(accessTokenEntity.ExpiresAt), "access_token"), nil
	}
	// Log error if any, or if token not found as access token
	if err != nil {
		s.logger.Debug("Error finding access token", zap.Error(err))
	} else {
		s.logger.Debug("Token not found as an access token")
	}

	// If not an access token, try to find it as a refresh token
	s.logger.Debug("Attempting to find token as a refresh token")
	refreshTokenEntity, err := s.refreshTokenRepository.FindByRefreshToken(command.Token)
	if err == nil && refreshTokenEntity != nil {
		s.logger.Debug("Refresh token found", zap.String("refreshTokenId", refreshTokenEntity.Id))
		return s.buildIntrospectionResponse(*refreshTokenEntity.UserId, *refreshTokenEntity.ClientId, refreshTokenEntity.Scope, refreshTokenEntity.CreatedAt, time.Until(refreshTokenEntity.ExpiresAt), "refresh_token"), nil
	}

	if err != nil {
		s.logger.Debug("Error finding refresh token", zap.Error(err))
	} else {
		s.logger.Debug("Token not found as a refresh token")
	}

	s.logger.Info("Introspection complete: token is inactive or not found",
		zap.String("token", command.Token),
		zap.Bool("active", false),
	)
	return &IntrospectionResponse{Active: false}, nil
}

func (s *introspectionService) buildIntrospectionResponse(userId, clientId, scope string, createdAt time.Time, expiresIn time.Duration, tokenType string) *IntrospectionResponse {
	s.logger.Debug("Building introspection response",
		zap.String("userId", userId),
		zap.String("clientId", clientId),
		zap.String("scope", scope),
		zap.String("tokenType", tokenType),
	)
	isActive := true
	expiresAt := createdAt.Add(expiresIn)

	if time.Now().After(expiresAt) {
		isActive = false
		s.logger.Debug("Token is expired", zap.String("tokenType", tokenType), zap.Time("expiresAt", expiresAt))
	}

	response := &IntrospectionResponse{
		Active:    isActive,
		Scope:     scope,
		ClientId:  clientId,
		ExpiresAt: expiresAt.Unix(),
		IssuedAt:  createdAt.Unix(),
		Subject:   userId,
		TokenType: tokenType,
	}
	s.logger.Debug("Introspection response built", zap.Any("response", response))
	return response
}
