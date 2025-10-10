package re

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"
)

// GetNextEpisodeFilename attempts to find a season/episode pattern (e.g., S01E06)
// in a filename and returns the filename for the next episode.
// It returns the new filename and a boolean indicating if the pattern was found.
func GetNextEpisodeFilename(currentFilename string) (string, bool) {
	// Separate the directory and the filename.
	dir := filepath.Dir(currentFilename)
	filename := filepath.Base(currentFilename)

	// Regex to find patterns like S01E06, s01e06, etc.
	// It captures the parts: (S##)(E)(##)
	re := regexp.MustCompile(`(?i)(S\d{2})(E)(\d{2})`)

	// Find the first match in the filename.
	submatches := re.FindStringSubmatch(filename)
	if len(submatches) != 4 {
		// Pattern not found.
		return "", false
	}

	episodeStr := submatches[3]
	episodeNum, err := strconv.Atoi(episodeStr)
	if err != nil {
		// This should not happen with our regex, but it's good practice to check.
		return "", false
	}

	// Increment the episode number.
	nextEpisodeNum := episodeNum + 1

	// Format the new episode number back to a two-digit string (e.g., "07").
	nextEpisodeStr := fmt.Sprintf("%02d", nextEpisodeNum)

	// Use the regex to replace only the episode number part of the matched string.
	// ${1} is the season part (e.g., "S01"), ${2} is "E".
	nextFilename := re.ReplaceAllString(filename, "${1}${2}"+nextEpisodeStr)

	// Combine the original directory with the new filename.
	nextPath := filepath.Join(dir, nextFilename)

	return nextPath, true
}
