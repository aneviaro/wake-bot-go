package user_service

import (
	"time"
	"wake-bot/storage"
	"wake-bot/user"
)

type UserService struct {
	store storage.Store
}

func NewUserService(store storage.Store) *UserService {
	return &UserService{store: store}
}

func (us UserService) GetByID(id int64) (*user.User, error) {
	return us.store.GetByID(id)
}

func (us UserService) NewUser(user user.User) error {
	return us.store.Save(user)
}

func (us UserService) Update(fromUser user.User) error {
	toUser, err := us.store.GetByID(fromUser.ChatID)
	if err != nil {
		return err
	}

	switch {
	case fromUser.TimeZone != "":
		toUser.TimeZone = fromUser.TimeZone
	case !fromUser.WakeUpTime.IsZero():
		toUser.WakeUpTime = fromUser.WakeUpTime
	case fromUser.TimeFormat != "":
		toUser.TimeFormat = fromUser.TimeFormat
	}

	toUser.ChatID = fromUser.ChatID

	return us.store.Save(*toUser)
}

func (us UserService) GetUserTime(chatID int64) string {
	defaultTime := time.Now().Format("3:04PM")

	u, err := us.GetByID(chatID)
	if err != nil {
		return defaultTime
	}

	var loc *time.Location

	switch {
	case u.UTCOffset != 0:
		loc = time.FixedZone("manual", u.UTCOffset*60*60)
	case u.TimeZone != "":
		loc, err = time.LoadLocation(u.TimeZone)
		if err != nil {
			return defaultTime
		}
	default:
		loc, err = time.LoadLocation("UTC")
		if err != nil {
			return defaultTime
		}
	}

	if u.TimeFormat != "" {
		return time.Now().In(loc).Format(string(u.TimeFormat))
	}

	return time.Now().In(loc).Format("3:04PM")
}
