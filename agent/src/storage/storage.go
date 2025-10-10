package storage

import (
	"vlc-tracker-agent/agent/src/config"
	"vlc-tracker-agent/agent/src/models"
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

func Initialize(conf *config.Config) error {
	var err error
	cache = NewCache[models.MediaFile]()
	db, err = NewDB(conf)
	return err
}
