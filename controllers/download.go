package controller

import (
	"../helpers"
	"../models"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"
)

func checkCaptcha(gRecaptchaResponse string) bool {

	secret := helpers.Config.CaptchaSecret
	response := gRecaptchaResponse

	postURL := "https://www.google.com/recaptcha/api/siteverify"

	resp, _ := http.PostForm(postURL,
		url.Values{"secret": {secret}, "response": {response}})

	if respPayload, err := ioutil.ReadAll(resp.Body); err == nil {
		jsonMap := make(map[string]interface{})
		e := json.Unmarshal(respPayload, &jsonMap)

		success, ok := jsonMap["success"]

		if e != nil || !ok {
			return false
		} else if success == true {
			return true
		}
	}

	return false
}

func DownloadFile(c *gin.Context) {
	db := models.Database{Filename: models.DbFilename, Bucket: models.Bucket}
	key := c.Param("key")
	fn := c.Param("filename")
	googleCaptchaCode := c.Query("g-recaptcha-response")

	if db.DoesKeyExist(key) == false {
		c.String(http.StatusForbidden, "Invalid key")
		return
	}

	expectedFilePath := fmt.Sprintf("%s/%s/%s",
		helpers.Config.StorageFolder,
		key,
		fn)

	if _, err := os.Stat(expectedFilePath); os.IsNotExist(err) {
		fmt.Print(expectedFilePath + " does not exist.")
		c.String(http.StatusNotFound, "Doesn't look like that file exists.")
		return
	}

	sf := db.ReadStoredFile(key)
	sf.LastAccess = time.Now().UTC()

	if sf.Downloads >= helpers.Config.MaxDownloadsBeforeInteraction && !checkCaptcha(googleCaptchaCode) {
		if helpers.IsWebBrowser(c.Request.Header.Get("User-Agent")) {
			c.HTML(http.StatusOK, "download.tmpl", gin.H{
				"filename": fn,
			})
			return
		}
		c.String(http.StatusForbidden, "This file has been download too many times.")
		return
	}

	sf.Downloads = sf.Downloads + 1
	db.WriteStoredFile(*sf)

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", sf.FileName))
	c.File(expectedFilePath)
}
