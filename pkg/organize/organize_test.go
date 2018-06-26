package organize

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/nilsbu/fastest"
	"github.com/nilsbu/lastfm/pkg/io"
	"github.com/nilsbu/lastfm/pkg/rsrc"
	"github.com/nilsbu/lastfm/pkg/store"
	"github.com/nilsbu/lastfm/pkg/unpack"
	"github.com/nilsbu/lastfm/test/mock"
)

func TestLoadAPIKey(t *testing.T) {
	ft := fastest.T{T: t}

	testCases := []struct {
		json   string
		apiKey rsrc.Key
		err    fastest.Code
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
				r, _, _ = mock.IO(map[rsrc.Locator][]byte{}, mock.Path)
			} else {
				r, _, _ = mock.IO(
					map[rsrc.Locator][]byte{rsrc.APIKey(): []byte(tc.json)},
					mock.Path)
			}
			apiKey, err := LoadAPIKey(r)
			ft.Equals(err != nil, tc.err == fastest.Fail)
			ft.Only(tc.err == fastest.OK)
			ft.Equals(apiKey, tc.apiKey)
		})
	}
}

func TestLoadSessionID(t *testing.T) {
	ft := fastest.T{T: t}

	testCases := []struct {
		json string
		sid  SessionID
		err  fastest.Code
	}{
		{"", "", fastest.Fail},
		{`{`, "", fastest.Fail},
		{`{}`, "", fastest.Fail},
		{`{"user":"asdf"}`, "asdf", fastest.OK},
	}

	for i, tc := range testCases {
		ft.Seq(fmt.Sprintf("#%v", i), func(ft fastest.T) {
			var r io.Reader
			if tc.json == "" {
				r, _, _ = mock.IO(map[rsrc.Locator][]byte{}, mock.Path)
			} else {
				r, _, _ = mock.IO(
					map[rsrc.Locator][]byte{rsrc.SessionID(): []byte(tc.json)},
					mock.Path)
			}
			sid, err := LoadSessionID(r)
			ft.Equals(err != nil, tc.err == fastest.Fail)
			ft.Only(tc.err == fastest.OK)
			ft.DeepEquals(sid, tc.sid)
		})
	}
}

func TestAllDayPlays(t *testing.T) {
	// also see TestAllDayPlaysFalseName below

	testCases := []struct {
		name    rsrc.Name
		plays   []unpack.DayPlays
		writeOK bool
		readOK  bool
	}{
		{
			"XX",
			[]unpack.DayPlays{unpack.DayPlays{"BTS": 2, "XX": 1, "12": 1}},
			true, true,
		},
		{
			"XX",
			[]unpack.DayPlays{
				unpack.DayPlays{"as": 42, "": 1, "12": 100},
				unpack.DayPlays{"ギルガメッシュ": 1000},
			},
			false, false,
		},
	}

	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			loc, _ := rsrc.AllDayPlays(tc.name)
			var files map[rsrc.Locator][]byte

			if tc.writeOK {
				files = map[rsrc.Locator][]byte{loc: nil}
			} else {
				files = map[rsrc.Locator][]byte{}
			}

			r, w, _ := mock.IO(files, mock.Path)
			err := WriteAllDayPlays(tc.plays, tc.name, w)
			if err != nil && tc.writeOK {
				t.Error("unexpected error during write:", err)
			} else if err == nil && !tc.writeOK {
				t.Error("expected error during write but none occured")
			}

			plays, err := ReadAllDayPlays(tc.name, r)
			if err != nil && tc.readOK {
				t.Error("unexpected error during read:", err)
			} else if err == nil && !tc.readOK {
				t.Error("expected error during read but none occurred")
			}
			if err == nil {
				if !reflect.DeepEqual(plays, tc.plays) {
					t.Errorf("read plays differ from written:\nread:    %v\nwritten: %v",
						plays, tc.plays)
				}
			}
		})
	}
}

