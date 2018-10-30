package command

import (
	"reflect"
	"testing"

	"github.com/nilsbu/lastfm/pkg/charts"
	"github.com/nilsbu/lastfm/pkg/format"
	"github.com/nilsbu/lastfm/pkg/rsrc"
	"github.com/nilsbu/lastfm/pkg/store"
	"github.com/nilsbu/lastfm/pkg/unpack"
	"github.com/nilsbu/lastfm/test/mock"
)

func TestTableTotal(t *testing.T) {
	user := "TestUser"

	cases := []struct {
		descr          string
		user           *unpack.User
		charts         *charts.Charts
		correctionsRaw []byte
		n              int
		step           int
		table          *format.Table
		ok             bool
	}{
		{
			"no user",
			nil,
			&charts.Charts{"X": []float64{1, 0, 1}},
			[]byte("{}"),
			10, 1,
			nil,
			false,
		},
		{
			"no charts",
			&unpack.User{Name: user, Registered: rsrc.ParseDay("2018-01-01")},
			nil,
			[]byte("{}"),
			10, 1,
			nil,
			false,
		},
		{
			"no corrections",
			&unpack.User{Name: user, Registered: rsrc.ParseDay("2018-01-01")},
			&charts.Charts{"X": []float64{1, 0, 1}},
			nil,
			10, 1,
			&format.Table{
				Charts: charts.Charts{"X": []float64{1, 1, 2}},
				First:  rsrc.ParseDay("2018-01-01"),
				Step:   1,
				Count:  10,
			},
			true,
		},
		{
			"ok", // TODO
			&unpack.User{Name: user, Registered: rsrc.ParseDay("2018-01-01")},
			&charts.Charts{"X": []float64{1, 0, 1}},
			[]byte("{}"),
			10, 1,
			&format.Table{
				Charts: charts.Charts{"X": []float64{1, 1, 2}},
				First:  rsrc.ParseDay("2018-01-01"),
				Step:   1,
				Count:  10,
			},
			true,
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			files, _ := mock.IO(
				map[rsrc.Locator][]byte{
					rsrc.AllDayPlays(user):       nil,
					rsrc.ArtistCorrections(user): c.correctionsRaw,
					rsrc.UserInfo(user):          nil},

				mock.Path)
			s, _ := store.New([][]rsrc.IO{[]rsrc.IO{files}})

			d := mock.NewDisplay()
			cmd := tableTotal{
				printCharts: printCharts{
					by:         "all",
					name:       "",
					percentage: false,
					normalized: false,
					n:          c.n,
				},
				step: c.step,
			}

			if c.charts != nil {
				err := unpack.WriteAllDayPlays(c.charts.UnravelDays(), user, s)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			}

			if c.user != nil {
				err := unpack.WriteUserInfo(c.user, s)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			}

			err := cmd.Execute(&unpack.SessionInfo{User: user}, s, d)
			if err != nil && c.ok {
				t.Fatalf("unexpected error: %v", err)
			} else if err == nil && !c.ok {
				t.Fatalf("expected error but none occurred")
			}

			if err == nil {
				if len(d.Msgs) == 0 {
					t.Fatalf("no message was printed")
				} else if len(d.Msgs) > 1 {
					t.Fatalf("got %v messages but expected 1", len(d.Msgs))
				} else {
					if !reflect.DeepEqual(d.Msgs[0], c.table) {
						t.Errorf("%v != %v", d.Msgs[0], c.table)
					}
				}
			}
		})
	}
}
