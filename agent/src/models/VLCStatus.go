package models

const (
	StateStopped = "stopped"
)

type StatusMessage interface {
	GetState() string
	GetTime() int
	GetLength() int
	GetFilename() string
	GetShowName() string
	GetTitle() string
	GetEpisodeNumber() string
	GetSeasonNumber() string
}

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

func (v VLCStatus) GetState() string {
	return v.State
}

func (v VLCStatus) GetTime() int {
	return v.Time
}

func (v VLCStatus) GetLength() int {
	return v.Length
}

func (v VLCStatus) GetFilename() string {
	return v.Information.Category.Meta.Filename
}

func (v VLCStatus) GetShowName() string {
	return v.Information.Category.Meta.ShowName
}

func (v VLCStatus) GetTitle() string {
	return v.Information.Category.Meta.Title
}

func (v VLCStatus) GetEpisodeNumber() string {
	return v.Information.Category.Meta.EpisodeNumber
}

func (v VLCStatus) GetSeasonNumber() string {
	return v.Information.Category.Meta.SeasonNumber
}
