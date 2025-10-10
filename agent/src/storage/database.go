package storage

import (
	"database/sql"
	_ "embed"
	"fmt"
	"time"
	"vlc-tracker-agent/agent/src/models"
	"vlc-tracker-agent/common/logger"

	_ "modernc.org/sqlite"
)

//go:embed queries/createTables.sql
var queryCreateTables string

//go:embed queries/setMediaFile.sql
var querySetMediaFile string

//go:embed queries/getMediaFile.sql
var queryGetMediaFile string

//go:embed queries/getLatestMediaFile.sql
var queryGetLatestMediaFile string

// DB represents a wrapper around the SQL database connection.
type DB struct {
	conn *sql.DB
}

// NewDB initializes a connection to an SQLite database file at the given path.
// It creates the file and the necessary table if they don't exist.
func NewDB(path string) (*DB, error) {
	// sql.Open() creates the database file if it doesn't exist.
	conn, err := sql.Open("sqlite", path)
	if err != nil {
		logger.Log.Error("could not open sqlite database", "error", err)
		return nil, err
	}

	if err := conn.Ping(); err != nil {
		logger.Log.Error("could not connect to sqlite database", "error", err)
		return nil, err
	}

	if _, err := conn.Exec(queryCreateTables); err != nil {
		logger.Log.Error("could not create table", "error", err)
		return nil, err
	}

	return &DB{conn: conn}, nil
}

// Close closes the database connection.
func (db *DB) Close() error {
	return db.conn.Close()
}

// SetMediaFile inserts or updates a record in the media_files table.
// It preserves the original created_at timestamp on updates.
func (db *DB) SetMediaFile(mf models.MediaFile) error {
	now := time.Now()
	// For an INSERT, both created_at and updated_at are `now`.
	// For an UPDATE, the new `updated_at` value from the `excluded` row is used,
	// and the `created_at` column is NOT mentioned in the `DO UPDATE` clause,
	// so it remains unchanged from the original record.
	_, err := db.conn.Exec(querySetMediaFile, mf.Filepath, mf.Filename, mf.TotalSeconds, mf.CurrentSecond, now, now)
	if err != nil {
		logger.Log.Error("failed to set media file for filepath", "Filepath", mf.Filepath)
		return fmt.Errorf("failed to set media file for filepath '%s': %w", mf.Filepath, err)
	}
	return nil
}

// GetMediaFile retrieves a media file record by its filepath.
// It returns sql.ErrNoRows if the filepath is not found.
func (db *DB) GetMediaFile(filepath string) (*models.MediaFile, error) {
	mf := &models.MediaFile{}

	err := db.conn.QueryRow(queryGetMediaFile, filepath).Scan(
		&mf.Filepath,
		&mf.Filename,
		&mf.TotalSeconds,
		&mf.CurrentSecond,
		&mf.CreatedAt,
		&mf.UpdatedAt,
	)
	if err != nil {
		logger.Log.Error("failed to get media file for filepath", "filepath", filepath, "error", err)
		return nil, fmt.Errorf("failed to get media file for filepath '%s': %w", filepath, err)
	}
	return mf, nil
}

// GetLatestUpdatedMediaFile retrieves the most recently updated record from the media_files table.
func (db *DB) GetLatestUpdatedMediaFile() (*models.MediaFile, error) {
	row := db.conn.QueryRow(queryGetLatestMediaFile)
	var mf models.MediaFile
	err := row.Scan(&mf.Filepath, &mf.Filename, &mf.TotalSeconds, &mf.CurrentSecond, &mf.CreatedAt, &mf.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			// This means the table is empty. It's not an application error.
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get latest media file: %w", err)
	}
	return &mf, nil
}
