package controller

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/GregorioDiStefano/go-file-storage/models"
	"github.com/GregorioDiStefano/go-file-storage/utils"
	"github.com/gin-gonic/gin"
)

type Uploader interface {
	upload(io.ReadCloser, string, string)
}

func checkUploadSize(size uint64) error {
	maxFileSize := uint64(utils.Config.GetInt("max_file_size"))
	if size > maxFileSize || size <= 0 {
		fmt.Printf("File upload was :%d, while max size allowed is: %d\n", maxFileSize)
		return errors.New("File too large")
	}
	return nil
}

func Upload(c *gin.Context) {

	s3upload := S3Upload{}

	uploadFileSize := c.Request.ContentLength

	if err := checkUploadSize(uint64(uploadFileSize)); err != nil {
		sendError(c, "Upload size either too larger or invalid")
		return
	}

	fn := c.Param("filename")
	key := models.DB.FindUnusedKey()
	deleteKey := utils.RandomString(uint8(utils.Config.GetInt("delete_key_size")))

	if err := s3upload.upload(c.Request.Body, key, fn); err != nil {
		utils.Log.Fatalln("Uploading file to S3 bucket failed.")
		sendError(c, "Uploading file to S3 bucket failed!")
		return
	}

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
	returnJSON["downloadURL"] = fmt.Sprintf("%s/%s/%s", utils.Config.GetString("domain"), key, fn)
	returnJSON["deleteURL"] = fmt.Sprintf("%s/%s/%s/%s", utils.Config.GetString("domain"), key, deleteKey, fn)
	utils.Log.WithFields(log.Fields{"key": key, "deleteKey": deleteKey, "fn": fn}).Infoln("Upload successful.")
	c.JSON(http.StatusCreated, returnJSON)
}
