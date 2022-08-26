package boltdb

import (
	"bytes"
	"encoding/json"
	"github.com/boltdb/bolt"
	"github.com/vadimpk/go-oxford-telegram-bot/internal/service"
	"strconv"
)

const (
	bucketName = "settings"
)

type SettingsRepository struct {
	db *bolt.DB
}

func NewSettingsRepository(db *bolt.DB) (*SettingsRepository, error) {
	err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		return err
	})
	return &SettingsRepository{db: db}, err
}

func (r *SettingsRepository) Save(chatID int64, settings service.Settings) error {
	return r.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		bt, err := structToBytes(settings)
		if err != nil {
			return err
		}
		return b.Put(intToBytes(chatID), bt)
	})
}

func (r *SettingsRepository) Get(chatID int64) (error, service.Settings) {
	var settings service.Settings

	return r.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		data := b.Get(intToBytes(chatID))
		return bytesToStruct(data, &settings)
	}), settings
}

func intToBytes(v int64) []byte {
	return []byte(strconv.FormatInt(v, 10))
}

func structToBytes(s service.Settings) ([]byte, error) {
	r := new(bytes.Buffer)
	err := json.NewEncoder(r).Encode(s)
	return r.Bytes(), err
}

func bytesToStruct(d []byte, s *service.Settings) error {
	return json.Unmarshal(d, &s)
}