func TestAllDayPlaysFalseName(t *testing.T) {
	r, w, _ := mock.IO(map[rsrc.Locator][]byte{}, mock.Path)

	if err := WriteAllDayPlays([]unpack.DayPlays{}, "I", w); err == nil {
		t.Error("expected error during write but non occurred")
	}

	if _, err := ReadAllDayPlays("I", r); err == nil {
		t.Error("expected error during read but non occurred")
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
			`{"unixtime":1529246468,"time":"2018-06-17 14:41:08 +0000 UTC"}`,
			true,
			fastest.OK,
		},
		{
			1529250983,
			`{"unixtime":1529250983,"time":"2018-06-17 15:56:23 +0000 UTC"}`,
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
			loc, _ := rsrc.Bookmark("Xx")
			var r io.Reader
			if tc.readOK {
				r, _, _ = mock.IO(
					map[rsrc.Locator][]byte{loc: []byte(tc.data)},
					mock.Path)
			} else {
				r, _, _ = mock.IO(map[rsrc.Locator][]byte{}, mock.Path)
			}
			bookmark, err := ReadBookmark("Xx", r)
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
		{1529246468, `{"unixtime":1529246468,"time":"2018-06-17 14:41:08 +0000 UTC"}`, fastest.OK},
		{1529250983, `{"unixtime":1529250983,"time":"2018-06-17 15:56:23 +0000 UTC"}`, fastest.Fail},
	}

	for i, tc := range testCases {
		ft.Seq(fmt.Sprintf("#%v", i), func(ft fastest.T) {
			loc, _ := rsrc.Bookmark("XX")
			var files map[rsrc.Locator][]byte
			if tc.err == fastest.OK {
				files = map[rsrc.Locator][]byte{loc: nil}
			} else {
				files = map[rsrc.Locator][]byte{}
			}

			r, w, _ := mock.IO(files, mock.Path)
			err := WriteBookmark(tc.timestamp, "XX", w)
			ft.Implies(err != nil, tc.err == fastest.Fail, err)
			ft.Implies(err == nil, tc.err == fastest.OK)
			ft.Only(err == nil)

			written, err := r.Read(loc)
			ft.Nil(err)
			ft.Equals(string(written), string(tc.data))
		})
	}
}

