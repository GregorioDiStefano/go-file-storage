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
	os.Exit(m.Run())
}

func TestDeleteUnusedFile_1(t *testing.T) {
	helpers.Config.FileCheckFrequency = 1
	helpers.Config.DeleteAfterSecondsLastAccessed = 10

	simpleStoredFiled := models.StoredFile{
		Key:        models.DB.FindUnsedKey(),
		DeleteKey:  helpers.RandomString(helpers.Config.DeleteKeySize),
		FileName:   "deleteFile",
		FileSize:   1024,
		UploadTime: time.Now().UTC(),
		LastAccess: time.Now().UTC(),
	}
	models.DB.WriteStoredFile(simpleStoredFiled)

	sf := models.DB.ReadStoredFile(simpleStoredFiled.Key)

	go deleteUnusedFile()

	time.Sleep(5 * time.Second)
	//after 5 seconds, the file will still be here
	assert.False(t, sf.Deleted)
	time.Sleep(15 * time.Second)
	//but now it should be gone.
	sf = models.DB.ReadStoredFile(simpleStoredFiled.Key)
	assert.True(t, sf.Deleted)
}
