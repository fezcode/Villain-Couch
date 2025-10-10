package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"vlc-tracker-agent/src/config"
	"vlc-tracker-agent/src/models"
)

func main() {

	// --- 1. Load Configuration ---
	config, err := loadConfig("src/config/config.json")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// --- 2. Get Media File from Command-Line Arguments ---
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s <path/to/mediafile>", os.Args[0])
	}
	mediaFile := os.Args[1]

	// Check if the media file exists before trying to launch VLC.
	if _, err := os.Stat(mediaFile); os.IsNotExist(err) {
		log.Fatalf("Media file not found: %s", mediaFile)
	}

	// --- 3. Launch VLC as a child process ---
	fmt.Println("Launching VLC...")
	cmd := exec.Command(
		config.VlcPath,
		mediaFile, // The file to open (from command-line arg)
		"--extraintf", "http",
		"--http-port", config.HttpPort,
		"--http-password", config.HttpPassword,
	)

	err = cmd.Start()
	if err != nil {
		log.Fatalf("Failed to start VLC. Is the vlc_path in config.json correct? Error: %v", err)
	}
	fmt.Printf("VLC process started with PID: %d\n", cmd.Process.Pid)

	// Ensure we kill the VLC process when our tracker exits.
	defer func() {
		fmt.Println("\nStopping VLC process...")
		if err := cmd.Process.Kill(); err != nil {
			log.Printf("Failed to kill VLC process: %v", err)
		}
		fmt.Println("VLC process stopped.")
	}()

	// --- 4. Setup graceful shutdown ---
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	fmt.Println("Waiting for VLC to initialize its web interface...")
	time.Sleep(3 * time.Second)

	// --- 5. Start the polling loop ---
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	fmt.Println("Starting VLC status polling. Press Ctrl+C to exit.")

	for {
		select {
		case <-ticker.C:
			status, err := pollVlcStatus(config)
			if err != nil {
				log.Printf("Error polling VLC: %v", err)
				continue
			}
			displayStatus(status)
		case <-sigChan:
			fmt.Println("\nShutdown signal received.")
			return
		}
	}
}

// --- HELPER FUNCTIONS ---

// loadConfig reads a JSON file and decodes it into a Config struct.
func loadConfig(filename string) (*config.Config, error) {
	configFile, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("could not open config file: %w", err)
	}
	defer configFile.Close()

	var config config.Config
	jsonParser := json.NewDecoder(configFile)
	if err = jsonParser.Decode(&config); err != nil {
		return nil, fmt.Errorf("could not decode config file: %w", err)
	}

	return &config, nil
}

// pollVlcStatus makes an HTTP request to the running VLC instance.
func pollVlcStatus(config *config.Config) (*models.VLCStatus, error) {
	const user = ""
	client := &http.Client{Timeout: 3 * time.Second}
	url := fmt.Sprintf("http://localhost:%s/requests/status.json", config.HttpPort)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("could not create request: %w", err)
	}

	req.SetBasicAuth(user, config.HttpPassword)

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not connect to VLC's web interface: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("vlc returned a non-200 status code: %s", res.Status)
	}

	var status models.VLCStatus
	if err := json.NewDecoder(res.Body).Decode(&status); err != nil {
		return nil, fmt.Errorf("could not decode VLC's JSON response: %w", err)
	}

	return &status, nil
}

// displayStatus clears the console and prints the current playback status.
func displayStatus(status *models.VLCStatus) {
	clearConsole()

	currentTime := fmt.Sprintf("%02d:%02d:%02d", status.Time/3600, (status.Time%3600)/60, status.Time%60)
	totalTime := fmt.Sprintf("%02d:%02d:%02d", status.Length/3600, (status.Length%3600)/60, status.Length%60)

	fmt.Println("--- VLC Tracker ---")
	fmt.Printf("File:    %s\n", status.Information.Category.Meta.Filename)
	fmt.Printf("Status:  %s\n", status.State)
	fmt.Printf("Time:    %s / %s\n", currentTime, totalTime)
	fmt.Println("-------------------")
	fmt.Println("(Watching for updates... Press Ctrl+C to exit)")
}

// clearConsole clears the terminal screen.
func clearConsole() {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		_ = cmd.Run() // We can ignore the error here
	} else {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		_ = cmd.Run()
	}
}
