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

var Version = "dev" // Overwritten by gobake build

func main() {
	// Basic flag handling
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "-v", "--version":
			fmt.Printf("villain.couch v%s\n", Version)
			return
		case "-h", "--help":
			fmt.Println("Villain Couch - VLC Progress Tracker")
			fmt.Println("\nUsage:")
			fmt.Println("  villain-couch [options]")
			fmt.Println("\nOptions:")
			fmt.Println("  -v, --version  Show version")
			fmt.Println("  -h, --help     Show this help")
			fmt.Println("\nControls (in TUI):")
			fmt.Println("  Space          Toggle Pause/Play")
			fmt.Println("  n              Next Episode")
			fmt.Println("  h/l            Seek -/+ 10s")
			fmt.Println("  +/-            Volume up/down")
			fmt.Println("  Tab            Switch to History view")
			fmt.Println("  Enter          Resume selected from History")
			fmt.Println("  q              Quit")
			return
		}
	}

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
