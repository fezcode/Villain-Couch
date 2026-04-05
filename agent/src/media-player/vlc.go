package media_player

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"
	"villain-couch/agent/src/cli"
	"villain-couch/agent/src/config"
	"villain-couch/agent/src/models"
	"villain-couch/agent/src/options"
	"villain-couch/common/encoding"
	"villain-couch/common/logger"
	re "villain-couch/common/regex"
)

type MediaPlayer interface {
	Build(*config.Config, *options.Options)
	Status() (models.StatusMessage, error)
	Playlist() (models.PlaylistMessage, error)
	PlayFile(filepath string) error
	TryNext(currentFilepath string) error
	LogStatus(s models.StatusMessage)
	TogglePause() error
	Stop() error
	Seek(val string) error
	Volume(val string) error
	Next() error
	Previous() error
}

var (
	ErrorMediaFileNotFound = errors.New("media file not found")
)

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

func (vlc *VLCMediaPlayer) sendCommand(command string, params map[string]string) error {
	client := &http.Client{Timeout: 3 * time.Second}
	u, err := url.Parse(vlc.StatusEndpoint)
	if err != nil {
		return err
	}
	q := u.Query()
	q.Set("command", command)
	for k, v := range params {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return err
	}
	req.SetBasicAuth("", vlc.Args.HttpPassword)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("VLC returned non-200 status: %s", resp.Status)
	}
	return nil
}

func (vlc *VLCMediaPlayer) TogglePause() error {
	return vlc.sendCommand("pl_pause", nil)
}

func (vlc *VLCMediaPlayer) Stop() error {
	return vlc.sendCommand("pl_stop", nil)
}

func (vlc *VLCMediaPlayer) Seek(val string) error {
	return vlc.sendCommand("seek", map[string]string{"val": val})
}

func (vlc *VLCMediaPlayer) Volume(val string) error {
	return vlc.sendCommand("volume", map[string]string{"val": val})
}

func (vlc *VLCMediaPlayer) Next() error {
	return vlc.sendCommand("pl_next", nil)
}

func (vlc *VLCMediaPlayer) Previous() error {
	return vlc.sendCommand("pl_previous", nil)
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

func (vlc *VLCMediaPlayer) PlayFile(filepathStr string) error {
	// 1. Convert the OS-specific file path to a valid URI.
	uri := encoding.FormatFileURI(filepathStr)

	// 2. Execute the request using sendCommand.
	err := vlc.sendCommand("in_play", map[string]string{"input": uri})
	if err != nil {
		return err
	}

	logger.Log.Info("played file", "file", filepathStr)

	// Wait for VLC to actually switch to the new file
	targetFilename := filepath.Base(filepathStr)
	timeout := time.After(5 * time.Second)
	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()

	// Loop to wait for the file to be loaded
	for {
		select {
		case <-timeout:
			logger.Log.Warn("Timeout waiting for VLC to load new file", "file", targetFilename)
			goto Seek
		case <-ticker.C:
			status, err := vlc.Status()
			if err == nil {
				// We check if the filename matches and if the state is playing
				if status.GetFilename() == targetFilename && status.GetState() == "playing" {
					goto Seek
				}
			}
		}
	}

Seek:
	err = vlc.Seek("1")
	if err != nil {
		logger.Log.Error("could not seek", "error", err.Error())
		return err
	}

	return nil
}

// TryNext can return following errors
// 1. Next episode name is in wrong format
// 2. Next episode media file not found
// 3. VLC API Error (PlayFile)
func (vlc *VLCMediaPlayer) TryNext(currentFilepath string) error {
	nextEpisodeName, ok := re.GetNextEpisodeFilename(currentFilepath)
	if !ok {
		logger.Log.Error("could not find next episode filename", "current", currentFilepath)
		return fmt.Errorf("could not find next episode filename")
	}

	// Check if the media file exists before trying to launch VLC.
	if _, err := os.Stat(nextEpisodeName); os.IsNotExist(err) {
		logger.Log.Warn("Media file not found", "Media File", nextEpisodeName)
		return ErrorMediaFileNotFound
	}

	return vlc.PlayFile(nextEpisodeName)
}

func (vlc *VLCMediaPlayer) LogStatus(s models.StatusMessage) {
	currentTime := fmt.Sprintf("%02d:%02d:%02d", s.GetTime()/3600, (s.GetTime()%3600)/60, s.GetTime()%60)
	totalTime := fmt.Sprintf("%02d:%02d:%02d", s.GetLength()/3600, (s.GetLength()%3600)/60, s.GetLength()%60)
	logger.Log.Info("Pinged", "Filename", s.GetFilename(), "State", s.GetState(), "Time", currentTime, "Total Time", totalTime)
}
