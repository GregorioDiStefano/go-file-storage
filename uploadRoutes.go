package main

import (
	"fmt"
	"net/http"
	"time"

	"./helpers"
	"github.com/gin-gonic/gin"
)

func simpleUpload(c *gin.Context) {

	if _, err := checkUploadSize(c); err != nil {
		return
	}

	db := database{filename: dbFilename, bucket: bucket}
	fn := c.Param("filename")
	key := findUnsedKey(db)

	processUpload(c.Request.Body, key, fn)

	simpleStoredFiled := StoredFile{
		MaxDownloads: 10,
		Key:          key,
		DeleteKey:    helpers.RandomString(helpers.Config.DeleteKeySize),
		FileName:     fn,
		UploadTime:   time.Now().UTC()}

	db.writeStoredFile(simpleStoredFiled)
	c.String(http.StatusOK, fmt.Sprintf("download key: %s, delete key: %s\n", key, simpleStoredFiled.DeleteKey))
}
