package unpack

import (
	"reflect"
	"testing"

	"github.com/nilsbu/lastfm/pkg/rsrc"
	"github.com/nilsbu/lastfm/test/mock"
)

func TestLoadAllDayPlays(t *testing.T) {
	cases := []struct {
		plays  []PlayCount
		write  bool
		readOK bool
	}{
		{
			[]PlayCount{PlayCount{"ABC": 34}},
			false, false,
		},
		{
			[]PlayCount{
				PlayCount{
					"ABC":    34,
					"|xöü#ß": 1},
				PlayCount{
					"<<><": 9999,
					"ABC":  8},
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
