package unpack

import (
	"reflect"
	"testing"

	"github.com/nilsbu/lastfm/pkg/charts2"
	"github.com/nilsbu/lastfm/pkg/rsrc"
	"github.com/nilsbu/lastfm/test/mock"
)

func TestLastfmError(t *testing.T) {
	err := &LastfmError{
		Code:    3,
		Message: "some error",
	}

	if err.Error() != "LastFM error (code = 3): some error" {
		t.Errorf("wrong error message: '%v'",
			err.Error())
	}
}

func TestLoadUserInfo(t *testing.T) {
	cases := []struct {
		json []byte
		name string
		user *User
		ok   bool
	}{
		{
			[]byte(`{"user":{"name":"What","playcount":1928,"registered":{"unixtime":114004225884}}}`),
			"What",
			&User{"What", rsrc.ToDay(114004195200)},
			true,
		},
		{
			[]byte(`{"user":{"name":"What","playcount":1928,`),
			"What",
			nil,
			false,
		},
		{
			nil,
			"What",
			nil,
			false,
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			io, err := mock.IO(
				map[rsrc.Locator][]byte{rsrc.UserInfo(c.name): c.json},
				mock.Path)
			if err != nil {
				t.Fatal("setup error")
			}

			user, err := LoadUserInfo(c.name, NewCacheless(io))
			if err != nil && c.ok {
				t.Error("unexpected error:", err)
			} else if err == nil && !c.ok {
				t.Error("expected error")
			}

			if err == nil {
				if user.Name != c.user.Name {
					t.Error("wrong name")
				}

				if user.Registered.Midnight() != c.user.Registered.Midnight() {
					t.Error("wrong registered")
				}
			}
		})
	}
}

func TestWriteUserInfo(t *testing.T) {
	cases := []struct {
		user *User
		json []byte
		ok   bool
	}{
		{
			&User{"What", rsrc.ToDay(114004195200)},
			[]byte(`{"user":{"name":"What","playcount":0,"registered":{"unixtime":114004195200}}}`),
			true,
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			io, err := mock.IO(
				map[rsrc.Locator][]byte{rsrc.UserInfo(c.user.Name): nil},
				mock.Path)
			if err != nil {
				t.Fatal("setup error")
			}

			err = WriteUserInfo(c.user, io)
			if err != nil && c.ok {
				t.Error("unexpected error:", err)
			} else if err == nil && !c.ok {
				t.Error("expected error")
			}

			if err == nil {
				json, err := io.Read(rsrc.UserInfo(c.user.Name))
				if err != nil {
					t.Fatalf("load error: %v", err)
				}

				if string(json) != string(c.json) {
					t.Errorf("wrong data: '%v' != '%v'", string(json), string(c.json))
				}
			}
		})
	}
}

