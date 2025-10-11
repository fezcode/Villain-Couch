package options

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"villain-couch/agent/src/cli"
	"villain-couch/agent/src/config"
	"villain-couch/agent/src/resolver"
	"villain-couch/agent/src/storage"
	"villain-couch/common/globals"
	"villain-couch/common/logger"
	"villain-couch/common/optional"
	"villain-couch/common/step"
)

type Options struct {
	// Generated Ones
	VLCPath               string
	DatabaseFilePath      string
	MediaFilePath         string
	MediaFileStartTime    string
	FuzzyFoundNextEpisode string
}

var opts *Options

func GetOptions() *Options {
	return opts
}

// Sets additional options after initalization
func SetOptions(db *storage.DB) error {
	if opts.MediaFilePath == "" {
		file, err := db.GetLatestUpdatedMediaFile()
		if err != nil {
			logger.Log.Error(err.Error(), "msg", "Error getting latest updated media file.")
			return err
		}

		if file == nil {
			logger.Log.Error("No file to play. Exiting.")
			os.Exit(1)
		}

		opts.MediaFilePath = file.Filepath
		opts.MediaFileStartTime = strconv.Itoa(file.CurrentSecond)
	}

	return nil
}

func ValidateOptions() {
	// Check if the media file exists before trying to launch VLC.
	if _, err := os.Stat(opts.MediaFilePath); os.IsNotExist(err) {
		logger.Log.Error("Media file not found", "Media File", opts.MediaFilePath)
		os.Exit(1)
	}
}

func Initialize(fl *cli.CLIFlags, conf *config.Config) error {
	opts = &Options{}
	steps := []step.Step{
		{F: putVLCPath, P: conf.VLCPath},
		{F: putDatabasePath},
		{F: putMediaFilePath, P: fl.MediaFile.String()},
	}
	return step.RunSteps(steps)
}

func putVLCPath(p ...string) error {
	optionalVLCPath := optional.FirstOrEmpty(p)
	location, found, err := resolver.GetVLCInstallLocation(optionalVLCPath)
	if err != nil {
		logger.Log.Error("could not get VLC location", "error", err)
		return err
	}

	if !found {
		logger.Log.Error("could not find media player install location")
		logger.Log.Warn("INSTALL VLC")
		logger.Log.Warn("INSTALL VLC")
		os.Exit(1)
	}

	opts.VLCPath = location
	return nil
}

func putDatabasePath(...string) error {
	dir, _, err := globals.GetConfigPaths()
	if err != nil {
		logger.Log.Error(err.Error(), "msg", "Error getting config file path")
		return err
	}

	opts.DatabaseFilePath = filepath.Join(dir, config.GetConfig().DatabaseFileName)
	return nil
}

func putMediaFilePath(params ...string) error {
	if len(params) == 0 {
		logger.Log.Error("no parameters provided for media file path location")
		return errors.New("no parameters provided for media file path location")
	}
	opts.MediaFilePath = params[0]
	return nil
}
