package repositories

import (
	"errors"
	"fmt"
	"github.com/manuelrojas19/go-oauth2-server/store"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// userRepository is a concrete implementation of UserRepository
type userRepository struct {
	db     *gorm.DB
	logger *zap.Logger
}

// NewUserRepository creates a new instance of userRepository
func NewUserRepository(db *gorm.DB, logger *zap.Logger) UserRepository {
	return &userRepository{db: db, logger: logger}
}

// Save creates or updates a user in the database
func (r *userRepository) Save(user *store.User) (*store.User, error) {
	r.logger.Info("Saving user", zap.String("userID", user.Id))

	// Perform the save operation
	result := r.db.Save(user)

	// Handle errors during the save operation
	if result.Error != nil {
		r.logger.Error("Error saving user",
			zap.String("userID", user.Id),
			zap.Error(result.Error),
			zap.Stack("stacktrace"),
		)
		return nil, fmt.Errorf("error saving user: %w", result.Error)
	}

	r.logger.Info("Successfully saved user", zap.String("userID", user.Id))
	return user, nil
}

// FindByUserId retrieves a user by ID from the database.
func (r *userRepository) FindByUserId(id string) (*store.User, error) {
	r.logger.Info("Searching for user", zap.String("userID", id))

	// Initialize a new User entity
	user := new(store.User)

	// Query the database for the user
	result := r.db.Where("id = ?", id).First(user)

	// Handle errors during the query
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			r.logger.Warn("User not found", zap.String("userID", id))
			return nil, fmt.Errorf("user not found")
		}
		r.logger.Error("Error finding user",
			zap.String("userID", id),
			zap.Error(result.Error),
			zap.Stack("stacktrace"),
		)
		return nil, fmt.Errorf("error finding user: %w", result.Error)
	}

	r.logger.Info("Successfully found user", zap.String("userID", id))
	return user, nil
}
