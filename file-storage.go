package main

import (
	"os"
	"time"

	redis "gopkg.in/redis.v3"

	"github.com/GregorioDiStefano/go-file-storage/controllers"
	"github.com/GregorioDiStefano/go-file-storage/models"
	"github.com/GregorioDiStefano/go-file-storage/utils"
	"github.com/etcinit/speedbump"
	"github.com/etcinit/speedbump/ginbump"
	"github.com/gin-gonic/gin"
)

func init() {
	configFile := os.Getenv("CONFIG_FILE")
	utils.ParseConfig(configFile)
	models.DB.OpenDatabaseFile()
}

func main() {
	defer models.DB.CloseDatabaseFile()

	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	utils.Log.Info("Starting.....")

	go deleteUnusedFile()

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	go func() {
		if client.Ping().Err() != nil {
			utils.Log.Warningf("Communication with Redis failed!")
		}
		time.Sleep(5 * time.Minute)
	}()

	router.GET("/", controller.IndexPage)
	router.GET("/:key/:filename", controller.DownloadFile)
	router.PUT("/:filename", ginbump.RateLimit(client, speedbump.PerHourHasher{}, 5), controller.Upload)
	router.DELETE("/:key/:delete_key/:filename", ginbump.RateLimit(client, speedbump.PerHourHasher{}, 5), controller.DeleteFile)

	router.Run(utils.Config.GetString("port"))
}
