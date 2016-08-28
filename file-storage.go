package main

import (
	"os"
	"time"

	redis "gopkg.in/redis.v3"

	"github.com/GregorioDiStefano/go-file-storage/controllers"
	"github.com/GregorioDiStefano/go-file-storage/log"
	"github.com/GregorioDiStefano/go-file-storage/models"
	"github.com/GregorioDiStefano/go-file-storage/utils"
	"github.com/etcinit/speedbump"
	"github.com/etcinit/speedbump/ginbump"
	"github.com/gin-gonic/gin"
)

func init() {

	log.SetLevel(log.DebugLevel)
	log.Debug("Starting..")

	configFile := os.Getenv("CONFIG_FILE")
	utils.LoadConfig(configFile)

	models.DB.Setup(utils.Config.GetInt("key_size"))
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

	downloader := controller.NewDownloader(utils.Config.GetString("CAPTCHA_SECRET"), utils.Config.GetInt("max_downloads"))
	uploader := controller.NewUploader(utils.Config.GetString("domain"), utils.Config.GetInt64("max_file_size"), utils.Config.GetInt("delete_key_size"), utils.Config.GetString("aws.bucket"), utils.Config.GetString("aws.region"))
	deleter := controller.NewDeleter(*uploader)

	router.GET("/", controller.IndexPage)
	router.GET("/:key/:filename", downloader.DownloadFile)
	router.PUT("/:filename", ginbump.RateLimit(client, speedbump.PerHourHasher{}, 5), uploader.UploadFile)
	router.DELETE("/:key/:delete_key/:filename", ginbump.RateLimit(client, speedbump.PerHourHasher{}, 5), deleter.DeleteFile)

	router.Run(utils.Config.GetString("port"))
}
