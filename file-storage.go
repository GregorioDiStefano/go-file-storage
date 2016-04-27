package main

import (
	"./controllers"
	"./helpers"
	"./models"
	"github.com/gin-gonic/gin"
)

func init() {
	db := models.Database{Filename: models.DbFilename, Bucket: models.Bucket}
	helpers.ParseConfig()
	db.OpenDatabaseFile()
}

func main() {
	router := gin.Default()

	router.GET("/", controller.IndexPage)
	router.PUT("/:filename", controller.SimpleUpload)
	router.GET("/:key/:filename", controller.FileDownloader)
	router.DELETE("/:key/:delete_key", controller.DeleteFile)
	router.Run(helpers.Config.ServerPort)
}
