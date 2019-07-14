package organize

import (
	"reflect"
	"testing"

	"github.com/nilsbu/lastfm/pkg/charts"
	"github.com/nilsbu/lastfm/pkg/rsrc"
	"github.com/nilsbu/lastfm/pkg/store"
	"github.com/nilsbu/lastfm/pkg/unpack"
	"github.com/nilsbu/lastfm/test/mock"
)

func TestLoadHistory(t *testing.T) {
	testCases := []struct {
		user  unpack.User
		until rsrc.Day
		data  [][]string
		dps   [][]charts.Song
		ok    bool
	}{
		{
			unpack.User{Name: "", Registered: rsrc.ParseDay("2018-01-10")},
			rsrc.ParseDay("2018-01-11"),
			[][]string{[]string{}, []string{}},
			nil,
			false,
		},
		{
			unpack.User{Name: "", Registered: rsrc.ParseDay("2018-01-10")},
			nil,
			[][]string{[]string{}, []string{}},
			nil,
			false,
		},
		{
			unpack.User{Name: "", Registered: nil},
			rsrc.ParseDay("2018-01-11"),
			[][]string{[]string{}, []string{}},
			nil,
			false,
		},
		{
			unpack.User{Name: "ASDF", Registered: rsrc.ParseDay("2018-01-11")},
			rsrc.ParseDay("2018-01-12"),
			[][]string{
				[]string{`{"recenttracks":{"track":[{"artist":{"#text":"ASDF"}}], "@attr":{"totalPages":"1"}}}`},
				[]string{`{"recenttracks":{"track":[{"artist":{"#text":"XXX"}}], "@attr":{"totalPages":"1"}}}`},
			},
			[][]charts.Song{
				{{Artist: "ASDF"}},
				{{Artist: "XXX"}},
			},
			true,
		},
		{
			unpack.User{Name: "ASDF", Registered: rsrc.ParseDay("2018-01-10")},
			rsrc.ParseDay("2018-01-10"),
			[][]string{
				[]string{
					`{"recenttracks":{"track":[{"artist":{"#text":"X"}}], "@attr":{"page":"1","totalPages":"3"}}}`,
					`{"recenttracks":{"track":[{"artist":{"#text":"Y"}}], "@attr":{"page":"2","totalPages":"3"}}}`,
					`{"recenttracks":{"track":[{"artist":{"#text":"Z"}}, {"artist":{"#text":"X"}}], "@attr":{"page":"3","totalPages":"3"}}}`,
				},
			},
			[][]charts.Song{
				{{Artist: "X"}, {Artist: "Y"}, {Artist: "Z"}, {Artist: "X"}},
			},
			true,
		},
		{
			unpack.User{Name: "ASDF", Registered: rsrc.ParseDay("2018-01-10")},
			rsrc.ParseDay("2018-01-10"),
			[][]string{
				[]string{
					`{"recenttracks":{"track":[{"artist":{"#text":"X"}}], "@attr":{"page":"1","totalPages":"3"}}}`,
					"", "",
				},
			},
			nil,
			false,
		},
	}

	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			files := make(map[rsrc.Locator][]byte)
			for j, day := range tc.data {
				for k, d := range day {
					time := tc.user.Registered.AddDate(0, 0, j)
					files[rsrc.History(tc.user.Name, k+1, time)] = []byte(d)
				}
			}
			io, _ := mock.IO(files, mock.Path)

			dps, err := LoadHistory(tc.user, tc.until, io)
			if err != nil && tc.ok {
				t.Error("unexpected error:", err)
			} else if err == nil && !tc.ok {
				t.Error("expected error but none occurred")
			}
			if err == nil {
				if !reflect.DeepEqual(dps, tc.dps) {
					t.Errorf("wrong data:\nhas:      %v\nexpected: %v",
						dps, tc.dps)
				}
			}
		})
	}
}

