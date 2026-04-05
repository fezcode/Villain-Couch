package main

import (
	"fmt"
	"os"
	"villain-couch/agent/src/bootstrap"
	"villain-couch/agent/src/cli"
	"villain-couch/agent/src/cli/operations"
	"villain-couch/agent/src/config"
	mediaplayer "villain-couch/agent/src/media-player"
	"villain-couch/agent/src/options"
	"villain-couch/agent/src/storage"
	"villain-couch/agent/src/tui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	bootstrap.Bootstrap()
	flags, db, opts, conf := cli.GetFlags(), storage.GetDB(), options.GetOptions(), config.GetConfig()
	operations.New().Build(flags, db, opts).Sort().Run().Finalize()
	vlc := mediaplayer.New(conf, opts)

	if err := vlc.CommandRunner.Start(); err != nil {
		fmt.Printf("Failed to start VLC: %v\n", err)
		os.Exit(1)
	}

	m := tui.NewModel(&vlc, opts)
	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}

	// Stop the background process.
	_ = vlc.CommandRunner.Stop()
	bootstrap.Teardown()
}
