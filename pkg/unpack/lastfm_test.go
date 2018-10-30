package unpack

import (
	"reflect"
	"testing"

	"github.com/nilsbu/lastfm/pkg/charts"
	"github.com/nilsbu/lastfm/pkg/rsrc"
	"github.com/nilsbu/lastfm/test/mock"
)

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

			user, err := LoadUserInfo(c.name, io)
			if err != nil && c.ok {
				t.Error("unexpected error:", err)
			} else if err == nil && !c.ok {
				t.Error("expected error")
			}

			if err == nil {
				if user.Name != c.user.Name {
					t.Error("wrong name")
				}

				hasMidn, _ := user.Registered.Midnight()
				wantMidn, _ := c.user.Registered.Midnight()
				if hasMidn != wantMidn {
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
			[]byte(`{"recenttracks":{"track":[{"artist":{"#text":"ASDF"}},{"artist":{"#text":"ASDF"}}], "@attr":{"totalPages":"1"}}}`),
			"user", rsrc.ToDay(86400), 1,
			&HistoryDayPage{charts.Charts{"ASDF": []float64{2}}, 1},
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

			hist, err := LoadHistoryDayPage(c.user, c.page, c.day, io)
			if err != nil && c.ok {
				t.Error("unexpected error:", err)
			} else if err == nil && !c.ok {
				t.Error("expected error")
			}

			if err != nil {
				if !reflect.DeepEqual(hist, c.hist) {
					t.Errorf("wrong data:\n has:  %v\nwant: %v",
						hist, c.hist)
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
			[]TagCount{TagCount{"bui", 100}, TagCount{"asdf", 12}},
			true,
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			io, err := mock.IO(c.files, mock.Path)
			if err != nil {
				t.Fatal("setup error")
			}

			tags, err := LoadArtistTags(c.artist, io)
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

func TestLoadTagInfo(t *testing.T) {
	cases := []struct {
		files map[rsrc.Locator][]byte
		names [][]string
		tags  []*charts.Tag
		ok    bool
	}{
		{
			map[rsrc.Locator][]byte{rsrc.TagInfo("african"): nil},
			[][]string{[]string{"african"}},
			[]*charts.Tag{nil},
			false,
		},
		{
			map[rsrc.Locator][]byte{rsrc.TagInfo("african"): []byte(`{"user":{"name":"xy","registered":{"unixtime":86400}}}`)},
			[][]string{[]string{"african"}},
			[]*charts.Tag{&charts.Tag{}},
			true, // no error is thrown, therefore this is acceppted
		},
		{
			map[rsrc.Locator][]byte{rsrc.TagInfo("african"): []byte(`{"tag":{"name":"african","total":55266,"reach":10493}}`)},
			[][]string{[]string{"african", "african"}},
			[]*charts.Tag{
				&charts.Tag{Name: "african", Total: 55266, Reach: 10493},
				&charts.Tag{Name: "african", Total: 55266, Reach: 10493},
			},
			true,
		},
		{
			map[rsrc.Locator][]byte{rsrc.TagInfo("african"): []byte(`{"tag":{"name":"african","total":55266,"reach":10493}}`)},
			[][]string{[]string{"african"}, []string{"african"}},
			[]*charts.Tag{
				&charts.Tag{Name: "african", Total: 55266, Reach: 10493},
				&charts.Tag{Name: "african", Total: 55266, Reach: 10493},
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

			buf := NewCachedTagLoader(io)

			n := 0
			for _, names := range c.names {
				n += len(names)
			}

			tags := make([]*charts.Tag, n)
			feedback := make(chan error)
			errs := []error{}

			n = 0
			for _, names := range c.names {
				for i := range names {
					go func(i int) {
						res, err := buf.LoadTagInfo(names[i])
						tags[i+n] = res
						feedback <- err
					}(i)
				}

				for range names {
					if err := <-feedback; err != nil {
						errs = append(errs, err)
						if c.ok {
							t.Error("unexpected error :", err)
						}
					}
				}

				n += len(names)
			}

			if len(errs) == 0 {
				if !c.ok {
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
