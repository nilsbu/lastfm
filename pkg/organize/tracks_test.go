package organize

import (
	"reflect"
	"testing"

	"github.com/nilsbu/lastfm/pkg/charts"
	"github.com/nilsbu/lastfm/pkg/rsrc"
	"github.com/nilsbu/lastfm/pkg/unpack"
	"github.com/nilsbu/lastfm/test/mock"
)

func TestLoadTrackInfos(t *testing.T) {
	cases := []struct {
		files    map[rsrc.Locator][]byte
		tracks   []charts.Song
		infos    []unpack.TrackInfo
		hasError bool
	}{
		{
			map[rsrc.Locator][]byte{},
			[]charts.Song{},
			[]unpack.TrackInfo{}, false,
		},
		{
			map[rsrc.Locator][]byte{
				rsrc.TrackInfo("a", "b"): nil,
			},
			[]charts.Song{{Artist: "a", Title: "b"}},
			nil, true,
		},
		{
			map[rsrc.Locator][]byte{
				rsrc.TrackInfo("a", "b"): nil,
				rsrc.TrackInfo("a", "c"): nil,
				rsrc.TrackInfo("1", "2"): nil,
			},
			[]charts.Song{
				{Artist: "a", Title: "b"},
				{Artist: "a", Title: "c"},
				{Artist: "1", Title: "2"},
			},
			[]unpack.TrackInfo{
				{Artist: "a", Track: "b", Duration: 2, Listeners: 4, Playcount: 6},
				{Artist: "a", Track: "c", Duration: 6, Listeners: 0, Playcount: 6},
				{Artist: "1", Track: "2", Duration: 2, Listeners: 4, Playcount: 9},
			},
			false,
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			io, err := mock.IO(c.files, mock.Path)
			if err != nil {
				t.Fatal("setup error")
			}

			for _, info := range c.infos {
				if err := unpack.WriteTrackInfo(info.Artist, info.Track, info, io); err != nil {
					t.Fatal(err)
				}
			}

			infos, err := LoadTrackInfos(c.tracks, io)
			if err != nil && !c.hasError {
				t.Error("unexpected error:", err)
			} else if err == nil && c.hasError {
				t.Error("expected error but none occurred")
			}

			if !c.hasError {
				if !reflect.DeepEqual(infos, c.infos) {
					t.Errorf("wrong data:\nhas:  %v\nwant: %v",
						infos, c.infos)
				}
			}
		})
	}
}
