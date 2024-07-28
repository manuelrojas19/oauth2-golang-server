package session

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/manuelrojas19/go-oauth2-server/services"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
)

type sessionService struct {
	redisClient *redis.Client
	ctx         context.Context
}

func NewSessionService(redisClient *redis.Client) services.SessionService {
	return &sessionService{
		redisClient: redisClient,
		ctx:         context.Background(),
	}
}

func (u *sessionService) CreateSession(userId, email string) (string, error) {
	// Generate a new session ID
	sessionID := uuid.New().String()
	sessionData := map[string]interface{}{
		"user_id": userId,
		"email":   email,
	}

	// Create a new context with a timeout to prevent hanging operations
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Start a transaction
	pipe := u.redisClient.TxPipeline()

	// Set session data
	pipe.HMSet(ctx, sessionID, sessionData)

	// Set expiration for the session
	pipe.Expire(ctx, sessionID, 1*time.Hour)

	// Execute the transaction
	_, err := pipe.Exec(ctx)
	if err != nil {
		return "", err
	}

	return sessionID, nil
}

func (u *sessionService) SessionExists(sessionID string) bool {
	if sessionID == "" {
		log.Println("Session ID should not be empty")
		return false
	}

	sessionKey := sessionID

	// Check if the session key exists in Redis
	existsCmd := u.redisClient.Exists(u.ctx, sessionKey)
	result, err := existsCmd.Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			// Session does not exist
			log.Println("Session does not exist")
		} else {
			// Error checking session
			log.Printf("Failed to check session existence: %v\n", err)
		}
		return false
	}

	// Check the result value
	if result == 0 {
		// Session does not exist
		log.Println("Session does not exist")
		return false
	}

	// Check the remaining TTL of the session key
	ttlCmd := u.redisClient.TTL(u.ctx, sessionKey)
	ttl, err := ttlCmd.Result()
	if err != nil {
		log.Printf("Failed to check session TTL: %v\n", err)
		return false
	}

	if ttl <= 0 {
		// Session has expired or does not exist
		log.Println("Session has expired")
		return false
	}

	// Session exists and is not expired
	return true
}

func (u *sessionService) GetUserIdFromSession(sessionID string) (string, error) {
	// Fetch user_id from session data in Redis
	userID, err := u.redisClient.HGet(context.Background(), sessionID, "user_id").Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			// Session does not exist or user_id is not found
			return "", fmt.Errorf("session not found or user_id not found in session")
		}
		// Error retrieving user_id
		return "", err
	}

	return userID, nil
}
