package unpack

import (
	"reflect"
	"testing"

	"github.com/nilsbu/lastfm/pkg/charts"
	"github.com/nilsbu/lastfm/pkg/rsrc"
	"github.com/nilsbu/lastfm/test/mock"
)

func TestLoadAllDayPlays(t *testing.T) {
	cases := []struct {
		plays  []charts.Charts
		write  bool
		readOK bool
	}{
		{
			[]charts.Charts{charts.Charts{"ABC": []float64{34}}},
			false, false,
		},
		{
			[]charts.Charts{
				charts.Charts{
					"ABC":    []float64{34},
					"|xöü#ß": []float64{1}},
				charts.Charts{
					"<<><": []float64{9999},
					"ABC":  []float64{8}},
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
