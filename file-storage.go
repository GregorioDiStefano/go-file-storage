package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"

	"./helpers"
)

const (
	dbFilename = "files.db"
	bucket     = "files"
)

type StoredFile struct {
	Key          string
	FileName     string
	FileSize     int64
	DeleteKey    string
	MaxDownloads int64
	Downloads    int64
	UploadTime   time.Time
}

func fileDownloader(c *gin.Context) {
	expectedFilePath := fmt.Sprintf("%s/%s/%s",
		helpers.Config.StorageFolder,
		c.Param("key"),
		c.Param("filename"))

	if _, err := os.Stat(expectedFilePath); os.IsNotExist(err) {
		fmt.Print(expectedFilePath + " does not exist.")
		c.String(http.StatusForbidden, "Doesn't look like that file exists.")
		return
	}
	c.File(expectedFilePath)
}

func init() {
	helpers.ParseConfig()
}

func main() {
	router := gin.Default()

	router.POST("/", upload)
	router.PUT("/:filename", simpleUpload)
	router.GET("/:key/:filename", fileDownloader)
	router.Run(helpers.Config.ServerPort)
}