func TestLoadHistoryDayPage(t *testing.T) {
	song1 := `{"artist":{"#text":"ASDF"},"name":"x","album":{"#text":"q"}}`
	song2 := `{"artist":{"#text":"ASDF"},"name":"y","album":{"#text":"q"}}`

	cases := []struct {
		json []byte
		user string
		day  rsrc.Day
		page int
		hist *HistoryDayPage
		ok   bool
	}{
		{
			[]byte{},
			"user", rsrc.ToDay(86400), 1,
			nil,
			false,
		},
		{
			[]byte(`{"recenttracks":{"track":[` + song1 + `,` + song2 + `], "@attr":{"totalPages":"1"}}}`),
			"user", rsrc.ToDay(86400), 1,
			&HistoryDayPage{
				[]charts2.Song{
					{
						Artist: "ASDF",
						Title:  "x",
						Album:  "q",
					},
					{
						Artist: "ASDF",
						Title:  "y",
						Album:  "q",
					},
				}, 1},
			true,
		},
		{
			[]byte(`{"recenttracks":{"@attr":{"page":"1","total":"0","user":"NBooN","perPage":"200","totalPages":"0"},"track":{"artist":{"mbid":"846e89f6-6257-4371-a26d-de960a60bec5","#text":"The Coup"},"@attr":{"nowplaying":"true"},"mbid":"293b4bc9-95c3-3032-a59f-53d6dfba5263","album":{"mbid":"e2f0f87f-763a-498e-9823-decef2cf62b3","#text":"Pick A Bigger Weapon"},"streamable":"0","url":"https:\/\/www.last.fm\/music\/The+Coup\/_\/My+Favorite+Mutiny","name":"My Favorite Mutiny"}}}`),
			"user", rsrc.ToDay(86400), 1,
			&HistoryDayPage{
				[]charts2.Song{}, 0},
			true,
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			io, err := mock.IO(
				map[rsrc.Locator][]byte{rsrc.History(c.user, c.page, c.day): c.json},
				mock.Path)
			if err != nil {
				t.Fatal("setup error")
			}

			hist, err := LoadHistoryDayPage(c.user, c.page, c.day, NewCacheless(io))
			if err != nil && c.ok {
				t.Error("unexpected error:", err)
			} else if err == nil && !c.ok {
				t.Error("expected error")
			}

			if err == nil {
				if !reflect.DeepEqual(hist, c.hist) {
					t.Errorf("wrong data:\n has:  %v\nwant: %v",
						hist, c.hist)
				}
			}
		})
	}
}

func TestLoadArtistInfo(t *testing.T) {
	cases := []struct {
		files     map[rsrc.Locator][]byte
		artist    string
		listeners int64
		playCount int64
		ok        bool
	}{
		{
			map[rsrc.Locator][]byte{rsrc.ArtistInfo("xy"): nil},
			"xy",
			0, 0,
			false,
		},
		{
			map[rsrc.Locator][]byte{rsrc.ArtistInfo("xy"): []byte(`{"artist":{"name":"xy","stats":{"listeners":"119","playcount":"4400"}}}`)},
			"xy",
			119, 4400,
			true,
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			io, err := mock.IO(c.files, mock.Path)
			if err != nil {
				t.Fatal("setup error")
			}

			info, err := LoadArtistInfo(c.artist, NewCacheless(io))
			if err != nil && c.ok {
				t.Error("unexpected error:", err)
			} else if err == nil && !c.ok {
				t.Error("expected error")
			}

			if err == nil {
				if info.Name != c.artist {
					t.Errorf("expected artist name '%v' but got '%v'",
						c.artist, info.Name)
				}
				if info.Listeners != c.listeners {
					t.Errorf("expected %v listeners but got %v",
						c.listeners, info.Listeners)
				}
				if info.PlayCount != c.playCount {
					t.Errorf("expected play count %v but got %v",
						c.playCount, info.PlayCount)
				}
			}
		})
	}
}

func TestLoadArtistTags(t *testing.T) {
	cases := []struct {
		files  map[rsrc.Locator][]byte
		artist string
		tags   []TagCount
		ok     bool
	}{
		{
			map[rsrc.Locator][]byte{rsrc.ArtistTags("xy"): nil},
			"xy",
			[]TagCount{},
			false,
		},
		{
			map[rsrc.Locator][]byte{rsrc.ArtistTags("xy"): []byte(`{"user":{"name":"xy","registered":{"unixtime":86400}}}`)},
			"xy",
			[]TagCount{},
			true, // no error thrown, we'll have to except that wrong data is accepted
		},
		{
			map[rsrc.Locator][]byte{rsrc.ArtistTags("xy"): []byte(`{"toptags":{"tag":[{"name":"bui", "count":100},{"count":12,"name":"asdf"}],"@attr":{"artist":"xy"}}}`)},
			"xy",
			[]TagCount{{"bui", 100}, {"asdf", 12}},
			true,
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			io, err := mock.IO(c.files, mock.Path)
			if err != nil {
				t.Fatal("setup error")
			}

			tags, err := LoadArtistTags(c.artist, NewCacheless(io))
			if err != nil && c.ok {
				t.Error("unexpected error:", err)
			} else if err == nil && !c.ok {
				t.Error("expected error")
			}

			if err == nil {
				if !reflect.DeepEqual(tags, c.tags) {
					t.Errorf("wrong data:\n has:  %v\nwant: %v",
						tags, c.tags)
				}
			}
		})
	}
}

