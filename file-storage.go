package main

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
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
	Filename     string
	FileSize     int64
	DeleteKey    string
	MaxDownloads int64
	Downloads    int64
	UploadTime   time.Time
}

func processUpload(file multipart.File, key string, fn string) {
	os.Mkdir("files/"+key, 0777)
	f, err := os.OpenFile("files/"+key+"/"+fn, os.O_CREATE|os.O_WRONLY, 0777)
	defer f.Close()

	if err != nil {
		panic(err)
	}

	for {
		tmp := make([]byte, 512*1024)
		count, _ := file.Read(tmp)

		if count > 0 {
			f.Write(tmp[0:count])
		} else {
			break
		}
	}
}

//func upload(rw http.ResponseWriter, req *http.Request) {
func upload(c *gin.Context) {
	FileSize, _ := strconv.ParseInt(c.Request.Header.Get("Content-Length"),
		10,
		64)

	if FileSize > helpers.Config.MaxSize {
		fmt.Printf("File upload was :%d, while max size allowed is: %d\n", FileSize, helpers.Config.MaxSize)
		c.String(http.StatusForbidden, helpers.Config.OverMaxSizeStr)
		return
	}

	file, headers, err := c.Request.FormFile("fileupload")
	DeleteKey := c.Request.FormValue("DeleteKey")
	MaxDownloads, _ := strconv.ParseInt(c.Request.FormValue("MaxDownloads"), 10, 64)
	Downloads := int64(0)
	UploadTime := time.Now().UTC()
	Key := helpers.RandomString(helpers.Config.KeySize)
	FileName := headers.Filename

	if file != nil {
		processUpload(file, Key, FileName)
	}

	if err != nil {
		panic(err)
	}

	sf := StoredFile{Key,
		FileName,
		FileSize,
		DeleteKey,
		MaxDownloads,
		Downloads,
		UploadTime}

	db := database{filename: dbFilename, bucket: bucket}
	fmt.Println(sf)
	db.writeStoredFile(sf)

	c.String(http.StatusOK, Key)
}

func init() {
	helpers.ParseConfig()
}

func main() {
	router := gin.Default()

	router.POST("/upload", upload)
	router.Run(":8080")
}
