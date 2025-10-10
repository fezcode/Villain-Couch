package models

// VLCStatus defines the structure of the JSON response from VLC's status endpoint.
type VLCStatus struct {
	State       string `json:"state"`
	Time        int    `json:"time"`
	Length      int    `json:"length"`
	Information struct {
		Category struct {
			Meta struct {
				ShowName      string `json:"showName"`
				Filename      string `json:"filename"`
				Title         string `json:"title"`
				EpisodeNumber string `json:"episodeNumber"`
				SeasonNumber  string `json:"seasonNumber"`
			} `json:"meta"`
		} `json:"category"`
	} `json:"information"`
}
