package config

import (
	"fmt"

	"golang.org/x/sys/windows/registry"
)

// GetVLCInstallLocation checks for VLC in the registry.
// It returns the installation path, a boolean indicating if it was found, and any error that occurred.
func GetVLCInstallLocation() (path string, found bool, err error) {
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
			return installLocation, true, nil
		}
		// If we couldn't get the string value, we'll just continue to the next path.
	}

	// If we looped through all paths and didn't find the key and its value, VLC is not installed.
	return "", false, nil
}
