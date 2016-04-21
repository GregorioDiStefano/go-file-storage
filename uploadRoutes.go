package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"./helpers"
	"github.com/gin-gonic/gin"
)

func simpleUpload(c *gin.Context) {
	checkUploadSize(c)
	FileName := c.Param("filename")
	Key := helpers.RandomString(helpers.Config.KeySize)
	processUpload(c.Request.Body, Key, FileName)
	db := database{filename: dbFilename, bucket: bucket}

	simpleStoredFiled := StoredFile{
		MaxDownloads: 10,
		Key:          Key,
		FileName:     FileName,
		UploadTime:   time.Now().UTC()}

	db.writeStoredFile(simpleStoredFiled)
	c.String(http.StatusOK, fmt.Sprintf(Key))
}

func upload(c *gin.Context) {
	FileSize, err := checkUploadSize(c)

	if err != nil {
		return
	}

	file, headers, err := c.Request.FormFile("upload")

	DeleteKey := c.Request.FormValue("DeleteKey")
	MaxDownloads, _ := strconv.ParseInt(c.Request.FormValue("MaxDownloads"), 10, 64)
	Downloads := int64(0)
	UploadTime := time.Now().UTC()
	Key := helpers.RandomString(helpers.Config.KeySize)
	FileName := headers.Filename

	if file != nil {
		processUpload(file, Key, FileName)
	}

	if err != nil {
		panic(err)
	}

	sf := StoredFile{Key,
		FileName,
		FileSize,
		DeleteKey,
		MaxDownloads,
		Downloads,
		UploadTime}

	db := database{filename: dbFilename, bucket: bucket}
	db.writeStoredFile(sf)

	c.String(http.StatusOK, fmt.Sprintf(Key))
}
