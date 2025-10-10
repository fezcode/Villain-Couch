package storage

import (
	"villian-couch/agent/src/models"
	"villian-couch/common/logger"
)

// Database
var db *DB // database object

func GetDB() *DB {
	return db
}

// Cache
var cache *Cache[models.MediaFile]

func GetCache() *Cache[models.MediaFile] {
	return cache
}

func Initialize(DatabaseFilePath string) error {
	var err error
	cache = NewCache[models.MediaFile]()
	db, err = NewDB(DatabaseFilePath)
	return err
}

func Shutdown() {
	if err := db.Close(); err != nil {
		logger.Log.Error("could not close database", "error", err.Error())
		return
	}
}
