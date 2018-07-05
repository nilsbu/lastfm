package organize

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/nilsbu/fastest"
	"github.com/nilsbu/lastfm/pkg/rsrc"
	"github.com/nilsbu/lastfm/pkg/store"
	"github.com/nilsbu/lastfm/pkg/unpack"
	"github.com/nilsbu/lastfm/test/mock"
)

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
			var io rsrc.IO
			if tc.json == "" {
				io, _ = mock.IO(map[rsrc.Locator][]byte{}, mock.Path)
			} else {
				io, _ = mock.IO(
					map[rsrc.Locator][]byte{rsrc.SessionID(): []byte(tc.json)},
					mock.Path)
			}
			sid, err := LoadSessionID(io)
			ft.Equals(err != nil, tc.err == fastest.Fail)
			ft.Only(tc.err == fastest.OK)
			ft.DeepEquals(sid, tc.sid)
		})
	}
}

func TestAllDayPlays(t *testing.T) {
	// also see TestAllDayPlaysFalseName below

	testCases := []struct {
		user    string
		plays   []HistoryDay
		writeOK bool
		readOK  bool
	}{
		{
			"XX",
			[]HistoryDay{HistoryDay{"BTS": 2, "XX": 1, "12": 1}},
			true, true,
		},
		{
			"XX",
			[]HistoryDay{
				HistoryDay{"as": 42, "": 1, "12": 100},
				HistoryDay{"ギルガメッシュ": 1000},
			},
			false, false,
		},
	}

	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			var files map[rsrc.Locator][]byte
			if tc.writeOK {
				files = map[rsrc.Locator][]byte{rsrc.AllDayPlays(tc.user): nil}
			} else {
				files = map[rsrc.Locator][]byte{}
			}

			io, _ := mock.IO(files, mock.Path)
			err := WriteAllDayPlays(tc.plays, tc.user, io)
			if err != nil && tc.writeOK {
				t.Error("unexpected error during write:", err)
			} else if err == nil && !tc.writeOK {
				t.Error("expected error during write but none occured")
			}

			plays, err := ReadAllDayPlays(tc.user, io)
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
	io, _ := mock.IO(map[rsrc.Locator][]byte{}, mock.Path)

	if err := WriteAllDayPlays([]HistoryDay{}, "I", io); err == nil {
		t.Error("expected error during write but non occurred")
	}

	if _, err := ReadAllDayPlays("I", io); err == nil {
		t.Error("expected error during read but non occurred")
	}
}

