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

func TestReadBookmark(t *testing.T) {
	ft := fastest.T{T: t}

	testCases := []struct {
		timestamp int64
		data      string
		readOK    bool
		err       fastest.Code
	}{
		{
			1529246468,
			`{"unixtime":"1529246468","time":"2018-06-17 14:41:08 +0000 UTC"}`,
			true,
			fastest.OK,
		},
		{
			1529250983,
			`{"unixtime":"1529250983","time":"2018-06-17 15:56:23 +0000 UTC"}`,
			false,
			fastest.Fail,
		},
		{
			1529250983,
			`{"unixtime":"1529250983xz`,
			true,
			fastest.Fail,
		},
	}

	for i, tc := range testCases {
		ft.Seq(fmt.Sprintf("#%v", i), func(ft fastest.T) {
			rsrc := io.NewBookmark("X")
			var r io.Reader
			if tc.readOK {
				r = testutils.Reader{*rsrc: []byte(tc.data)}
			} else {
				r = testutils.Reader{}
			}
			bookmark, err := ReadBookmark("X", r)
			ft.Implies(err != nil, tc.err == fastest.Fail, err)
			ft.Implies(err == nil, tc.err == fastest.OK)
			ft.Only(err == nil)

			ft.Equals(bookmark, tc.timestamp)
		})
	}
}

func TestWriteBookmark(t *testing.T) {
	ft := fastest.T{T: t}

	testCases := []struct {
		timestamp int64
		data      string
		err       fastest.Code
	}{
		{1529246468, `{"unixtime":"1529246468","time":"2018-06-17 14:41:08 +0000 UTC"}`, fastest.OK},
		{1529250983, `{"unixtime":"1529250983","time":"2018-06-17 15:56:23 +0000 UTC"}`, fastest.Fail},
	}

	for i, tc := range testCases {
		ft.Seq(fmt.Sprintf("#%v", i), func(ft fastest.T) {
			rsrc := io.NewBookmark("X")
			w := testutils.NewWriter(map[io.Resource][]bool{
				*rsrc: []bool{tc.err == fastest.OK}})
			err := WriteBookmark(tc.timestamp, "X", w)
			ft.Implies(err != nil, tc.err == fastest.Fail)
			ft.Implies(err == nil, tc.err == fastest.OK, err)
			ft.Only(err == nil)

			written, ok := w.Data[*rsrc]
			ft.True(ok)
			ft.Equals(len(written), 1)
			ft.Equals(string(written[0]), string(tc.data))
		})
	}
}
