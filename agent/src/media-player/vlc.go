package media_player

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"vlc-tracker-agent/agent/src/cli"
	"vlc-tracker-agent/agent/src/config"
	"vlc-tracker-agent/agent/src/models"
	"vlc-tracker-agent/common/logger"
)

type VLCMediaPlayer struct {
	Args           cli.VLCRunnerArguments
	CommandRunner  *cli.CommandRunner
	StatusEndpoint string
}

type MediaPlayer interface {
	Build(conf *config.Config, flags *cli.CLIFlags)
	Status()
	PrintStatus(models.StatusMessage)
}

func New(conf *config.Config, flags *cli.CLIFlags) VLCMediaPlayer {
	vlc := VLCMediaPlayer{}
	vlc.Build(conf, flags)
	return vlc
}

func (vlc *VLCMediaPlayer) Build(conf *config.Config, flags *cli.CLIFlags) {
	vlc.Args = cli.PrepareRunnerArguments(conf.VlcPath, flags.MediaFile, conf.ExtraIntf, conf.HttpPort, conf.HttpPassword)
	vlc.CommandRunner = cli.NewCommandRunnerForVLC(vlc.Args)
	vlc.StatusEndpoint = fmt.Sprintf("%s:%s/%s", conf.WebUrl, conf.HttpPort, conf.StatusEndpoint)
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

func (vlc *VLCMediaPlayer) PrintStatus(s models.StatusMessage) {
	currentTime := fmt.Sprintf("%02d:%02d:%02d", s.GetTime()/3600, (s.GetTime()%3600)/60, s.GetTime()%60)
	totalTime := fmt.Sprintf("%02d:%02d:%02d", s.GetLength()/3600, (s.GetLength()%3600)/60, s.GetLength()%60)
	logger.Log.Info("Pinged", "Filename", s.GetFilename(), "Status", s.GetState(), "Time", currentTime, "Total Time", totalTime)
}
