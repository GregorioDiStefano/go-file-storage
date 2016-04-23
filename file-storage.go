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
	key := c.Param("key")
	fn := c.Param("filename")
	expectedFilePath := fmt.Sprintf("%s/%s/%s",
		helpers.Config.StorageFolder,
		key,
		fn)

	if _, err := os.Stat(expectedFilePath); os.IsNotExist(err) {
		fmt.Print(expectedFilePath + " does not exist.")
		c.String(http.StatusForbidden, "Doesn't look like that file exists.")
		return
	}
	db := database{filename: dbFilename, bucket: bucket}
	sf := db.readStoredFile(key)
    sf.Downloads = sf.Downloads + 1
    db.writeStoredFile(*sf)
	fmt.Println(sf)

	c.File(expectedFilePath)
}

func init() {
	helpers.ParseConfig()
}

func main() {
	router := gin.Default()

	router.PUT("/:filename", simpleUpload)
	router.GET("/:key/:filename", fileDownloader)
	router.Run(helpers.Config.ServerPort)
}
