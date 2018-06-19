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
			w := testutils.NewWriter(map[io.Resource]bool{})
			err := WriteAllDayPlays(tc.plays, tc.name, w)
			ft.Nil(err, err)

			rsrc := io.NewAllDayPlays(tc.name)
			written, ok := w.Data[*rsrc]
			ft.True(ok)

			var r io.Reader
			if tc.failRead {
				r = testutils.Reader{}
			} else {
				r = testutils.Reader{*rsrc: written}
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
			// TODO check tat implies of this kind are handled propperly everywhere
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
			w := testutils.NewWriter(map[io.Resource]bool{
				*rsrc: tc.err == fastest.OK})
			err := WriteBookmark(tc.timestamp, "X", w)
			ft.Implies(err != nil, tc.err == fastest.Fail, err)
			ft.Implies(err == nil, tc.err == fastest.OK)
			ft.Only(err == nil)

			written, ok := w.Data[*rsrc]
			ft.True(ok)
			ft.Equals(string(written), string(tc.data))
		})
	}
}

func TestUpdateAllDayPlays(t *testing.T) {
	ft := fastest.T{T: t}

	testCases := []struct {
		user           unpack.User
		until          io.Midnight
		saved          []unpack.DayPlays
		tracksFile     map[io.Resource][]byte
		tracksDownload map[io.Resource][]byte
		plays          []unpack.DayPlays
		err            fastest.Code
	}{
		{ // No data
			unpack.User{Name: "A", Registered: 0},
			0,
			nil,
			map[io.Resource][]byte{},
			map[io.Resource][]byte{},
			[]unpack.DayPlays{},
			fastest.Fail,
		},
		{ // download one day
			unpack.User{Name: "A", Registered: 300}, // registered at 0:05
			0,
			[]unpack.DayPlays{},
			map[io.Resource][]byte{},
			map[io.Resource][]byte{
				*io.NewUserRecentTracks("A", 1, 0): []byte(`{"recenttracks":{"track":[{"artist":{"#text":"ASDF"}}], "@attr":{"totalPages":"1"}}}`),
			},
			[]unpack.DayPlays{
				unpack.DayPlays{"ASDF": 1},
			},
			fastest.OK,
		},
		{ // download some, have some
			unpack.User{Name: "A", Registered: 86400},
			3 * 86400,
			[]unpack.DayPlays{
				unpack.DayPlays{"XX": 4},
				unpack.DayPlays{}, // will be overwritten
			},
			map[io.Resource][]byte{
				*io.NewUserRecentTracks("A", 1, 1*86400): []byte(`{"recenttracks":{"track":[{"artist":{"#text":"XX"}},{"artist":{"#text":"XX"}},{"artist":{"#text":"XX"}},{"artist":{"#text":"XX"}}], "@attr":{"totalPages":"1"}}}`),
				*io.NewUserRecentTracks("A", 1, 2*86400): []byte(`{"recenttracks":{"track":[], "@attr":{"totalPages":"1"}}}`),
			},
			map[io.Resource][]byte{
				*io.NewUserRecentTracks("A", 1, 2*86400): []byte(`{"recenttracks":{"track":[{"artist":{"#text":"ASDF"}}], "@attr":{"totalPages":"1"}}}`),
				*io.NewUserRecentTracks("A", 1, 3*86400): []byte(`{"recenttracks":{"track":[{"artist":{"#text":"B"}}], "@attr":{"totalPages":"1"}}}`),
			},
			[]unpack.DayPlays{
				unpack.DayPlays{"XX": 4},
				unpack.DayPlays{"ASDF": 1},
				unpack.DayPlays{"B": 1},
			},
			fastest.OK,
		},
		{ // have more than want
			unpack.User{Name: "A", Registered: 0},
			86400,
			[]unpack.DayPlays{
				unpack.DayPlays{"XX": 2},
				unpack.DayPlays{"A": 1},
				unpack.DayPlays{"DropMe": 1},
				unpack.DayPlays{"DropMeToo": 100},
			},
			map[io.Resource][]byte{
				*io.NewUserRecentTracks("A", 1, 0):       []byte(`{"recenttracks":{"track":[{"artist":{"#text":"XX"}},{"artist":{"#text":"XX"}}], "@attr":{"totalPages":"1"}}}`),
				*io.NewUserRecentTracks("A", 1, 1*86400): []byte(`{"recenttracks":{"track":[{"artist":{"#text":"A"}}], "@attr":{"totalPages":"1"}}}`),
			},
			map[io.Resource][]byte{},
			[]unpack.DayPlays{
				unpack.DayPlays{"XX": 2},
				unpack.DayPlays{"A": 1},
			},
			fastest.OK,
		},
		{ // download error
			unpack.User{Name: "A", Registered: 0},
			0,
			[]unpack.DayPlays{},
			map[io.Resource][]byte{},
			map[io.Resource][]byte{},
			[]unpack.DayPlays{},
			fastest.Fail,
		},
	}

	for i, tc := range testCases {
		ft.Seq(fmt.Sprintf("#%v", i), func(ft fastest.T) {
			w := testutils.NewWriter(map[io.Resource]bool{})
			w.Data = tc.tracksFile
			if tc.saved != nil {
				err := WriteAllDayPlays(tc.saved, tc.user.Name, w)
				ft.Nil(err)
			}

			pool := io.NewPool(
				[]io.Reader{testutils.Reader(tc.tracksDownload)},
				[]io.Reader{testutils.Reader(w.Data)},
				[]io.Writer{w})

			plays, err := UpdateAllDayPlays(tc.user, tc.until, pool)
			ft.Implies(err != nil, tc.err == fastest.Fail, err)
			ft.Implies(err == nil, tc.err == fastest.OK)
			ft.Only(err == nil)
			ft.DeepEquals(plays, tc.plays)
		})
	}
}
