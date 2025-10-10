package bootstrap

import (
	"os"
	"villian-couch/agent/src/cli"
	"villian-couch/agent/src/config"
	"villian-couch/agent/src/options"
	"villian-couch/agent/src/storage"
	"villian-couch/common/logger"
)

func Bootstrap() {
	cli.Initialize()
	logger.Initialize(cli.GetFlags().Verbose)
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
