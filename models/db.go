package models

import (
	"../helpers"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/boltdb/bolt"
	"time"
)

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

func FindUnsedKey(d Database) string {
	count := 0
	possibleKey := helpers.RandomString(helpers.Config.KeySize)
	for d.DoesKeyExist(possibleKey) {
		possibleKey = helpers.RandomString(helpers.Config.KeySize + uint8(count/10))
		count++
	}
	return possibleKey
}

func (database *Database) DoesKeyExist(key string) bool {
	db, _ := bolt.Open(database.Filename, 0600, nil)
	defer db.Close()

	err := db.View(func(tx *bolt.Tx) error {
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
	db, err := bolt.Open(database.Filename, 0600, nil)
	defer db.Close()

	if err != nil {
		panic("Unable to open database for writing")
	}

	err = db.Update(func(tx *bolt.Tx) error {
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
	db, err := bolt.Open(database.Filename, 0666, &bolt.Options{ReadOnly: true})
	defer db.Close()

	if err != nil {
		fmt.Println(err.Error())
		panic("Unable to open database for reading:" + err.Error())
	}

	var sf *StoredFile

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(database.Bucket))

		sf, err = decode(b.Get([]byte(key)))
		fmt.Println(sf, err)

		return nil
	})

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
