package charts

import "testing"

func TestTitles(t *testing.T) {
	for _, c := range []struct {
		name                      string
		title                     Title
		string, key, artist, song string
	}{
		{
			"key title",
			KeyTitle("a"),
			"a", "a", "", "",
		},
		{
			"artist title",
			ArtistTitle("x"),
			"x", "x", "x", "",
		},
		{
			"song title",
			SongTitle(Song{Artist: "x", Title: "y"}),
			"x - y", "x\ny", "x", "y",
		},
		{
			"string title",
			StringTitle("a"),
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
