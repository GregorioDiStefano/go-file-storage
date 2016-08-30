package controller

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/GregorioDiStefano/go-file-storage/log"
	"github.com/GregorioDiStefano/go-file-storage/models"
	"github.com/GregorioDiStefano/go-file-storage/utils"
	"github.com/gin-gonic/gin"
)

type Download struct {
	CaptchaSecret string
	MaxDownloads  int
}

func NewDownloader(captchaSecret string, maxDownloads int) *Download {
	return &Download{captchaSecret, maxDownloads}
}

func (download Download) checkCaptcha(gRecaptchaResponse string) bool {
	secret := utils.Config.GetString("CAPTCHA_SECRET")
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

func (download Download) DownloadFile(c *gin.Context) {
	key := c.Param("key")
	fn := c.Param("filename")
	googleCaptchaCode := c.Query("g-recaptcha-response")

	log.WithFields(log.Fields{"key": key, "fn": fn}).Info("Incoming download.")

	sf := models.DB.ReadStoredFile(key)

	if !models.DB.DoesKeyExist(key) || sf == nil || sf.Deleted || sf.FileName != fn {
		if sf == nil {
			log.WithFields(log.Fields{"key": key, "fn": fn}).Error("Download failed since key doesn't exist in database")
		} else {
			log.WithFields(log.Fields{"key": key, "fn": fn, "delete": sf.Deleted}).Error("Download failed")
		}
		sendError(c, "Invalid filename, key, or file is deleted")
		return
	}

	sf.LastAccess = time.Now().UTC()

	//if file has been download too many times, show this page, with requires a captcha to be solved
	if sf.Downloads >= download.MaxDownloads && !download.checkCaptcha(googleCaptchaCode) {
		if utils.IsWebBrowser(c.Request.Header.Get("User-Agent")) {
			c.HTML(http.StatusOK, "download.tmpl", gin.H{
				"filename": fn,
			})
			return
		}

		log.WithFields(log.Fields{"key": key, "filename": fn}).Info("File downloads exceeded - not web browser")
		sendError(c, "This file has been download too many times, visit the URL with a Browser")
		return
	}

	sf.Downloads = sf.Downloads + 1
	log.WithFields(log.Fields{"key": key, "filename": fn}).Info("Downloads set to ", sf.Downloads)
	models.DB.WriteStoredFile(*sf)

	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")

	ip := c.ClientIP()
	log.WithFields(log.Fields{"ip": ip, "key": key, "fn": fn}).Info("S3 download started.")
	c.Redirect(http.StatusMovedPermanently, utils.GetS3SignedURL(sf.Key, sf.FileName, ip))
	return
}
