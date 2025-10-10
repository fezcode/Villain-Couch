package config

import (
	_ "embed"
	"encoding/json"
	"errors"
	"os"
	"villian-couch/common/globals"
	"villian-couch/common/logger"
	"villian-couch/common/step"
)

//go:embed config.json
var defaultSettings string

var appConfig *Config

func GetConfig() *Config {
	return appConfig
}

// Config holds the application configuration loaded from config.json.
type Config struct {
	WebUrl           string `json:"web_url"`
	StatusEndpoint   string `json:"status_endpoint"`
	PlaylistEndpoint string `json:"playlist_endpoint"`
	HttpPort         string `json:"http_port"`
	ExtraIntf        string `json:"extra_intf"`
	HttpPassword     string `json:"http_password"`
	DatabaseFileName string `json:"database_file_name"`
}

// setupConfig ensures the required configuration directory and the config file exist.
func setupConfig(...string) error {
	// 1. Get the user's home directory to resolve the '~' character.
	configDir, configFile, err := globals.GetConfigPaths()
	if err != nil {
		logger.Log.Error(err.Error(), "msg", "Error getting config paths")
		return err
	}

	logger.Log.Info("Checking for directory", "config dir", configDir)

	// 3. Check if the directory exists.
	// os.Stat returns an error if the path does not exist.
	if _, err := os.Stat(configDir); errors.Is(err, os.ErrNotExist) {
		// The directory does not exist, so we create it.
		logger.Log.Info("Directory not found, creating it...")

		// os.MkdirAll creates the directory including any necessary parents.
		// 0755 is a common permission mode (owner: rwx, group: r-x, other: r-x).
		if err := os.MkdirAll(configDir, 0755); err != nil {
			logger.Log.Error("failed to create directory", "config dir", configDir, "error", err)
			return err
		}
		logger.Log.Info("Directory created successfully.")
	} else if err != nil {
		// An error other than "not exist" occurred (e.g., permissions).
		logger.Log.Error("error checking directory %s: %w", configDir, err)
		return err
	} else {
		// The directory already exists, no action needed.
		logger.Log.Info("Directory already exists.")
	}

	// 4. Construct the full path for the configuration file.
	logger.Log.Info("Checking for config file", "path", configFile)

	// 5. Check if the file exists.
	if _, err := os.Stat(configFile); errors.Is(err, os.ErrNotExist) {
		// The file does not exist, so we create it.
		logger.Log.Info("Config file not found, creating it...")

		if err := os.WriteFile(configFile, []byte(defaultSettings), 0644); err != nil {
			logger.Log.Error("failed to create file", "config file", configFile, "error", err)
			return err
		}

		logger.Log.Info("Config file created successfully.")
	} else if err != nil {
		// Another error occurred.
		logger.Log.Error("error checking file", "config file", configFile, "error", err)
		return err
	} else {
		// The file already exists.
		logger.Log.Info("Config file already exists.")
	}
	return nil
}

func loadConfig(...string) error {
	_, configFilePath, err := globals.GetConfigPaths()
	if err != nil {
		logger.Log.Error(err.Error(), "msg", "Error getting config file path")
		return err
	}

	configFile, err := os.Open(configFilePath)
	if err != nil {
		logger.Log.Error("could not open config file", "error", err)
		return err
	}
	defer configFile.Close()

	var config Config
	jsonParser := json.NewDecoder(configFile)
	if err = jsonParser.Decode(&config); err != nil {
		logger.Log.Error("could not decode config file", "error", err)
		return err
	}

	appConfig = &config
	return nil
}

func Initialize() error {
	steps := []step.Step{
		{F: setupConfig},
		{F: loadConfig},
	}
	return step.RunSteps(steps)
}
