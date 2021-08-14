package repository

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"wake-bot/user"

	"github.com/boltdb/bolt"
)

// Storage is deprecated
type Storage struct {
	db *bolt.DB
}

func NewBoltStore() (*Storage, error) {
	db, err := bolt.Open("user.db", 0600, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to open db connection, err: %v", err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		_, e := tx.CreateBucketIfNotExists([]byte("users"))
		if e != nil {
			return fmt.Errorf("failed to create users bucket, err: %v", e)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &Storage{db: db}, nil
}

func (b Storage) GetByID(id int64) (*user.User, error) {
	var u user.User
	err := b.db.View(func(tx *bolt.Tx) error {
		userBkt := tx.Bucket([]byte("users"))

		v := userBkt.Get(itob(id))

		if len(v) > 0 {
			err := json.Unmarshal(v, &u)
			if err != nil {
				return fmt.Errorf("failed to unmarshal user")
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (b Storage) Save(user user.User) error {
	err := b.db.Update(func(tx *bolt.Tx) error {
		userBkt := tx.Bucket([]byte("users"))

		buf, err := json.Marshal(user)
		if err != nil {
			return err
		}

		return userBkt.Put(itob(user.ChatID), buf)
	})

	if err != nil {
		return fmt.Errorf("failed to create user, err: %v", err)
	}

	return nil
}

func itob(v int64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}
