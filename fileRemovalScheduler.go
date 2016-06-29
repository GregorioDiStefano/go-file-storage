package main

import (
	"time"

	log "github.com/Sirupsen/logrus"

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
				helpers.Log.WithFields(log.Fields{"filename": sf.FileName, "last_accessed": sf.LastAccess, "upload_time": sf.UploadTime}).Infoln("Removed old file.")
				sf.Deleted = true
				db.WriteStoredFile(*sf)
			}
		}
	}
}
