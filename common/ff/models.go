package ff

// EpisodeInfo holds structured data parsed from a media filename.
type EpisodeInfo struct {
	ShowName string
	Season   int
	Episode  int
	FilePath string // Keep track of the original file path
}
