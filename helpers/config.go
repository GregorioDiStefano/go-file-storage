package helpers

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
)

var Config Configuration

type Configuration struct {
	KeySize                        uint8
	DeleteKeySize                  uint8
	MaxDownloadsBeforeInteraction  int64
	MaxSize                        int64
	DeleteAfterSecondsLastAccessed int64
	FileCheckFrequency             uint
	OverMaxSizeStr                 string

	Domain        string
	ServerPort    string
	CaptchaSecret string

	AccessKey string
	SecretKey string

	StorageMethod string
	StorageFolder string
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
			//		} else if field == "StorageMethod" && value.(string) != "S3" || value.(string) != "local" {
			//			panic("StorageMethod in configuration file is incorrect, specify S3 or local")
		} else if value == nil || value == reflect.Zero(reflect.TypeOf(value)).Interface() {
			fmt.Printf("Warning: no value set for field: %s\n", field)
		}
	}
}

func ParseConfig(filename string) {
	file, err := os.Open(filename)

	if err != nil {
		panic("quitting, due to problem reading config.json:" + err.Error())
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&Config)
	fmt.Println(Config)

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
