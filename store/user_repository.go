package store

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"log"
)

// userRepository is a concrete implementation of UserRepository
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new instance of userRepository
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// Save creates or updates a user in the database
func (r *userRepository) Save(user *User) (*User, error) {
	log.Printf("Saving user with ID %s", user.Id)

	// Perform the save operation
	result := r.db.Save(user)

	// Handle errors during the save operation
	if result.Error != nil {
		log.Printf("Error saving user with ID %s: %v", user.Id, result.Error)
		return nil, fmt.Errorf("error saving user: %w", result.Error)
	}

	log.Printf("Successfully saved user with ID %s", user.Id)
	return user, nil
}

// FindByUserId retrieves a user by ID from the database.
func (r *userRepository) FindByUserId(id string) (*User, error) {
	log.Printf("Searching for user with ID %s", id)

	// Initialize a new User entity
	user := new(User)

	// Query the database for the user
	result := r.db.Where("id = ?", id).First(user)

	// Handle errors during the query
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Printf("User with ID %s not found", id)
			return nil, fmt.Errorf("user not found")
		}
		log.Printf("Error finding user with ID %s: %v", id, result.Error)
		return nil, fmt.Errorf("error finding user: %w", result.Error)
	}

	log.Printf("Successfully found user with ID %s", id)
	return user, nil
}
