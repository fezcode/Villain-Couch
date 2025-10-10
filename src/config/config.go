package config

// Config holds the application configuration loaded from config.json.
type Config struct {
	VlcPath      string `json:"vlc_path"`
	HttpPort     string `json:"http_port"`
	HttpPassword string `json:"http_password"`
}
