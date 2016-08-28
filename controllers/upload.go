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

type GenericUploader interface {
	UploadFile(c *gin.Context)
	doActualUpload(data io.ReadCloser, key string, fn string) error

	doActualDelete(deleteKey, fileKey, filename string) error
}

type Upload struct {
	awsBucket     string
	awsRegion     string
	maxUploadSize int64
	deleteKeySize int

	uploadDomain string
}

func NewUploader(uploadDomain string, maxUploadSize int64, deleteKeySize int, awsBucket, awsRegion string) *Upload {
	return &Upload{uploadDomain: uploadDomain,
		maxUploadSize: maxUploadSize,
		deleteKeySize: deleteKeySize,
		awsBucket:     awsBucket,
		awsRegion:     awsRegion}
}

func checkUploadSize(actualSize, maxFileSize int64) error {
	if actualSize > maxFileSize || actualSize <= 0 {
		fmt.Printf("File upload was :%d, while max size allowed is: %d\n", actualSize, maxFileSize)
		return errors.New("File too large")
	}
	return nil
}

func (upload Upload) UploadFile(c *gin.Context) {
	uploadFileSize := c.Request.ContentLength

	if err := checkUploadSize(uploadFileSize, upload.maxUploadSize); err != nil {
		sendError(c, "Upload size either too larger or invalid")
		return
	}

	fn := c.Param("filename")
	key := models.DB.FindUnusedKey()
	deleteKey := utils.RandomString(utils.Config.GetInt("delete_key_size"))

	fmt.Print("do we get here")

	if err := upload.doActualUpload(c.Request.Body, key, fn); err != nil {
		utils.Log.Fatalln("Uploading file to S3 bucket failed.")
		sendError(c, "Uploading file to S3 bucket failed!")
		fmt.Println("sent error")
		return
	}

	fmt.Print("do we get here")

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
	returnJSON["downloadURL"] = fmt.Sprintf("%s/%s/%s", upload.uploadDomain, key, fn)
	returnJSON["deleteURL"] = fmt.Sprintf("%s/%s/%s/%s", upload.uploadDomain, key, deleteKey, fn)
	utils.Log.WithFields(log.Fields{"key": key, "deleteKey": deleteKey, "fn": fn}).Infoln("Upload successful.")
	c.JSON(http.StatusCreated, returnJSON)
}
