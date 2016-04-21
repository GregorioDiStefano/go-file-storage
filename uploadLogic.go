package main

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"

	"./helpers"
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

func processUpload(data interface{}, key string, fn string) {
	directoryToCreate := fmt.Sprintf("%s/%s/", helpers.Config.StorageFolder, key)
	fileToCreate := fmt.Sprintf("%s/%s/%s", helpers.Config.StorageFolder, key, fn)

	os.Mkdir(directoryToCreate, 0777)
	f, err := os.OpenFile(fileToCreate, os.O_CREATE|os.O_WRONLY, 0777)
	defer f.Close()

	if err != nil {
		panic(err)
	}

	for {
		tmp := make([]byte, 512*1024)
		var count int

		switch data.(type) {
		case multipart.File:
			count, _ = data.(multipart.File).Read(tmp)
		case io.ReadCloser:
			count, _ = data.(io.ReadCloser).Read(tmp)
		}

		if count > 0 {
			f.Write(tmp[0:count])
		} else {
			break
		}
	}
}
