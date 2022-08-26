package boltdb

import (
	"errors"
	"github.com/boltdb/bolt"
)

const (
	statesBucket = "states"
)

type StatesRepository struct {
	db *bolt.DB
}

func NewStatesRepository(db *bolt.DB) (*StatesRepository, error) {
	err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(statesBucket))
		return err
	})
	return &StatesRepository{db: db}, err
}

func (r *StatesRepository) Save(chatID int64, state string) error {
	return r.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(statesBucket))
		return b.Put(intToBytes(chatID), []byte(state))
	})
}

func (r *StatesRepository) Get(chatID int64) (error, string) {
	var state string
	err := r.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(statesBucket))
		data := b.Get(intToBytes(chatID))
		if data == nil {
			return errors.New("user not found")
		}
		state = string(data)
		return nil
	})
	if err != nil {
		return err, ""
	}
	return err, state
}
