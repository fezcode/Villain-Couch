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
)

// --- CONFIGURATION ---
// IMPORTANT: Change these values to match your system.
const (
	// On Windows, this is the typical path. Adjust if yours is different.
	vlcPath = `D:\Apps\VideoLAN\VLC\vlc.exe`

	// The full path to the movie or episode you want to play.
	mediaFile = `D:\Downloads\The.Bear.S01.COMPLETE.1080p.HULU.WEB.H264-CAKES[TGx]\The.Bear.S01E06.1080p.WEB.H264-CAKES.mkv`

	// The port for VLC's web interface. 8080 is the default, but we'll use 9090 to be safe.
	httpPort = "6845"

	// The password for VLC's web interface. Keep it simple for this local tool.
	httpPassword = "my_secret_password"
)

// VlcStatus defines the structure of the JSON response from VLC's status endpoint.
// We only map the fields we care about.
type VlcStatus struct {
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

func main() {
	// --- 1. Launch VLC as a child process ---
	fmt.Println("Launching VLC...")
	cmd := exec.Command(
		vlcPath,
		" ",
		mediaFile,             // The file to open
		"--extraintf", "http", // Enable the web interface
		"--http-port", httpPort, // Set the custom port
		"--http-password", httpPassword, // Set the password
	)

	// Start the command. We use Start() instead of Run() because we don't
	// want to block the Go program while VLC is running.
	err := cmd.Start()
	if err != nil {
		log.Fatalf("Failed to start VLC. Is the vlcPath correct? Error: %v", err)
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

	// --- 2. Setup graceful shutdown ---
	// We want to catch Ctrl+C so we can clean up properly.
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Give VLC a couple of seconds to start its web server.
	fmt.Println("Waiting for VLC to initialize its web interface...")
	time.Sleep(3 * time.Second)

	// --- 3. Start the polling loop ---
	// A ticker is a channel that sends a "tick" at a specified interval.
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	fmt.Println("Starting VLC status polling. Press Ctrl+C to exit.")

	for {
		select {
		case <-ticker.C:
			// Every 5 seconds, poll VLC
			status, err := pollVlcStatus()
			if err != nil {
				log.Printf("Error polling VLC: %v", err)
				// We continue here, as VLC might just be starting up or temporarily unresponsive.
				continue
			}
			// If polling was successful, display the status.
			displayStatus(status)
		case <-sigChan:
			// The user pressed Ctrl+C. The deferred function will handle cleanup.
			fmt.Println("\nShutdown signal received.")
			return
		}
	}
}

// pollVlcStatus makes an HTTP request to the running VLC instance to get its status.
func pollVlcStatus() (*VlcStatus, error) {
	// The username for VLC's basic auth is empty.
	const user = ""

	client := &http.Client{Timeout: 3 * time.Second}
	url := fmt.Sprintf("http://localhost:%s/requests/status.json", httpPort)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("could not create request: %w", err)
	}

	// VLC's web interface uses Basic Authentication.
	req.SetBasicAuth(user, httpPassword)

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not connect to VLC's web interface: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("vlc returned a non-200 status code: %s", res.Status)
	}

	var status VlcStatus
	if err := json.NewDecoder(res.Body).Decode(&status); err != nil {
		return nil, fmt.Errorf("could not decode VLC's JSON response: %w", err)
	}

	return &status, nil
}

// displayStatus clears the console and prints the current playback status.
func displayStatus(status *VlcStatus) {
	// Simple console clear for a cleaner look.
	clearConsole()

	// Format time from seconds to MM:SS format
	currentTime := fmt.Sprintf("%02d:%02d", status.Time/60, status.Time%60)
	totalTime := fmt.Sprintf("%02d:%02d", status.Length/60, status.Length%60)

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
		cmd.Run()
	} else {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}
