package unpack

import (
	"reflect"
	"testing"

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
			[]byte(`{"recenttracks":{"track":[{"artist":{"#text":"ASDF"}}], "@attr":{"totalPages":"1"}}}`),
			"user", rsrc.ToDay(86400), 1,
			&HistoryDayPage{map[string]int{"ASDF": 1}, 1},
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
