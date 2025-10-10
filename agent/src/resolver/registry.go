package resolver

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"villian-couch/common/fs"

	"golang.org/x/sys/windows/registry"
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

func findVlcOnWindows() (path string, found bool, err error) {
	// Paths to check for VLC's Uninstall registry key.
	// We check both the standard path and the path for 32-bit applications on 64-bit Windows.
	registryPaths := []string{
		`SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\VLC media player`,
		`SOFTWARE\WOW6432Node\Microsoft\Windows\CurrentVersion\Uninstall\VLC media player`,
	}

	for _, path := range registryPaths {
		// Attempt to open the key with read-only access.
		key, err := registry.OpenKey(registry.LOCAL_MACHINE, path, registry.QUERY_VALUE)
		if err != nil {
			// If the error is ErrNotExist, this specific key isn't present,
			// which is normal. We can continue to the next path.
			if err == registry.ErrNotExist {
				continue
			}
			// For any other error, something went wrong reading the registry.
			return "", false, fmt.Errorf("failed to open registry key %s: %w", path, err)
		}
		// Defer closing the key handle to ensure it's closed before the function returns.
		defer key.Close()

		// If the key was opened, try to read the "InstallLocation" string value.
		installLocation, _, err := key.GetStringValue("InstallLocation")
		if err == nil {
			// Success! We found the key and the installation path.

			if !strings.Contains(installLocation, "vlc.exe") {
				installLocation = filepath.Join(installLocation, "vlc.exe")
			}
			return installLocation, true, nil
		}
		// If we couldn't get the string value, we'll just continue to the next path.
	}

	// If we looped through all paths and didn't find the key and its value, VLC is not installed.
	return "", false, nil
}

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
