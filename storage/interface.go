package storage

import "wake-bot/user"

type Store interface {
	GetByID(id int64) (*user.User, error)
	Save(user user.User) error
}