func TestUpdateAllDayPlays(t *testing.T) {
	h0 := rsrc.History("AA", 1, rsrc.ToDay(0*86400))
	h1 := rsrc.History("AA", 1, rsrc.ToDay(1*86400))
	h2 := rsrc.History("AA", 1, rsrc.ToDay(2*86400))
	h3 := rsrc.History("AA", 1, rsrc.ToDay(3*86400))

	testCases := []struct {
		user           unpack.User
		until          rsrc.Day
		saved          []HistoryDay
		tracksFile     map[rsrc.Locator][]byte
		tracksDownload map[rsrc.Locator][]byte
		plays          []HistoryDay
		ok             bool
	}{
		{ // No data
			unpack.User{Name: "AA", Registered: rsrc.ToDay(0)},
			rsrc.ToDay(0),
			nil,
			map[rsrc.Locator][]byte{},
			map[rsrc.Locator][]byte{},
			[]HistoryDay{},
			false,
		},
		{ // Registration day invalid
			unpack.User{Name: "AA", Registered: rsrc.NoDay()},
			rsrc.ToDay(0),
			nil,
			map[rsrc.Locator][]byte{},
			map[rsrc.Locator][]byte{},
			[]HistoryDay{},
			false,
		},
		{ // Begin no valid day
			unpack.User{Name: "AA", Registered: rsrc.ToDay(0)},
			rsrc.NoDay(),
			nil,
			map[rsrc.Locator][]byte{},
			map[rsrc.Locator][]byte{},
			[]HistoryDay{},
			false,
		},
		{ // download one day
			unpack.User{Name: "AA", Registered: rsrc.ToDay(300)}, // registered at 0:05
			rsrc.ToDay(0),
			[]HistoryDay{},
			map[rsrc.Locator][]byte{h0: nil},
			map[rsrc.Locator][]byte{
				h0: []byte(`{"recenttracks":{"track":[{"artist":{"#text":"ASDF"}}], "@attr":{"totalPages":"1"}}}`),
			},
			[]HistoryDay{
				HistoryDay{"ASDF": 1},
			},
			true,
		},
		{ // download some, have some
			unpack.User{Name: "AA", Registered: rsrc.ToDay(86400)},
			rsrc.ToDay(3 * 86400),
			[]HistoryDay{
				HistoryDay{"XX": 4},
				HistoryDay{}, // will be overwritten
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
			[]HistoryDay{
				HistoryDay{"XX": 4},
				HistoryDay{"ASDF": 1},
				HistoryDay{"B": 1},
			},
			true,
		},
		{ // have more than want
			unpack.User{Name: "AA", Registered: rsrc.ToDay(0)},
			rsrc.ToDay(86400),
			[]HistoryDay{
				HistoryDay{"XX": 2},
				HistoryDay{"A": 1},
				HistoryDay{"DropMe": 1},
				HistoryDay{"DropMeToo": 100},
			},
			map[rsrc.Locator][]byte{
				h0: []byte(`{"recenttracks":{"track":[{"artist":{"#text":"XX"}},{"artist":{"#text":"XX"}}], "@attr":{"totalPages":"1"}}}`),
				h1: []byte(`{"recenttracks":{"track":[{"artist":{"#text":"A"}}], "@attr":{"totalPages":"1"}}}`),
			},
			map[rsrc.Locator][]byte{},
			[]HistoryDay{
				HistoryDay{"XX": 2},
				HistoryDay{"A": 1},
			},
			true,
		},
		{ // download error
			unpack.User{Name: "AA", Registered: rsrc.ToDay(0)},
			rsrc.ToDay(0),
			[]HistoryDay{},
			map[rsrc.Locator][]byte{},
			map[rsrc.Locator][]byte{},
			[]HistoryDay{},
			false,
		},
	}

	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			tc.tracksFile[rsrc.AllDayPlays(tc.user.Name)] = nil
			io1, _ := mock.IO(tc.tracksFile, mock.Path)
			if tc.saved != nil {
				if err := WriteAllDayPlays(tc.saved, tc.user.Name, io1); err != nil {
					t.Error("unexpected error during write of all day plays:", err)
				}

			}

			io0, _ := mock.IO(tc.tracksDownload, mock.URL)

			pool, _ := store.NewCache([][]rsrc.IO{[]rsrc.IO{io0}, []rsrc.IO{io1}})

			plays, err := UpdateAllDayPlays(&tc.user, tc.until, pool)
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
		})
	}
}

func TestReadARtistTags(t *testing.T) {
	artistTags := rsrc.ArtistTags("xy")
	cases := []struct {
		files  map[rsrc.Locator][]byte
		artist string
		tags   []TagCount
		ok     bool
	}{
		{ // no file
			map[rsrc.Locator][]byte{artistTags: nil},
			"xy",
			nil,
			false,
		},
		{ // invalid user
			map[rsrc.Locator][]byte{artistTags: nil},
			"",
			nil,
			false,
		},
		{ // broken file
			map[rsrc.Locator][]byte{artistTags: []byte(`{"user":{"name":"x`)},
			"xy",
			nil,
			false,
		},
		{ // wrong content
			map[rsrc.Locator][]byte{artistTags: []byte(`{"user":{"name":"xy","registered":{"unixtime":86400}}}`)},
			"xy",
			nil,
			false,
		},
		{ // ok
			map[rsrc.Locator][]byte{artistTags: []byte(`{"toptags":{"tag":[{"name":"bui", "count":100},{"count":12,"name":"asdf"}]}}`)},
			"xy",
			[]TagCount{TagCount{"bui", 100}, TagCount{"asdf", 12}},
			true,
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			io, err := mock.IO(c.files, mock.Path)
			if err != nil {
				t.Fatal("setup error:", err)
			}

			tags, err := ReadArtistTags(c.artist, io)
			if err != nil && c.ok {
				t.Error("unexpected error:", err)
			} else if err == nil && !c.ok {
				t.Error("expected error but none occurred")
			}
			if err == nil {
				if !reflect.DeepEqual(tags, c.tags) {
					t.Errorf("read user faulty:\nhas:      %v\nexpected: %v",
						tags, c.tags)
				}
			}
		})
	}
}
