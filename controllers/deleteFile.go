package controller

import (
	"net/http"

	log "github.com/Sirupsen/logrus"

	"github.com/GregorioDiStefano/go-file-storage/models"
	"github.com/gin-gonic/gin"
)

func deleteS3File(key, filename string) error {
	//TODO: implement this later
	return nil
}

func DeleteFile(c *gin.Context) {
	key := c.Param("key")
	deleteKey := c.Param("delete_key")
	fileName := c.Param("filename")

	if models.DB.DoesKeyExist(key) {
		sf := models.DB.ReadStoredFile(key)

		if sf.Deleted {
			log.WithFields(log.Fields{"key": key, "delete_key": deleteKey, "fn": fileName}).Infoln("File is already deleted")
			sendError(c, "File is already deleted.")
			return
		}

		if sf.DeleteKey != deleteKey || sf.FileName != fileName {
			log.WithFields(log.Fields{"key": key, "delete_key": deleteKey, "fn": fileName}).Infoln("Delete key or filename was incorrect")
			sendError(c, "Delete key or filename was incorrect")
			return
		}

		if err := deleteS3File(key, fileName); err != nil {
			//log.WithFields(log.Fields{"key" : key, "delete_key": deleteKey, "fn": fileName}).Infoln("Delete key or filename was incorrect")
			sendError(c, "Failed to delete file")
			return
		}

		sf.Deleted = true
		models.DB.WriteStoredFile(*sf)
		log.WithFields(log.Fields{"key": key, "delete_key": deleteKey, "fn": fileName}).Infoln("File set to deleted.")
		c.JSON(http.StatusOK, map[string]string{"msg": "File deleted"})
		return
	}

	log.WithFields(log.Fields{"key": key, "delete_key": deleteKey, "fn": fileName}).Infoln("Failed to delete file.")
	sendError(c, "Something went wrong deleting the requested file")
}
