package main

import (
	"./controllers"
	"./helpers"
	"./models"
	"github.com/gin-gonic/gin"
)

func init() {
	db := models.Database{Filename: models.DbFilename, Bucket: models.Bucket}
	helpers.ParseConfig("config.json")
	db.OpenDatabaseFile()
}

func main() {
	router := gin.Default()

	router.GET("/", controller.IndexPage)
	router.PUT("/:filename", controller.SimpleUpload)
	router.GET("/:key/:filename", controller.DownloadFile)
	router.DELETE("/:key/:delete_key/:filename", controller.DeleteFile)
	router.Run(helpers.Config.ServerPort)
}
