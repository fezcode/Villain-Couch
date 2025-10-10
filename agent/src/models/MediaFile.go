package models

import "time"

// MediaFile represents a row in the media_files table.
type MediaFile struct {
	Filepath      string
	Filename      string
	TotalSeconds  int
	CurrentSecond int
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func FromStatus(v StatusMessage) MediaFile {
	return MediaFile{
		Filepath:      v.GetFilename(), // TODO find a better way.
		Filename:      v.GetFilename(),
		TotalSeconds:  v.GetLength(),
		CurrentSecond: v.GetTime(),
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
}