func TestUpdateHistory(t *testing.T) {
	h0 := rsrc.History("AA", 1, rsrc.ParseDay("2018-01-10"))
	h1 := rsrc.History("AA", 1, rsrc.ParseDay("2018-01-11"))
	h2 := rsrc.History("AA", 1, rsrc.ParseDay("2018-01-12"))
	h3 := rsrc.History("AA", 1, rsrc.ParseDay("2018-01-13"))
	bm := rsrc.Bookmark("AA")

	testCases := []struct {
		user           unpack.User
		until          rsrc.Day
		bookmark       rsrc.Day
		saved          [][]charts.Song
		tracksFile     map[rsrc.Locator][]byte
		tracksDownload map[rsrc.Locator][]byte
		plays          [][]charts.Song
		ok             bool
	}{
		{ // No data
			unpack.User{Name: "AA", Registered: rsrc.ParseDay("2018-01-10")},
			rsrc.ParseDay("2018-01-10"),
			nil,
			nil,
			map[rsrc.Locator][]byte{},
			map[rsrc.Locator][]byte{},
			[][]charts.Song{},
			false,
		},
		{ // Registration day invalid
			unpack.User{Name: "AA", Registered: nil},
			rsrc.ParseDay("2018-01-10"),
			nil,
			nil,
			map[rsrc.Locator][]byte{},
			map[rsrc.Locator][]byte{},
			[][]charts.Song{},
			false,
		},
		{ // Begin no valid day
			unpack.User{Name: "AA", Registered: rsrc.ParseDay("2018-01-10")},
			nil,
			nil,
			nil,
			map[rsrc.Locator][]byte{},
			map[rsrc.Locator][]byte{},
			[][]charts.Song{},
			false,
		},
		{ // download one day
			unpack.User{Name: "AA", Registered: rsrc.ParseDay("2018-01-10")},
			rsrc.ParseDay("2018-01-10"),
			rsrc.ParseDay("2018-01-10"),
			[][]charts.Song{},
			map[rsrc.Locator][]byte{h0: nil, bm: nil},
			map[rsrc.Locator][]byte{
				h0: []byte(`{"recenttracks":{"track":[{"artist":{"#text":"ASDF"}}], "@attr":{"totalPages":"1"}}}`),
			},
			[][]charts.Song{
				{{Artist: "ASDF"}},
			},
			true,
		},
		{ // download some, have some
			unpack.User{Name: "AA", Registered: rsrc.ParseDay("2018-01-11")},
			rsrc.ParseDay("2018-01-13"),
			rsrc.ParseDay("2018-01-12"),
			[][]charts.Song{
				{{Artist: "XX"}, {Artist: "XX"}, {Artist: "XX"}, {Artist: "XX"}},
				{}, // will be overwritten
			},
			map[rsrc.Locator][]byte{
				h1: []byte(`{"recenttracks":{"track":[{"artist":{"#text":"XX"}},{"artist":{"#text":"XX"}},{"artist":{"#text":"XX"}},{"artist":{"#text":"XX"}}], "@attr":{"totalPages":"1"}}}`),
				h2: []byte(`{"recenttracks":{"track":[], "@attr":{"totalPages":"1"}}}`),
				h3: nil,
				bm: nil,
			},
			map[rsrc.Locator][]byte{
				h1: nil,
				h2: []byte(`{"recenttracks":{"track":[{"artist":{"#text":"ASDF"}}], "@attr":{"totalPages":"1"}}}`),
				h3: []byte(`{"recenttracks":{"track":[{"artist":{"#text":"B"}}], "@attr":{"totalPages":"1"}}}`),
			},
			[][]charts.Song{
				{{Artist: "XX"}, {Artist: "XX"}, {Artist: "XX"}, {Artist: "XX"}},
				{{Artist: "ASDF"}},
				{{Artist: "B"}},
			},
			true,
		},
		{ // have more than want
			unpack.User{Name: "AA", Registered: rsrc.ParseDay("2018-01-10")},
			rsrc.ParseDay("2018-01-11"),
			rsrc.ParseDay("2018-01-13"),
			[][]charts.Song{
				{{Artist: "XX"}, {Artist: "XX"}},
				{{Artist: "A"}},
				{{Artist: "DropMe"}},
				{{Artist: "DropMeToo"}, {Artist: "DropMeToo"}, {Artist: "DropMeToo"}},
			},
			map[rsrc.Locator][]byte{
				h0: []byte(`{"recenttracks":{"track":[{"artist":{"#text":"XX"}},{"artist":{"#text":"XX"}}], "@attr":{"totalPages":"1"}}}`),
				h1: []byte(`{"recenttracks":{"track":[{"artist":{"#text":"A"}}], "@attr":{"totalPages":"1"}}}`),
				bm: nil,
			},
			map[rsrc.Locator][]byte{},
			[][]charts.Song{
				{{Artist: "XX"}, {Artist: "XX"}},
				{{Artist: "A"}},
			},
			true,
		},
		{ // saved days ahead of bookmark
			unpack.User{Name: "AA", Registered: rsrc.ParseDay("2018-01-10")},
			rsrc.ParseDay("2018-01-13"),
			rsrc.ParseDay("2018-01-12"),
			[][]charts.Song{
				{{Artist: "XX"}, {Artist: "XX"}},
				{{Artist: "A"}},
				{{Artist: "hui"}},
				{{Artist: "hui"}},
			},
			map[rsrc.Locator][]byte{
				h0: []byte(`{"recenttracks":{"track":[{"artist":{"#text":"XX"}},{"artist":{"#text":"XX"}}], "@attr":{"totalPages":"1"}}}`),
				h1: []byte(`{"recenttracks":{"track":[{"artist":{"#text":"A"}}], "@attr":{"totalPages":"1"}}}`),
				h2: []byte(`{"recenttracks":{"track":[{"artist":{"#text":"hui"}},{"artist":{"#text":"hui"}}], "@attr":{"totalPages":"1"}}}`),
				h3: []byte(`{"recenttracks":{"track":[{"artist":{"#text":"hui"}},{"artist":{"#text":"hui"}}], "@attr":{"totalPages":"1"}}}`),
				bm: nil,
			},
			map[rsrc.Locator][]byte{},
			[][]charts.Song{
				{{Artist: "XX"}, {Artist: "XX"}},
				{{Artist: "A"}},
				{{Artist: "hui"}, {Artist: "hui"}},
				{{Artist: "hui"}, {Artist: "hui"}},
			},
			true,
		},
		{ // saved days but no bookmark
			unpack.User{Name: "AA", Registered: rsrc.ParseDay("2018-01-10")},
			rsrc.ParseDay("2018-01-13"),
			nil,
			[][]charts.Song{
				{{Artist: "XX"}, {Artist: "XX"}},
				{{Artist: "A"}},
				{{Artist: "hui"}, {Artist: "hui"}},
				{{Artist: "hui"}},
			},
			map[rsrc.Locator][]byte{
				h0: nil,
				h1: nil,
				h2: nil, // ensure that these days aren't read
				h3: []byte(`{"recenttracks":{"track":[{"artist":{"#text":"hui"}},{"artist":{"#text":"hui"}}], "@attr":{"totalPages":"1"}}}`),
				bm: nil,
			},
			map[rsrc.Locator][]byte{},
			[][]charts.Song{
				{{Artist: "XX"}, {Artist: "XX"}},
				{{Artist: "A"}},
				{{Artist: "hui"}, {Artist: "hui"}},
				{{Artist: "hui"}, {Artist: "hui"}},
			},
			true,
		},
		{ // download error
			unpack.User{Name: "AA", Registered: rsrc.ParseDay("2018-01-10")},
			rsrc.ParseDay("2018-01-10"),
			rsrc.ParseDay("2018-01-10"),
			[][]charts.Song{},
			map[rsrc.Locator][]byte{bm: nil},
			map[rsrc.Locator][]byte{},
			[][]charts.Song{},
			false,
		},
	}

	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			tc.tracksFile[rsrc.SongHistory(tc.user.Name)] = nil
			io1, _ := mock.IO(tc.tracksFile, mock.Path)
			if tc.saved != nil {
				if err := unpack.WriteSongHistory(tc.saved, tc.user.Name, io1); err != nil {
					t.Fatal("unexpected error during write of all day plays:", err)
				}
			}

			if tc.bookmark != nil {
				dt := int(tc.bookmark.Midnight()-tc.user.Registered.Midnight()) / 86400
				sd := len(tc.saved)
				if dt > sd {
					t.Fatalf("bookmark is %vd after registered but must not be more "+
						"than number of days saved (%v)",
						dt, sd)
				}

				if err := unpack.WriteBookmark(tc.bookmark, tc.user.Name, io1); err != nil {
					t.Fatal("unexpected error during write of bookmark:", err)
				}
			}

			io0, _ := mock.IO(tc.tracksDownload, mock.URL)

			store, _ := store.New([][]rsrc.IO{[]rsrc.IO{io0}, []rsrc.IO{io1}})

			plays, err := UpdateHistory(&tc.user, tc.until, store)
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
