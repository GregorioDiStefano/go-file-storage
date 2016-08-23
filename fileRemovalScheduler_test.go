package main

import (
	"os"
	"testing"
	"time"

	"github.com/GregorioDiStefano/go-file-storage/models"
	"github.com/GregorioDiStefano/go-file-storage/utils"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	utils.ParseConfig("config/config.testing.yaml")

	os.Exit(m.Run())
}

func TestDeleteUnusedFile_1(t *testing.T) {

	utils.Config.Set("FileCheckFrequency", 1)
	utils.Config.Set("DeleteAfterSecondsLastAccessed", 30)
	utils.Config.Set("DeleteKeySize", 6)

	simpleStoredFiled := models.StoredFile{
		Key:        models.DB.FindUnusedKey(),
		DeleteKey:  utils.RandomString(uint8(utils.Config.GetInt("DeleteKeySize"))),
		FileName:   "deleteFile",
		FileSize:   1024,
		UploadTime: time.Now().UTC(),
		LastAccess: time.Now().UTC(),
	}
	models.DB.WriteStoredFile(simpleStoredFiled)

	sf := models.DB.ReadStoredFile(simpleStoredFiled.Key)

	go deleteUnusedFile()

	time.Sleep(15 * time.Second)
	//after 5 seconds, the file will still be here
	assert.False(t, sf.Deleted)
	time.Sleep(35 * time.Second)
	//but now it should be gone.
	sf = models.DB.ReadStoredFile(simpleStoredFiled.Key)
	assert.True(t, sf.Deleted)
}
