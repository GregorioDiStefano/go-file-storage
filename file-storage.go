package main

import (
	"github.com/GregorioDiStefano/go-file-storage/controllers"
	"github.com/GregorioDiStefano/go-file-storage/helpers"
	"github.com/GregorioDiStefano/go-file-storage/models"
	"github.com/gin-gonic/gin"
)

func init() {
	helpers.ParseConfig("config/config.json")
	models.DB.OpenDatabaseFile()
}

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")

	go deleteUnusedFile()

	router.GET("/", controller.IndexPage)
	router.PUT("/:filename", controller.SimpleUpload)
	router.GET("/:key/:filename", controller.DownloadFile)
	router.DELETE("/:key/:delete_key/:filename", controller.DeleteFile)
	router.Run(helpers.Config.ServerPort)
}