func TestWriteLoadArtistTags(t *testing.T) {
	// WriteArtistTags only tested in combination with loading for simplicity.
	cases := []struct {
		artist string
		tags   []TagCount
	}{
		{
			"xy",
			[]TagCount{{"bui", 100}, {"asdf", 12}},
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			io, err := mock.IO(
				map[rsrc.Locator][]byte{rsrc.ArtistTags(c.artist): nil},
				mock.Path)
			if err != nil {
				t.Fatal("setup error")
			}

			err = WriteArtistTags(c.artist, c.tags, io)
			if err != nil {
				t.Fatal("unexpected error:", err)
			}

			tags, err := LoadArtistTags(c.artist, NewCacheless(io))
			if err != nil {
				t.Fatal("unexpected error:", err)
			}

			if !reflect.DeepEqual(tags, c.tags) {
				t.Errorf("wrong data:\n has:  %v\nwant: %v",
					tags, c.tags)
			}
		})
	}
}

func TestLoadTagInfo(t *testing.T) {
	cases := []struct {
		files map[rsrc.Locator][]byte
		names [][]string
		tags  []*charts2.Tag
		ok    bool
	}{
		{
			map[rsrc.Locator][]byte{rsrc.TagInfo("african"): nil},
			[][]string{{"african"}},
			[]*charts2.Tag{nil},
			false,
		},
		{
			map[rsrc.Locator][]byte{rsrc.TagInfo("african"): []byte(`{"user":{"name":"xy","registered":{"unixtime":86400}}}`)},
			[][]string{{"african"}},
			[]*charts2.Tag{{}},
			true, // no error is thrown, therefore this is acceppted
		},
		{
			map[rsrc.Locator][]byte{rsrc.TagInfo("african"): []byte(`{"tag":{"name":"african","total":55266,"reach":10493}}`)},
			[][]string{{"african", "african"}},
			[]*charts2.Tag{
				{Name: "african", Total: 55266, Reach: 10493},
				{Name: "african", Total: 55266, Reach: 10493},
			},
			true,
		},
		{
			map[rsrc.Locator][]byte{rsrc.TagInfo("african"): []byte(`{"tag":{"name":"african","total":55266,"reach":10493}}`)},
			[][]string{{"african"}, {"african"}},
			[]*charts2.Tag{
				{Name: "african", Total: 55266, Reach: 10493},
				{Name: "african", Total: 55266, Reach: 10493},
			},
			true,
		},
		{
			map[rsrc.Locator][]byte{
				rsrc.TagInfo("error"):   []byte(`{"error":29,"message":"Rate Limit Exceeded"}`),
				rsrc.TagInfo("african"): []byte(`{"tag":{"name":"african","total":55266,"reach":10493}}`),
			},
			[][]string{{"error"}, {"african"}},
			[]*charts2.Tag{},
			false,
		},
		{
			map[rsrc.Locator][]byte{
				rsrc.TagInfo("error"):   []byte(`{"error":6,"message":"Invalid parameters"}`),
				rsrc.TagInfo("african"): []byte(`{"tag":{"name":"african","total":55266,"reach":10493}}`),
			},
			[][]string{{"error"}, {"african"}},
			[]*charts2.Tag{
				nil,
				{Name: "african", Total: 55266, Reach: 10493}},
			false,
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			io, err := mock.IO(c.files, mock.Path)
			if err != nil {
				t.Fatal("setup error")
			}

			buf := NewCached(io)

			n := 0
			for _, names := range c.names {
				n += len(names)
			}

			tags := make([]*charts2.Tag, n)
			feedback := make(chan error)
			errs := []error{}

			n = 0
			for _, names := range c.names {
				for i := range names {
					go func(i int) {
						res, err := LoadTagInfo(names[i], buf)
						tags[i+n] = res
						feedback <- err
					}(i)
				}

				for range names {
					if err := <-feedback; err != nil {
						errs = append(errs, err)
						if c.ok {
							t.Error("unexpected error:", err)
						}
					}
				}

				n += len(names)
			}

			if !c.ok {
				if len(errs) == 0 {
					t.Error("expected error but none occurred")
				}

				for i, want := range c.tags {
					if !reflect.DeepEqual(tags[i], want) {
						t.Errorf("wrong data at position %v\nhas:  %v\nwant: %v",
							i, tags[i], want)
					}
				}
			}
		})
	}
}

