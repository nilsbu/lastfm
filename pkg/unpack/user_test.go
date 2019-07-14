package unpack

import (
	"reflect"
	"testing"

	"github.com/nilsbu/lastfm/pkg/charts"
	"github.com/nilsbu/lastfm/pkg/rsrc"
	"github.com/nilsbu/lastfm/test/mock"
)

func TestBookmarks(t *testing.T) {
	cases := []struct {
		bookmark rsrc.Day
		write    bool
		readOK   bool
	}{
		{
			rsrc.ParseDay("2019-02-01"),
			false, false,
		},
		{
			rsrc.ParseDay("2019-02-01"),
			true, true,
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			io, err := mock.IO(
				map[rsrc.Locator][]byte{rsrc.Bookmark("user"): nil}, mock.Path)
			if err != nil {
				t.Fatal("setup error")
			}

			if c.write {
				err = WriteBookmark(c.bookmark, "user", io)
				if err != nil {
					t.Error("unexpected error during write:", err)
				}
			}

			bookmark, err := LoadBookmark("user", io)
			if err != nil && c.readOK {
				t.Error("unexpected error:", err)
			} else if err == nil && !c.readOK {
				t.Error("expected error but none occurred")
			}

			if err == nil {
				if c.bookmark.Midnight() != bookmark.Midnight() {
					t.Errorf("wrong data\nwant: '%v'\nhas:  '%v'",
						c.bookmark, bookmark)
				}
			}
		})
	}
}

func TestAllDayPlays(t *testing.T) {
	cases := []struct {
		plays  []map[string]float64
		write  bool
		readOK bool
	}{
		{
			[]map[string]float64{{"ABC": 34}},
			false, false,
		},
		{
			[]map[string]float64{
				{
					"ABC":    34,
					"|xöü#ß": 1,
				},
				{
					"<<><": 9999,
					"ABC":  8,
				},
			},
			true, true,
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			io, err := mock.IO(
				map[rsrc.Locator][]byte{rsrc.AllDayPlays("user"): nil}, mock.Path)
			if err != nil {
				t.Fatal("setup error")
			}

			if c.write {
				err = WriteAllDayPlays(c.plays, "user", io)
				if err != nil {
					t.Error("unexpected error during write:", err)
				}
			}

			plays, err := LoadAllDayPlays("user", io)
			if err != nil && c.readOK {
				t.Error("unexpected error:", err)
			} else if err == nil && !c.readOK {
				t.Error("expected error but none occurred")
			}

			if err == nil {
				if !reflect.DeepEqual(plays, c.plays) {
					t.Errorf("wrong data\nhas:  '%v'\nwant: '%v'", plays, c.plays)
				}
			}
		})
	}
}

func TestSongHistory(t *testing.T) {
	cases := []struct {
		plays  [][]charts.Song
		write  bool
		readOK bool
	}{
		{
			[][]charts.Song{},
			true, true,
		},
		{
			[][]charts.Song{
				{{Artist: "ABC", Title: "a", Album: "y"}}},
			false, false,
		},
		{
			[][]charts.Song{
				{
					{Artist: "ABC", Title: "a", Album: "y"},
					{Artist: "ABC", Title: "|xöü#ß", Album: ""},
				},
				{
					{Artist: "<<><", Title: "22", Album: "y"},
				},
			},
			true, true,
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			io, err := mock.IO(
				map[rsrc.Locator][]byte{rsrc.SongHistory("user"): nil}, mock.Path)
			if err != nil {
				t.Fatal("setup error")
			}

			if c.write {
				err = WriteSongHistory(c.plays, "user", io)
				if err != nil {
					t.Error("unexpected error during write:", err)
				}
			}

			plays, err := LoadSongHistory("user", io)
			if err != nil && c.readOK {
				t.Error("unexpected error:", err)
			} else if err == nil && !c.readOK {
				t.Error("expected error but none occurred")
			}

			if err == nil {
				if !reflect.DeepEqual(plays, c.plays) {
					t.Errorf("wrong data\nhas:  '%v'\nwant: '%v'", plays, c.plays)
				}
			}
		})
	}
}

func TestLoadArtistCorrections(t *testing.T) {
	cases := []struct {
		json        []byte
		corrections map[string]string
		ok          bool
	}{
		{
			nil, nil, false,
		},
		{
			[]byte(`{"corrections":{"abc":"x","yy":"x"}}`),
			map[string]string{"abc": "x", "yy": "x"},
			true,
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			io, err := mock.IO(
				map[rsrc.Locator][]byte{rsrc.ArtistCorrections("user"): c.json},
				mock.Path)
			if err != nil {
				t.Fatal("setup error")
			}

			corrections, err := LoadArtistCorrections("user", io)
			if err != nil && c.ok {
				t.Error("unexpected error:", err)
			} else if err == nil && !c.ok {
				t.Error("expected error but none occurred")
			}

			if err == nil {
				if !reflect.DeepEqual(corrections, c.corrections) {
					t.Errorf("wrong data\nhas:  '%v'\nwant: '%v'",
						corrections, c.corrections)
				}
			}
		})
	}
}

func TestLoadSupertagCorrections(t *testing.T) {
	cases := []struct {
		json        []byte
		corrections map[string]string
		ok          bool
	}{
		{
			nil, nil, false,
		},
		{
			[]byte(`{"corrections":{"abc":"x","yy":"x"}}`),
			map[string]string{"abc": "x", "yy": "x"},
			true,
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			io, err := mock.IO(
				map[rsrc.Locator][]byte{rsrc.SupertagCorrections("user"): c.json},
				mock.Path)
			if err != nil {
				t.Fatal("setup error")
			}

			corrections, err := LoadSupertagCorrections("user", io)
			if err != nil && c.ok {
				t.Error("unexpected error:", err)
			} else if err == nil && !c.ok {
				t.Error("expected error but none occurred")
			}

			if err == nil {
				if !reflect.DeepEqual(corrections, c.corrections) {
					t.Errorf("wrong data\nhas:  '%v'\nwant: '%v'",
						corrections, c.corrections)
				}
			}
		})
	}
}
