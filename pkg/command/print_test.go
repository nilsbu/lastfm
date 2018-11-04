package command

import (
	"reflect"
	"testing"
	"time"

	"github.com/nilsbu/lastfm/pkg/charts"
	"github.com/nilsbu/lastfm/pkg/format"
	"github.com/nilsbu/lastfm/pkg/rsrc"
	"github.com/nilsbu/lastfm/pkg/store"
	"github.com/nilsbu/lastfm/pkg/unpack"
	"github.com/nilsbu/lastfm/test/mock"
)

func date(str string) time.Time {
	date, _ := time.Parse("2006-01-02", str)
	return date
}

func TestPrintTotal(t *testing.T) {
	user := "TestUser"

	tagsX := []unpack.TagCount{{Name: "pop", Count: 100}}
	tagsY := []unpack.TagCount{{Name: "rock", Count: 100}}

	tagPop := &charts.Tag{
		Name:  "pop",
		Total: 100,
		Reach: 25,
	}
	tagRock := &charts.Tag{
		Name:  "rock",
		Total: 100,
		Reach: 25,
	}

	cases := []struct {
		descr     string
		user      *unpack.User
		charts    *charts.Charts
		cmd       command
		formatter *format.Charts
		ok        bool
	}{
		{
			"no user",
			nil,
			&charts.Charts{"X": []float64{1, 0, 1}},
			printTotal{
				printCharts: printCharts{
					by:         "all",
					name:       "",
					percentage: false,
					normalized: false,
					n:          10,
				},
				date: date("2018-01-01"),
			},
			nil,
			false,
		},
		{
			"no charts",
			&unpack.User{Name: user, Registered: rsrc.ParseDay("2018-01-01")},
			nil,
			printTotal{
				printCharts: printCharts{
					by:         "all",
					name:       "",
					percentage: false,
					normalized: false,
					n:          10,
				},
				date: date("2018-01-01"),
			},
			nil,
			false,
		},
		{
			"no corrections",
			&unpack.User{Name: user, Registered: rsrc.ParseDay("2018-01-01")},
			&charts.Charts{"X": []float64{1, 0, 1}},
			printTotal{
				printCharts: printCharts{
					by:         "all",
					name:       "",
					percentage: false,
					normalized: false,
					n:          10,
				},
				date: date("2018-01-01"),
			},
			&format.Charts{
				Charts:     charts.Charts{"X": []float64{1, 1, 2}},
				Column:     -1,
				Count:      10,
				Numbered:   true,
				Precision:  0,
				Percentage: false,
			},
			true,
		},
		{
			"by super",
			&unpack.User{Name: user, Registered: rsrc.ParseDay("2018-01-01")},
			&charts.Charts{
				"X": []float64{1, 0, 1},
				"Y": []float64{0, 1, 0}},
			printTotal{
				printCharts: printCharts{
					by:         "super",
					name:       "",
					percentage: false,
					normalized: false,
					n:          10,
				},
				date: date("2018-01-01"),
			},
			&format.Charts{
				Charts: charts.Charts{
					"classical":  []float64{0, 0, 0},
					"electronic": []float64{0, 0, 0},
					"folk":       []float64{0, 0, 0},
					"gothic":     []float64{0, 0, 0},
					"hip-hop":    []float64{0, 0, 0},
					"jazz":       []float64{0, 0, 0},
					"metal":      []float64{0, 0, 0},
					"pop":        []float64{1, 1, 2},
					"reggae":     []float64{0, 0, 0},
					"rock":       []float64{0, 1, 1},
					"":           []float64{0, 0, 0},
				},
				Column:     -1,
				Count:      10,
				Numbered:   true,
				Precision:  0,
				Percentage: false,
			},
			true,
		},
		{
			"rock bucket",
			&unpack.User{Name: user, Registered: rsrc.ParseDay("2018-01-01")},
			&charts.Charts{
				"X": []float64{1, 0, 1},
				"Y": []float64{1, 1, 0}},
			printTotal{
				printCharts: printCharts{
					by:         "super",
					name:       "rock",
					percentage: true,
					normalized: false,
					n:          10,
				},
				date: date("2018-01-01"),
			},
			&format.Charts{
				Charts:     charts.Charts{"Y": []float64{1, 2, 2}},
				Column:     -1,
				Count:      10,
				Numbered:   true,
				Precision:  2,
				Percentage: true,
			},
			true,
		},
		{
			"'all' with name invalid",
			&unpack.User{Name: user, Registered: rsrc.ParseDay("2018-01-01")},
			&charts.Charts{"X": []float64{1, 0, 1}},
			printTotal{
				printCharts: printCharts{
					by:   "all",
					name: "rock",
				},
			},
			nil,
			false,
		},
		{
			"by invalid with name empty",
			&unpack.User{Name: user, Registered: rsrc.ParseDay("2018-01-01")},
			&charts.Charts{"X": []float64{1, 0, 1}},
			printTotal{
				printCharts: printCharts{
					by:   "invalid",
					name: "",
				},
			},
			nil,
			false,
		},
		{
			"by invalid with name non-empty",
			&unpack.User{Name: user, Registered: rsrc.ParseDay("2018-01-01")},
			&charts.Charts{"X": []float64{1, 0, 1}},
			printTotal{
				printCharts: printCharts{
					by:   "invalid",
					name: "rock",
				},
			},
			nil,
			false,
		},
		{
			"by super with name invalid",
			&unpack.User{Name: user, Registered: rsrc.ParseDay("2018-01-01")},
			&charts.Charts{"X": []float64{1, 0, 1}},
			printTotal{
				printCharts: printCharts{
					by:   "super",
					name: "what is this?",
				},
			},
			nil,
			false,
		},
		{
			"by year",
			&unpack.User{Name: user, Registered: rsrc.ParseDay("2017-12-31")},
			&charts.Charts{
				"X": []float64{100, 1, 1},
				"Y": []float64{99, 1, 0}},
			printTotal{
				printCharts: printCharts{
					by:         "year",
					name:       "",
					percentage: false,
					normalized: false,
					n:          10,
				},
			},
			&format.Charts{
				Charts: charts.Charts{
					"2017": []float64{100, 101, 102},
					"2018": []float64{99, 100, 100},
					"":     []float64{0, 0, 0},
				},
				Column:     -1,
				Count:      10,
				Numbered:   true,
				Precision:  0,
				Percentage: false,
			},
			true,
		},
		{
			"by year 2017",
			&unpack.User{Name: user, Registered: rsrc.ParseDay("2017-12-31")},
			&charts.Charts{
				"X": []float64{100, 1, 1},
				"Y": []float64{99, 1, 0}},
			printTotal{
				printCharts: printCharts{
					by:         "year",
					name:       "2017",
					percentage: false,
					normalized: false,
					n:          10,
				},
			},
			&format.Charts{
				Charts: charts.Charts{
					"X": []float64{100, 101, 102},
				},
				Column:     -1,
				Count:      10,
				Numbered:   true,
				Precision:  0,
				Percentage: false,
			},
			true,
		},
		// TODO test corrections (in other test)
		// TODO test normalized (in other test)
	}

	for _, c := range cases {
		t.Run(c.descr, func(t *testing.T) {
			files, _ := mock.IO(
				map[rsrc.Locator][]byte{
					rsrc.AllDayPlays(user):       nil,
					rsrc.ArtistCorrections(user): nil,
					rsrc.UserInfo(user):          nil,
					rsrc.ArtistTags("X"):         nil,
					rsrc.ArtistTags("Y"):         nil,
					rsrc.TagInfo("pop"):          nil,
					rsrc.TagInfo("rock"):         nil},
				mock.Path)
			s, _ := store.New([][]rsrc.IO{[]rsrc.IO{files}})

			d := mock.NewDisplay()

			unpack.WriteArtistTags("X", tagsX, s)
			unpack.WriteArtistTags("Y", tagsY, s)
			unpack.WriteTagInfo(tagPop, s)
			unpack.WriteTagInfo(tagRock, s)

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
					if !reflect.DeepEqual(d.Msgs[0], c.formatter) {
						t.Errorf("%v != %v", d.Msgs[0], c.formatter)
					}
				}
			}
		})
	}
}
