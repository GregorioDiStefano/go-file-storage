package main

import (
	"encoding/json"
	"fmt"

	"github.com/boltdb/bolt"
)

type database struct {
	filename string
	bucket   string
}

func (database *database) writeStoredFile(sf StoredFile) error {
	db, err := bolt.Open(database.filename, 0600, nil)
	defer db.Close()

	if err != nil {
		panic("Unable to open database for writing")
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(database.bucket))
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

func (database *database) readStoredFile(key string) error {
	db, err := bolt.Open(database.filename, 0666, &bolt.Options{ReadOnly: true})
	defer db.Close()

	if err != nil {
		panic("Unable to open database for reading:" + err.Error())
	}

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(database.bucket))
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			fmt.Printf("key=%s, value=%s\n", k, v)
		}

		return nil
	})

	return nil
}
