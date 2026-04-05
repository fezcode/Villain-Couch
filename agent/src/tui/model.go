package tui

import (
	"errors"
	"fmt"
	"time"
	mediaplayer "villain-couch/agent/src/media-player"
	"villain-couch/agent/src/models"
	"villain-couch/agent/src/options"
	"villain-couch/agent/src/storage"
	"villain-couch/common/logger"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type TickMsg time.Time

func Tick() tea.Cmd {
	return tea.Every(500*time.Millisecond, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

type Model struct {
	VLC      *mediaplayer.VLCMediaPlayer
	Opts     *options.Options
	Status   models.StatusMessage
	Playlist models.PlaylistMessage
	History  []models.MediaFile
	Progress progress.Model
	Error    error
	LastTick time.Time
	ViewMode string // "now-playing" or "history"
	Cursor   int    // For history navigation
}

func NewModel(vlc *mediaplayer.VLCMediaPlayer, opts *options.Options) Model {
	return Model{
		VLC:      vlc,
		Opts:     opts,
		Progress: progress.New(progress.WithDefaultGradient()),
		ViewMode: "now-playing",
	}
}

func (m Model) Init() tea.Cmd {
	return Tick()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.saveMediaStates()
			return m, tea.Quit
		case "tab":
			if m.ViewMode == "now-playing" {
				m.ViewMode = "history"
				m.History, _ = storage.GetDB().GetMediaHistory()
			} else {
				m.ViewMode = "now-playing"
			}
		case "up", "k":
			if m.ViewMode == "history" && m.Cursor > 0 {
				m.Cursor--
			}
		case "down", "j":
			if m.ViewMode == "history" && m.Cursor < len(m.History)-1 {
				m.Cursor++
			}
		case "enter":
			if m.ViewMode == "history" && len(m.History) > 0 {
				_ = m.VLC.PlayFile(m.History[m.Cursor].Filepath)
				m.ViewMode = "now-playing"
			}
		case " ":
			_ = m.VLC.TogglePause()
		case "n":
			current, _ := m.Playlist.GetCurrent()
			_ = m.VLC.TryNext(current)
		case "s":
			_ = m.VLC.Stop()
		case "l": // Seek forward 10s
			if m.Status != nil {
				newTime := m.Status.GetTime() + 10
				_ = m.VLC.Seek(fmt.Sprintf("%d", newTime))
			}
		case "h": // Seek backward 10s
			if m.Status != nil {
				newTime := m.Status.GetTime() - 10
				if newTime < 0 {
					newTime = 0
				}
				_ = m.VLC.Seek(fmt.Sprintf("%d", newTime))
			}
		case "+":
			_ = m.VLC.Volume("+10")
		case "-":
			_ = m.VLC.Volume("-10")
		}

	case TickMsg:
		m.LastTick = time.Time(msg)
		m.handleTick()
		return m, Tick()

	case progress.FrameMsg:
		newProgressModel, progressCmd := m.Progress.Update(msg)
		m.Progress = newProgressModel.(progress.Model)
		return m, progressCmd
	}

	return m, cmd
}

func (m *Model) handleTick() {
	status, err := m.VLC.Status()
	if err != nil {
		m.Error = err
		return
	}
	m.Status = status

	playlist, err := m.VLC.Playlist()
	if err != nil {
		m.Error = err
		return
	}
	m.Playlist = playlist

	currentFilepath, err := playlist.GetCurrent()
	if err != nil {
		return
	}

	// If file is completed then stop VLC.
	if status.GetState() == models.StateStopped {
		m.saveMediaStates()
		storage.GetCache().Delete(currentFilepath)
		if err := m.VLC.TryNext(currentFilepath); err != nil {
			if errors.Is(err, mediaplayer.ErrorMediaFileNotFound) {
				if m.Opts.FuzzyFoundNextEpisode != "" {
					err := m.VLC.PlayFile(m.Opts.FuzzyFoundNextEpisode)
					if err != nil {
						logger.Log.Error("VLC PlayFile Error on Fuzzy Found Next Episode", "error", err, "path", m.Opts.FuzzyFoundNextEpisode)
						_ = m.VLC.CommandRunner.Stop()
					}
				} else {
					_ = m.VLC.CommandRunner.Stop()
				}
			} else {
				_ = m.VLC.CommandRunner.Stop()
			}
		}
	} else {
		mf := models.NewMediaFileFromStatus(status, currentFilepath)
		storage.GetCache().Set(mf.Filepath, mf)
	}
}

func (m *Model) saveMediaStates() {
	db := storage.GetDB()
	cache := storage.GetCache()
	for _, key := range cache.Keys() {
		if key == "" {
			continue
		}
		val, _ := cache.Get(key)
		_ = db.SetMediaFile(val)
	}
}

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1)

	statusStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#04B575"))

	infoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#3C3C3C"))

	cursorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7D56F4")).
			Bold(true)

	selectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7D56F4"))
)

func (m Model) View() string {
	s := titleStyle.Render("Villain Couch") + "\n\n"

	if m.ViewMode == "history" {
		s += "--- Watching History ---\n\n"
		for i, mf := range m.History {
			cursor := " "
			if m.Cursor == i {
				cursor = cursorStyle.Render(">")
			}
			s += fmt.Sprintf("%s %s\n", cursor, mf.Filename)
		}
		s += "\n" + infoStyle.Render("tab: back | enter: resume | q: quit")
		return s
	}

	if m.Status == nil {
		return s + "Connecting to VLC..."
	}

	state := m.Status.GetState()
	filename := m.Status.GetFilename()
	if filename == "" {
		filename = "No file playing"
	}

	currentTime := m.Status.GetTime()
	totalTime := m.Status.GetLength()
	percent := 0.0
	if totalTime > 0 {
		percent = float64(currentTime) / float64(totalTime)
	}

	currStr := fmt.Sprintf("%02d:%02d:%02d", currentTime/3600, (currentTime%3600)/60, currentTime%60)
	totalStr := fmt.Sprintf("%02d:%02d:%02d", totalTime/3600, (totalTime%3600)/60, totalTime%60)

	s += statusStyle.Render(fmt.Sprintf("State: %s", state)) + "\n"
	s += fmt.Sprintf("Playing: %s\n", filename)
	s += fmt.Sprintf("Time: %s / %s\n", currStr, totalStr)
	s += fmt.Sprintf("Volume: %d%%\n\n", m.Status.GetVolume()*100/256)

	s += m.Progress.ViewAs(percent) + "\n\n"

	s += infoStyle.Render("Controls: space: pause | n: next | h/l: seek | +/-: vol | tab: history | q: quit")

	return s
}
