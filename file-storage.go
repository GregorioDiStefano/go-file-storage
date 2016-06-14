package main

import (
	"os"
	"time"

	"gopkg.in/redis.v3"

	"github.com/GregorioDiStefano/go-file-storage/controllers"
	"github.com/GregorioDiStefano/go-file-storage/helpers"
	"github.com/GregorioDiStefano/go-file-storage/models"
	"github.com/etcinit/speedbump"
	"github.com/etcinit/speedbump/ginbump"
	"github.com/gin-gonic/gin"
)

func init() {

	configFile := os.Getenv("CONFIG_FILE")
	if configFile == "" {
		panic("No CONFIG_FILE set")
	}
	helpers.ParseConfig(configFile)
	models.DB.OpenDatabaseFile()
}

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	defer models.DB.CloseDatabaseFile()

	helpers.Log.Info("Starting.....")
	go deleteUnusedFile()

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	go func() {
		if client.Ping().Err() != nil {
			panic("Communication with Redis failed!")
		}
		time.Sleep(5 * time.Minute)
	}()

	router.GET("/", controller.IndexPage)
	router.GET("/:key/:filename", controller.DownloadFile)
	router.PUT("/:filename", ginbump.RateLimit(client, speedbump.PerHourHasher{}, 5), controller.SimpleUpload)
	router.DELETE("/:key/:delete_key/:filename", ginbump.RateLimit(client, speedbump.PerHourHasher{}, 5), controller.DeleteFile)

	router.Run(helpers.Config.ServerPort)
}
