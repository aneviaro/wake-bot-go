package repository

import (
	"context"
	"wake-bot/user"

	"cloud.google.com/go/datastore"
)

type UserRepository struct {
	client *datastore.Client
}

const datastoreKind = "User"

func NewDatastore(client *datastore.Client) UserRepository {
	return UserRepository{client: client}
}

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

func (d UserRepository) Save(u user.User) error {
	ctx := context.Background()

	if u.K == nil {
		u.K = datastore.IncompleteKey(datastoreKind, nil)
	}

	_, err := d.client.Put(ctx, u.K, &u)

	return err
}
