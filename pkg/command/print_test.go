package command

import (
	"bytes"
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

func breakUp(plays map[string][]float64) (days []map[string]float64) {
	days = []map[string]float64{}

	size := 0
	for _, values := range plays {
		size = len(values)
	}

	if size == 0 {
		return
	}

	for i := 0; i < size; i++ {
		day := map[string]float64{}
		for key, values := range plays {
			day[key] = values[i]
		}
		days = append(days, day)
	}

	return
}

func iotaF(base float64, n int) []float64 {
	nums := make([]float64, n)
	for i := range nums {
		nums[i] = base + float64(i)
	}

	return nums
}

func repeat(x float64, n int) []float64 {
	nums := make([]float64, n)
	for i := range nums {
		nums[i] = float64(x)
	}

	return nums
}

func repeatSongs(songs []charts.Song, n int) [][]charts.Song {
	days := make([][]charts.Song, n)
	for i := range days {
		days[i] = songs
	}
	return days
}

func times(song charts.Song, n int) []charts.Song {
	songs := make([]charts.Song, n)
	for i := range songs {
		songs[i] = song
	}
	return songs
}

func TestPrint(t *testing.T) {
	user := "TestUser"

	tagsX := []unpack.TagCount{
		{Name: "pop", Count: 100},
		{Name: "french", Count: 88}}
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
	tagFrench := &charts.Tag{
		Name:  "french",
		Total: 100,
		Reach: 25,
	}

	cases := []struct {
		descr     string
		user      *unpack.User
		history   [][]charts.Song
		cmd       command
		formatter format.Formatter
		ok        bool
	}{
		{
			"no user",
			nil,
			[][]charts.Song{
				{{Artist: "X", Title: "x"}},
				{},
				{{Artist: "X", Title: "x"}},
			},
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
			"no user (year)",
			nil,
			[][]charts.Song{
				{{Artist: "X", Title: "x"}},
				{},
				{{Artist: "X", Title: "x"}},
			},
			printTotal{
				printCharts: printCharts{
					by:         "year",
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
			[][]charts.Song{
				{{Artist: "X", Title: "x"}},
				{},
				{{Artist: "X", Title: "x"}},
			},
			printTotal{
				printCharts: printCharts{
					by:         "all",
					name:       "",
					percentage: false,
					normalized: false,
					n:          10,
				},
				date: time.Time{},
			},
			&format.Charts{
				Charts: charts.CompileArtists(
					[]map[string]float64{
						map[string]float64{"X": 1},
						map[string]float64{"X": 1},
						map[string]float64{"X": 2},
					}, rsrc.ParseDay("2018-01-01")),
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
			[][]charts.Song{
				{{Artist: "X", Title: "x"}},
				{{Artist: "Y", Title: "y"}},
				{{Artist: "X", Title: "x"}},
			},
			printTotal{
				printCharts: printCharts{
					by:         "super",
					name:       "",
					percentage: false,
					normalized: false,
					n:          10,
				},
				date: time.Time{},
			},
			&format.Charts{
				Charts: charts.CompileArtists(
					[]map[string]float64{
						map[string]float64{"pop": 1, "rock": 0},
						map[string]float64{"pop": 1, "rock": 1},
						map[string]float64{"pop": 2, "rock": 1},
					}, rsrc.ParseDay("2018-01-01")),
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
			[][]charts.Song{
				{{Artist: "X", Title: "x"}, {Artist: "Y", Title: "y"}},
				{{Artist: "Y", Title: "y"}},
				{{Artist: "X", Title: "x"}},
			},
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
				Charts: charts.CompileArtists(
					[]map[string]float64{
						map[string]float64{"Y": 1},
						map[string]float64{"Y": 2},
						map[string]float64{"Y": 2},
					}, rsrc.ParseDay("2018-01-01")),
				Column:     -1,
				Count:      10,
				Numbered:   true,
				Precision:  2,
				Percentage: true,
			},
			true,
		},
		{
			"by country",
			&unpack.User{Name: user, Registered: rsrc.ParseDay("2018-01-01")},
			[][]charts.Song{
				{{Artist: "X", Title: "x"}},
				{{Artist: "Y", Title: "y"}},
				{{Artist: "X", Title: "x"}},
			},
			printTotal{
				printCharts: printCharts{
					by:         "country",
					name:       "",
					percentage: false,
					normalized: false,
					n:          10,
				},
				date: time.Time{},
			},
			&format.Charts{
				Charts: charts.CompileArtists(
					[]map[string]float64{
						map[string]float64{"France": 1, "-": 0},
						map[string]float64{"France": 1, "-": 1},
						map[string]float64{"France": 2, "-": 1},
					}, rsrc.ParseDay("2018-01-01")),
				Column:     -1,
				Count:      10,
				Numbered:   true,
				Precision:  0,
				Percentage: false,
			},
			true,
		},
		{
			"'all' with name invalid",
			&unpack.User{Name: user, Registered: rsrc.ParseDay("2018-01-01")},
			[][]charts.Song{
				{{Artist: "X", Title: "x"}},
				{},
				{{Artist: "X", Title: "x"}},
			},
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
			[][]charts.Song{
				{{Artist: "X", Title: "x"}},
				{},
				{{Artist: "X", Title: "x"}},
			},
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
			[][]charts.Song{
				{{Artist: "X", Title: "x"}},
				{},
				{{Artist: "X", Title: "x"}},
			},
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
			[][]charts.Song{
				{{Artist: "X", Title: "x"}},
				{},
				{{Artist: "X", Title: "x"}},
			},
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
			&unpack.User{Name: user, Registered: rsrc.ParseDay("2017-11-01")},
			append(
				repeatSongs([]charts.Song{{Artist: "X", Title: "x"}}, 30+31+31),
				repeatSongs([]charts.Song{{Artist: "Y", Title: "y"}}, 30+28)...),
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
				Charts: charts.CompileArtists(breakUp(map[string][]float64{
					"2017": append(iotaF(1, 30+31+31), repeat(92, 30+28)...),
					"2018": append(repeat(0, 30+31+31), iotaF(1, 30+28)...)}),
					rsrc.ParseDay("2017-12-30")),
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
			&unpack.User{Name: user, Registered: rsrc.ParseDay("2017-12-01")},
			append(append(
				[][]charts.Song{{{Artist: "X", Title: "x"}, {Artist: "Y", Title: "y"}}},
				repeatSongs([]charts.Song{{Artist: "X", Title: "x"}}, 30)...),
				repeatSongs([]charts.Song{}, 30)...),
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
				Charts: charts.CompileArtists(breakUp(map[string][]float64{
					"X": append(iotaF(1, 31), repeat(31, 30)...)}),
					rsrc.ParseDay("2017-12-30")),
				Column:     -1,
				Count:      10,
				Numbered:   true,
				Precision:  0,
				Percentage: false,
			},
			true,
		},
		{
			"super withno tags",
			&unpack.User{Name: user, Registered: rsrc.ParseDay("2017-12-31")},
			[][]charts.Song{
				{{Artist: "Z", Title: "z"}},
			},
			printTotal{
				printCharts: printCharts{
					by:         "super",
					name:       "",
					percentage: false,
					normalized: false,
					n:          10,
				},
			},
			nil,
			false,
		},
		{
			"country with no tags",
			&unpack.User{Name: user, Registered: rsrc.ParseDay("2017-12-31")},
			[][]charts.Song{
				{{Artist: "Z", Title: "z"}},
			},
			printTotal{
				printCharts: printCharts{
					by:         "country",
					name:       "",
					percentage: false,
					normalized: false,
					n:          10,
				},
			},
			nil,
			false,
		},
		{
			"all regular",
			&unpack.User{Name: user, Registered: rsrc.ParseDay("2018-01-01")},
			[][]charts.Song{
				{{Artist: "X", Title: "x"}},
				{},
				{{Artist: "X", Title: "x"}},
			},
			printTotal{
				printCharts: printCharts{
					by:         "all",
					name:       "",
					percentage: true,
					normalized: false,
					n:          10,
				},
				date: date("2018-01-01"),
			},
			&format.Charts{
				Charts: charts.CompileArtists(breakUp(map[string][]float64{
					"X": []float64{1, 1, 2}}),
					rsrc.ParseDay("2017-12-30")),
				Column:     -1,
				Count:      10,
				Numbered:   true,
				Precision:  2,
				Percentage: true,
			},
			true,
		},
		{
			"total total",
			&unpack.User{Name: user, Registered: rsrc.ParseDay("2018-01-01")},
			[][]charts.Song{
				{{Artist: "X", Title: "x"}},
				{{Artist: "X", Title: "x"}},
				{{Artist: "X", Title: "x"}},
			},
			printTotal{
				printCharts: printCharts{
					by:         "total",
					name:       "",
					percentage: false,
					normalized: false,
					n:          10,
				},
				date: date("2018-01-01"),
			},
			&format.Charts{
				Charts: charts.CompileArtists(breakUp(map[string][]float64{
					"total": []float64{1, 2, 3}}),
					rsrc.ParseDay("2018-01-01")),
				Column:     0,
				Count:      10,
				Numbered:   true,
				Precision:  0,
				Percentage: false,
			},
			true,
		},
		// Fade
		{
			"fade regular",
			&unpack.User{Name: user, Registered: rsrc.ParseDay("2018-01-01")},
			[][]charts.Song{
				{{Artist: "X", Title: "x"}},
				{},
				{{Artist: "X", Title: "x"}},
			},
			printFade{
				printCharts: printCharts{
					by:         "all",
					name:       "",
					percentage: true,
					normalized: false,
					n:          10,
				},
				hl:   1,
				date: date("2018-01-01"),
			},
			&format.Charts{
				Charts: charts.CompileArtists(breakUp(map[string][]float64{
					"X": []float64{1, 0.5, 0.25}}),
					rsrc.ParseDay("2017-12-30")),
				Column:     -1,
				Count:      10,
				Numbered:   true,
				Precision:  2,
				Percentage: true,
			},
			true,
		},
		{
			"fade fail",
			&unpack.User{Name: user, Registered: rsrc.ParseDay("2018-01-01")},
			[][]charts.Song{
				{{Artist: "X", Title: "x"}},
				{},
				{{Artist: "X", Title: "x"}},
			},
			printFade{
				printCharts: printCharts{
					by:         "year",
					name:       "9",
					percentage: true,
					normalized: false,
					n:          10,
				},
				hl:   1,
				date: date("2018-01-01"),
			},
			nil,
			false,
		},
		{
			"fade false user",
			&unpack.User{Name: "no user", Registered: rsrc.ParseDay("2018-01-01")},
			[][]charts.Song{
				{{Artist: "X", Title: "x"}},
				{},
				{{Artist: "X", Title: "x"}},
			},
			printFade{
				printCharts: printCharts{
					by:         "all",
					name:       "",
					percentage: true,
					normalized: false,
					n:          10,
				},
				hl:   1,
				date: date("2018-01-01"),
			},
			nil,
			false,
		},
		{
			"period functional",
			&unpack.User{Name: user, Registered: rsrc.ParseDay("2017-12-31")},
			[][]charts.Song{
				times(charts.Song{Artist: "X", Title: "x"}, 7),
				times(charts.Song{Artist: "X", Title: "x"}, 1),
				times(charts.Song{Artist: "X", Title: "x"}, 8),
			},
			printPeriod{
				printCharts: printCharts{
					by:         "all",
					name:       "",
					percentage: false,
					normalized: false,
					n:          10,
				},
				period: "2018",
			},
			&format.Column{
				Column:     charts.Column{charts.Score{Name: "X", Score: 9}},
				Numbered:   true,
				Precision:  0,
				Percentage: false,
				SumTotal:   9,
			},
			true,
		},
		{
			"period; no charts",
			&unpack.User{Name: user, Registered: rsrc.ParseDay("2017-12-31")},
			[][]charts.Song{
				times(charts.Song{Artist: "X", Title: "x"}, 7),
				times(charts.Song{Artist: "X", Title: "x"}, 1),
				times(charts.Song{Artist: "X", Title: "x"}, 8),
			},
			printPeriod{
				printCharts: printCharts{by: "xx", n: 10},
				period:      "2018",
			},
			nil, false,
		},
		{
			"period; user",
			&unpack.User{Name: "nop", Registered: rsrc.ParseDay("2017-12-31")},
			[][]charts.Song{
				times(charts.Song{Artist: "X", Title: "x"}, 7),
				times(charts.Song{Artist: "X", Title: "x"}, 1),
				times(charts.Song{Artist: "X", Title: "x"}, 8),
			},
			printPeriod{
				printCharts: printCharts{by: "all", n: 10},
				period:      "2018",
			},
			nil, false,
		},
		{
			"period; broken period",
			&unpack.User{Name: user, Registered: rsrc.ParseDay("2017-12-31")},
			[][]charts.Song{
				times(charts.Song{Artist: "X", Title: "x"}, 7),
				times(charts.Song{Artist: "X", Title: "x"}, 1),
				times(charts.Song{Artist: "X", Title: "x"}, 8),
			},
			printPeriod{
				printCharts: printCharts{by: "all", n: 10},
				period:      "I don't work",
			},
			nil, false,
		},
		{
			"period; percentage",
			&unpack.User{Name: user, Registered: rsrc.ParseDay("2017-12-31")},
			[][]charts.Song{
				times(charts.Song{Artist: "X", Title: "x"}, 7),
				times(charts.Song{Artist: "X", Title: "x"}, 1),
				times(charts.Song{Artist: "X", Title: "x"}, 8),
			},
			printPeriod{
				printCharts: printCharts{
					by:         "all",
					name:       "",
					percentage: true,
					normalized: false,
					n:          10,
				},
				period: "2018",
			},
			&format.Column{
				Column:     charts.Column{charts.Score{Name: "X", Score: 9}},
				Numbered:   true,
				Precision:  2,
				Percentage: true,
				SumTotal:   9,
			},
			true,
		},
		{
			"interval; basic",
			&unpack.User{Name: user, Registered: rsrc.ParseDay("2017-12-31")},
			[][]charts.Song{
				times(charts.Song{Artist: "X", Title: "x"}, 7),
				times(charts.Song{Artist: "X", Title: "x"}, 1),
				times(charts.Song{Artist: "X", Title: "x"}, 8),
				times(charts.Song{Artist: "X", Title: "x"}, 99),
			},
			printInterval{
				printCharts: printCharts{
					by:         "all",
					name:       "",
					percentage: false,
					normalized: false,
					n:          10,
				},
				begin:  date("2018-01-01"),
				before: date("2018-01-03"),
			},
			&format.Column{
				Column:     charts.Column{charts.Score{Name: "X", Score: 9}},
				Numbered:   true,
				Precision:  0,
				Percentage: false,
				SumTotal:   9,
			},
			true,
		},
		{
			"interval; no charts",
			&unpack.User{Name: user, Registered: rsrc.ParseDay("2017-12-31")},
			[][]charts.Song{
				times(charts.Song{Artist: "X", Title: "x"}, 7),
				times(charts.Song{Artist: "X", Title: "x"}, 1),
				times(charts.Song{Artist: "X", Title: "x"}, 8),
				times(charts.Song{Artist: "X", Title: "x"}, 99),
			},
			printInterval{
				printCharts: printCharts{
					by:         "sss",
					name:       "",
					percentage: false,
					normalized: false,
					n:          10,
				},
				begin:  date("2018-01-01"),
				before: date("2018-01-03"),
			},
			nil, false,
		},
		{
			"interval; no user",
			&unpack.User{Name: "", Registered: rsrc.ParseDay("2017-12-31")},
			[][]charts.Song{
				times(charts.Song{Artist: "X", Title: "x"}, 7),
				times(charts.Song{Artist: "X", Title: "x"}, 1),
				times(charts.Song{Artist: "X", Title: "x"}, 8),
				times(charts.Song{Artist: "X", Title: "x"}, 99),
			},
			printInterval{
				printCharts: printCharts{
					by:         "all",
					name:       "",
					percentage: false,
					normalized: false,
					n:          10,
				},
				begin:  date("2018-01-01"),
				before: date("2018-01-03"),
			},
			nil, false,
		},
		{
			"interval; percentage",
			&unpack.User{Name: user, Registered: rsrc.ParseDay("2017-12-31")},
			[][]charts.Song{
				times(charts.Song{Artist: "X", Title: "x"}, 7),
				times(charts.Song{Artist: "X", Title: "x"}, 1),
				times(charts.Song{Artist: "X", Title: "x"}, 8),
				times(charts.Song{Artist: "X", Title: "x"}, 99),
			},
			printInterval{
				printCharts: printCharts{
					by:         "all",
					name:       "",
					percentage: true,
					normalized: false,
					n:          10,
				},
				begin:  date("2018-01-01"),
				before: date("2018-01-03"),
			},
			&format.Column{
				Column:     charts.Column{charts.Score{Name: "X", Score: 9}},
				Numbered:   true,
				Precision:  2,
				Percentage: true,
				SumTotal:   9,
			},
			true,
		},
		{
			"songs",
			&unpack.User{Name: user, Registered: rsrc.ParseDay("2017-12-31")},
			[][]charts.Song{
				{{Artist: "A", Title: "d"}, {Artist: "A", Title: "d"}, {Artist: "B", Title: "c"}},
				{{Artist: "A", Title: "d"}, {Artist: "B", Title: "c"}},
			},
			printTotal{
				printCharts: printCharts{
					keys:       "song",
					by:         "all",
					name:       "",
					percentage: false,
					normalized: false,
					n:          10,
				},
			},
			&format.Column{
				Column: charts.Column{
					charts.Score{Name: "A - d", Score: 3},
					charts.Score{Name: "B - c", Score: 2},
				},
				Numbered:   true,
				Precision:  0,
				Percentage: false,
				SumTotal:   5,
			},
			true,
		},
		{
			"songs by super",
			&unpack.User{Name: user, Registered: rsrc.ParseDay("2017-12-31")},
			[][]charts.Song{
				{{Artist: "X", Title: "d"}, {Artist: "X", Title: "d"}, {Artist: "Y", Title: "c"}},
				{{Artist: "X", Title: "d"}, {Artist: "Y", Title: "c"}},
			},
			printTotal{
				printCharts: printCharts{
					keys:       "song",
					by:         "super",
					name:       "pop",
					percentage: false,
					normalized: false,
					n:          10,
				},
			},
			&format.Column{
				Column: charts.Column{
					charts.Score{Name: "X - d", Score: 3},
				},
				Numbered:   true,
				Precision:  0,
				Percentage: false,
				SumTotal:   3,
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
					rsrc.SongHistory(user):       nil,
					rsrc.ArtistCorrections(user): nil,
					rsrc.UserInfo(user):          nil,
					rsrc.ArtistTags("X"):         nil,
					rsrc.ArtistTags("Y"):         nil,
					rsrc.TagInfo("pop"):          nil,
					rsrc.TagInfo("rock"):         nil,
					rsrc.TagInfo("french"):       nil},
				mock.Path)
			s, _ := store.New([][]rsrc.IO{[]rsrc.IO{files}})
			d := mock.NewDisplay()

			unpack.WriteArtistTags("X", tagsX, s)
			unpack.WriteArtistTags("Y", tagsY, s)
			unpack.WriteTagInfo(tagPop, s)
			unpack.WriteTagInfo(tagRock, s)
			unpack.WriteTagInfo(tagFrench, s)

			if c.history != nil {
				err := unpack.WriteSongHistory(c.history, user, s)
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

					buf0, buf1 := new(bytes.Buffer), new(bytes.Buffer)
					c.formatter.Plain(buf0)
					d.Msgs[0].Plain(buf1)

					// TODO checking Plain() is no a sufficient test
					if buf0.String() != buf1.String() {
						t.Errorf("formatter does not match expected: %v != %v", c.formatter, d.Msgs[0])
					}
				}
			}
		})
	}
}

