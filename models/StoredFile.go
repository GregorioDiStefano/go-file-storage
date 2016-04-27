package models

import (
	"../helpers"
	"encoding/json"
	"errors"
	_ "fmt"
	"github.com/boltdb/bolt"
	"time"
)

var boltdb *bolt.DB

const (
	DbFilename = "files.db"
	Bucket     = "files"
)

type Database struct {
	Filename string
	Bucket   string
}

type StoredFile struct {
	Key        string
	FileName   string
	FileSize   int64
	DeleteKey  string
	Downloads  int64
	UploadTime time.Time
}

func (database *Database) OpenDatabaseFile() {
	var err error
	boltdb, err = bolt.Open(database.Filename, 0600, nil)
	if err != nil {
		panic(err)
	}
}

func (database *Database) FindUnsedKey() string {
	count := 0
	possibleKey := helpers.RandomString(helpers.Config.KeySize)
	for database.DoesKeyExist(possibleKey) {
		possibleKey = helpers.RandomString(helpers.Config.KeySize + uint8(count/10))
		count++
	}
	return possibleKey
}

func (database *Database) DoesKeyExist(key string) bool {

	err := boltdb.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(database.Bucket))
		record := b.Get([]byte(key))
		if len(record) > 0 {
			return errors.New("Key exists")
		}
		return nil
	})

	if err != nil {
		return true
	}
	return false
}

func (database *Database) WriteStoredFile(sf StoredFile) error {
	err := boltdb.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(database.Bucket))
		if err != nil {
			return err
		}

		encoded, err := json.Marshal(sf)

		if err != nil {
			return err
		}

		err = b.Put([]byte(sf.Key), []byte(encoded))
		return err
	})
	return err
}

func (database *Database) ReadStoredFile(key string) *StoredFile {

	var sf *StoredFile

	err := boltdb.View(func(tx *bolt.Tx) error {
		var err error
		b := tx.Bucket([]byte(database.Bucket))
		sf, err = decode(b.Get([]byte(key)))

		return err
	})

	if err != nil || sf == nil {
		panic(err)
	}

	return sf
}

func decode(data []byte) (*StoredFile, error) {
	var p *StoredFile
	err := json.Unmarshal(data, &p)
	if err != nil {
		return nil, err
	}
	return p, nil
}
