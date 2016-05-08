package controller

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/GregorioDiStefano/go-file-storage/helpers"
	"github.com/gin-gonic/gin"
	"github.com/rlmcpherson/s3gof3r"
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
	S3Keys := s3gof3r.Keys{AccessKey: helpers.Config.AccessKey,
		SecretKey: helpers.Config.SecretKey}

	// Open bucket to put file into
	s3 := s3gof3r.New("s3-eu-west-1.amazonaws.com", S3Keys)
	b := s3.Bucket("greg-filestorage")

	w, err := b.PutWriter(fmt.Sprintf("%s/%s", key, fn), nil, nil)

	if err != nil {
		fmt.Println("1: ", err)
		return err
	}

	// Copy into S3
	if _, err = io.Copy(w, data); err != nil {
		fmt.Println(err)
		return err
	}

	if err = w.Close(); err != nil {
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
