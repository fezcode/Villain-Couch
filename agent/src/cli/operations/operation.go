package operations

import (
	"villain-couch/agent/src/options"
	"villain-couch/agent/src/storage"
)

const (
	OrderNone   = -1
	OrderHigh   = 2
	OrderMedium = 1
	OrderLow    = 0
)

type Runnable interface {
	Name() string
	Run() error
	Priority() int
	DefaultError() string
	Finalize()
}

type Operation struct {
	Database *storage.DB
	Options  *options.Options
}
