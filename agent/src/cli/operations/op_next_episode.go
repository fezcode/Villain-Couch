package operations

import (
	"errors"
	"fmt"
	"villain-couch/common/ff"
	"villain-couch/common/logger"
	"villain-couch/common/optional"
)

type NextEpisode struct {
	Operation
}

func (a NextEpisode) Priority() int {
	return OrderMedium
}

func (a NextEpisode) Finalize() {}

func (a NextEpisode) DefaultError() string {
	return fmt.Sprintf("Cannot run %s operation", a.Name())
}

func (a NextEpisode) Name() string {
	return "Try to Find Next Operation"
}

func (a NextEpisode) Run() error {
	file, err := a.Database.GetLatestUpdatedMediaFile()
	if err != nil {
		logger.Log.Error("Error getting latest updated media file", "error", err)
		return err
	}

	if file == nil {
		logger.Log.Error("no media file found")
		return errors.New("no media file found")
	}

	workspaces, err := a.Database.GetWorkspaces()
	if err != nil {
		logger.Log.Error("Error getting workspaces", "error", err)
		return err
	}

	ws := optional.FirstOrEmpty(workspaces)
	if ws.IsEmpty() {
		logger.Log.Error("no workspace found")
		return errors.New("no workspace found")
	}

	relatedFiles, err := ff.FindRelatedFiles(ws.DirectoryPath, file.Filename)
	if err != nil {
		logger.Log.Error("Error finding related files", "error", err)
		return err
	}

	info, err := ff.ParseEpisodeInfo(file.Filename)
	if err != nil {
		logger.Log.Error("Could not parse target filename to find next episode", "error", err)
		return err
	}

	nextEpisode, found := ff.FindNextEpisode(info, relatedFiles)
	if !found {
		logger.Log.Warn("Could not find a subsequent episode in the related files.")
		logger.Log.Warn("You might be on the last available episode.")
		return errors.New("could not find a subsequent episode in the related files")
	}

	a.Options.FuzzyFoundNextEpisode = nextEpisode.FilePath
	return nil
}
