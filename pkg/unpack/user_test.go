package unpack

import (
	"testing"

	"github.com/nilsbu/lastfm/pkg/rsrc"
	"github.com/nilsbu/lastfm/test/mock"
)

func TestUserInfo(t *testing.T) {
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
