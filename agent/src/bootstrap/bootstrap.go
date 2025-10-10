package bootstrap

import (
	"os"
	"vlc-tracker-agent/agent/src/cli"
	"vlc-tracker-agent/agent/src/config"
	"vlc-tracker-agent/agent/src/options"
	"vlc-tracker-agent/agent/src/storage"
	"vlc-tracker-agent/common/logger"
)

func Bootstrap() {
	logger.Initialize()
	cli.Initialize()
	if err := config.Initialize(); err != nil {
		logger.Log.Error(err.Error(), "msg", "Error setting up config.")
		os.Exit(1)
	}

	if err := options.Initialize(cli.GetFlags()); err != nil {
		logger.Log.Error(err.Error(), "msg", "Error setting up options.")
		os.Exit(1)
	}

	if err := storage.Initialize(options.GetOptions().DatabaseFilePath); err != nil {
		logger.Log.Error(err.Error(), "msg", "Error setting up storage.")
		os.Exit(1)
	}

	if err := options.SetOptions(storage.GetDB()); err != nil {
		logger.Log.Error(err.Error(), "msg", "Error setting up options.")
		os.Exit(1)
	}

	options.ValidateOptions()

}

func Teardown() {
	storage.Shutdown()
}
