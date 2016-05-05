package main

import (
	"os"
	"testing"
	"time"

	"github.com/GregorioDiStefano/go-file-storage/helpers"
	"github.com/GregorioDiStefano/go-file-storage/models"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	helpers.ParseConfig("config/config.testing.json")
	//models.DB.OpenDatabaseFile()
	os.Exit(m.Run())
}

func TestDeleteUnusedFile(t *testing.T) {
	helpers.Config.FileCheckFrequency = 1
	helpers.Config.DeleteAfterSecondsLastAccessed = 10

	key := models.DB.FindUnsedKey()
	deleteKey := helpers.RandomString(helpers.Config.DeleteKeySize)

	simpleStoredFiled := models.StoredFile{
		Key:        key,
		DeleteKey:  deleteKey,
		FileName:   "deleteFile",
		FileSize:   1024,
		UploadTime: time.Now().UTC(),
		LastAccess: time.Now().UTC(),
	}
	models.DB.WriteStoredFile(simpleStoredFiled)

	sf := models.DB.ReadStoredFile(simpleStoredFiled.Key)

	go deleteUnusedFile()

	time.Sleep(1 * time.Second)
	assert.False(t, sf.Deleted)
	time.Sleep(15 * time.Second)
	sf = models.DB.ReadStoredFile(simpleStoredFiled.Key)
	assert.True(t, sf.Deleted)
}
