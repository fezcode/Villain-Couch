package operations

import (
	"villain-couch/agent/src/options"
	"villain-couch/agent/src/storage"
)

type Runnable interface {
	Name() string
	Run() error
	DefaultError() string
	Finalize()
}
type Operation struct {
	Database *storage.DB
	Options  *options.Options
}
