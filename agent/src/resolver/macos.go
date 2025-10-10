//go:build darwin

package resolver

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// findVlcOnDarwin searches for VLC in standard locations on macOS.
func findVlcOnDarwin() (string, bool, error) {
	// 1. Check the standard /Applications path first.
	standardPath := "/Applications/VLC.app/Contents/MacOS/VLC"
	if _, err := os.Stat(standardPath); err == nil {
		return standardPath, true, nil
	}

	// 2. Fallback to Spotlight search if not in the standard path.
	// This finds VLC even if it's in ~/Applications or elsewhere.
	cmd := exec.Command("mdfind", "kMDItemCFBundleIdentifier == 'org.videolan.vlc'")
	output, err := cmd.Output()
	if err == nil && len(output) > 0 {
		// mdfind returns the path to the .app bundle, so we append the executable path.
		appPath := strings.TrimSpace(string(output))
		vlcPath := filepath.Join(appPath, "Contents", "MacOS", "VLC")
		if _, err := os.Stat(vlcPath); err == nil {
			return vlcPath, true, nil
		}
	}

	return "", false, fmt.Errorf("could not find VLC in /Applications or via Spotlight")
}

func findVlcOnWindows() (string, bool, error) {
	return "", false, nil
}
