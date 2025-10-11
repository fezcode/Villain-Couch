package models

import "time"

type Workspace struct {
	DirectoryPath string
	DirectoryName string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     time.Time
}

func (w Workspace) IsEmpty() bool {
	return w.DirectoryPath == "" && w.DirectoryName == ""
}
