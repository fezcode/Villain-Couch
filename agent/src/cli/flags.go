package cli

import (
	"flag"
	"villain-couch/common/str"
)

type CLIFlags struct {
	Verbose      bool
	FindNext     bool
	MediaFile    str.Str
	AddWorkspace str.Str
}

var cliFlags *CLIFlags

func Initialize() {
	cliFlags = parseFlags()
}

func GetFlags() *CLIFlags {
	return cliFlags
}

func parseFlags() *CLIFlags {
	var Verbose, FindNext bool
	var MF, AW string

	flag.BoolVar(&Verbose, "verbose", false, "prints info level logs")
	flag.BoolVar(&FindNext, "find-next", false, "tries to find next episode when there is nothing else to play.")
	flag.StringVar(&MF, "file", "", "media file to play")
	flag.StringVar(&AW, "ws", "", "add directory to workspace, will close agent afterwards")
	flag.Parse()

	return &CLIFlags{
		FindNext:     FindNext,
		Verbose:      Verbose,
		MediaFile:    str.Str(MF),
		AddWorkspace: str.Str(AW),
	}
}
