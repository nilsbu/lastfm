package organize

import (
	"fmt"
	"testing"

	"github.com/nilsbu/fastest"
	"github.com/nilsbu/lastfm/io"
	"github.com/nilsbu/lastfm/testutils"
	"github.com/nilsbu/lastfm/unpack"
)

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

			dpRes := <-LoadDayPlays(tc.user, tc.time, r)
			ft.Implies(dpRes.Err != nil, tc.err == fastest.Fail, dpRes.Err)
			ft.Implies(dpRes.Err == nil, tc.err == fastest.OK)

			ft.Implies(dpRes.DayPlays == nil, tc.dp == nil)
			ft.DeepEquals(dpRes.DayPlays, tc.dp)
		})
	}
}
