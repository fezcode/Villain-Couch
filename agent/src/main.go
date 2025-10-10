package main

import (
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"syscall"
	"time"
	"vlc-tracker-agent/agent/src/bootstrap"
	"vlc-tracker-agent/agent/src/config"
	mediaplayer "vlc-tracker-agent/agent/src/media-player"
	"vlc-tracker-agent/agent/src/models"
	"vlc-tracker-agent/agent/src/options"
	"vlc-tracker-agent/agent/src/storage"
	"vlc-tracker-agent/common/logger"
)

func main() {
	bootstrap.Bootstrap()
	vlc := mediaplayer.New(config.GetConfig(), options.GetOptions())

	logger.Log.Info("Starting VilLain Couch [VLC Tracker]")
	if err := vlc.CommandRunner.Start(); err != nil {
		logger.Log.Error("Failed to start command", "error", err)
		os.Exit(1)
	}

	run(&vlc)
}

func run(vlc *mediaplayer.VLCMediaPlayer) {
	// Graceful Shutdown Setup
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	for {
		select {
		case err := <-vlc.CommandRunner.Done():
			// The Done channel is closed, and the final error state is received.
			if err != nil {
				// The error "signal: interrupt" is expected here because we stopped it.
				logger.Log.Warn("Background command finished with an error (as expected).", "error", err)
			}

			// TODO Post Close handle here
			saveMediaStates()
			os.Exit(0)
			// Exit the for loop somehow

		case sig := <-sigChan:
			logger.Log.Info("Received signal, initiating graceful shutdown.", "signal", sig.String())

			// Stop the background process.
			if err := vlc.CommandRunner.Stop(); err != nil {
				logger.Log.Error("Failed to send stop signal to command", "error", err)
			}

			// Wait for the command to fully terminate.
			<-vlc.CommandRunner.Done()
			logger.Log.Info("Background command stopped successfully.")

		case <-time.After(5 * time.Second):
			// This case executes if the Done channel is not ready yet.
			handleTick(vlc)
		}
	}
}

func handleTick(vlc *mediaplayer.VLCMediaPlayer) {
	status, err := vlc.Status()
	if err != nil {
		logger.Log.Error("VLC GetStatus Error", "error", err)
		// ignore error
	}

	playlist, err := vlc.Playlist()
	if err != nil {
		logger.Log.Error("VLC GetPlaylist Error", "error", err)
		// ignore error
	}

	currentFilepath, err := playlist.GetCurrent()
	if err != nil {
		logger.Log.Error("VLC GetCurrent Error", "error", err)
		// ignore error
	}

	vlc.LogStatus(status)
	mf := models.NewMediaFileFromStatus(status, currentFilepath)
	storage.GetCache().Set(mf.Filename, mf)
}

func saveMediaStates() {
	logger.Log.Info("Saving media states...")
	db := storage.GetDB()
	cache := storage.GetCache()
	for _, key := range cache.Keys() {
		val, _ := cache.Get(key)
		err := db.SetMediaFile(val)
		if err != nil {
			logger.Log.Error("could not save state to database", "error", err.Error())
			return
		}
	}

	if err := db.Close(); err != nil {
		logger.Log.Error("could not close database", "error", err.Error())
		return
	}
}

// clearConsole clears the terminal screen.
func clearConsole() {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		_ = cmd.Run() // We can ignore the error here
	} else {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		_ = cmd.Run()
	}
}
