package user

import (
	"errors"
	"strconv"
	"time"
	"wake-bot/repository"
	"wake-bot/user"
)

// IService sets a user service behavior.
type IService interface {
	GetByID(id int64) (*user.User, error)
	NewUser(user *user.User) error
	Update(newUser *user.User) error
	GetUserTimeWithFormat(chatID int64) (*user.TimeWithFormat, error)
}

// Service is responsible for creating, updating and getting users. Is using repository.Store implementation.
type Service struct {
	store repository.Store
}

// NewService returns a new Service instance.
func NewService(store repository.Store) *Service {
	return &Service{store: store}
}

// GetByID select a user by id from repository.
func (u *Service) GetByID(id int64) (*user.User, error) {
	return u.store.GetByID(id)
}

// NewUser creates a new user from given one. Is using repository.Store
func (u *Service) NewUser(newUser *user.User) error {
	return u.store.Save(newUser)
}

// Update updates a user with a new one. Is using repository.Store.
func (u *Service) Update(fromUser *user.User) error {
	toUser, err := u.store.GetByID(fromUser.ChatID)
	if err != nil {
		return err
	}

	if fromUser.TimeZone != "" {
		toUser.TimeZone = fromUser.TimeZone
	}

	if fromUser.UTCOffset != "" {
		toUser.UTCOffset = fromUser.UTCOffset
	}

	if !fromUser.WakeUpTime.IsZero() {
		toUser.WakeUpTime = fromUser.WakeUpTime
	}

	if fromUser.TimeFormat != "" {
		toUser.TimeFormat = fromUser.TimeFormat
	}

	if fromUser.TZDefinedBy != "" {
		toUser.TZDefinedBy = fromUser.TZDefinedBy
	}

	toUser.ChatID = fromUser.ChatID

	return u.store.Save(toUser)
}

func (u *Service) GetUserTimeWithFormat(chatID int64) (*user.TimeWithFormat, error) {
	var tf user.TimeWithFormat

	us, err := u.GetByID(chatID)
	if err != nil {
		return nil, err
	}

	var loc *time.Location

	switch us.TZDefinedBy {
	case user.Timezone:
		loc, err = time.LoadLocation(us.TimeZone)
		if err != nil {
			return nil, err
		}
	case user.Offset:
		offset, err := strconv.Atoi(us.UTCOffset)
		if err != nil {
			return nil, err
		}

		loc = time.FixedZone("manual", offset*60*60)
	default:
		return nil, errors.New("timezone not fulfilled")
	}

	tf.Time = time.Now().In(loc)

	if us.TimeFormat != "" {
		tf.Format = us.TimeFormat
	}

	return &tf, nil
}
