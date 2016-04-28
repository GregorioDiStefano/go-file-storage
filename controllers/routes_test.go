package controller

import (
	"../helpers"
	"../models"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

var r *gin.Engine

func performRequest(r http.Handler, method, path string, data string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, strings.NewReader(data))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestMain(m *testing.M) {
	db := models.Database{Filename: models.DbFilename, Bucket: models.Bucket}
	helpers.ParseConfig()
	db.OpenDatabaseFile()

	r = gin.Default()

	fmt.Println("r : ", r)
	r.PUT("/:filename", func(c *gin.Context) {
		SimpleUpload(c)
	})

	r.DELETE("/:key/:delete_key/:filename", func(c *gin.Context) {
		DeleteFile(c)
	})

	r.GET("/:key/:filename", func(c *gin.Context) {
		DownloadFile(c)
	})
	os.Exit(m.Run())
}

func TestSimpleUpload(t *testing.T) {

	fmt.Println("r : ", r)
	var dl map[string]string

	superLongString := strings.Repeat("A", 1024*1024*10)

	w := performRequest(r, "PUT", "/test", superLongString)

	if err := json.Unmarshal(w.Body.Bytes(), &dl); err != nil {
		panic(err)
		t.Fail()
	}

	downloadURL := dl["downloadURL"]
	deleteURL := dl["deleteURL"]

	for i := int64(0); i <= helpers.Config.MaxDownloadsBeforeInteraction; i++ {
		w = performRequest(r, "GET", downloadURL, "")
		assert.Equal(t, string(w.Body.Bytes()), superLongString)
	}

	w = performRequest(r, "GET", downloadURL, "")
	assert.Equal(t, w.Code, http.StatusForbidden)

	deletePath := strings.Replace(deleteURL, helpers.Config.Domain, "", 1)
	w = performRequest(r, "DELETE", deletePath, "TestData")

	w = performRequest(r, "GET", downloadURL, "")
	assert.Equal(t, w.Code, http.StatusNotFound)

}
