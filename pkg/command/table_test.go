package command

import (
	"strings"
	"testing"

	"github.com/nilsbu/lastfm/pkg/charts"
	"github.com/nilsbu/lastfm/pkg/format"
	"github.com/nilsbu/lastfm/pkg/rsrc"
	"github.com/nilsbu/lastfm/pkg/store"
	"github.com/nilsbu/lastfm/pkg/unpack"
	"github.com/nilsbu/lastfm/test/mock"
)

// TODO replace trustRanges with trusted function from charts
func trustRanges(s string, registered rsrc.Day, l int) charts.Ranges {
	ranges, _ := charts.ParseRanges(s, registered, l)
	return ranges
}

func TestTable(t *testing.T) {
	user := "TestUser"

	cases := []struct {
		descr     string
		user      *unpack.User
		history   [][]charts.Song
		hasCharts bool
		cmd       command
		table     *format.Table
		ok        bool
	}{
		{
			"no user",
			nil,
			[][]charts.Song{
				{{Artist: "X"}},
				{},
				{{Artist: "X"}},
			}, true,
			tableTotal{
				printCharts: printCharts{by: "all", n: 10},
				step:        1,
			},
			nil,
			false,
		},
		{
			"no charts",
			&unpack.User{Name: user, Registered: rsrc.ParseDay("2018-01-01")},
			[][]charts.Song{}, false,
			tableTotal{
				printCharts: printCharts{by: "all", n: 10},
				step:        1,
			},
			nil,
			false,
		},
		{
			"no corrections",
			&unpack.User{Name: user, Registered: rsrc.ParseDay("2018-01-01")},
			[][]charts.Song{
				{{Artist: "X"}},
				{},
				{{Artist: "X"}},
			}, true,
			tableTotal{
				printCharts: printCharts{by: "all", n: 10},
				step:        1,
			},
			&format.Table{
				Charts: charts.FromMap(map[string][]float64{
					"X": {1, 1, 2},
				}),
				Ranges: trustRanges("1d", rsrc.ParseDay("2018-01-01"), 3),
			},
			true,
		},
		{
			"ok", // TODO
			&unpack.User{Name: user, Registered: rsrc.ParseDay("2018-01-01")},
			[][]charts.Song{
				{{Artist: "X"}},
				{},
				{{Artist: "X"}},
			}, true,
			tableTotal{
				printCharts: printCharts{by: "all", n: 10},
				step:        1,
			},
			&format.Table{
				Charts: charts.FromMap(map[string][]float64{
					"X": {1, 1, 2},
				}),
				Ranges: trustRanges("1d", rsrc.ParseDay("2018-01-01"), 3),
			},
			true,
		},
		{
			"ok, every other day",
			&unpack.User{Name: user, Registered: rsrc.ParseDay("2018-01-01")},
			[][]charts.Song{
				{{Artist: "X"}},
				{{Artist: "X"}},
			}, true,
			tableTotal{
				printCharts: printCharts{by: "all", n: 3},
				step:        2,
			},
			&format.Table{
				Charts: charts.FromMap(map[string][]float64{
					"X": {1, 2},
				}),
				Ranges: trustRanges("2d", rsrc.ParseDay("2018-01-01"), 3),
			},
			true,
		},
		{
			"table period; years",
			&unpack.User{Name: user, Registered: rsrc.ParseDay("2017-12-30")},
			[][]charts.Song{
				{{Artist: "X"}},
				{},
				{{Artist: "X"}},
				{{Artist: "X"}, {Artist: "X"}, {Artist: "X"}, {Artist: "X"}, {Artist: "X"}},
			}, true,
			tablePeriods{
				printCharts: printCharts{by: "all", n: 10},
				period:      "1y",
			},
			&format.Table{
				Charts: charts.FromMap(map[string][]float64{
					"X": {1, 6},
				}),
				Ranges: trustRanges("1d", rsrc.ParseDay("2017-01-01"), 3),
			},
			true,
		},
		{
			"table period; charts broken",
			&unpack.User{Name: user, Registered: rsrc.ParseDay("2017-12-30")},
			[][]charts.Song{
				{{Artist: "X"}},
				{},
				{{Artist: "X"}},
				{{Artist: "X"}, {Artist: "X"}, {Artist: "X"}, {Artist: "X"}, {Artist: "X"}},
			}, true,
			tablePeriods{
				printCharts: printCharts{by: "allxxx", n: 10},
				period:      "y",
			},
			nil, false,
		},
		{
			"table period; no user",
			&unpack.User{Name: "no one", Registered: rsrc.ParseDay("2017-12-30")},
			[][]charts.Song{
				{{Artist: "X"}},
				{},
				{{Artist: "X"}},
				{{Artist: "X"}, {Artist: "X"}, {Artist: "X"}, {Artist: "X"}, {Artist: "X"}},
			}, true,
			tablePeriods{
				printCharts: printCharts{by: "all", n: 10},
				period:      "y",
			},
			nil, false,
		},
		{
			"table period; false period",
			&unpack.User{Name: user, Registered: rsrc.ParseDay("2017-12-30")},
			[][]charts.Song{
				{{Artist: "X"}},
				{},
				{{Artist: "X"}},
				{{Artist: "X"}, {Artist: "X"}, {Artist: "X"}, {Artist: "X"}, {Artist: "X"}},
			}, true,
			tablePeriods{
				printCharts: printCharts{by: "all", n: 10},
				period:      "invalid",
			},
			nil, false,
		},
	}

	for _, c := range cases {
		t.Run(c.descr, func(t *testing.T) {
			expectedFiles := map[rsrc.Locator][]byte{
				rsrc.SongHistory(user):       nil,
				rsrc.Bookmark(user):          nil,
				rsrc.ArtistCorrections(user): []byte(`{"corrections": {}}`),
				rsrc.UserInfo(user):          nil}

			if c.user != nil && c.hasCharts {
				for i := 0; i < len(c.history); i++ {
					expectedFiles[rsrc.DayHistory(user, c.user.Registered.AddDate(0, 0, i))] = nil
				}
			}

			files, _ := mock.IO(expectedFiles, mock.Path)
			s, _ := store.New([][]rsrc.IO{{files}})

			d := mock.NewDisplay()

			if c.hasCharts && c.user != nil {
				unpack.WriteBookmark(c.user.Registered.AddDate(0, 0, len(c.history)-1), user, s)

				for i, day := range c.history {
					err := unpack.WriteDayHistory(day, user, c.user.Registered.AddDate(0, 0, i), s)
					if err != nil {
						t.Fatalf("unexpected error: %v", err)
					}
				}
			}

			if c.user != nil {
				err := unpack.WriteUserInfo(c.user, s)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			}

			err := c.cmd.Execute(&unpack.SessionInfo{User: user}, s, d)
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
					var sb0 strings.Builder
					c.table.Plain(&sb0)
					var sb1 strings.Builder
					d.Msgs[0].Plain(&sb1)
					if sb0.String() != sb1.String() {
						t.Errorf("actual does not match expected:\n%v----------\n%v", sb1.String(), sb0.String())
					}
				}
			}
		})
	}
}
