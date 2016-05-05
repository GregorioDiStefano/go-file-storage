package controller

import (
	"fmt"
	"net/http"
	"os"

	"github.com/GregorioDiStefano/go-file-storage/helpers"
	"github.com/GregorioDiStefano/go-file-storage/models"
	"github.com/gin-gonic/gin"
)

func DeleteFile(c *gin.Context) {
	key := c.Param("key")
	deleteKey := c.Param("delete_key")
	fileName := c.Param("filename")

	if models.DB.DoesKeyExist(key) {
		sf := models.DB.ReadStoredFile(key)

		if sf.DeleteKey == deleteKey && sf.FileName == fileName {

			filePath := fmt.Sprintf("%s/%s/%s",
				helpers.Config.StorageFolder,
				key,
				sf.FileName)

			if _, err := os.Stat(filePath); err == nil {
				os.Remove(filePath)
				c.String(http.StatusOK, "Deleted file: "+sf.FileName+"\n")
				return
			}
		}
	}
	sendError(c, "That key or filename doesn't exist")
}
