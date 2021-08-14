package user

import (
	"wake-bot/repository"
	"wake-bot/user"
)

type Service struct {
	store repository.Store
}

func NewUserService(store repository.Store) *Service {
	return &Service{store: store}
}

func (u Service) GetByID(id int64) (*user.User, error) {
	return u.store.GetByID(id)
}

func (u Service) NewUser(user user.User) error {
	return u.store.Save(user)
}

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
