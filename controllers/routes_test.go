package controller

import (
	"../helpers"
	"../models"
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

var r *gin.Engine

type uploadedFileInfo struct {
	payload     []byte
	downloadURL string
	deleteURL   string
}

func performRequest(r http.Handler, method, path string, data []byte) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, bytes.NewReader(data))
	cl := strconv.Itoa(len(data))
	req.Header.Add("Content-Length", cl)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func downloadPathToMap(path string) map[string]string {
	splitStr := strings.Split(path, "/")

	key := splitStr[1]
	filename := splitStr[2]

	return map[string]string{
		"key":      key,
		"filename": filename,
	}
}

func deletePathToMap(path string) map[string]string {
	splitStr := strings.Split(path, "/")

	key := splitStr[1]
	deleteKey := splitStr[2]
	filename := splitStr[3]

	return map[string]string{
		"key":       key,
		"deleteKey": deleteKey,
		"filename":  filename,
	}
}

func TestMain(m *testing.M) {
	db := models.Database{Filename: models.DbFilename, Bucket: models.Bucket}
	helpers.ParseConfig("config/config.testing.json")
	db.OpenDatabaseFile()

	r = gin.Default()

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

func uploadFile(data []byte, filename string) uploadedFileInfo {
	var dl map[string]string

	w := performRequest(r, "PUT", "/"+filename, data)

	if err := json.Unmarshal(w.Body.Bytes(), &dl); err != nil {
		panic(err)
	}

	return uploadedFileInfo{data, dl["downloadURL"], dl["deleteURL"]}
}

func TestUpload_1byte(t *testing.T) {
	data := []byte{0x01}
	uf := uploadFile(data, "test123")

	w := performRequest(r, "GET", uf.downloadURL, nil)
	assert.Equal(t, uf.payload, w.Body.Bytes())
}

func TestUpload_10MB(t *testing.T) {
	data := []byte{0x01}
	uf := uploadFile(bytes.Repeat(data, 10000000), "test1234")

	w := performRequest(r, "GET", uf.downloadURL, nil)
	assert.Equal(t, uf.payload, w.Body.Bytes())
}

func TestUploadEqualMaxSize(t *testing.T) {
	maxUploadSize := helpers.Config.MaxSize

	if maxUploadSize >= (1 * 1024 * 1024 * 1024) {
		t.Skip("maxUploadSize is huge, skipping")
	}
	randomData := make([]byte, maxUploadSize)
	_, err := rand.Read(randomData)

	if err != nil {
		panic(err)
	}

	w := performRequest(r, "PUT", "/hugepass", randomData)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUploadExceedingMaxSize(t *testing.T) {
	maxUploadSize := helpers.Config.MaxSize

	if maxUploadSize >= (1 * 1024 * 1024 * 1024) {
		t.Skip("maxUploadSize is huge, skipping")
	}
	randomData := make([]byte, maxUploadSize+1)
	_, err := rand.Read(randomData)

	if err != nil {
		panic(err)
	}

	w := performRequest(r, "PUT", "/hugefail", randomData)
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestInvalidDownload_1(t *testing.T) {
	ufl := uploadFile(bytes.Repeat([]byte{0xff}, 1024*1024), "test")
	downloadURL := downloadPathToMap(strings.Replace(ufl.downloadURL, helpers.Config.Domain, "", 1))

	key := downloadURL["key"]
	invalidFilename := downloadURL["filename"] + ".x"

	invalidDownloadURL := fmt.Sprintf("%s/%s/%s", helpers.Config.Domain, key, invalidFilename)
	w := performRequest(r, "GET", invalidDownloadURL, nil)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestInvalidDownload_2(t *testing.T) {
	ufl := uploadFile(bytes.Repeat([]byte{0x00, 0x11}, 1024*1024), "test")
	downloadURL := downloadPathToMap(strings.Replace(ufl.downloadURL, helpers.Config.Domain, "", 1))

	key := downloadURL["key"] + ".x"
	invalidFilename := downloadURL["filename"]

	invalidDownloadURL := fmt.Sprintf("%s/%s/%s", helpers.Config.Domain, key, invalidFilename)
	w := performRequest(r, "GET", invalidDownloadURL, nil)
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestMaxDownloads(t *testing.T) {
	db := models.Database{Filename: models.DbFilename, Bucket: models.Bucket}
	ufl := uploadFile(bytes.Repeat([]byte{0x11, 0x22, 0x33, 0xff}, 1024), "test")

	for i := int64(1); i <= helpers.Config.MaxDownloadsBeforeInteraction; i++ {
		w := performRequest(r, "GET", ufl.downloadURL, nil)
		assert.Equal(t, ufl.payload, w.Body.Bytes())
		assert.Equal(t, http.StatusOK, w.Code)
	}

	downloadURLPath := downloadPathToMap(strings.Replace(ufl.downloadURL, helpers.Config.Domain, "", 1))
	downloadsStoredInDB := db.ReadStoredFile(downloadURLPath["key"]).Downloads
	assert.Equal(t, helpers.Config.MaxDownloadsBeforeInteraction, downloadsStoredInDB)

	w := performRequest(r, "GET", ufl.downloadURL, nil)
	downloadPathToMap(strings.Replace(ufl.downloadURL, helpers.Config.Domain, "", 1))
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestDeleteFile(t *testing.T) {
	ufl := uploadFile(bytes.Repeat([]byte{0x11, 0x22, 0x33, 0xff}, 1024), "test")

	w := performRequest(r, "GET", ufl.downloadURL, nil)
	assert.Equal(t, http.StatusOK, w.Code)

	w = performRequest(r, "DELETE", ufl.deleteURL, nil)
	assert.Equal(t, http.StatusOK, w.Code)

	w = performRequest(r, "DELETE", ufl.deleteURL, nil)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDownloadLastAccess(t *testing.T) {
	db := models.Database{Filename: models.DbFilename, Bucket: models.Bucket}
	ufl := uploadFile(bytes.Repeat([]byte{0xff}, 1024), "test")
	downloadURLPath := downloadPathToMap(strings.Replace(ufl.downloadURL, helpers.Config.Domain, "", 1))
	performRequest(r, "GET", ufl.downloadURL, nil)
	expectedTime := db.ReadStoredFile(downloadURLPath["key"]).LastAccess.Unix()

	assert.Equal(t, expectedTime, time.Now().Unix())
}

func TestUploadDownloadUnicode(t *testing.T) {
	ufl := uploadFile(bytes.Repeat([]byte{0xff}, 1024), "Шဩ.ext")
	assert.True(t, strings.Contains(ufl.downloadURL, "Шဩ.ext"))

	w := performRequest(r, "GET", ufl.downloadURL, nil)
	assert.Equal(t, ufl.payload, w.Body.Bytes())
}

func TestDeleteInvalid_1(t *testing.T) {
	ufl := uploadFile(bytes.Repeat([]byte{0x11, 0x22, 0x33}, 1024), "test")
	deleteURL := deletePathToMap(strings.Replace(ufl.deleteURL, helpers.Config.Domain, "", 1))

	key := deleteURL["key"]
	deleteKey := deleteURL["deleteKey"]
	invalidFilename := deleteURL["filename"] + ".x"

	invalidDeleteURL := fmt.Sprintf("%s/%s/%s/%s", helpers.Config.Domain, key, deleteKey, invalidFilename)
	w := performRequest(r, "DELETE", invalidDeleteURL, nil)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestDeleteInvalid_2(t *testing.T) {
	ufl := uploadFile(bytes.Repeat([]byte{0x11, 0x22, 0x33}, 1024), "test")
	deleteURL := deletePathToMap(strings.Replace(ufl.deleteURL, helpers.Config.Domain, "", 1))

	key := deleteURL["key"]
	invalidDeleteKey := deleteURL["deleteKey"] + ".x"
	filename := deleteURL["filename"]

	invalidDeleteURL := fmt.Sprintf("%s/%s/%s/%s", helpers.Config.Domain, key, invalidDeleteKey, filename)
	w := performRequest(r, "DELETE", invalidDeleteURL, nil)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestDeleteInvalid_3(t *testing.T) {
	ufl := uploadFile(bytes.Repeat([]byte{0x11, 0x22, 0x33}, 1024), "test")
	deleteURL := deletePathToMap(strings.Replace(ufl.deleteURL, helpers.Config.Domain, "", 1))

	invalidKey := deleteURL["key"] + ".x"
	deleteKey := deleteURL["deleteKey"]
	filename := deleteURL["filename"]

	invalidDeleteURL := fmt.Sprintf("%s/%s/%s/%s", helpers.Config.Domain, invalidKey, deleteKey, filename)
	w := performRequest(r, "DELETE", invalidDeleteURL, nil)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
