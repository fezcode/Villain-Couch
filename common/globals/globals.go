package globals

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	CONFIG_FOLDER_NAME = ".villain_couch"
	CONFIG_NAME        = "settings.json"
)

func GetConfigPaths() (dir string, filename string, err error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", "", fmt.Errorf("could not get user home directory: %w", err)
	}

	configDir := filepath.Join(homeDir, ".", CONFIG_FOLDER_NAME)
	configFile := filepath.Join(configDir, CONFIG_NAME)

	return configDir, configFile, nil
}
