package organize

import (
	"reflect"
	"testing"

	"github.com/nilsbu/lastfm/pkg/charts"
	"github.com/nilsbu/lastfm/pkg/rsrc"
	"github.com/nilsbu/lastfm/test/mock"
)

func TestLoadArtistSimilar(t *testing.T) {
	for _, c := range []struct {
		descr   string
		files   map[rsrc.Locator][]byte
		artists []charts.Key
		similar map[string]map[string]float32
		hasErrs bool
	}{
		{
			"empty input",
			map[rsrc.Locator][]byte{},
			[]charts.Key{},
			map[string]map[string]float32{},
			false,
		},
		{
			"correct data",
			map[rsrc.Locator][]byte{
				rsrc.ArtistSimilar("X"): []byte(`{"similarartists":{"artist":[{"name":"Kylie Minogue","match":"1"},{"name":"Sido","match":"0.5"}]}}`),
				rsrc.ArtistSimilar("Y"): []byte(`{"similarartists":{"artist":[{"name":"H","match":"1"},{"name":"Sido","match":"0.25"}]}}`)},
			[]charts.Key{
				charts.NewCustomKey("Y", "Y", ""),
				charts.NewCustomKey("X", "X", "")},
			map[string]map[string]float32{
				"X": {"Kylie Minogue": 1, "Sido": .5},
				"Y": {"H": 1, "Sido": .25}},
			false,
		},
		{
			"non-fatal error and correct data",
			map[rsrc.Locator][]byte{
				rsrc.ArtistSimilar("X"): []byte(`{"error":6,"message":"Invalid parameters"}`),
				rsrc.ArtistSimilar("Y"): []byte(`{"similarartists":{"artist":[{"name":"H","match":"1"},{"name":"Sido","match":"0.25"}]}}`)},
			[]charts.Key{
				charts.NewCustomKey("Y", "Y", ""),
				charts.NewCustomKey("X", "X", "")},
			map[string]map[string]float32{
				"X": nil,
				"Y": {"H": 1, "Sido": .25}},
			true,
		},
	} {
		t.Run(c.descr, func(t *testing.T) {
			io, err := mock.IO(c.files, mock.Path)
			if err != nil {
				t.Fatal("setup error")
			}

			similar, err := LoadArtistSimilar(c.artists, io)
			if c.hasErrs && err == nil {
				t.Error("expected error but none occurred")
			} else if !c.hasErrs && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if !reflect.DeepEqual(c.similar, similar) {
				t.Errorf("unexpected data: %v != %v", c.similar, similar)
			}
		})
	}
}
