package controller

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/GregorioDiStefano/go-file-storage/helpers"
	"github.com/gin-gonic/gin"
)

func checkUploadSize(c *gin.Context) (int64, error) {
	FileSize, _ := strconv.ParseInt(c.Request.Header.Get("Content-Length"),
		10,
		64)

	if FileSize > helpers.Config.MaxSize {
		fmt.Printf("File upload was :%d, while max size allowed is: %d\n",
			FileSize,
			helpers.Config.MaxSize)
		c.String(http.StatusForbidden, helpers.Config.OverMaxSizeStr)
		return FileSize, errors.New("File too large")
	}
	return FileSize, nil
}

func processUploadS3(data io.ReadCloser, key string, fn string) error {
	uploader := s3manager.NewUploader(session.New(&aws.Config{Region: aws.String(helpers.Config.AWSRegion)}))
	_, err := uploader.Upload(&s3manager.UploadInput{
		Body:   data,
		Bucket: aws.String(helpers.Config.S3BucketName),
		Key:    aws.String(fmt.Sprintf("%s/%s", key, fn)),
	})

	if err != nil {
		return err
	}

	return nil
}

func processUpload(data io.ReadCloser, key string, fn string) {
	directoryToCreate := fmt.Sprintf("%s/%s/", helpers.Config.StorageFolder, key)
	fileToCreate := fmt.Sprintf("%s/%s/%s", helpers.Config.StorageFolder, key, fn)

	os.Mkdir(directoryToCreate, 0777)

	f, err := os.OpenFile(fileToCreate, os.O_CREATE|os.O_WRONLY, 0777)
	defer f.Close()

	if err != nil {
		panic(err)
	}

	io.Copy(f, data)
}
