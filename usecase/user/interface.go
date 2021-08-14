package user

import "wake-bot/user"

type IService interface {
	GetByID(id int64) (*user.User, error)
	NewUser(user user.User) error
	Update(newUser user.User) error
}
