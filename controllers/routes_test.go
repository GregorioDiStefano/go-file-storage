package controller

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/GregorioDiStefano/go-file-storage/models"
	"github.com/GregorioDiStefano/go-file-storage/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var r *gin.Engine

type uploadedFileInfo struct {
	payload     []byte
	downloadURL string
	deleteURL   string
}

func HTTPResponseToBytes(resp http.Response) []byte {
	b, _ := ioutil.ReadAll(resp.Body)
	return b
}

func performRequest(r http.Handler, method, path string, data []byte) ([]byte, http.Header, int) {
	req, _ := http.NewRequest(method, path, bytes.NewReader(data))
	req.Header.Add("Content-Length", strconv.Itoa(len(data)))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code == http.StatusMovedPermanently {
		finalURL := w.Header().Get("Location")
		resp, err := http.Get(finalURL)
		if err != nil {
			panic(err)
		}
		respBytes, _ := ioutil.ReadAll(resp.Body)
		return respBytes, resp.Header, resp.StatusCode
	}

	return w.Body.Bytes(), w.Header(), w.Code
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

	if err := os.Chdir("../"); err != nil {
		panic(err)
	}

	utils.ParseConfig("config/config.testing.yaml")
	models.DB.OpenDatabaseFile()

	r = gin.Default()

	r.LoadHTMLGlob("templates/*")

	r.PUT("/:filename", func(c *gin.Context) {
		Upload(c)
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

	respBytes, _, _ := performRequest(r, "PUT", "/"+filename, data)

	if err := json.Unmarshal(respBytes, &dl); err != nil {
		panic(err)
	}

	return uploadedFileInfo{data, dl["downloadURL"], dl["deleteURL"]}
}

func TestUpload_1byte(t *testing.T) {
	data := []byte{0x01}
	uf := uploadFile(data, "test123")

	respBytes, _, statusCode := performRequest(r, "GET", uf.downloadURL, nil)

	assert.Equal(t, uf.payload, respBytes)
	assert.Equal(t, http.StatusOK, statusCode)
}

func TestUpload_10MB(t *testing.T) {
	data := []byte{0x01}
	uf := uploadFile(bytes.Repeat(data, 10000000), "test1234")

	respBytes, _, _ := performRequest(r, "GET", uf.downloadURL, nil)
	assert.Equal(t, uf.payload, respBytes)
}

func TestUploadEqualMaxSize(t *testing.T) {
	maxUploadSize := utils.Config.GetInt("max_file_size")

	if maxUploadSize >= (1 * 1024 * 1024 * 1024) {
		t.Skip("maxUploadSize is huge, skipping")
	}

	randomData := make([]byte, maxUploadSize)
	_, err := rand.Read(randomData)

	if err != nil {
		panic(err)
	}

	_, _, statusCode := performRequest(r, "PUT", "/hugepass", randomData)
	assert.Equal(t, http.StatusCreated, statusCode)
}

func TestUploadExceedingMaxSize(t *testing.T) {
	maxUploadSize := utils.Config.GetInt("max_file_size")

	if maxUploadSize >= (1 * 1024 * 1024 * 1024) {
		t.Skip("maxUploadSize is huge, skipping")
	}
	randomData := make([]byte, maxUploadSize+1)
	_, err := rand.Read(randomData)

	if err != nil {
		panic(err)
	}

	_, _, statusCode := performRequest(r, "PUT", "/hugefail", randomData)
	assert.Equal(t, http.StatusForbidden, statusCode)
}

func TestInvalidDownload_1(t *testing.T) {
	ufl := uploadFile(bytes.Repeat([]byte{0xff}, 1024*1024), "test")
	downloadURL := downloadPathToMap(strings.Replace(ufl.downloadURL, utils.Config.GetString("domain"), "", 1))

	key := downloadURL["key"]
	invalidFilename := downloadURL["filename"] + ".x"

	invalidDownloadURL := fmt.Sprintf("%s/%s/%s", utils.Config.GetString("domain"), key, invalidFilename)

	_, _, statusCode := performRequest(r, "GET", invalidDownloadURL, nil)
	assert.Equal(t, http.StatusForbidden, statusCode)
}

func TestInvalidDownload_2(t *testing.T) {
	ufl := uploadFile(bytes.Repeat([]byte{0x00, 0x11}, 1024*1024), "test")
	downloadURL := downloadPathToMap(strings.Replace(ufl.downloadURL, utils.Config.GetString("domain"), "", 1))

	key := downloadURL["key"] + "foo"
	invalidFilename := downloadURL["filename"]

	invalidDownloadURL := fmt.Sprintf("%s/%s/%s", utils.Config.GetString("domain"), key, invalidFilename)
	_, _, statusCode := performRequest(r, "GET", invalidDownloadURL, nil)
	assert.Equal(t, http.StatusForbidden, statusCode)
}

func TestMaxDownloads(t *testing.T) {
	db := models.Database{Filename: models.DbFilename, Bucket: models.Bucket}
	ufl := uploadFile(bytes.Repeat([]byte{0x11, 0x22, 0x33, 0xff}, 1024), "test")

	for i := int64(1); i <= int64(utils.Config.GetInt("max_downloads")); i++ {
		respBytes, _, statusCode := performRequest(r, "GET", ufl.downloadURL, nil)
		assert.Equal(t, ufl.payload, respBytes)
		assert.Equal(t, http.StatusOK, statusCode)
	}

	downloadURLPath := downloadPathToMap(strings.Replace(ufl.downloadURL, utils.Config.GetString("domain"), "", 1))
	downloadsStoredInDB := db.ReadStoredFile(downloadURLPath["key"]).Downloads
	assert.Equal(t, utils.Config.GetInt("max_downloads"), int(downloadsStoredInDB))

	_, _, statusCode := performRequest(r, "GET", ufl.downloadURL, nil)
	downloadPathToMap(strings.Replace(ufl.downloadURL, utils.Config.GetString("domain"), "", 1))
	assert.Equal(t, http.StatusForbidden, statusCode)
}

func TestDeleteFile(t *testing.T) {
	ufl := uploadFile(bytes.Repeat([]byte{0x11, 0x22, 0x33, 0xff}, 1024), "test")

	_, _, statusCode := performRequest(r, "GET", ufl.downloadURL, nil)
	assert.Equal(t, http.StatusOK, statusCode)

	_, _, statusCode = performRequest(r, "DELETE", ufl.deleteURL, nil)
	assert.Equal(t, http.StatusOK, statusCode)

	_, _, statusCode = performRequest(r, "DELETE", ufl.deleteURL, nil)
	assert.Equal(t, http.StatusForbidden, statusCode)
}

func TestDownloadLastAccess(t *testing.T) {
	ufl := uploadFile(bytes.Repeat([]byte{0xff}, 10), "test")
	downloadURLPath := downloadPathToMap(strings.Replace(ufl.downloadURL, utils.Config.GetString("domain"), "", 1))
	performRequest(r, "GET", ufl.downloadURL, nil)
	expectedTime := models.DB.ReadStoredFile(downloadURLPath["key"]).LastAccess.Unix()
	assert.True(t, expectedTime == time.Now().Unix() || expectedTime+1 == time.Now().Unix())
}

func TestDownloadDeletedFile(t *testing.T) {
	ufl := uploadFile(bytes.Repeat([]byte{0xff}, 1024), "random.ext")
	downloadURLPath := downloadPathToMap(strings.Replace(ufl.downloadURL, utils.Config.GetString("domain"), "", 1))

	storedFile := models.DB.ReadStoredFile(downloadURLPath["key"])
	storedFile.Deleted = true
	models.DB.WriteStoredFile(*storedFile)

	_, _, statusCode := performRequest(r, "GET", ufl.downloadURL, nil)
	assert.Equal(t, http.StatusForbidden, statusCode)
}

func TestUploadDownloadUnicode(t *testing.T) {
	ufl := uploadFile(bytes.Repeat([]byte{0xff}, 1024), "testҖ.ext")
	assert.True(t, strings.Contains(ufl.downloadURL, "testҖ.ext"))

	respBytes, _, _ := performRequest(r, "GET", ufl.downloadURL, nil)
	assert.Equal(t, ufl.payload, respBytes)
}

func TestDeleteInvalid_1(t *testing.T) {
	ufl := uploadFile(bytes.Repeat([]byte{0x11, 0x22, 0x33}, 1024), "test")
	deleteURL := deletePathToMap(strings.Replace(ufl.deleteURL, utils.Config.GetString("domain"), "", 1))

	key := deleteURL["key"]
	deleteKey := deleteURL["deleteKey"]
	invalidFilename := deleteURL["filename"] + ".x"

	invalidDeleteURL := fmt.Sprintf("%s/%s/%s/%s", utils.Config.GetString("domain"), key, deleteKey, invalidFilename)
	_, _, statusCode := performRequest(r, "DELETE", invalidDeleteURL, nil)
	assert.Equal(t, http.StatusForbidden, statusCode)
}

func TestDeleteInvalid_2(t *testing.T) {
	ufl := uploadFile(bytes.Repeat([]byte{0x11, 0x22, 0x33}, 1024), "test")
	deleteURL := deletePathToMap(strings.Replace(ufl.deleteURL, utils.Config.GetString("domain"), "", 1))

	key := deleteURL["key"]
	invalidDeleteKey := deleteURL["deleteKey"] + ".x"
	filename := deleteURL["filename"]

	invalidDeleteURL := fmt.Sprintf("%s/%s/%s/%s", utils.Config.GetString("domain"), key, invalidDeleteKey, filename)

	_, _, statusCode := performRequest(r, "DELETE", invalidDeleteURL, nil)
	assert.Equal(t, http.StatusForbidden, statusCode)
}

func TestDeleteInvalid_3(t *testing.T) {
	ufl := uploadFile(bytes.Repeat([]byte{0x11, 0x22, 0x33}, 1024), "test")
	deleteURL := deletePathToMap(strings.Replace(ufl.deleteURL, utils.Config.GetString("domain"), "", 1))

	invalidKey := deleteURL["key"] + ".x"
	deleteKey := deleteURL["deleteKey"]
	filename := deleteURL["filename"]

	invalidDeleteURL := fmt.Sprintf("%s/%s/%s/%s", utils.Config.GetString("domain"), invalidKey, deleteKey, filename)
	_, _, statusCode := performRequest(r, "DELETE", invalidDeleteURL, nil)
	assert.Equal(t, http.StatusForbidden, statusCode)
}
