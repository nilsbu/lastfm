package organize

import (
	"fmt"
	"testing"

	"github.com/nilsbu/fastest"
	"github.com/nilsbu/lastfm/pkg/rsrc"
	"github.com/nilsbu/lastfm/pkg/unpack"
	"github.com/nilsbu/lastfm/test/mock"
)

func TestLoadAllDayPlays(t *testing.T) {
	ft := fastest.T{T: t}

	testCases := []struct {
		user  unpack.User
		until rsrc.Day
		data  [][]string
		dps   []unpack.DayPlays
		err   fastest.Code
	}{
		{
			unpack.User{Name: "", Registered: rsrc.ToDay(0)},
			rsrc.ToDay(86400),
			[][]string{[]string{}, []string{}},
			nil,
			fastest.Fail,
		},
		{
			unpack.User{Name: "ASDF", Registered: rsrc.ToDay(86400)},
			rsrc.ToDay(2 * 86400),
			[][]string{
				[]string{`{"recenttracks":{"track":[{"artist":{"#text":"ASDF"}}], "@attr":{"totalPages":"1"}}}`},
				[]string{`{"recenttracks":{"track":[{"artist":{"#text":"XXX"}}], "@attr":{"totalPages":"1"}}}`},
			},
			[]unpack.DayPlays{unpack.DayPlays{"ASDF": 1}, unpack.DayPlays{"XXX": 1}},
			fastest.OK,
		},
		{
			unpack.User{Name: "ASDF", Registered: rsrc.ToDay(0)},
			rsrc.ToDay(0),
			[][]string{
				[]string{
					`{"recenttracks":{"track":[{"artist":{"#text":"X"}}], "@attr":{"page":"1","totalPages":"3"}}}`,
					`{"recenttracks":{"track":[{"artist":{"#text":"Y"}}], "@attr":{"page":"2","totalPages":"3"}}}`,
					`{"recenttracks":{"track":[{"artist":{"#text":"Z"}}, {"artist":{"#text":"X"}}], "@attr":{"page":"3","totalPages":"3"}}}`,
				},
			},
			[]unpack.DayPlays{unpack.DayPlays{"X": 2, "Y": 1, "Z": 1}},
			fastest.OK,
		},
	}

	for i, tc := range testCases {
		ft.Seq(fmt.Sprintf("#%v", i), func(ft fastest.T) {
			files := make(map[rsrc.Locator][]byte)
			for j, day := range tc.data {
				for k, d := range day {
					reg, _ := tc.user.Registered.Midnight()
					time := reg + int64(j*86400)
					files[rsrc.History(tc.user.Name, k+1, rsrc.ToDay(time))] = []byte(d)
				}
			}
			io, _ := mock.IO(files, mock.Path)

			dps, err := LoadAllDayPlays(tc.user, tc.until, io)
			ft.Implies(err != nil, tc.err == fastest.Fail, err)
			ft.Implies(err == nil, tc.err == fastest.OK)
			ft.DeepEquals(dps, tc.dps)
		})
	}
}

func TestLoadDayPlays(t *testing.T) {
	ft := fastest.T{T: t}

	testCases := []struct {
		user string
		time rsrc.Day
		data []string
		dp   unpack.DayPlays
		err  fastest.Code
	}{
		{
			"", rsrc.ToDay(86400),
			[]string{},
			nil,
			fastest.Fail,
		},
		{
			"NOP", rsrc.ToDay(86400),
			[]string{
				"FAIL",
			},
			nil,
			fastest.Fail,
		},
		{
			"ASDF", rsrc.ToDay(86400),
			[]string{
				`{"recenttracks":{"track":[{"artist":{"#text":"ASDF"}}], "@attr":{"totalPages":"1"}}}`,
			},
			unpack.DayPlays{"ASDF": 1},
			fastest.OK,
		},
		{
			"XX", rsrc.ToDay(86400),
			[]string{
				`{"recenttracks":{"track":[{"artist":{"#text":"X"}}], "@attr":{"page":"1","totalPages":"3"}}}`,
				`{"recenttracks":{"track":[{"artist":{"#text":"Y"}}], "@attr":{"page":"2","totalPages":"3"}}}`,
				`{"recenttracks":{"track":[{"artist":{"#text":"Z"}}, {"artist":{"#text":"X"}}], "@attr":{"page":"3","totalPages":"3"}}}`,
			},
			unpack.DayPlays{"X": 2, "Y": 1, "Z": 1},
			fastest.OK,
		},
		{
			"XX", rsrc.ToDay(86400),
			[]string{
				`{"recenttracks":{"track":[{"artist":{"#text":"X"}}], "@attr":{"page":"1","totalPages":"2"}}}`,
				`{"recenttracks":{"track":[], "@attr":{"page":"2","totalPages":"2"}}}`,
			},
			unpack.DayPlays{"X": 1},
			fastest.OK,
		},
		{
			"XX", rsrc.ToDay(86400),
			[]string{
				`{"recenttracks":{"track":[{"artist":{"#text":"X"}}], "@attr":{"page":"1","totalPages":"3"}}}`,
				`FAIL`,
			},
			nil,
			fastest.Fail,
		},
	}

	for i, tc := range testCases {
		ft.Seq(fmt.Sprintf("#%v", i), func(ft fastest.T) {
			files := make(map[rsrc.Locator][]byte)
			for j, d := range tc.data {
				time, _ := tc.time.Midnight()
				files[rsrc.History(tc.user, j+1, rsrc.ToDay(time))] = []byte(d)
			}
			io, _ := mock.IO(files, mock.Path)

			dp, err := loadDayPlays(tc.user, tc.time, io)
			ft.Implies(err != nil, tc.err == fastest.Fail, err)
			ft.Implies(err == nil, tc.err == fastest.OK)

			ft.Implies(dp == nil, tc.dp == nil)
			ft.DeepEquals(dp, tc.dp)
		})
	}
}
