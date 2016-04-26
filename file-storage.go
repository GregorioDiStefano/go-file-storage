package main

import (
	"./controllers"
	"./helpers"
	"github.com/gin-gonic/gin"
)

func init() {
	helpers.ParseConfig()
}

func main() {
	router := gin.Default()

	router.GET("/", controller.IndexPage)
	router.PUT("/:filename", controller.SimpleUpload)
	router.GET("/:key/:filename", controller.FileDownloader)
	router.DELETE("/:key/:delete_key", controller.DeleteFile)
	router.Run(helpers.Config.ServerPort)
}
