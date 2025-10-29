package media_player

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
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
	TryNext() error
	LogStatus(s models.StatusMessage)
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

func (vlc *VLCMediaPlayer) PlayFile(filepath string) error {
	// 1. Convert the OS-specific file path to a valid URI.
	uri := encoding.FormatFileURI(filepath)
	const user = ""
	client := &http.Client{Timeout: 3 * time.Second}

	// 2. Construct the full URL with the correct command and input.
	// The input parameter must be URL-encoded.
	requestURL := fmt.Sprintf("%s?command=in_play&input=%s", vlc.StatusEndpoint, url.QueryEscape(uri))

	// 3. Create the HTTP GET request.
	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// 4. Set the Basic Authentication header. The username is blank.
	req.SetBasicAuth(user, vlc.Args.HttpPassword)

	// 5. Execute the request.
	resp, err := client.Do(req)
	if err != nil {
		logger.Log.Error("could not connect to VLC's web interface", "error", err.Error())
		return err
	}
	defer resp.Body.Close()

	// 6. Check for a successful status code.
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("VLC returned non-200 status: %s", resp.Status)
	}

	logger.Log.Info("played file", "file", filepath)

	err = vlc.SeekSecond("1")
	if err != nil {
		logger.Log.Error("could not seek second", "error", err.Error())
		return err
	}

	return nil
}

func (vlc *VLCMediaPlayer) SeekSecond(second string) error {
	client := &http.Client{Timeout: 3 * time.Second}

	requestURL := fmt.Sprintf("%s?command=seek&val=%s", vlc.StatusEndpoint, second)

	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.SetBasicAuth("", vlc.Args.HttpPassword)

	resp, err := client.Do(req)
	if err != nil {
		logger.Log.Error("could not connect to VLC's web interface", "error", err.Error())
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("VLC returned non-200 status: %s", resp.Status)
	}

	logger.Log.Info("seeked", "second", second)

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
