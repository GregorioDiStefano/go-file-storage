package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/GregorioDiStefano/go-file-storage/helpers"
	"github.com/GregorioDiStefano/go-file-storage/models"
	"github.com/gin-gonic/gin"
)

func checkCaptcha(gRecaptchaResponse string) bool {
	secret := helpers.Config.CAPTCHA_SECRET
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

	helpers.Log.WithFields(log.Fields{"key": key, "fn": fn}).Info("Incoming download.")

	var expectedFilePath string
	sf := models.DB.ReadStoredFile(key)

	if !models.DB.DoesKeyExist(key) || sf == nil || sf.Deleted || sf.FileName != fn {
		if sf == nil {
			helpers.Log.WithFields(log.Fields{"key": key, "fn": fn}).Error("Download failed since key doesn't exist in database")
		} else {
			helpers.Log.WithFields(log.Fields{"key": key, "fn": fn, "delete": sf.Deleted}).Error("Download failed")
		}
		sendError(c, "Invalid filename, key, or file is deleted")
		return
	}

	if sf.StorageMethod == LOCAL {
		expectedFilePath = fmt.Sprintf("%s/%s/%s",
			helpers.Config.StorageFolder,
			key,
			fn)

		if _, err := os.Stat(expectedFilePath); os.IsNotExist(err) {
			sendError(c, "File does not exist")
			return
		}
	}

	sf.LastAccess = time.Now().UTC()

	//if file has been download too many times, show this page, with requires a captcha to be solved
	if sf.Downloads >= helpers.Config.MaxDownloadsBeforeInteraction && !checkCaptcha(googleCaptchaCode) {
		if helpers.IsWebBrowser(c.Request.Header.Get("User-Agent")) {
			c.HTML(http.StatusOK, "download.tmpl", gin.H{
				"filename": fn,
			})
			return
		}

		sendError(c, "This file has been download too many times, visit the URL with a Browser")
		return
	}

	sf.Downloads = sf.Downloads + 1
	helpers.Log.WithFields(log.Fields{"key": key, "filename": fn}).Info("Downloads set to ", sf.Downloads)
	models.DB.WriteStoredFile(*sf)

	if sf.StorageMethod == S3 {
		ip := helpers.GetXFF(c.Request.Header)
		c.Redirect(http.StatusMovedPermanently, helpers.GetS3SignedURL(sf.Key, sf.FileName, ip))
		return
	} else if sf.StorageMethod == LOCAL {
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", sf.FileName))
		c.File(expectedFilePath)
		return
	}
}
