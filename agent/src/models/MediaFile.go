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

func NewMediaFileFromStatus(v StatusMessage, s string) MediaFile {
	return MediaFile{
		Filepath:      s,
		Filename:      v.GetFilename(),
		TotalSeconds:  v.GetLength(),
		CurrentSecond: v.GetTime(),
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
}

func (mf MediaFile) IsEmpty() bool {
	return mf.Filepath == ""
}
