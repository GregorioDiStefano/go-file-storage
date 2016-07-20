package main

import (
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/GregorioDiStefano/go-file-storage/models"
	"github.com/GregorioDiStefano/go-file-storage/utils"
)

//Iterate over all file stored, and delete files that have not been accessed since DeleteAfterSecondsLastAccessed
func deleteUnusedFile() {
	for {
		time.Sleep(time.Duration(utils.Config.GetInt("file_check_freq")) * time.Second)
		db := models.Database{Filename: models.DbFilename, Bucket: models.Bucket}
		for _, v := range *db.GetAllKeys() {
			sf := db.ReadStoredFile(v)
			delta := time.Now().Sub(sf.LastAccess)
			if sf.Deleted == false && int64(delta.Seconds()) > int64(utils.Config.GetInt("delete_after_seconds")) {
				utils.Log.WithFields(log.Fields{"filename": sf.FileName, "last_accessed": sf.LastAccess, "upload_time": sf.UploadTime}).Infoln("Removed old file.")
				sf.Deleted = true
				db.WriteStoredFile(*sf)
			}
		}
	}
}
