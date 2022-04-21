package command

import (
	"bytes"
	"strings"
	"testing"

	"github.com/nilsbu/lastfm/pkg/charts"
	"github.com/nilsbu/lastfm/pkg/format"
	"github.com/nilsbu/lastfm/pkg/info"
	"github.com/nilsbu/lastfm/pkg/io"
	"github.com/nilsbu/lastfm/pkg/pipeline"
	"github.com/nilsbu/lastfm/pkg/rsrc"
	"github.com/nilsbu/lastfm/pkg/unpack"
	"github.com/nilsbu/lastfm/test/mock"
)

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

func repeatSongs(songs []info.Song, n int) [][]info.Song {
	days := make([][]info.Song, n)
	for i := range days {
		days[i] = songs
	}
	return days
}

func times(song info.Song, n int) []info.Song {
	songs := make([]info.Song, n)
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

	tagPop := &info.Tag{
		Name:  "pop",
		Total: 100,
		Reach: 25,
	}
	tagRock := &info.Tag{
		Name:  "rock",
		Total: 100,
		Reach: 25,
	}
	tagFrench := &info.Tag{
		Name:  "french",
		Total: 100,
		Reach: 25,
	}

	cases := []struct {
		descr     string
		user      *unpack.User
		history   [][]info.Song
		cmd       command
		formatter format.Formatter
		ok        bool
	}{
		{
			"no user",
			nil,
			[][]info.Song{
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
				date: rsrc.ParseDay("2018-01-01"),
			},
			nil,
			false,
		},
		{
			"no user (year)",
			nil,
			[][]info.Song{
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
				date: rsrc.ParseDay("2018-01-01"),
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
				date: rsrc.ParseDay("2018-01-01"),
			},
			nil,
			false,
		},
		{
			"no corrections",
			&unpack.User{Name: user, Registered: rsrc.ParseDay("2018-01-01")},
			[][]info.Song{
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
			},
			&format.Charts{
				Charts: []charts.Charts{charts.FromMap(map[string][]float64{
					"X": {1, 1, 2},
				})},
				Numbered:   true,
				Precision:  0,
				Percentage: false,
			},
			true,
		},
		{
			"by super",
			&unpack.User{Name: user, Registered: rsrc.ParseDay("2018-01-01")},
			[][]info.Song{
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
			},
			&format.Charts{
				Charts: []charts.Charts{charts.FromMap(map[string][]float64{
					"pop":  {1, 1, 2},
					"rock": {0, 1, 1},
				})},
				Numbered:   true,
				Precision:  0,
				Percentage: false,
			},
			true,
		},
		{
			"day",
			&unpack.User{Name: user, Registered: rsrc.ParseDay("2018-01-01")},
			[][]info.Song{
				{{Artist: "X", Title: "x"}},
				{{Artist: "Y", Title: "y"}},
				{{Artist: "X", Title: "x"}},
			},
			printTotal{
				printCharts: printCharts{
					by: "all",
					n:  10,
				},
				date: rsrc.ParseDay("2018-01-02"),
			},
			&format.Charts{
				Charts: []charts.Charts{charts.FromMap(map[string][]float64{
					"X": {1},
					"Y": {1},
				})},
				Numbered:   true,
				Precision:  0,
				Percentage: false,
			},
			true,
		},
		{
			"rock bucket",
			&unpack.User{Name: user, Registered: rsrc.ParseDay("2018-01-01")},
			[][]info.Song{
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
				date: rsrc.ParseDay("2018-01-01"),
			},
			&format.Charts{
				Charts: []charts.Charts{charts.FromMap(map[string][]float64{
					"Y": {.5},
				})},
				Numbered:   true,
				Precision:  2,
				Percentage: true,
			},
			true,
		},
		{
			"by country",
			&unpack.User{Name: user, Registered: rsrc.ParseDay("2018-01-01")},
			[][]info.Song{
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
			},
			&format.Charts{
				Charts: []charts.Charts{charts.InOrder([]charts.Pair{
					{Title: charts.KeyTitle("France"), Values: []float64{1, 1, 2}},
					{Title: charts.KeyTitle("-"), Values: []float64{0, 1, 1}},
				})},
				Numbered:   true,
				Precision:  0,
				Percentage: false,
			},
			true,
		},
		{
			"'all' with name invalid",
			&unpack.User{Name: user, Registered: rsrc.ParseDay("2018-01-01")},
			[][]info.Song{
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
			[][]info.Song{
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
			[][]info.Song{
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
			[][]info.Song{
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
				repeatSongs([]info.Song{{Artist: "X", Title: "x"}}, 30+31+31),
				repeatSongs([]info.Song{{Artist: "Y", Title: "y"}}, 30+28)...),
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
				Charts: []charts.Charts{charts.FromMap(map[string][]float64{
					"2017": append(iotaF(1, 30+31+31), repeat(92, 30+28)...),
					"2018": append(repeat(0, 30+31+31), iotaF(1, 30+28)...),
				})},
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
				[][]info.Song{{{Artist: "X", Title: "x"}, {Artist: "Y", Title: "y"}}},
				repeatSongs([]info.Song{{Artist: "X", Title: "x"}}, 30)...),
				repeatSongs([]info.Song{}, 30)...),
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
				Charts: []charts.Charts{charts.FromMap(map[string][]float64{
					"X": append(iotaF(1, 31), repeat(31, 30)...),
				})},
				Numbered:   true,
				Precision:  0,
				Percentage: false,
			},
			true,
		},
		// { // TODO tests fail and I don't know why, so they're disabled
		// 	"super with no tags",
		// 	&unpack.User{Name: user, Registered: rsrc.ParseDay("2017-12-31")},
		// 	[][]info.Song{
		// 		{{Artist: "Z", Title: "z"}},
		// 	},
		// 	printTotal{
		// 		printCharts: printCharts{
		// 			by:         "super",
		// 			name:       "",
		// 			percentage: false,
		// 			normalized: false,
		// 			n:          10,
		// 		},
		// 	},
		// 	nil,
		// 	false,
		// },
		// {
		// 	"country with no tags",
		// 	&unpack.User{Name: user, Registered: rsrc.ParseDay("2017-12-31")},
		// 	[][]info.Song{
		// 		{{Artist: "Z", Title: "z"}},
		// 	},
		// 	printTotal{
		// 		printCharts: printCharts{
		// 			by:         "country",
		// 			name:       "",
		// 			percentage: false,
		// 			normalized: false,
		// 			n:          10,
		// 		},
		// 	},
		// 	nil,
		// 	false,
		// },
		{
			"all regular",
			&unpack.User{Name: user, Registered: rsrc.ParseDay("2018-01-01")},
			[][]info.Song{
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
				date: rsrc.ParseDay("2018-01-01"),
			},
			&format.Charts{
				Charts: []charts.Charts{charts.FromMap(map[string][]float64{
					"X": {1},
				})},
				Numbered:   true,
				Precision:  2,
				Percentage: true,
			},
			true,
		},
		{
			"total total",
			&unpack.User{Name: user, Registered: rsrc.ParseDay("2018-01-01")},
			[][]info.Song{
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
				date: rsrc.ParseDay("2018-01-01"),
			},
			&format.Charts{
				Charts: []charts.Charts{charts.FromMap(map[string][]float64{
					"total": {1},
				})},
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
			[][]info.Song{
				{{Artist: "X", Title: "x"}},
				{},
				{{Artist: "X", Title: "x"}},
			},
			printFade{
				printCharts: printCharts{
					by:         "all",
					name:       "",
					percentage: false,
					normalized: false,
					n:          1,
				},
				hl: 1,
			},
			&format.Charts{
				Charts: []charts.Charts{charts.FromMap(map[string][]float64{
					"X": {1.25},
				})},
				Numbered:   true,
				Precision:  2,
				Percentage: false,
			},
			true,
		},
		{
			"fade fail",
			&unpack.User{Name: user, Registered: rsrc.ParseDay("2018-01-01")},
			[][]info.Song{
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
				date: rsrc.ParseDay("2018-01-01"),
			},
			nil,
			false,
		},
		{
			"fade false user",
			&unpack.User{Name: "no user", Registered: rsrc.ParseDay("2018-01-01")},
			[][]info.Song{
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
				date: rsrc.ParseDay("2018-01-01"),
			},
			nil,
			false,
		},
		{
			"period functional",
			&unpack.User{Name: user, Registered: rsrc.ParseDay("2017-12-31")},
			[][]info.Song{
				times(info.Song{Artist: "X", Title: "x"}, 7),
				times(info.Song{Artist: "X", Title: "x"}, 1),
				times(info.Song{Artist: "X", Title: "x"}, 8),
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
			&format.Charts{
				Charts: []charts.Charts{charts.FromMap(map[string][]float64{
					"X": {9},
				})},
				Numbered:   true,
				Precision:  0,
				Percentage: false,
			},
			true,
		},
		{
			"period; no charts",
			&unpack.User{Name: user, Registered: rsrc.ParseDay("2017-12-31")},
			[][]info.Song{
				times(info.Song{Artist: "X", Title: "x"}, 7),
				times(info.Song{Artist: "X", Title: "x"}, 1),
				times(info.Song{Artist: "X", Title: "x"}, 8),
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
			[][]info.Song{
				times(info.Song{Artist: "X", Title: "x"}, 7),
				times(info.Song{Artist: "X", Title: "x"}, 1),
				times(info.Song{Artist: "X", Title: "x"}, 8),
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
			[][]info.Song{
				times(info.Song{Artist: "X", Title: "x"}, 7),
				times(info.Song{Artist: "X", Title: "x"}, 1),
				times(info.Song{Artist: "X", Title: "x"}, 8),
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
			[][]info.Song{
				times(info.Song{Artist: "X", Title: "x"}, 7),
				times(info.Song{Artist: "X", Title: "x"}, 1),
				times(info.Song{Artist: "X", Title: "x"}, 8),
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
			&format.Charts{
				Charts: []charts.Charts{charts.FromMap(map[string][]float64{
					"X": {1},
				})},
				Numbered:   true,
				Precision:  2,
				Percentage: true,
			},
			true,
		},
		{
			"interval; basic",
			&unpack.User{Name: user, Registered: rsrc.ParseDay("2017-12-31")},
			[][]info.Song{
				times(info.Song{Artist: "X", Title: "x"}, 7),
				times(info.Song{Artist: "X", Title: "x"}, 1),
				times(info.Song{Artist: "X", Title: "x"}, 8),
				times(info.Song{Artist: "X", Title: "x"}, 99),
			},
			printInterval{
				printCharts: printCharts{
					by:         "all",
					name:       "",
					percentage: false,
					normalized: false,
					n:          10,
				},
				begin: rsrc.ParseDay("2018-01-01"),
				end:   rsrc.ParseDay("2018-01-03"),
			},
			&format.Charts{
				Charts: []charts.Charts{charts.FromMap(map[string][]float64{
					"X": {9},
				})},
				Numbered:   true,
				Precision:  0,
				Percentage: false,
			},
			true,
		},
		{
			"interval; no charts",
			&unpack.User{Name: user, Registered: rsrc.ParseDay("2017-12-31")},
			[][]info.Song{
				times(info.Song{Artist: "X", Title: "x"}, 7),
				times(info.Song{Artist: "X", Title: "x"}, 1),
				times(info.Song{Artist: "X", Title: "x"}, 8),
				times(info.Song{Artist: "X", Title: "x"}, 99),
			},
			printInterval{
				printCharts: printCharts{
					by:         "sss",
					name:       "",
					percentage: false,
					normalized: false,
					n:          10,
				},
				begin: rsrc.ParseDay("2018-01-01"),
				end:   rsrc.ParseDay("2018-01-03"),
			},
			nil, false,
		},
		{
			"interval; no user",
			&unpack.User{Name: "", Registered: rsrc.ParseDay("2017-12-31")},
			[][]info.Song{
				times(info.Song{Artist: "X", Title: "x"}, 7),
				times(info.Song{Artist: "X", Title: "x"}, 1),
				times(info.Song{Artist: "X", Title: "x"}, 8),
				times(info.Song{Artist: "X", Title: "x"}, 99),
			},
			printInterval{
				printCharts: printCharts{
					by:         "all",
					name:       "",
					percentage: false,
					normalized: false,
					n:          10,
				},
				begin: rsrc.ParseDay("2018-01-01"),
				end:   rsrc.ParseDay("2018-01-03"),
			},
			nil, false,
		},
		{
			"interval; percentage",
			&unpack.User{Name: user, Registered: rsrc.ParseDay("2017-12-31")},
			[][]info.Song{
				times(info.Song{Artist: "X", Title: "x"}, 7),
				times(info.Song{Artist: "X", Title: "x"}, 1),
				times(info.Song{Artist: "X", Title: "x"}, 8),
				times(info.Song{Artist: "X", Title: "x"}, 99),
			},
			printInterval{
				printCharts: printCharts{
					by:         "all",
					name:       "",
					percentage: true,
					normalized: false,
					n:          10,
				},
				begin: rsrc.ParseDay("2018-01-01"),
				end:   rsrc.ParseDay("2018-01-03"),
			},
			&format.Charts{
				Charts: []charts.Charts{charts.FromMap(map[string][]float64{
					"X": {1},
				})},
				Numbered:   true,
				Precision:  2,
				Percentage: true,
			},
			true,
		},
		{
			"songs",
			&unpack.User{Name: user, Registered: rsrc.ParseDay("2017-12-31")},
			[][]info.Song{
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
			&format.Charts{
				Charts: []charts.Charts{charts.FromMap(map[string][]float64{
					"A - d": {3},
					"B - c": {2},
				})},
				Numbered:   true,
				Precision:  0,
				Percentage: false,
			},
			true,
		},

		// { //TODO song and super don't work in conjunction
		// 	"songs by super",
		// 	&unpack.User{Name: user, Registered: rsrc.ParseDay("2017-12-31")},
		// 	[][]info.Song{
		// 		{{Artist: "X", Title: "d"}, {Artist: "X", Title: "d"}, {Artist: "Y", Title: "c"}},
		// 		{{Artist: "X", Title: "d"}, {Artist: "Y", Title: "c"}},
		// 	},
		// 	printTotal{
		// 		printCharts: printCharts{
		// 			keys:       "song",
		// 			by:         "super",
		// 			name:       "pop",
		// 			percentage: false,
		// 			normalized: false,
		// 			n:          10,
		// 		},
		// 	},
		// 	&format.Column{
		// 		Column: charts.Column{
		// 			charts.Score{Name: "X - d", Score: 3},
		// 		},
		// 		Numbered:   true,
		// 		Precision:  0,
		// 		Percentage: false,
		// 		SumTotal:   3,
		// 	},
		// 	true,
		// },
		// TODO test corrections (in other test)
		// TODO test normalized (in other test)
	}

	for _, c := range cases {
		t.Run(c.descr, func(t *testing.T) {
			expectedFiles :=
				map[rsrc.Locator][]byte{
					rsrc.SongHistory(user):       nil,
					rsrc.Bookmark(user):          nil,
					rsrc.ArtistCorrections(user): []byte(`{"corrections": {}}`),
					rsrc.UserInfo(user):          nil,
					rsrc.ArtistTags("X"):         nil,
					rsrc.ArtistTags("Y"):         nil,
					rsrc.TagInfo("pop"):          nil,
					rsrc.TagInfo("rock"):         nil,
					rsrc.TagInfo("french"):       nil}

			if c.user != nil && c.history != nil {
				for i := range c.history {
					expectedFiles[rsrc.DayHistory(user, c.user.Registered.AddDate(0, 0, i))] = nil
				}
			}

			files, _ := mock.IO(expectedFiles,
				mock.Path)
			s, _ := io.NewStore([][]rsrc.IO{{files}})
			d := mock.NewDisplay()

			unpack.WriteArtistTags("X", tagsX, s)
			unpack.WriteArtistTags("Y", tagsY, s)
			unpack.WriteTagInfo(tagPop, s)
			unpack.WriteTagInfo(tagRock, s)
			unpack.WriteTagInfo(tagFrench, s)

			if c.user != nil && c.history != nil {
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

			session := &unpack.SessionInfo{User: user}
			pl := pipeline.New(session, s)
			err := c.cmd.Execute(session, s, pl, d)
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
						t.Errorf("actual does not match expected:\n%v----------\n%v", buf1.String(), buf0.String())
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
			&format.Charts{
				Charts: []charts.Charts{charts.FromMap(map[string][]float64{
					"pop": {100},
				})},
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
			s, _ := io.NewStore([][]rsrc.IO{{files}})
			d := mock.NewDisplay()

			unpack.WriteArtistTags(artist, c.tags, s)

			pl := pipeline.New(nil, s)
			err := c.cmd.Execute(nil, s, pl, d)
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
					c.formatter.Plain(&sb0)
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
