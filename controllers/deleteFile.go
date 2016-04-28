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
	fileName := c.Param("filename")

	db := models.Database{Filename: models.DbFilename, Bucket: models.Bucket}

	if db.DoesKeyExist(key) {
		sf := db.ReadStoredFile(key)
		if sf.DeleteKey == deleteKey && sf.FileName == fileName {
			filePath := fmt.Sprintf("%s/%s/%s",
				helpers.Config.StorageFolder,
				key,
				sf.FileName)

			if _, err := os.Stat(filePath); err == nil {
				os.Remove(filePath)
				c.String(http.StatusOK, "Deleted file: "+sf.FileName+"\n")
				return
			} else {
				fmt.Println(filePath, err)
				c.String(http.StatusBadRequest, "File does not exist")
				return
			}

		} else {
			c.String(http.StatusUnauthorized, "Wrong delete key or filename!")
			return
		}
	}
	c.String(http.StatusUnauthorized, "That key doesn't exist!")
}
