package controller

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/GregorioDiStefano/go-file-storage/helpers"
	"github.com/GregorioDiStefano/go-file-storage/models"
	"github.com/gin-gonic/gin"
)

func deleteS3File(key, filename string) error {
	//TODO: implement this later
	return nil
}

func deleteLocalFile(filePath string) error {
	if _, err := os.Stat(filePath); err == nil {
		os.Remove(filePath)
		return nil
	} else {
		return errors.New("Failed to delete file")
	}
}

func DeleteFile(c *gin.Context) {
	key := c.Param("key")
	deleteKey := c.Param("delete_key")
	fileName := c.Param("filename")

	if models.DB.DoesKeyExist(key) {
		sf := models.DB.ReadStoredFile(key)

		if sf.Deleted {
			sendError(c, "File is already deleted.")
			return
		}

		if sf.DeleteKey != deleteKey || sf.FileName != fileName {
			sendError(c, "Delete key or filename was incorrect")
			return
		}

		var deleteErr error
		if sf.StorageMethod == LOCAL {
			filePath := fmt.Sprintf("%s/%s/%s",
				helpers.Config.StorageFolder,
				key,
				sf.FileName)
			deleteErr = deleteLocalFile(filePath)
		} else if sf.StorageMethod == S3 {
			deleteErr = deleteS3File(key, fileName)
		}

		if deleteErr != nil {
			sendError(c, "Failed to delete file")
			return
		}

		sf.Deleted = true
		models.DB.WriteStoredFile(*sf)
		c.JSON(http.StatusOK, map[string]string{"msg": "File deleted"})
		return
	}
	sendError(c, "Something went wrong deleting the requested file")
}