func TestWriteLoadTagInfo(t *testing.T) {
	// WriteTagInfo only tested in combination with loading for simplicity.
	cases := []struct {
		tag *charts2.Tag
	}{
		{
			&charts2.Tag{Name: "african", Total: 55266, Reach: 10493},
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			io, err := mock.IO(
				map[rsrc.Locator][]byte{
					rsrc.TagInfo(c.tag.Name): nil},
				mock.Path)
			if err != nil {
				t.Fatal("setup error")
			}

			ctl := NewCached(io)

			err = WriteTagInfo(c.tag, io)
			if err != nil {
				t.Fatal("unexpected error:", err)
			}

			tag, err := LoadTagInfo(c.tag.Name, ctl)
			if err != nil {
				t.Fatal("unexpected error:", err)
			}

			if !reflect.DeepEqual(tag, c.tag) {
				t.Errorf("wrong data:\n has:  %v\nwant: %v",
					tag, c.tag)
			}
		})
	}
}

func TestLoadTrackInfo(t *testing.T) {
	cases := []struct {
		files  map[rsrc.Locator][]byte
		artist string
		track  string
		info   TrackInfo
		ok     bool
	}{
		{
			map[rsrc.Locator][]byte{rsrc.TrackInfo("xy", "a"): nil},
			"xy", "a",
			TrackInfo{},
			false,
		},
		{
			map[rsrc.Locator][]byte{rsrc.TrackInfo("xy", "a"): []byte(`{"track":{"duration":"123000","listeners":"2","playcount":"3"}}`)},
			"xy", "a",
			TrackInfo{
				Artist:    "xy",
				Track:     "a",
				Duration:  123,
				Listeners: 2,
				Playcount: 3,
			},
			true,
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			io, err := mock.IO(c.files, mock.Path)
			if err != nil {
				t.Fatal("setup error")
			}

			buf := NewCached(io)

			info, err := LoadTrackInfo(c.artist, c.track, buf)
			if err != nil && c.ok {
				t.Error("unexpected error:", err)
			} else if err == nil && !c.ok {
				t.Error("expected error")
			}

			if err == nil {
				if !reflect.DeepEqual(info, c.info) {
					t.Errorf("wrong data:\n has:  %v\nwant: %v",
						info, c.info)
				}
			}
		})
	}
}

func TestWriteLoadTrackInfo(t *testing.T) {
	// WritedTrackInfo only tested in combination with loading for simplicity.
	cases := []struct {
		artist string
		track  string
		info   TrackInfo
	}{
		{
			"xy", "a",
			TrackInfo{
				Artist:    "xy",
				Track:     "a",
				Duration:  123,
				Listeners: 2,
				Playcount: 3,
			},
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			io, err := mock.IO(
				map[rsrc.Locator][]byte{rsrc.TrackInfo(c.artist, c.track): nil},
				mock.Path)
			if err != nil {
				t.Fatal("setup error")
			}

			err = WriteTrackInfo(c.artist, c.track, c.info, io)
			if err != nil {
				t.Fatal("unexpected error:", err)
			}

			info, err := LoadTrackInfo(c.artist, c.track, NewCached(io))
			if err != nil {
				t.Fatal("unexpected error:", err)
			}

			if !reflect.DeepEqual(info, c.info) {
				t.Errorf("wrong data:\n has:  %v\nwant: %v",
					info, c.info)
			}
		})
	}
}