func TestPrintTags(t *testing.T) {
	artist := "X"

	cases := []struct {
		descr     string
		tags      []unpack.TagCount
		cmd       command
		formatter format.Formatter
		ok        bool
	}{
		{
			"artist not available",
			[]unpack.TagCount{{Name: "pop", Count: 100}},
			printTags{artist: "nope"},
			nil,
			false,
		},
		{
			"with tags",
			[]unpack.TagCount{{Name: "pop", Count: 100}},
			printTags{artist: artist},
			&format.Column{
				Column:   charts.Column{charts.Score{Name: "pop", Score: 100}},
				Numbered: true,
			},
			true,
		},
	}

	for _, c := range cases {
		t.Run(c.descr, func(t *testing.T) {
			files, _ := mock.IO(
				map[rsrc.Locator][]byte{
					rsrc.ArtistTags(artist): nil},
				mock.Path)
			s, _ := store.New([][]rsrc.IO{[]rsrc.IO{files}})
			d := mock.NewDisplay()

			unpack.WriteArtistTags(artist, c.tags, s)

			err := c.cmd.Execute(nil, s, d)
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
					if !reflect.DeepEqual(c.formatter, d.Msgs[0]) {
						t.Errorf("formatter does not match expected: %v != %v", c.formatter, d.Msgs[0])
					}
				}
			}
		})
	}
}
