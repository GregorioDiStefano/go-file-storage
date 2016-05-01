package controller

import (
	"../helpers"
	"../models"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func SimpleUpload(c *gin.Context) {

	if _, err := checkUploadSize(c); err != nil {
		c.String(http.StatusForbidden, "Upload too large")
		return
	}

	db := models.Database{Filename: models.DbFilename, Bucket: models.Bucket}
	fn := c.Param("filename")
	key := db.FindUnsedKey()
	deleteKey := helpers.RandomString(helpers.Config.DeleteKeySize)

	processUpload(c.Request.Body, key, fn)

	simpleStoredFiled := models.StoredFile{
		Key:        key,
		DeleteKey:  deleteKey,
		FileName:   fn,
		FileSize:   c.Request.ContentLength,
		UploadTime: time.Now().UTC()}

	db.WriteStoredFile(simpleStoredFiled)

	returnJSON := make(map[string]string)
	returnJSON["downloadURL"] = fmt.Sprintf("%s/%s/%s", helpers.Config.Domain, key, fn)
	returnJSON["deleteURL"] = fmt.Sprintf("%s/%s/%s/%s", helpers.Config.Domain, key, deleteKey, fn)

	jsonString, _ := json.Marshal(returnJSON)

	c.String(http.StatusOK, string(jsonString))
}
