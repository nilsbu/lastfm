package organize

import (
	"fmt"
	"testing"

	"github.com/nilsbu/fastest"
	"github.com/nilsbu/lastfm/io"
	"github.com/nilsbu/lastfm/testutils"
	"github.com/nilsbu/lastfm/unpack"
)

func TestLoadAPIKey(t *testing.T) {
	ft := fastest.T{T: t}

	testCases := []struct {
		json string
		key  io.APIKey
		err  fastest.Code
	}{
		{"", "", fastest.Fail},
		{`{`, "", fastest.Fail},
		{`{}`, "", fastest.Fail},
		{`{"apikey":"asdf"}`, "asdf", fastest.OK},
	}

	for i, tc := range testCases {
		ft.Seq(fmt.Sprintf("#%v", i), func(ft fastest.T) {
			var r io.Reader
			if tc.json == "" {
				r = testutils.Reader(map[io.Resource][]byte{})
			} else {
				r = testutils.Reader(map[io.Resource][]byte{
					*io.NewAPIKey(): []byte(tc.json)})
			}
			key, err := LoadAPIKey(r)
			ft.Equals(err != nil, tc.err == fastest.Fail)
			ft.Only(tc.err == fastest.OK)
			ft.Equals(key, tc.key)
		})
	}
}

func TestAllDayPlays(t *testing.T) {
	ft := fastest.T{T: t}

	testCases := []struct {
		name     io.Name
		plays    []unpack.DayPlays
		failRead bool
	}{
		{
			"X",
			[]unpack.DayPlays{},
			false,
		},
		{
			"X",
			[]unpack.DayPlays{unpack.DayPlays{"BTS": 2, "XX": 1, "12": 1}},
			false,
		},
		{
			"X",
			[]unpack.DayPlays{
				unpack.DayPlays{"as": 42, "": 1, "12": 100},
				unpack.DayPlays{"ギルガメッシュ": 1000},
			},
			false,
		},

		{
			"X",
			[]unpack.DayPlays{unpack.DayPlays{"BTS": 2, "XX": 1, "12": 1}},
			true,
		},
	}

	for i, tc := range testCases {
		ft.Seq(fmt.Sprintf("#%v", i), func(ft fastest.T) {
			w := testutils.NewWriter(map[io.Resource][]bool{})
			err := WriteAllDayPlays(tc.plays, tc.name, w)
			ft.Nil(err, err)

			rsrc := io.NewAllDayPlays(tc.name)
			written, ok := w.Data[*rsrc]
			ft.True(ok)
			ft.Equals(len(written), 1)

			var r io.Reader
			if tc.failRead {
				r = testutils.Reader{}
			} else {
				r = testutils.Reader{*rsrc: written[0]}
			}
			plays, err := ReadAllDayPlays(tc.name, r)
			ft.Implies(err != nil, tc.failRead, err)
			ft.Implies(err == nil, !tc.failRead)
			ft.Only(err == nil)
			ft.DeepEquals(plays, tc.plays)
		})
	}
}
