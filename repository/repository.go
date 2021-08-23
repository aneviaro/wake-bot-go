package repository

import (
	"context"
	"wake-bot/user"

	"cloud.google.com/go/datastore"
)

// UserRepository represents a user datastore repository, implements the Store interface.
type UserRepository struct {
	client *datastore.Client
}

const datastoreKind = "User"

// NewRepository creates a new UserRepository.
func NewRepository(client *datastore.Client) UserRepository {
	return UserRepository{client: client}
}

// GetById gets a user from datastore by id.
func (d UserRepository) GetByID(id int64) (*user.User, error) {
	ctx := context.Background()
	q := datastore.NewQuery(datastoreKind).Filter("chat_id =", id)

	var u []user.User
	_, err := d.client.GetAll(ctx, q, &u)

	if err != nil || len(u) == 0 {
		return &user.User{}, err
	}

	return &u[0], err
}

// Save saves the given user into the datastore.
func (d UserRepository) Save(u user.User) error {
	ctx := context.Background()

	if u.K == nil {
		u.K = datastore.IncompleteKey(datastoreKind, nil)
	}

	_, err := d.client.Put(ctx, u.K, &u)

	return err
}
