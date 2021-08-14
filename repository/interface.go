package repository

import "wake-bot/user"

// Store sets behavior for user repo
type Store interface {
	// GetByID select users with provided id
	GetByID(id int64) (*user.User, error)
	// Save saves user entity somewhere
	Save(user user.User) error
}
