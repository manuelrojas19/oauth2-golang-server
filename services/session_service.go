package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
)

type sessionService struct {
	redisClient *redis.Client
	ctx         context.Context
}

func NewSessionService(redisClient *redis.Client) SessionService {
	return &sessionService{
		redisClient: redisClient,
		ctx:         context.Background(),
	}
}

func (u *sessionService) CreateSession(userId, email string) (string, error) {
	sessionID := uuid.New().String()
	sessionData := map[string]string{
		"user_id": userId,
		"email":   email,
	}

	err := u.redisClient.HSet(context.Background(), sessionID, sessionData).Err()
	if err != nil {
		return "", err
	}

	// Set expiration for the session
	err = u.redisClient.Expire(context.Background(), sessionID, 1*time.Hour).Err()
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

	// Session exists
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
