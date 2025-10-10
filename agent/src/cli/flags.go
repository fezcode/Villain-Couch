package cli

import (
	"flag"
)

type CLIFlags struct {
	Continue  bool
	Verbose   bool
	MediaFile string
}

var cliFlags *CLIFlags

func Initialize() {
	a := CLIFlags{}
	flag.BoolVar(&a.Continue, "continue", false, "continue where you left off")
	flag.BoolVar(&a.Verbose, "verbose", false, "prints info level logs")
	flag.StringVar(&a.MediaFile, "file", "", "media file to play")
	flag.Parse()
	cliFlags = &a
}

func GetFlags() *CLIFlags {
	return cliFlags
}
