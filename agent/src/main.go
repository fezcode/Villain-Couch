package main

import (
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"syscall"
	"time"
	"vlc-tracker-agent/agent/src/bootstrap"
	"vlc-tracker-agent/agent/src/cli"
	"vlc-tracker-agent/agent/src/config"
	mediaplayer "vlc-tracker-agent/agent/src/media-player"
	"vlc-tracker-agent/common/logger"
)

func main() {
	bootstrap.Bootstrap()

	conf := config.GetConfig()
	flags := cli.GetFlags()
	vlc := mediaplayer.New(conf, flags)

	logger.Log.Info("Starting Villain Couch [VLC Tracker]")
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
			os.Exit(0)
		//break waitingLoop // Exit the for loop.

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
			// logger.Log.Info("...still waiting for process to terminate...")
			// You can perform other periodic tasks here.
			status, err := vlc.Status()
			if err != nil {
				logger.Log.Error("VLC GetStatus Error", "error", err)
				// ignore error
			}

			vlc.PrintStatus(status)
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
