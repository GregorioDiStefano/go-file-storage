package models

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func init() {
	const (
		DbFilename = "/tmp/files.testing.db"
		Bucket     = "files"
	)

	//make sure the test db is clean
	os.Remove("/tmp/files.testing.db")

	var DB = Database{Filename: DbFilename, Bucket: Bucket}
	DB.Setup(10)
	DB.OpenDatabaseFile()
}

func TestWriteStoredFile(t *testing.T) {
	simpleStoredFiled := StoredFile{
		Key:        "a",
		DeleteKey:  "b",
		FileName:   "filename",
		FileSize:   2345,
		UploadTime: time.Now().UTC(),
		LastAccess: time.Now().UTC(),
	}
	DB.WriteStoredFile(simpleStoredFiled)
}

func TestReadStoredFile(t *testing.T) {
	simpleStoredFiled := StoredFile{
		Key:        "a1",
		DeleteKey:  "b",
		FileName:   "filename",
		FileSize:   2345,
		UploadTime: time.Now().UTC(),
		LastAccess: time.Now().UTC(),
	}
	DB.WriteStoredFile(simpleStoredFiled)
	assert.Equal(t, simpleStoredFiled, *DB.ReadStoredFile("a1"))
}

func TestGetAllKeys(t *testing.T) {
	assert.Equal(t, []string{"a", "a1"}, *DB.GetAllKeys())
}

func TestDoesKeyExist(t *testing.T) {
	for _, s := range *DB.GetAllKeys() {
		assert.True(t, DB.DoesKeyExist(s))
	}
}
