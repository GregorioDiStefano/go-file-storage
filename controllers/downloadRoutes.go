package controller

import (
	"../helpers"
	"../models"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

func FileDownloader(c *gin.Context) {
	db := models.Database{Filename: models.DbFilename, Bucket: models.Bucket}
	key := c.Param("key")
	fn := c.Param("filename")

	if db.DoesKeyExist(key) == false {
		c.String(http.StatusForbidden, "Invalid key")
		return
	}

	expectedFilePath := fmt.Sprintf("%s/%s/%s",
		helpers.Config.StorageFolder,
		key,
		fn)

	if _, err := os.Stat(expectedFilePath); os.IsNotExist(err) {
		fmt.Print(expectedFilePath + " does not exist.")
		c.String(http.StatusForbidden, "Doesn't look like that file exists.")
		return
	}

	sf := db.ReadStoredFile(key)

	if sf.Downloads >= helpers.Config.MaxDownloadsBeforeInteraction {
		c.String(http.StatusForbidden, "This file has been download too many times.")
		return
	}

	sf.Downloads = sf.Downloads + 1
	db.WriteStoredFile(*sf)

	c.File(expectedFilePath)
}
