package helpers

import (
	"encoding/json"
	"fmt"
	"os"
)

var Config Configuration

type Configuration struct {
	KeySize                       uint8
	MaxDownloadsBeforeInteraction int64
	MaxSize                       int64
	OverMaxSizeStr                string
}

func ParseConfig() {
	file, err := os.Open("config.json")
	if err != nil {
		panic("quitting, due to problem reading config.json:" + err.Error())
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&Config)

	if err != nil {
		panic("quitting, due to problem parsing configuration: " + err.Error())
	}

	fmt.Println(Config.MaxSize)
}
