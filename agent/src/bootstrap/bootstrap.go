package bootstrap

import (
	"os"
	"vlc-tracker-agent/agent/src/cli"
	"vlc-tracker-agent/agent/src/config"
	"vlc-tracker-agent/common/logger"
)

func Bootstrap() {
	var err error
	logger.Initialize()
	cli.Initialize()
	err = config.Initialize()

	if err != nil {
		logger.Log.Error(err.Error(), "msg", "Error setting up config.")
		os.Exit(1)
	}

}
