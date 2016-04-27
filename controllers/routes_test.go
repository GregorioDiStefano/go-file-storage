package controller

import (
	"../helpers"
	"../models"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

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
	os.Exit(m.Run())
}

func TestSimpleUpload(t *testing.T) {

	r := gin.Default()

	r.PUT("/:filename", func(c *gin.Context) {
		SimpleUpload(c)
	})

	r.DELETE("/:key/:delete_key", func(c *gin.Context) {
		DeleteFile(c)
	})

	var dl map[string]string
	w := performRequest(r, "PUT", "/test", "TestData")
	fmt.Println(string(w.Body.Bytes()))
	if err := json.Unmarshal(w.Body.Bytes(), &dl); err != nil {
		panic(err)
		t.Fail()
	}

	deleteURL := dl["deleteURL"]
	downloadURL := dl["downloadURL"]

	deletePath := strings.Replace(deleteURL, helpers.Config.Domain, "", 1)

	w = performRequest(r, "DELETE", deletePath, "TestData")
	fmt.Println(string(w.Body.Bytes()))

	os.Setenv("DOWNLOAD_URL", downloadURL)
	os.Setenv("DELETE_URL", deleteURL)
}
