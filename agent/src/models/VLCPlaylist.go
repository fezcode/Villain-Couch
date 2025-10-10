package models

import (
	"fmt"
	"net/url"
	"runtime"
	"strings"
)

type PlaylistMessage interface {
	GetCurrent() (string, error)
}

// VLCPlaylistNode represents a node in the VLC playlist tree (either a folder or the root).
type VLCPlaylistNode struct {
	ID       string            `json:"id"`
	Name     string            `json:"name"`
	Type     string            `json:"type"`
	Children []VLCPlaylistNode `json:"children"`
	Current  string            `json:"current,omitempty"` // omitempty because it only appears on the current item
	URI      string            `json:"uri,omitempty"`     // omitempty because it only appears on leaves
}

func (node VLCPlaylistNode) GetCurrent() (string, error) {
	// findCurrent recursively searches the tree for the current item.
	uri, found := findCurrent(node)
	if !found {
		return "", nil
	}

	// The URI is in the format "file:///C:/Path/To/File.mkv".
	// We need to parse it to get a clean, OS-specific path.
	parsedURL, err := url.Parse(uri)
	if err != nil {
		return "", fmt.Errorf("failed to parse file URI: %w", err)
	}

	// On Windows, url.Parse leaves a leading slash (e.g., "/C:/...").
	// We need to remove it.
	cleanPath := parsedURL.Path
	if runtime.GOOS == "windows" {
		cleanPath = strings.TrimPrefix(cleanPath, "/")
		cleanPath = strings.ReplaceAll(cleanPath, "/", "\\")
	}

	return cleanPath, nil
}

// findCurrent is a recursive helper function to search the playlist tree.
func findCurrent(node VLCPlaylistNode) (uri string, found bool) {
	// Check if the current node is the one we're looking for.
	if node.Type == "leaf" && node.Current == "current" {
		return node.URI, true
	}

	// If not, recursively search its children.
	for _, child := range node.Children {
		uri, found := findCurrent(child)
		if found {
			return uri, true
		}
	}

	// Not found in this branch of the tree.
	return "", false
}
