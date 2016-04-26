package controller

import (
	"../helpers"
	"../models"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

func DeleteFile(c *gin.Context) {
	key := c.Param("key")
	deleteKey := c.Param("delete_key")

	db := models.Database{Filename: models.DbFilename, Bucket: models.Bucket}

	if db.DoesKeyExist(key) {
		sf := db.ReadStoredFile(key)
		if sf.DeleteKey == deleteKey {
			filePath := fmt.Sprintf("%s/%s/%s",
				helpers.Config.StorageFolder,
				key,
				sf.FileName)

			if _, err := os.Stat(filePath); err == nil {
				os.Remove(filePath)
				c.String(http.StatusOK, "Deleted file: "+sf.FileName+"\n")
			} else {
				fmt.Println(filePath, err)
				c.String(http.StatusBadRequest, "File does not exist")
			}

		} else {
			c.String(http.StatusUnauthorized, "Wrong delete key!")
		}
	}
}
