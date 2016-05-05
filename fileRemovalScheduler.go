package main

import (
	"fmt"
	"time"

	"github.com/GregorioDiStefano/go-file-storage/helpers"
	"github.com/GregorioDiStefano/go-file-storage/models"
)

//Iterate over all file stored, and delete files that have not been accessed since DeleteAfterSecondsLastAccessed
func deleteUnusedFile() {
	for {
		time.Sleep(time.Duration(helpers.Config.FileCheckFrequency) * time.Second)
		db := models.Database{Filename: models.DbFilename, Bucket: models.Bucket}
		for _, v := range *db.GetAllKeys() {
			sf := db.ReadStoredFile(v)
			delta := time.Now().Sub(sf.LastAccess)
			if sf.Deleted == false && int64(delta.Seconds()) > helpers.Config.DeleteAfterSecondsLastAccessed {
				fmt.Println("Deleting: ", sf.Key)
				sf.Deleted = true
				db.WriteStoredFile(*sf)
			}
		}
	}
}
