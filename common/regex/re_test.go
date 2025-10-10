package re

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestValidateAndUnwrapErrors tests that our validate function works as expected
// and that we can correctly iterate through the joined errors.
func TestGetNextEpisodeFilename(t *testing.T) {
	current := `The.Best.TV.Show.Ever.S01E06.1080p.WEB.H264-GOLANG.mp4`
	expected := `The.Best.TV.Show.Ever.S01E07.1080p.WEB.H264-GOLANG.mp4`

	got, ok := GetNextEpisodeFilename(current)

	assert.True(t, ok)
	assert.Equal(t, expected, got)

	fullPath := "D:\\Downloads\\The.Bear.S01.COMPLETE.1080p.HULU.WEB.H264-CAKES[TGx]\\The.Bear.S01E05.1080p.WEB.H264-CAKES.mkv"
	expected = "D:\\Downloads\\The.Bear.S01.COMPLETE.1080p.HULU.WEB.H264-CAKES[TGx]\\The.Bear.S01E06.1080p.WEB.H264-CAKES.mkv"

	got, ok = GetNextEpisodeFilename(fullPath)

	assert.True(t, ok)
	assert.Equal(t, expected, got)

}
