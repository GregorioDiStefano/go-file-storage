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
	StorageFolder                 string
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

	if len(Config.StorageFolder) > 0 {
		if _, err := os.Stat(Config.StorageFolder); os.IsNotExist(err) {
			panic(fmt.Sprintf("quitting, storage dir (%s) does not exist.", Config.StorageFolder))
		}
	}

	fmt.Println(Config.MaxSize)
}
