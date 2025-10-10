package media_player

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"vlc-tracker-agent/agent/src/cli"
	"vlc-tracker-agent/agent/src/config"
	"vlc-tracker-agent/agent/src/models"
	"vlc-tracker-agent/agent/src/options"
	"vlc-tracker-agent/common/logger"
)

type MediaPlayer interface {
	Build(*config.Config, *options.Options)
	Status()
	Playlist()
	PrintStatus(models.StatusMessage)
}

type VLCMediaPlayer struct {
	Args             cli.VLCRunnerArguments
	CommandRunner    *cli.CommandRunner
	StatusEndpoint   string
	PlaylistEndpoint string
}

func New(conf *config.Config, opts *options.Options) VLCMediaPlayer {
	vlc := VLCMediaPlayer{}
	vlc.Build(conf, opts)
	return vlc
}

func (vlc *VLCMediaPlayer) Build(conf *config.Config, opts *options.Options) {
	vlc.Args = cli.PrepareRunnerArguments(opts.VLCPath, opts.MediaFilePath, opts.MediaFileStartTime, conf.ExtraIntf, conf.HttpPort, conf.HttpPassword)
	vlc.CommandRunner = cli.NewCommandRunnerForVLC(vlc.Args)
	vlc.StatusEndpoint = fmt.Sprintf("%s:%s/%s", conf.WebUrl, conf.HttpPort, conf.StatusEndpoint)
	vlc.PlaylistEndpoint = fmt.Sprintf("%s:%s/%s", conf.WebUrl, conf.HttpPort, conf.PlaylistEndpoint)
}

func (vlc *VLCMediaPlayer) Status() (models.StatusMessage, error) {
	var status models.VLCStatus
	const user = ""
	client := &http.Client{Timeout: 3 * time.Second}

	req, err := http.NewRequest("GET", vlc.StatusEndpoint, nil)
	if err != nil {
		logger.Log.Error("Error building http request", "error", err.Error())
		return status, err
	}

	req.SetBasicAuth(user, vlc.Args.HttpPassword)

	res, err := client.Do(req)
	if err != nil {
		logger.Log.Error("could not connect to VLC's web interface", "error", err.Error())
		return status, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		logger.Log.Error("vlc returned a non-200 status code", "status_code", res.StatusCode)
		return status, err
	}

	if err := json.NewDecoder(res.Body).Decode(&status); err != nil {
		logger.Log.Error("could not decode VLC's JSON response", "error", err.Error())
		return nil, err
	}

	return &status, nil
}

func (vlc *VLCMediaPlayer) Playlist() (models.PlaylistMessage, error) {
	var playlist models.VLCPlaylistNode
	const user = ""
	client := &http.Client{Timeout: 3 * time.Second}

	req, err := http.NewRequest("GET", vlc.PlaylistEndpoint, nil)
	if err != nil {
		logger.Log.Error("Error building http request", "error", err.Error())
		return playlist, err
	}

	req.SetBasicAuth(user, vlc.Args.HttpPassword)

	res, err := client.Do(req)
	if err != nil {
		logger.Log.Error("could not connect to VLC's web interface", "error", err.Error())
		return playlist, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		logger.Log.Error("vlc returned a non-200 status code", "status_code", res.StatusCode)
		return playlist, err
	}

	if err := json.NewDecoder(res.Body).Decode(&playlist); err != nil {
		logger.Log.Error("could not decode VLC's JSON response", "error", err.Error())
		return nil, err
	}

	return &playlist, nil
}

func (vlc *VLCMediaPlayer) LogStatus(s models.StatusMessage) {
	currentTime := fmt.Sprintf("%02d:%02d:%02d", s.GetTime()/3600, (s.GetTime()%3600)/60, s.GetTime()%60)
	totalTime := fmt.Sprintf("%02d:%02d:%02d", s.GetLength()/3600, (s.GetLength()%3600)/60, s.GetLength()%60)
	logger.Log.Info("Pinged", "Filename", s.GetFilename(), "Status", s.GetState(), "Time", currentTime, "Total Time", totalTime)
}
