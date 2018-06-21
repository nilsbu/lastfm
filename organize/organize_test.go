package organize

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/nilsbu/fastest"
	"github.com/nilsbu/lastfm/io"
	"github.com/nilsbu/lastfm/mock"
	"github.com/nilsbu/lastfm/rsrc"
	"github.com/nilsbu/lastfm/unpack"
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
				r, _ = mock.FileIO(map[string][]byte{})
			} else {
				path, _ := rsrc.APIKey().Path()
				r, _ = mock.FileIO(map[string][]byte{
					path: []byte(tc.json)})
			}
			apiKey, err := LoadAPIKey(r)
			ft.Equals(err != nil, tc.err == fastest.Fail)
			ft.Only(tc.err == fastest.OK)
			ft.Equals(apiKey, tc.apiKey)
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
		{
			"XX",
			[]unpack.DayPlays{unpack.DayPlays{"BTS": 2, "XX": 1, "12": 1}},
			true, false,
		},
	}

	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			rs, _ := rsrc.AllDayPlays(tc.name)
			path, _ := rs.Path()
			var files map[string][]byte

			if tc.writeOK {
				files = map[string][]byte{path: nil}
			} else {
				files = map[string][]byte{}
			}

			r, w := mock.FileIO(files)
			err := WriteAllDayPlays(tc.plays, tc.name, w)
			if err != nil && tc.writeOK {
				t.Error("unexpected error during write:", err)
			} else if err == nil && !tc.writeOK {
				t.Error("expected error during write but none occured")
			}

			if !tc.readOK {
				delete(files, path)
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
	r, w := mock.FileIO(map[string][]byte{})

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
			rs, _ := rsrc.Bookmark("Xx")
			var r io.Reader
			if tc.readOK {
				path, _ := rs.Path()
				r, _ = mock.FileIO(map[string][]byte{path: []byte(tc.data)})
			} else {
				r, _ = mock.FileIO(map[string][]byte{})
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
			rs, _ := rsrc.Bookmark("XX")
			path, _ := rs.Path()
			var files map[string][]byte
			if tc.err == fastest.OK {
				files = map[string][]byte{path: nil}
			} else {
				files = map[string][]byte{}
			}

			_, w := mock.FileIO(files)
			err := WriteBookmark(tc.timestamp, "XX", w)
			ft.Implies(err != nil, tc.err == fastest.Fail)
			ft.Implies(err == nil, tc.err == fastest.OK, err)
			ft.Only(err == nil)

			written, ok := files[path]
			ft.True(ok)
			ft.Equals(string(written), string(tc.data))
		})
	}
}
