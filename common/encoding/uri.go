package encoding

import (
	"path/filepath"
	"runtime"
)

func FormatFileURI(path string) string {
	// Convert backslashes to forward slashes for the URI.
	path = filepath.ToSlash(path)

	// On Windows, a path like "C:/Users/..." needs a leading slash for the URI.
	if runtime.GOOS == "windows" {
		path = "/" + path
	}

	return "file://" + path
}
