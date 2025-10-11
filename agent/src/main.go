package main

import (
	"errors"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"syscall"
	"time"
	"villain-couch/agent/src/bootstrap"
	"villain-couch/agent/src/cli"
	"villain-couch/agent/src/cli/operations"
	"villain-couch/agent/src/config"
	mediaplayer "villain-couch/agent/src/media-player"
	"villain-couch/agent/src/models"
	"villain-couch/agent/src/options"
	"villain-couch/agent/src/storage"
	"villain-couch/common/logger"
)

func main() {
	bootstrap.Bootstrap()
	flags, db, opts, conf := cli.GetFlags(), storage.GetDB(), options.GetOptions(), config.GetConfig()
	operations.New().Build(flags, db, opts).Sort().Run().Finalize() // What in the Java...
	vlc := mediaplayer.New(conf, opts)
	run(&vlc, opts)
}

func run(vlc *mediaplayer.VLCMediaPlayer, opts *options.Options) {
	logger.Log.Info("------ Starting VilLain Couch [VLC Tracker] ------")

	if err := vlc.CommandRunner.Start(); err != nil {
		logger.Log.Error("Failed to start command", "error", err)
		os.Exit(1)
	}

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
			bootstrap.Teardown()
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

		case <-time.After(500 * time.Millisecond):
			// This case executes if the Done channel is not ready yet.
			handleTick(vlc, opts)
		}
	}
}

func handleTick(vlc *mediaplayer.VLCMediaPlayer, opts *options.Options) {
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

	// If file is completed then stop VLC.
	// You don't need to update cache since it will be empty string.
	if status.GetState() == models.StateStopped {
		saveMediaStates()
		storage.GetCache().Delete(currentFilepath)
		if err := vlc.TryNext(currentFilepath); err != nil {
			if errors.Is(err, mediaplayer.ErrorMediaFileNotFound) {
				if opts.FuzzyFoundNextEpisode != "" {
					err := vlc.PlayFile(opts.FuzzyFoundNextEpisode)
					if err != nil {
						logger.Log.Error("VLC PlayFile Error on Fuzzy Found Next Episode", "error", err, "path", opts.FuzzyFoundNextEpisode)
						_ = vlc.CommandRunner.Stop()
					}
				}
			} else {
				logger.Log.Warn("cannot play next file", "error", err)
				_ = vlc.CommandRunner.Stop()
			}

		}
	} else {
		mf := models.NewMediaFileFromStatus(status, currentFilepath)
		storage.GetCache().Set(mf.Filepath, mf)
	}
}

func saveMediaStates() {
	logger.Log.Info("Saving media states...")
	db := storage.GetDB()
	cache := storage.GetCache()
	for _, key := range cache.Keys() {
		// If empty filepath then do not add it to db
		if key == "" {
			continue
		}
		val, _ := cache.Get(key)
		err := db.SetMediaFile(val)
		if err != nil {
			logger.Log.Error("could not save state to database", "error", err.Error())
			return
		}
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
