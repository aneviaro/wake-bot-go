package repository

import (
	"context"
	"wake-bot/user"

	"cloud.google.com/go/datastore"
)

// Store sets behavior for user repo
type Store interface {
	// GetByID select users with provided id
	GetByID(id int64) (*user.User, error)
	// Save saves user entity somewhere
	Save(user *user.User) error
}

// UserRepository represents a user datastore repository, implements the Store interface.
type UserRepository struct {
	client *datastore.Client
}

const datastoreKind = "User"

// NewRepository creates a new UserRepository.
func NewRepository(client *datastore.Client) *UserRepository {
	return &UserRepository{client: client}
}

// GetByID gets a user from datastore by id.
func (d *UserRepository) GetByID(id int64) (*user.User, error) {
	q := datastore.NewQuery(datastoreKind).FilterField("chat_id", "=", id)

	var u []user.User
	_, err := d.client.GetAll(context.TODO(), q, &u)

	if err != nil || len(u) == 0 {
		return &user.User{}, err
	}

	return &u[0], err
}

// Save saves the given user into the datastore.
func (d *UserRepository) Save(u *user.User) error {
	ctx := context.TODO()

	if u.K == nil {
		u.K = datastore.IncompleteKey(datastoreKind, nil)
	}

	_, err := d.client.Put(ctx, u.K, u)

	return err
}
