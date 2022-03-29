package organize

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/nilsbu/lastfm/pkg/info"
	"github.com/nilsbu/lastfm/pkg/io"
	"github.com/nilsbu/lastfm/pkg/rsrc"
	"github.com/nilsbu/lastfm/pkg/unpack"
	"github.com/nilsbu/lastfm/test/mock"
)

func TestLoadHistory(t *testing.T) {
	p01 := rsrc.History("ASDF", 1, rsrc.ParseDay("2018-01-10"))
	p02 := rsrc.History("ASDF", 2, rsrc.ParseDay("2018-01-10"))
	p03 := rsrc.History("ASDF", 3, rsrc.ParseDay("2018-01-10"))
	p11 := rsrc.History("ASDF", 1, rsrc.ParseDay("2018-01-11"))

	tASDF := rsrc.TrackInfo("ASDF", "")
	tXXX := rsrc.TrackInfo("XXX", "")
	tX := rsrc.TrackInfo("X", "")
	tY := rsrc.TrackInfo("Y", "")
	tZ := rsrc.TrackInfo("Z", "")

	testCases := []struct {
		user  unpack.User
		until rsrc.Day
		files map[rsrc.Locator][]byte
		dps   [][]info.Song
		ok    bool
	}{
		{
			unpack.User{Name: "", Registered: rsrc.ParseDay("2018-01-10")},
			rsrc.ParseDay("2018-01-11"),
			map[rsrc.Locator][]byte{
				p01: nil,
				p11: nil,
			},
			nil,
			false,
		},
		{
			unpack.User{Name: "", Registered: rsrc.ParseDay("2018-01-10")},
			nil,
			map[rsrc.Locator][]byte{
				p01: nil,
			},
			nil,
			false,
		},
		{
			unpack.User{Name: "", Registered: nil},
			rsrc.ParseDay("2018-01-11"),
			map[rsrc.Locator][]byte{
				p01: nil,
			},
			nil,
			false,
		},
		{
			unpack.User{Name: "ASDF", Registered: rsrc.ParseDay("2018-01-10")},
			rsrc.ParseDay("2018-01-11"),
			map[rsrc.Locator][]byte{
				p01:   []byte(`{"recenttracks":{"track":[{"artist":{"#text":"ASDF"}}], "@attr":{"totalPages":"1"}}}`),
				p11:   []byte(`{"recenttracks":{"track":[{"artist":{"#text":"XXX"}}], "@attr":{"totalPages":"1"}}}`),
				tASDF: []byte(`{"track":{"duration":"120000"}}`),
				tXXX:  []byte(`{"track":{"duration":"60000"}}`),
			},
			[][]info.Song{
				{{Artist: "ASDF", Duration: 2}},
				{{Artist: "XXX", Duration: 1}},
			},
			true,
		},
		{
			unpack.User{Name: "ASDF", Registered: rsrc.ParseDay("2018-01-10")},
			rsrc.ParseDay("2018-01-10"),
			map[rsrc.Locator][]byte{
				p01: []byte(`{"recenttracks":{"track":[{"artist":{"#text":"X"}}], "@attr":{"page":"1","totalPages":"3"}}}`),
				p02: []byte(`{"recenttracks":{"track":[{"artist":{"#text":"Y"}}], "@attr":{"page":"2","totalPages":"3"}}}`),
				p03: []byte(`{"recenttracks":{"track":[{"artist":{"#text":"Z"}}, {"artist":{"#text":"X"}}], "@attr":{"page":"3","totalPages":"3"}}}`),
				tX:  []byte(`{"track":{"duration":"60000"}}`),
				tY:  []byte(`{"track":{"duration":"60000"}}`),
				tZ:  []byte(`{"track":{"duration":"60000"}}`),
			},
			[][]info.Song{
				{{Artist: "X", Duration: 1}, {Artist: "Y", Duration: 1}, {Artist: "Z", Duration: 1}, {Artist: "X", Duration: 1}},
			},
			true,
		},
		{
			unpack.User{Name: "ASDF", Registered: rsrc.ParseDay("2018-01-10")},
			rsrc.ParseDay("2018-01-10"),
			map[rsrc.Locator][]byte{
				p01: []byte(`{"recenttracks":{"track":[{"artist":{"#text":"X"}}], "@attr":{"page":"1","totalPages":"3"}}}`),
				p02: []byte(""),
				p03: []byte(""),
				tX:  []byte(`{"track":{"duration":"60000"}}`),
			},
			nil,
			false,
		},
	}

	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			io, _ := mock.IO(tc.files, mock.Path)

			dps, err := loadHistory(tc.user, tc.until, io, unpack.NewCached(io))
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

	b0 := rsrc.DayHistory("AA", rsrc.ParseDay("2018-01-10"))
	b1 := rsrc.DayHistory("AA", rsrc.ParseDay("2018-01-11"))
	b2 := rsrc.DayHistory("AA", rsrc.ParseDay("2018-01-12"))
	b3 := rsrc.DayHistory("AA", rsrc.ParseDay("2018-01-13"))

	bm := rsrc.Bookmark("AA")

	tASDF := rsrc.TrackInfo("ASDF", "")
	tXX := rsrc.TrackInfo("XX", "")
	thui := rsrc.TrackInfo("hui", "")
	tB := rsrc.TrackInfo("B", "")

	testCases := []struct {
		name           string
		user           unpack.User
		until          rsrc.Day
		bookmark       rsrc.Day
		saved          [][]info.Song
		tracksFile     map[rsrc.Locator][]byte
		tracksDownload map[rsrc.Locator][]byte
		plays          [][]info.Song
		ok             bool
	}{
		{
			"no data",
			unpack.User{Name: "AA", Registered: rsrc.ParseDay("2018-01-10")},
			rsrc.ParseDay("2018-01-10"),
			nil,
			nil,
			map[rsrc.Locator][]byte{},
			map[rsrc.Locator][]byte{},
			[][]info.Song{},
			false,
		},
		{
			"registration day invalid",
			unpack.User{Name: "AA", Registered: nil},
			rsrc.ParseDay("2018-01-10"),
			nil,
			nil,
			map[rsrc.Locator][]byte{},
			map[rsrc.Locator][]byte{},
			[][]info.Song{},
			false,
		},
		{
			"begin no valid day",
			unpack.User{Name: "AA", Registered: rsrc.ParseDay("2018-01-10")},
			nil,
			nil,
			nil,
			map[rsrc.Locator][]byte{},
			map[rsrc.Locator][]byte{},
			[][]info.Song{},
			false,
		},
		{
			"download one day",
			unpack.User{Name: "AA", Registered: rsrc.ParseDay("2018-01-10")},
			rsrc.ParseDay("2018-01-10"),
			rsrc.ParseDay("2018-01-10"),
			[][]info.Song{},
			map[rsrc.Locator][]byte{h0: nil, bm: nil},
			map[rsrc.Locator][]byte{
				h0:    []byte(`{"recenttracks":{"track":[{"artist":{"#text":"ASDF"}}], "@attr":{"totalPages":"1"}}}`),
				tASDF: []byte(`{"track":{"duration":"120000"}}`),
			},
			[][]info.Song{
				{{Artist: "ASDF", Duration: 2}},
			},
			true,
		},
		{
			"download some, have some",
			unpack.User{Name: "AA", Registered: rsrc.ParseDay("2018-01-11")},
			rsrc.ParseDay("2018-01-13"),
			rsrc.ParseDay("2018-01-12"),
			[][]info.Song{
				{{Artist: "XX", Duration: 1}, {Artist: "XX", Duration: 1}, {Artist: "XX", Duration: 1}, {Artist: "XX", Duration: 1}},
				{}, // will be overwritten
			},
			map[rsrc.Locator][]byte{
				h1:    []byte(`{"recenttracks":{"track":[{"artist":{"#text":"XX"}},{"artist":{"#text":"XX"}},{"artist":{"#text":"XX"}},{"artist":{"#text":"XX"}}], "@attr":{"totalPages":"1"}}}`),
				h2:    []byte(`{"recenttracks":{"track":[], "@attr":{"totalPages":"1"}}}`),
				h3:    nil,
				b1:    nil,
				b2:    nil,
				bm:    nil,
				tASDF: []byte(`{"track":{"duration":"120000"}}`),
				tXX:   []byte(`{"track":{"duration":"60000"}}`),
				tB:    []byte(`{"track":{"duration":"60000"}}`),
			},
			map[rsrc.Locator][]byte{
				h1: nil,
				h2: []byte(`{"recenttracks":{"track":[{"artist":{"#text":"ASDF"}}], "@attr":{"totalPages":"1"}}}`),
				h3: []byte(`{"recenttracks":{"track":[{"artist":{"#text":"B"}}], "@attr":{"totalPages":"1"}}}`),
			},
			[][]info.Song{
				{{Artist: "XX", Duration: 1}, {Artist: "XX", Duration: 1}, {Artist: "XX", Duration: 1}, {Artist: "XX", Duration: 1}},
				{{Artist: "ASDF", Duration: 2}},
				{{Artist: "B", Duration: 1}},
			},
			true,
		},
		{
			"don't reload what you don't need",
			unpack.User{Name: "AA", Registered: rsrc.ParseDay("2018-01-10")},
			rsrc.ParseDay("2018-01-12"),
			rsrc.ParseDay("2018-01-12"),
			[][]info.Song{
				{{Artist: "XX", Duration: 1}, {Artist: "XX", Duration: 1}, {Artist: "XX", Duration: 1}, {Artist: "XX", Duration: 1}},
				{{Artist: "XX", Duration: 1}, {Artist: "XX", Duration: 1}, {Artist: "XX", Duration: 1}, {Artist: "XX", Duration: 1}},
				{{Artist: "XX", Duration: 1}},
			},
			map[rsrc.Locator][]byte{
				// h1 is not read
				h2:    nil,
				b0:    nil,
				b1:    nil,
				b2:    nil,
				bm:    nil,
				tASDF: []byte(`{"track":{"duration":"120000"}}`),
				tXX:   []byte(`{"track":{"duration":"60000"}}`),
				tB:    []byte(`{"track":{"duration":"60000"}}`),
			},
			map[rsrc.Locator][]byte{
				h2: []byte(`{"recenttracks":{"track":[{"artist":{"#text":"XX"}},{"artist":{"#text":"XX"}},{"artist":{"#text":"XX"}},{"artist":{"#text":"XX"}}], "@attr":{"totalPages":"1"}}}`),
			},
			[][]info.Song{
				{{Artist: "XX", Duration: 1}, {Artist: "XX", Duration: 1}, {Artist: "XX", Duration: 1}, {Artist: "XX", Duration: 1}},
				{{Artist: "XX", Duration: 1}, {Artist: "XX", Duration: 1}, {Artist: "XX", Duration: 1}, {Artist: "XX", Duration: 1}},
				{{Artist: "XX", Duration: 1}, {Artist: "XX", Duration: 1}, {Artist: "XX", Duration: 1}, {Artist: "XX", Duration: 1}},
			},
			true,
		},
		{
			"have more than want",
			unpack.User{Name: "AA", Registered: rsrc.ParseDay("2018-01-10")},
			rsrc.ParseDay("2018-01-11"),
			rsrc.ParseDay("2018-01-13"),
			[][]info.Song{
				{{Artist: "XX", Duration: 1}, {Artist: "XX", Duration: 1}},
				{{Artist: "A", Duration: 4}},
				{{Artist: "DropMe"}},
				{{Artist: "DropMeToo"}, {Artist: "DropMeToo"}, {Artist: "DropMeToo"}},
			},
			map[rsrc.Locator][]byte{
				h0: []byte(`{"recenttracks":{"track":[{"artist":{"#text":"XX"}},{"artist":{"#text":"XX"}}], "@attr":{"totalPages":"1"}}}`),
				h1: []byte(`{"recenttracks":{"track":[{"artist":{"#text":"A"}}], "@attr":{"totalPages":"1"}}}`),
				bm: nil,
				b0: nil,
				b1: nil,
				b2: nil,
				b3: nil,
			},
			map[rsrc.Locator][]byte{},
			[][]info.Song{
				{{Artist: "XX", Duration: 1}, {Artist: "XX", Duration: 1}},
				{{Artist: "A", Duration: 4}},
			},
			true,
		},
		{
			"saved days ahead of bookmark",
			unpack.User{Name: "AA", Registered: rsrc.ParseDay("2018-01-10")},
			rsrc.ParseDay("2018-01-13"),
			rsrc.ParseDay("2018-01-12"),
			[][]info.Song{
				{{Artist: "XX", Duration: 1}, {Artist: "XX", Duration: 1}},
				{{Artist: "ASDF", Duration: 2}},
				{{Artist: "hui", Duration: 1}},
				{{Artist: "hui", Duration: 1}},
			},
			map[rsrc.Locator][]byte{
				h0:    []byte(`{"recenttracks":{"track":[{"artist":{"#text":"XX"}},{"artist":{"#text":"XX"}}], "@attr":{"totalPages":"1"}}}`),
				h1:    []byte(`{"recenttracks":{"track":[{"artist":{"#text":"ASDF"}}], "@attr":{"totalPages":"1"}}}`),
				h2:    []byte(`{"recenttracks":{"track":[{"artist":{"#text":"hui"}},{"artist":{"#text":"hui"}}], "@attr":{"totalPages":"1"}}}`),
				h3:    []byte(`{"recenttracks":{"track":[{"artist":{"#text":"hui"}},{"artist":{"#text":"hui"}}], "@attr":{"totalPages":"1"}}}`),
				bm:    nil,
				b0:    nil,
				b1:    nil,
				b2:    nil,
				b3:    nil,
				tXX:   []byte(`{"track":{"duration":"60000"}}`),
				thui:  []byte(`{"track":{"duration":"60000"}}`),
				tASDF: []byte(`{"track":{"duration":"120000"}}`),
			},
			map[rsrc.Locator][]byte{},
			[][]info.Song{
				{{Artist: "XX", Duration: 1}, {Artist: "XX", Duration: 1}},
				{{Artist: "ASDF", Duration: 2}},
				{{Artist: "hui", Duration: 1}, {Artist: "hui", Duration: 1}},
				{{Artist: "hui", Duration: 1}, {Artist: "hui", Duration: 1}},
			},
			true,
		},
		{
			"saved days but no bookmark",
			unpack.User{Name: "AA", Registered: rsrc.ParseDay("2018-01-10")},
			rsrc.ParseDay("2018-01-13"),
			nil,
			[][]info.Song{
				{{Artist: "XX"}, {Artist: "XX"}},
				{{Artist: "A"}},
				{{Artist: "hui"}, {Artist: "hui"}},
				{{Artist: "hui"}},
			},
			map[rsrc.Locator][]byte{
				h0:    []byte(`{"recenttracks":{"track":[{"artist":{"#text":"XX"}},{"artist":{"#text":"XX"}}], "@attr":{"totalPages":"1"}}}`),
				h1:    []byte(`{"recenttracks":{"track":[{"artist":{"#text":"ASDF"}}], "@attr":{"totalPages":"1"}}}`),
				h2:    []byte(`{"recenttracks":{"track":[{"artist":{"#text":"hui"}},{"artist":{"#text":"hui"}}], "@attr":{"totalPages":"1"}}}`),
				h3:    []byte(`{"recenttracks":{"track":[{"artist":{"#text":"hui"}},{"artist":{"#text":"hui"}}], "@attr":{"totalPages":"1"}}}`),
				bm:    nil,
				b0:    nil,
				b1:    nil,
				b2:    nil,
				b3:    nil,
				tASDF: []byte(`{"track":{"duration":"120000"}}`),
				tXX:   []byte(`{"track":{"duration":"60000"}}`),
				thui:  []byte(`{"track":{"duration":"60000"}}`),
			},
			map[rsrc.Locator][]byte{},
			[][]info.Song{
				{{Artist: "XX", Duration: 1}, {Artist: "XX", Duration: 1}},
				{{Artist: "ASDF", Duration: 2}},
				{{Artist: "hui", Duration: 1}, {Artist: "hui", Duration: 1}},
				{{Artist: "hui", Duration: 1}, {Artist: "hui", Duration: 1}},
			},
			true,
		},
		{
			"download error",
			unpack.User{Name: "AA", Registered: rsrc.ParseDay("2018-01-10")},
			rsrc.ParseDay("2018-01-10"),
			rsrc.ParseDay("2018-01-10"),
			[][]info.Song{},
			map[rsrc.Locator][]byte{bm: nil},
			map[rsrc.Locator][]byte{},
			[][]info.Song{},
			false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.tracksFile[rsrc.SongHistory(tc.user.Name)] = nil
			io1, _ := mock.IO(tc.tracksFile, mock.Path)
			if tc.saved != nil {
				for i, songs := range tc.saved {
					if err := unpack.WriteDayHistory(songs, tc.user.Name, tc.user.Registered.AddDate(0, 0, i), io1); err != nil {
						t.Fatal("unexpected error during write of all day plays:", err)
					}
				}
				// if err := unpack.WriteSongHistory(tc.saved, tc.user.Name, io1); err != nil {
				// 	t.Fatal("unexpected error during write of all day plays:", err)
				// }
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

			store, _ := io.NewStore([][]rsrc.IO{{io0}, {io1}})

			plays, err := UpdateHistory(&tc.user, tc.until, store)
			if err != nil && tc.ok {
				t.Error("unexpected error:", err)
			} else if err == nil && !tc.ok {
				t.Error("expected error but none occurred")
			}
			if err == nil {
				if !reflect.DeepEqual(plays, tc.plays) {
					t.Errorf("updated plays faulty:\nhas:      %v\nexpected: %v",
						printSongs(plays), printSongs(tc.plays))
				}
			}
		})
	}
}

type printSongs [][]info.Song

func (s printSongs) String() string {
	var sb strings.Builder
	sb.WriteString("[")

	for _, songs := range s {
		sb.WriteString("[")
		for _, song := range songs {
			sb.WriteString(fmt.Sprintf("A:%v T:%v A:%v D:%v', ",
				song.Artist, song.Title, song.Album, song.Duration))
		}
		sb.WriteString("]")
	}

	sb.WriteString("]")
	return sb.String()
}
