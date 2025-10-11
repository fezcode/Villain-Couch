package operations

import (
	"villian-couch/agent/src/options"
	"villian-couch/agent/src/storage"
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
