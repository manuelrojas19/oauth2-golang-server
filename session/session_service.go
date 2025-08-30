package session

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/manuelrojas19/go-oauth2-server/services"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type sessionService struct {
	redisClient *redis.Client
	ctx         context.Context
	logger      *zap.Logger
}

func NewSessionService(redisClient *redis.Client, logger *zap.Logger) services.SessionService {
	return &sessionService{
		redisClient: redisClient,
		ctx:         context.Background(),
		logger:      logger,
	}
}

func (u *sessionService) CreateSession(userId, email string) (string, error) {
	start := time.Now()
	sessionID := uuid.New().String()
	sessionData := map[string]interface{}{
		"user_id": userId,
		"email":   email,
	}
	u.logger.Info("Attempting to create session", zap.String("userId", userId), zap.String("email", email))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pipe := u.redisClient.TxPipeline()
	pipe.HMSet(ctx, sessionID, sessionData)
	pipe.Expire(ctx, sessionID, 1*time.Hour)
	u.logger.Debug("Executing Redis pipeline to set session data and expiry", zap.String("sessionId", sessionID))

	_, err := pipe.Exec(ctx)
	if err != nil {
		u.logger.Error("Error executing Redis pipeline for session creation",
			zap.String("sessionId", sessionID),
			zap.String("userId", userId),
			zap.String("email", email),
			zap.Error(err),
			zap.Duration("duration", time.Since(start)),
			zap.Stack("stacktrace"),
		)
		return "", fmt.Errorf("failed to create session: %w", err)
	}

	u.logger.Info("Successfully created session",
		zap.String("sessionId", sessionID),
		zap.String("userId", userId),
		zap.String("email", email),
		zap.Duration("duration", time.Since(start)),
	)
	u.logger.Debug("Session creation complete", zap.String("sessionId", sessionID))

	return sessionID, nil
}

func (u *sessionService) SessionExists(sessionID string) bool {
	start := time.Now()
	u.logger.Info("Checking if session exists", zap.String("sessionId", sessionID))
	if sessionID == "" {
		u.logger.Warn("Session ID is empty during existence check")
		return false
	}

	sessionKey := sessionID
	existsCmd := u.redisClient.Exists(u.ctx, sessionKey)
	result, err := existsCmd.Result()
	if err != nil {
		if err == redis.Nil {
			u.logger.Info("Session does not exist for ID",
				zap.String("sessionId", sessionID),
				zap.Duration("duration", time.Since(start)),
			)
		} else {
			u.logger.Error("Error checking session existence in Redis",
				zap.String("sessionId", sessionID),
				zap.Error(err),
				zap.Duration("duration", time.Since(start)),
				zap.Stack("stacktrace"),
			)
		}
		return false
	}

	if result == 0 {
		u.logger.Info("Session does not exist for ID",
			zap.String("sessionId", sessionID),
			zap.Duration("duration", time.Since(start)),
		)
		return false
	}

	ttlCmd := u.redisClient.TTL(u.ctx, sessionKey)
	ttl, err := ttlCmd.Result()
	if err != nil {
		u.logger.Error("Error checking session TTL in Redis",
			zap.String("sessionId", sessionID),
			zap.Error(err),
			zap.Duration("duration", time.Since(start)),
			zap.Stack("stacktrace"),
		)
		return false
	}

	if ttl <= 0 {
		u.logger.Info("Session has expired",
			zap.String("sessionId", sessionID),
			zap.Duration("duration", time.Since(start)),
		)
		return false
	}

	u.logger.Info("Session exists and is not expired",
		zap.String("sessionId", sessionID),
		zap.Duration("duration", time.Since(start)),
	)
	return true
}

func (u *sessionService) GetUserIdFromSession(sessionID string) (string, error) {
	start := time.Now()
	u.logger.Info("Attempting to retrieve user ID from session", zap.String("sessionId", sessionID))
	userID, err := u.redisClient.HGet(context.Background(), sessionID, "user_id").Result()
	if err != nil {
		if err == redis.Nil {
			u.logger.Info("Session not found or user_id not found in session",
				zap.String("sessionId", sessionID),
				zap.Duration("duration", time.Since(start)),
			)
			return "", fmt.Errorf("session not found or user_id not found in session")
		}
		u.logger.Error("Error retrieving user_id from session",
			zap.String("sessionId", sessionID),
			zap.Error(err),
			zap.Duration("duration", time.Since(start)),
			zap.Stack("stacktrace"),
		)
		return "", fmt.Errorf("failed to retrieve user ID from session: %w", err)
	}

	u.logger.Info("Successfully retrieved user_id from session",
		zap.String("sessionId", sessionID),
		zap.String("userId", userID),
		zap.Duration("duration", time.Since(start)),
	)
	u.logger.Debug("User ID retrieved from session", zap.String("userId", userID))

	return userID, nil
}

func (u *sessionService) DeleteSession(sessionID string) error {
	start := time.Now()
	u.logger.Info("Attempting to delete session", zap.String("sessionId", sessionID))

	if sessionID == "" {
		u.logger.Warn("Session ID is empty for deletion attempt")
		return fmt.Errorf("session ID cannot be empty")
	}

	err := u.redisClient.Del(u.ctx, sessionID).Err()
	if err != nil {
		u.logger.Error("Error deleting session from Redis",
			zap.String("sessionId", sessionID),
			zap.Error(err),
			zap.Duration("duration", time.Since(start)),
			zap.Stack("stacktrace"),
		)
		return fmt.Errorf("failed to delete session: %w", err)
	}

	u.logger.Info("Session deleted successfully",
		zap.String("sessionId", sessionID),
		zap.Duration("duration", time.Since(start)),
	)
	return nil
}
