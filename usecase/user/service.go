package user

import (
	"wake-bot/repository"
	"wake-bot/user"
)

// Service is responsible for creating, updating and getting users. Is using repository.Store implementation.
type Service struct {
	store repository.Store
}

// NewService returns a new Service instance.
func NewService(store repository.Store) *Service {
	return &Service{store: store}
}

// GetByID select a user by id from repository.
func (u Service) GetByID(id int64) (*user.User, error) {
	return u.store.GetByID(id)
}

// NewUser creates a new user from given one. Is using repository.Store
func (u Service) NewUser(newUser user.User) error {
	return u.store.Save(newUser)
}

// Update updates a user with a new one. Is using repository.Store.
func (u Service) Update(fromUser user.User) error {
	toUser, err := u.store.GetByID(fromUser.ChatID)
	if err != nil {
		return err
	}

	switch {
	case fromUser.TimeZone != nil:
		toUser.TimeZone = fromUser.TimeZone
	case !fromUser.WakeUpTime.IsZero():
		toUser.WakeUpTime = fromUser.WakeUpTime
	case fromUser.TimeFormat != "":
		toUser.TimeFormat = fromUser.TimeFormat
	}

	toUser.ChatID = fromUser.ChatID

	return u.store.Save(*toUser)
}
