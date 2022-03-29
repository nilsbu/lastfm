package charts_test

import (
	"testing"

	"github.com/nilsbu/lastfm/pkg/charts"
	"github.com/nilsbu/lastfm/pkg/info"
)

func TestTitles(t *testing.T) {
	for _, c := range []struct {
		name                      string
		title                     charts.Title
		string, key, artist, song string
	}{
		{
			"key title",
			charts.KeyTitle("a"),
			"a", "a", "", "",
		},
		{
			"artist title",
			charts.ArtistTitle("x"),
			"x", "x", "x", "",
		},
		{
			"song title",
			charts.SongTitle(info.Song{Artist: "x", Title: "y"}),
			"x - y", "x\ny", "x", "y",
		},
		{
			"string title",
			charts.StringTitle("a"),
			"a", "a", "", "",
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			if c.string != c.title.String() {
				t.Errorf("String: %v != %v", c.string, c.title.String())
			}
			if c.key != c.title.Key() {
				t.Errorf("Key: %v != %v", c.key, c.title.Key())
			}
			if c.artist != c.title.Artist() {
				t.Errorf("Artist: %v != %v", c.artist, c.title.Artist())
			}
		})
	}
}
