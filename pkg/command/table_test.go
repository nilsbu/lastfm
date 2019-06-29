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

func TestTable(t *testing.T) {
	user := "TestUser"

	cases := []struct {
		descr          string
		user           *unpack.User
		charts         *charts.Charts
		correctionsRaw []byte
		cmd            command
		table          *format.Table
		ok             bool
	}{
		{
			"no user",
			nil,
			&charts.Charts{"X": []float64{1, 0, 1}},
			[]byte("{}"),
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
			nil,
			[]byte("{}"),
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
			&charts.Charts{"X": []float64{1, 0, 1}},
			nil,
			tableTotal{
				printCharts: printCharts{by: "all", n: 10},
				step:        1,
			},
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
			tableTotal{
				printCharts: printCharts{by: "all", n: 10},
				step:        1,
			},
			&format.Table{
				Charts: charts.Charts{"X": []float64{1, 1, 2}},
				First:  rsrc.ParseDay("2018-01-01"),
				Step:   1,
				Count:  10,
			},
			true,
		},
		{
			"ok, different values", // TODO
			&unpack.User{Name: user, Registered: rsrc.ParseDay("2018-01-01")},
			&charts.Charts{"X": []float64{1, 0, 1}},
			[]byte("{}"),
			tableTotal{
				printCharts: printCharts{by: "all", n: 3},
				step:        2,
			},
			&format.Table{
				Charts: charts.Charts{"X": []float64{1, 1, 2}},
				First:  rsrc.ParseDay("2018-01-01"),
				Step:   2,
				Count:  3,
			},
			true,
		}, {
			"table period; years",
			&unpack.User{Name: user, Registered: rsrc.ParseDay("2017-12-30")},
			&charts.Charts{"X": []float64{1, 0, 1, 5}},
			[]byte("{}"),
			tablePeriods{
				printCharts: printCharts{by: "all", n: 10},
				period:      "y",
			},
			&format.Table{
				Charts: charts.Charts{"X": []float64{1, 6}},
				First:  rsrc.ParseDay("2017-01-01"),
				Step:   1,
				Count:  10,
			},
			true,
		}, {
			"table period; charts broken",
			&unpack.User{Name: user, Registered: rsrc.ParseDay("2017-12-30")},
			&charts.Charts{"X": []float64{1, 0, 1, 5}},
			[]byte("{}"),
			tablePeriods{
				printCharts: printCharts{by: "allxxx", n: 10},
				period:      "y",
			},
			nil, false,
		}, {
			"table period; no user",
			&unpack.User{Name: "no one", Registered: rsrc.ParseDay("2017-12-30")},
			&charts.Charts{"X": []float64{1, 0, 1, 5}},
			[]byte("{}"),
			tablePeriods{
				printCharts: printCharts{by: "all", n: 10},
				period:      "y",
			},
			nil, false,
		}, {
			"table period; false period",
			&unpack.User{Name: user, Registered: rsrc.ParseDay("2017-12-30")},
			&charts.Charts{"X": []float64{1, 0, 1, 5}},
			[]byte("{}"),
			tablePeriods{
				printCharts: printCharts{by: "all", n: 10},
				period:      "invalid",
			},
			nil, false,
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
					if !reflect.DeepEqual(d.Msgs[0], c.table) {
						t.Errorf("%v != %v", d.Msgs[0], c.table)
					}
				}
			}
		})
	}
}

func TestParsePeriod(t *testing.T) {
	cases := []struct {
		descr      string
		registered rsrc.Day
		n          int
		intervals  []charts.Interval
		ok         bool
	}{
		{
			"y", rsrc.ParseDay("2007-01-01"), 600,
			[]charts.Interval{
				{Begin: rsrc.ParseDay("2007-01-01"), Before: rsrc.ParseDay("2008-01-01")},
				{Begin: rsrc.ParseDay("2008-01-01"), Before: rsrc.ParseDay("2009-01-01")},
			}, true,
		},
		{
			"y", rsrc.ParseDay("2007-02-01"), 3,
			[]charts.Interval{
				{Begin: rsrc.ParseDay("2007-01-01"), Before: rsrc.ParseDay("2008-01-01")},
			}, true,
		},
		{
			"m", rsrc.ParseDay("2007-02-01"), 30,
			[]charts.Interval{
				{Begin: rsrc.ParseDay("2007-02-01"), Before: rsrc.ParseDay("2007-03-01")},
				{Begin: rsrc.ParseDay("2007-03-01"), Before: rsrc.ParseDay("2007-04-01")},
			}, true,
		},
		{
			"asdasd", rsrc.ParseDay("2007-02-01"), 30,
			[]charts.Interval{},
			false,
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			days := []float64{}
			for i := 0; i < c.n; i++ {
				days = append(days, 0)
			}
			cha := charts.Charts{"x": days}

			intervals, err := parsePeriod(cha, c.registered, c.descr)
			if err != nil && c.ok {
				t.Fatalf("unexpected error: %v", err)
			} else if err == nil && !c.ok {
				t.Fatalf("expected error but none occurred")
			}

			if c.ok {
				if !reflect.DeepEqual(c.intervals, intervals) {
					t.Errorf("%v != %v", c.intervals, intervals)
				}
			}
		})
	}
}
