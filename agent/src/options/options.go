package options

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"vlc-tracker-agent/agent/src/cli"
	"vlc-tracker-agent/agent/src/config"
	"vlc-tracker-agent/agent/src/registry"
	"vlc-tracker-agent/agent/src/storage"
	"vlc-tracker-agent/common/globals"
	"vlc-tracker-agent/common/logger"
	"vlc-tracker-agent/common/step"
)

type Options struct {
	// Generated Ones
	VLCPath            string
	DatabaseFilePath   string
	MediaFilePath      string
	MediaFileStartTime string
}

var opts *Options

func GetOptions() *Options {
	return opts
}

func SetOptions(db *storage.DB) error {
	if opts.MediaFilePath == "" {
		file, err := db.GetLatestUpdatedMediaFile()
		if err != nil {
			logger.Log.Error(err.Error(), "msg", "Error getting latest updated media file.")
			return err
		}

		if file == nil {
			logger.Log.Error(err.Error(), "msg", "No file to play. Exiting.")
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

func Initialize(fl *cli.CLIFlags) error {
	opts = &Options{}
	steps := []step.Step{
		{F: putVLCPath},
		{F: putDatabasePath},
		{F: putMediaFilePath, P: &fl.MediaFile},
	}
	return step.RunSteps(steps)
}

func putVLCPath(...string) error {
	location, found, err := registry.GetVLCInstallLocation()
	if err != nil {
		logger.Log.Error("could not read registry", "error", err)
		return err
	}

	if !found {
		logger.Log.Error("could not find media player install location")
		logger.Log.Warn("INSTALL VLC")
		os.Exit(1)
	}

	opts.VLCPath = filepath.Join(location, "vlc.exe")
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