func TestUpdateAllDayPlays(t *testing.T) {
	h0, _ := rsrc.History("AA", 1, rsrc.ToDay(0*86400))
	h1, _ := rsrc.History("AA", 1, rsrc.ToDay(1*86400))
	h2, _ := rsrc.History("AA", 1, rsrc.ToDay(2*86400))
	h3, _ := rsrc.History("AA", 1, rsrc.ToDay(3*86400))

	testCases := []struct {
		user           unpack.User
		until          rsrc.Day
		saved          []unpack.DayPlays
		tracksFile     map[rsrc.Locator][]byte
		tracksDownload map[rsrc.Locator][]byte
		plays          []unpack.DayPlays
		ok             bool
	}{
		{ // No data
			unpack.User{Name: "AA", Registered: rsrc.ToDay(0)},
			rsrc.ToDay(0),
			nil,
			map[rsrc.Locator][]byte{},
			map[rsrc.Locator][]byte{},
			[]unpack.DayPlays{},
			false,
		},
		{ // Registration day invalid
			unpack.User{Name: "AA", Registered: rsrc.NoDay()},
			rsrc.ToDay(0),
			nil,
			map[rsrc.Locator][]byte{},
			map[rsrc.Locator][]byte{},
			[]unpack.DayPlays{},
			false,
		},
		{ // Begin no valid day
			unpack.User{Name: "AA", Registered: rsrc.ToDay(0)},
			rsrc.NoDay(),
			nil,
			map[rsrc.Locator][]byte{},
			map[rsrc.Locator][]byte{},
			[]unpack.DayPlays{},
			false,
		},
		{ // download one day
			unpack.User{Name: "AA", Registered: rsrc.ToDay(300)}, // registered at 0:05
			rsrc.ToDay(0),
			[]unpack.DayPlays{},
			map[rsrc.Locator][]byte{h0: nil},
			map[rsrc.Locator][]byte{
				h0: []byte(`{"recenttracks":{"track":[{"artist":{"#text":"ASDF"}}], "@attr":{"totalPages":"1"}}}`),
			},
			[]unpack.DayPlays{
				unpack.DayPlays{"ASDF": 1},
			},
			true,
		},
		{ // download some, have some
			unpack.User{Name: "AA", Registered: rsrc.ToDay(86400)},
			rsrc.ToDay(3 * 86400),
			[]unpack.DayPlays{
				unpack.DayPlays{"XX": 4},
				unpack.DayPlays{}, // will be overwritten
			},
			map[rsrc.Locator][]byte{
				h1: []byte(`{"recenttracks":{"track":[{"artist":{"#text":"XX"}},{"artist":{"#text":"XX"}},{"artist":{"#text":"XX"}},{"artist":{"#text":"XX"}}], "@attr":{"totalPages":"1"}}}`),
				h2: []byte(`{"recenttracks":{"track":[], "@attr":{"totalPages":"1"}}}`),
				h3: nil,
			},
			map[rsrc.Locator][]byte{
				h1: nil,
				h2: []byte(`{"recenttracks":{"track":[{"artist":{"#text":"ASDF"}}], "@attr":{"totalPages":"1"}}}`),
				h3: []byte(`{"recenttracks":{"track":[{"artist":{"#text":"B"}}], "@attr":{"totalPages":"1"}}}`),
			},
			[]unpack.DayPlays{
				unpack.DayPlays{"XX": 4},
				unpack.DayPlays{"ASDF": 1},
				unpack.DayPlays{"B": 1},
			},
			true,
		},
		{ // have more than want
			unpack.User{Name: "AA", Registered: rsrc.ToDay(0)},
			rsrc.ToDay(86400),
			[]unpack.DayPlays{
				unpack.DayPlays{"XX": 2},
				unpack.DayPlays{"A": 1},
				unpack.DayPlays{"DropMe": 1},
				unpack.DayPlays{"DropMeToo": 100},
			},
			map[rsrc.Locator][]byte{
				h0: []byte(`{"recenttracks":{"track":[{"artist":{"#text":"XX"}},{"artist":{"#text":"XX"}}], "@attr":{"totalPages":"1"}}}`),
				h1: []byte(`{"recenttracks":{"track":[{"artist":{"#text":"A"}}], "@attr":{"totalPages":"1"}}}`),
			},
			map[rsrc.Locator][]byte{},
			[]unpack.DayPlays{
				unpack.DayPlays{"XX": 2},
				unpack.DayPlays{"A": 1},
			},
			true,
		},
		{ // download error
			unpack.User{Name: "AA", Registered: rsrc.ToDay(0)},
			rsrc.ToDay(0),
			[]unpack.DayPlays{},
			map[rsrc.Locator][]byte{},
			map[rsrc.Locator][]byte{},
			[]unpack.DayPlays{},
			false,
		},
	}

	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			loc, _ := rsrc.AllDayPlays(tc.user.Name)
			tc.tracksFile[loc] = nil
			r, w, _ := mock.IO(tc.tracksFile, mock.Path)
			if tc.saved != nil {
				if err := WriteAllDayPlays(tc.saved, tc.user.Name, w); err != nil {
					t.Error("unexpected error during write of all day plays:", err)
				}

			}

			d, _, _ := mock.IO(tc.tracksDownload, mock.URL)

			pool, _ := store.New(
				[][]io.Reader{[]io.Reader{d}, []io.Reader{r}},
				[][]io.Writer{[]io.Writer{io.FailIO{}}, []io.Writer{w}})

			plays, err := UpdateAllDayPlays(tc.user, tc.until, pool)
			if err != nil && tc.ok {
				t.Error("unexpected error:", err)
			} else if err == nil && !tc.ok {
				t.Error("expected error but none occurred")
			}
			if err == nil {
				if !reflect.DeepEqual(plays, tc.plays) {
					t.Errorf("updated plays faulty:\nhas:      %v\nexpected: %v",
						plays, tc.plays)
				}
			}
			// ft.DeepEquals(plays, tc.plays)
		})
	}
}
