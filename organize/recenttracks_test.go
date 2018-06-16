package organize

import (
	"fmt"
	"testing"

	"github.com/nilsbu/fastest"
	"github.com/nilsbu/lastfm/io"
	"github.com/nilsbu/lastfm/testutils"
	"github.com/nilsbu/lastfm/unpack"
)

func TestLoadAllDayPlays(t *testing.T) {
	ft := fastest.T{T: t}

	testCases := []struct {
		user  unpack.User
		until io.Midnight
		data  [][]string
		dps   []unpack.DayPlays
		err   fastest.Code
	}{
		{
			unpack.User{Name: "", Registered: 0},
			86400,
			[][]string{[]string{}, []string{}},
			nil,
			fastest.Fail,
		},
		{
			unpack.User{Name: "ASDF", Registered: 86400},
			2 * 86400,
			[][]string{
				[]string{`{"recenttracks":{"track":[{"artist":{"#text":"ASDF"}}], "@attr":{"totalPages":"1"}}}`},
				[]string{`{"recenttracks":{"track":[{"artist":{"#text":"XXX"}}], "@attr":{"totalPages":"1"}}}`},
			},
			[]unpack.DayPlays{unpack.DayPlays{"ASDF": 1}, unpack.DayPlays{"XXX": 1}},
			fastest.OK,
		},
		{
			unpack.User{Name: "ASDF", Registered: 0},
			0,
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
			files := make(map[io.Resource][]byte)
			for j, day := range tc.data {
				for k, d := range day {
					time := tc.user.Registered + io.Midnight(j*86400)
					urt := io.NewUserRecentTracks(tc.user.Name, io.Page(k+1), time)
					files[*urt] = []byte(d)
				}
			}
			r := testutils.AsyncReader(files)

			dps, err := LoadAllDayPlays(tc.user, tc.until, r)
			ft.Implies(err != nil, tc.err == fastest.Fail, err)
			ft.Implies(err == nil, tc.err == fastest.OK)
			ft.DeepEquals(dps, tc.dps)
		})
	}
}

func TestLoadDayPlays(t *testing.T) {
	ft := fastest.T{T: t}

	testCases := []struct {
		user io.Name
		time io.Midnight
		data []string
		dp   unpack.DayPlays
		err  fastest.Code
	}{
		{
			"", 86400,
			[]string{},
			nil,
			fastest.Fail,
		},
		{
			"NOP", 86400,
			[]string{
				"FAIL",
			},
			nil,
			fastest.Fail,
		},
		{
			"ASDF", 86400,
			[]string{
				`{"recenttracks":{"track":[{"artist":{"#text":"ASDF"}}], "@attr":{"totalPages":"1"}}}`,
			},
			unpack.DayPlays{"ASDF": 1},
			fastest.OK,
		},
		{
			"XX", 86400,
			[]string{
				`{"recenttracks":{"track":[{"artist":{"#text":"X"}}], "@attr":{"page":"1","totalPages":"3"}}}`,
				`{"recenttracks":{"track":[{"artist":{"#text":"Y"}}], "@attr":{"page":"2","totalPages":"3"}}}`,
				`{"recenttracks":{"track":[{"artist":{"#text":"Z"}}, {"artist":{"#text":"X"}}], "@attr":{"page":"3","totalPages":"3"}}}`,
			},
			unpack.DayPlays{"X": 2, "Y": 1, "Z": 1},
			fastest.OK,
		},
		{
			"XX", 86400,
			[]string{
				`{"recenttracks":{"track":[{"artist":{"#text":"X"}}], "@attr":{"page":"1","totalPages":"2"}}}`,
				`{"recenttracks":{"track":[], "@attr":{"page":"2","totalPages":"2"}}}`,
			},
			unpack.DayPlays{"X": 1},
			fastest.OK,
		},
		{
			"XX", 86400,
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
			files := make(map[io.Resource][]byte)
			for j, d := range tc.data {
				urt := io.NewUserRecentTracks(tc.user, io.Page(j+1), tc.time)
				files[*urt] = []byte(d)
			}
			r := testutils.AsyncReader(files)

			dp, err := loadDayPlays(tc.user, tc.time, r)
			ft.Implies(err != nil, tc.err == fastest.Fail, err)
			ft.Implies(err == nil, tc.err == fastest.OK)

			ft.Implies(dp == nil, tc.dp == nil)
			ft.DeepEquals(dp, tc.dp)
		})
	}
}
