package user_service

import (
	"wake-bot/storage"
	"wake-bot/user"
)

type UserService struct {
	store *storage.Storage
}

func NewUserService(store *storage.Storage) *UserService {
	return &UserService{store: store}
}

func (u UserService) GetByID(id int64) (*user.User, error) {
	return u.store.GetByID(id)
}

func (u UserService) NewUser(user user.User) error {
	return u.store.Save(user)
}

func (u UserService) Update(fromUser user.User) error {
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
