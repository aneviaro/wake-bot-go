package user_service

import "wake-bot/user"

type Service interface {
	GetByID(id int64) (*user.User, error)
	NewUser(user user.User) error
	Update(newUser user.User) error
}
