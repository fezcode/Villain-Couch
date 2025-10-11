package operations

import (
	"errors"
	"fmt"
	"os"
	"time"
	"villain-couch/agent/src/models"
	"villain-couch/common/fs"
	"villain-couch/common/logger"
)

type AddWorkspace struct {
	Operation
	DirPath, DirName string
	AlreadyExists    bool
}

func (a AddWorkspace) Priority() int {
	return OrderMedium
}

func (a AddWorkspace) DefaultError() string {
	return fmt.Sprintf("Cannot run %s operation", a.Name())
}

func (a AddWorkspace) Name() string {
	return "Add Workspace"
}

func (a AddWorkspace) Run() error {
	now := time.Now()
	ws := models.Workspace{DirectoryPath: a.DirPath, DirectoryName: a.DirName, CreatedAt: now, UpdatedAt: now}

	if !fs.DirectoryExists(a.DirPath) {
		logger.Log.Error("Given directory does not exist")
		return errors.New("given directory does not exist")
	}

	err := a.Database.InsertWorkspace(ws)
	if err != nil {
		logger.Log.Error(a.DefaultError(), "error", err)
		return err
	}
	return nil
}

func (a AddWorkspace) Finalize() {
	msg := fmt.Sprintf("COMMAND '%s' WILL NOW CLOSE THE AGENT.", a.Name())
	logger.Log.Warn(msg)
	os.Exit(0)
}
