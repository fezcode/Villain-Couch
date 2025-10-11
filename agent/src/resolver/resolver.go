package resolver

import (
	"fmt"
	"runtime"
	"villain-couch/common/fs"
)

// GetVLCInstallLocation checks for VLC in the registry.
// It returns the installation path, a boolean indicating if it was found, and any error that occurred.
func GetVLCInstallLocation(p string) (path string, found bool, err error) {
	// Handle Linux and custom provided path
	if p != "" {
		if !fs.FileExists(p) {
			return "", false, nil
		}
		return p, true, nil
	}

	switch runtime.GOOS {
	case "windows":
		return findVlcOnWindows()
	case "darwin": // "darwin" is the Go identifier for macOS
		return findVlcOnDarwin()
	default:
		// Cannot see here if
		// For Linux or other OSes, we can add detection later.
		// For now, it will require the config.json entry.
		return "", false, fmt.Errorf("auto-detection not supported on %s", runtime.GOOS)
	}
}
