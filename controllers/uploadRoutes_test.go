package controller

import (
	"../helpers"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func performRequest(r http.Handler, method, path string, data []byte) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestSimpleUpload(t *testing.T) {
	helpers.ParseConfig()

	r := gin.Default()
	r.PUT("/:filename", func(c *gin.Context) {
		SimpleUpload(c)
	})

	w := performRequest(r, "PUT", "/test", []byte("TestData"))

	fmt.Println(w.Body)
}
