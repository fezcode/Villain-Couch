//	 ff - Fuzzy Finder
//
//	 Example usage:
//
//	 	baseDirectory := "D:\\Downloads"
//		targetFilename := "The.Awesome.Show.S01E09.1080p.WEB.H264-WESTSIDEGUNN.mkv"
//
//		log.Println("Starting fuzzy search...")
//
//		// --- EXECUTION ---
//		relatedFiles, err := FindRelatedFiles(baseDirectory, targetFilename)
//		if err != nil {
//			log.Fatalf("Error: %v", err)
//		}
//
//		// --- DISPLAY RESULTS ---
//		fmt.Println("\n--- Search Results ---")
//		if len(relatedFiles) == 0 {
//			fmt.Println("No related directories found.")
//		} else {
//			for dir, files := range relatedFiles {
//				fmt.Printf("\n[+] Files in: %s\n", dir)
//				if len(files) == 0 {
//					fmt.Println("  (No files found in this directory)")
//				} else {
//					for _, file := range files {
//						// Print just the filename for cleaner output
//						fmt.Printf("  - %s\n", filepath.Base(file))
//					}
//				}
//			}
//		}
//		fmt.Println("\n--- Search Complete ---")
//
//		// --- FIND NEXT EPISODE ---
//		log.Println("\nSearching for the next episode...")
//		targetInfo, err := ParseEpisodeInfo(targetFilename)
//		if err != nil {
//			log.Fatalf("Could not parse target filename to find next episode: %v", err)
//		}
//
//		nextEp, found := FindNextEpisode(targetInfo, relatedFiles)
//
//		// --- DISPLAY NEXT EPISODE RESULT ---
//		fmt.Println("\n--- Next Episode ---")
//		if found {
//			fmt.Printf("Next episode found!\n")
//			fmt.Printf("  Show:    %s\n", nextEp.ShowName)
//			fmt.Printf("  Season:  %d\n", nextEp.Season)
//			fmt.Printf("  Episode: %d\n", nextEp.Episode)
//			fmt.Printf("  File:    %s\n", nextEp.FilePath)
//		} else {
//			fmt.Println("Could not find a subsequent episode in the related files.")
//			fmt.Println("You might be on the last available episode.")
//		}
//		fmt.Println("\n--- Search Complete ---")
package ff

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"villian-couch/common/logger"
)

// Enhanced regex to capture Show Name (Group 1), Season (Group 3), and Episode (Group 5).
var episodeRegex = regexp.MustCompile(`^(.*?)[\._ ]?(S|s)?(\d{1,2})(E|e|X|x)(\d{1,2})`)
var bracketRegex = regexp.MustCompile(`\[.*?\]`)

// ParseEpisodeInfo attempts to parse a filename into an EpisodeInfo struct.
func ParseEpisodeInfo(filePath string) (EpisodeInfo, error) {
	filename := filepath.Base(filePath)
	matches := episodeRegex.FindStringSubmatch(filename)

	if len(matches) < 6 {
		return EpisodeInfo{}, fmt.Errorf("could not parse episode info from: %s", filename)
	}

	showName := matches[1]
	seasonStr := matches[3]
	episodeStr := matches[5]

	season, err := strconv.Atoi(seasonStr)
	if err != nil {
		return EpisodeInfo{}, fmt.Errorf("could not parse season number '%s'", seasonStr)
	}

	episode, err := strconv.Atoi(episodeStr)
	if err != nil {
		return EpisodeInfo{}, fmt.Errorf("could not parse episode number '%s'", episodeStr)
	}

	showName = strings.ReplaceAll(showName, ".", " ")
	showName = strings.ReplaceAll(showName, "_", " ")
	showName = strings.TrimSpace(showName)

	return EpisodeInfo{
		ShowName: showName,
		Season:   season,
		Episode:  episode,
		FilePath: filePath,
	}, nil
}

// FindNextEpisode searches through a list of all found files to find the next episode.
// This version has been corrected with more robust logic.
func FindNextEpisode(targetInfo EpisodeInfo, allFiles map[string][]string) (EpisodeInfo, bool) {
	var potentialEpisodes []EpisodeInfo
	normalizedTargetName := normalizeString(targetInfo.ShowName)

	for _, files := range allFiles {
		for _, file := range files {
			info, err := ParseEpisodeInfo(file)
			if err != nil {
				continue
			}

			// Add any episode from the same show to our list of candidates.
			if normalizeString(info.ShowName) == normalizedTargetName {
				potentialEpisodes = append(potentialEpisodes, info)
			}
		}
	}

	// Sort all found episodes for the show chronologically.
	sort.Slice(potentialEpisodes, func(i, j int) bool {
		if potentialEpisodes[i].Season != potentialEpisodes[j].Season {
			return potentialEpisodes[i].Season < potentialEpisodes[j].Season
		}
		return potentialEpisodes[i].Episode < potentialEpisodes[j].Episode
	})

	// --- CORRECTED LOGIC ---
	// Instead of looking FOR the target episode in the list, we now look for the
	// FIRST episode in the sorted list that comes chronologically AFTER the target.
	for _, candidate := range potentialEpisodes {
		isLaterSeason := candidate.Season > targetInfo.Season
		isLaterEpisodeInSameSeason := candidate.Season == targetInfo.Season && candidate.Episode > targetInfo.Episode

		if isLaterSeason || isLaterEpisodeInSameSeason {
			// Because the list is sorted, the first candidate that meets this
			// condition is guaranteed to be the next episode.
			return candidate, true
		}
	}

	// If the loop completes, no episode was found after the target.
	return EpisodeInfo{}, false
}

// --- Functions from previous version (normalizeString, FindRelatedFiles) ---
func normalizeString(s string) string {
	s = strings.ToLower(s)
	s = bracketRegex.ReplaceAllString(s, "")
	s = strings.ReplaceAll(s, ".", "")
	s = strings.ReplaceAll(s, "_", "")
	s = strings.ReplaceAll(s, " ", "")
	return strings.TrimSpace(s)
}

func FindRelatedFiles(baseDir, targetFilename string) (map[string][]string, error) {
	targetInfo, err := ParseEpisodeInfo(targetFilename)
	if err != nil {
		return nil, fmt.Errorf("could not parse target filename: %w", err)
	}
	normalizedKey := normalizeString(targetInfo.ShowName)

	logger.Log.Info("extracted", "show name", targetInfo.ShowName, "normalized search key", normalizedKey)

	results := make(map[string][]string)
	entries, err := os.ReadDir(baseDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read base directory %s: %w", baseDir, err)
	}
	logger.Log.Info("entries", "size", len(entries), "directory", baseDir)
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		dirName := entry.Name()
		normalizedDirName := normalizeString(dirName)
		if strings.Contains(normalizedDirName, normalizedKey) {
			matchedDirPath := filepath.Join(baseDir, dirName)
			logger.Log.Info("found matchings", "directory", matchedDirPath)
			var filesInDir []string
			err := filepath.WalkDir(matchedDirPath, func(path string, d os.DirEntry, err error) error {
				if err != nil {
					return err
				}
				if !d.IsDir() {
					filesInDir = append(filesInDir, path)
				}
				return nil
			})
			if err != nil {
				logger.Log.Warn("could not read files", "directory", matchedDirPath, "error", err)
				continue
			}
			results[matchedDirPath] = filesInDir
		}
	}
	return results, nil
}
