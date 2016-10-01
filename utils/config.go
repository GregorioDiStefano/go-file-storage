package utils

import (
	"fmt"
	"os"

	"errors"

	"github.com/spf13/viper"
)

var Config viper.Viper

func LoadConfig(filename string) {
	Config = *viper.New()
	fmt.Println(os.Getwd())
	fmt.Println(filename)

	if _, err := os.Stat(filename); err != nil {
		panic("Unable to load config file:" + err.Error())
	}

	Config.SetConfigFile(filename)
	Config.BindEnv("CAPTCHA_SECRET")
	Config.BindEnv("AWS_ACCESS_KEY_ID")
	Config.BindEnv("AWS_SECRET_ACCESS_KEY")

	if err := Config.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	if err := checkConfig(); err != nil {
		panic(err)
	}

}

func checkConfig() error {
	for _, key := range Config.AllKeys() {
		switch key {
		case "file_check_freq":
		case "delete_key_size":
		case "delete_after_seconds":
		case "aws":
		case "max_file_size":
		case "captcha_secret":
		case "aws_access_key_id":
		case "aws_secret_access_key":
		case "max_downloads":
		case "key_size":
		case "domain":
		case "port":

		default:
			return errors.New("Missing key in config file:" + key)
		}
	}

	return nil
}
