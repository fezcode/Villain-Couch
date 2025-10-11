package cli

import (
	"flag"
	"villain-couch/common/str"
)

type CLIFlags struct {
	Version      bool
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
	var Version, Verbose, FindNext bool
	var MF, AW string

	flag.BoolVar(&Version, "version", false, "prints version (terminates immediately)")
	flag.BoolVar(&Verbose, "verbose", false, "prints info level logs")
	flag.BoolVar(&FindNext, "find-next", false, "tries to find next episode when there is nothing else to play.")
	flag.StringVar(&MF, "file", "", "media file to play")
	flag.StringVar(&AW, "ws", "", "add directory to workspace, will close agent after all operations.")
	flag.Parse()

	return &CLIFlags{
		Version:      Version,
		FindNext:     FindNext,
		Verbose:      Verbose,
		MediaFile:    str.Str(MF),
		AddWorkspace: str.Str(AW),
	}
}
