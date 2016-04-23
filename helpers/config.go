package helpers

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
)

var Config Configuration

type Configuration struct {
	KeySize                       uint8
	DeleteKeySize                 uint8
	MaxDownloadsBeforeInteraction int64
	MaxSize                       int64
	OverMaxSizeStr                string
	StorageFolder                 string
	Domain                        string
	ServerPort                    string
}

func checkParsedValues() {
	v := reflect.ValueOf(Config)

	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i).Name
		value := v.Field(i).Interface()

		if field == "ServerPort" {
			if len(string(value.(string))) == 0 {
				fmt.Println("ServerPort not set, using 8080 as default")
				Config.ServerPort = "127.0.0.1:8080"
			}
		} else if value == nil || value == reflect.Zero(reflect.TypeOf(value)).Interface() {
			fmt.Printf("Warning: no value set for field: %s\n", field)
		}
	}
}

func ParseConfig() {
	file, err := os.Open("config.json")
	if err != nil {
		panic("quitting, due to problem reading config.json:" + err.Error())
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&Config)

	checkParsedValues()

	if err != nil {
		panic("quitting, due to problem parsing configuration: " + err.Error())
	}

	if len(Config.StorageFolder) > 0 {
		if _, err := os.Stat(Config.StorageFolder); os.IsNotExist(err) {
			panic(fmt.Sprintf("quitting, storage dir (%s) does not exist.", Config.StorageFolder))
		}
	}
}
