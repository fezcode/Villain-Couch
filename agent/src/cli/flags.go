package cli

import (
	"flag"
	"os"
	"vlc-tracker-agent/common/logger"
)

type CLIFlags struct {
	Continue  bool
	MediaFile string
}

var cliFlags *CLIFlags

func Initialize() {
	logger.Log.Info("Initializing Flags")

	a := CLIFlags{}
	flag.BoolVar(&a.Continue, "continue", false, "continue where you left off")
	flag.StringVar(&a.MediaFile, "file", "", "media file to play")
	flag.Parse()
	cliFlags = &a

	// Validate Flags
	// Check if the media file exists before trying to launch VLC.
	if _, err := os.Stat(a.MediaFile); os.IsNotExist(err) {
		logger.Log.Error("Media file not found", "Media File", a.MediaFile)
		os.Exit(1)
	}
}

func GetFlags() *CLIFlags {
	return cliFlags
}
