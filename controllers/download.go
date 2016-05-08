package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/GregorioDiStefano/go-file-storage/helpers"
	"github.com/GregorioDiStefano/go-file-storage/models"
	"github.com/gin-gonic/gin"
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

	key := c.Param("key")
	fn := c.Param("filename")
	googleCaptchaCode := c.Query("g-recaptcha-response")

	if models.DB.DoesKeyExist(key) == false {
		sendError(c, "Invalid file key")
		return
	}

	expectedFilePath := fmt.Sprintf("%s/%s/%s",
		helpers.Config.StorageFolder,
		key,
		fn)

	sf := models.DB.ReadStoredFile(key)
	if _, err := os.Stat(expectedFilePath); os.IsNotExist(err) && sf.StorageMethod == LOCAL {
		sendError(c, "File does not exist")
		return
	}

	if sf.Deleted {
		sendError(c, "Sorry, this file has been deleted")
	}

	sf.LastAccess = time.Now().UTC()

	if sf.Downloads >= helpers.Config.MaxDownloadsBeforeInteraction && !checkCaptcha(googleCaptchaCode) {
		if helpers.IsWebBrowser(c.Request.Header.Get("User-Agent")) {
			c.HTML(http.StatusOK, "download.tmpl", gin.H{
				"filename": fn,
				"url":      "http://cloudfrontURL",
			})
			return
		}

		sendError(c, "This file has been download too many times, visit the URL with a Browser")
		return
	}

	sf.Downloads = sf.Downloads + 1
	models.DB.WriteStoredFile(*sf)

	if sf.StorageMethod == S3 {
		c.Redirect(http.StatusTemporaryRedirect, helpers.GetS3SignedURL(sf.Key, sf.FileName))
		return
	} else if sf.StorageMethod == LOCAL {
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", sf.FileName))
		c.File(expectedFilePath)
		return
	}
	sendError(c, "Error locating the file you requested!")
}
