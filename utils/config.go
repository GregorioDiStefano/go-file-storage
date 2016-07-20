package utils

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

var Config viper.Viper

func ParseConfig(filename string) {
	Config = *viper.New()
	Config.Debug()

	if _, err := os.Stat(filename); err != nil {
		panic(err)
	}

	Config.SetConfigFile(filename)
	Config.BindEnv("CAPTCHA_SECRET")
	Config.BindEnv("AWS_ACCESS_KEY_ID")
	Config.BindEnv("AWS_SECRET_ACCESS_KEY")

	err := Config.ReadInConfig()

	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}
