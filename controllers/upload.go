package controller

import (
	"fmt"
	"net/http"
	"time"

	"github.com/GregorioDiStefano/go-file-storage/helpers"
	"github.com/GregorioDiStefano/go-file-storage/models"
	"github.com/gin-gonic/gin"
)

func SimpleUpload(c *gin.Context) {

	if _, err := checkUploadSize(c); err != nil {
		sendError(c, "Upload too large")
		return
	}

	fn := c.Param("filename")
	key := models.DB.FindUnsedKey()
	deleteKey := helpers.RandomString(helpers.Config.DeleteKeySize)

	processUpload(c.Request.Body, key, fn)

	simpleStoredFiled := models.StoredFile{
		Key:        key,
		DeleteKey:  deleteKey,
		FileName:   fn,
		FileSize:   c.Request.ContentLength,
		UploadTime: time.Now().UTC(),
		LastAccess: time.Now().UTC(),
	}

	models.DB.WriteStoredFile(simpleStoredFiled)

	returnJSON := make(map[string]string)
	returnJSON["downloadURL"] = fmt.Sprintf("%s/%s/%s", helpers.Config.Domain, key, fn)
	returnJSON["deleteURL"] = fmt.Sprintf("%s/%s/%s/%s", helpers.Config.Domain, key, deleteKey, fn)

	c.JSON(http.StatusOK, returnJSON)
}
