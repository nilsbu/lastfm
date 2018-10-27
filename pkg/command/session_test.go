package command

import (
	"reflect"
	"testing"

	"github.com/nilsbu/lastfm/pkg/format"
	"github.com/nilsbu/lastfm/pkg/rsrc"
	"github.com/nilsbu/lastfm/pkg/store"
	"github.com/nilsbu/lastfm/pkg/unpack"
	"github.com/nilsbu/lastfm/test/mock"
)

func TestSessionInfo(t *testing.T) {
	cases := []struct {
		session *unpack.SessionInfo
		ok      bool
	}{
		{&unpack.SessionInfo{User: "U"}, true},
		{nil, true},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			d := mock.NewDisplay()
			cmd := sessionInfo{}

			err := cmd.Execute(c.session, nil, d)
			if err != nil && c.ok {
				t.Fatalf("unexpected error: %v", err)
			} else if err == nil && !c.ok {
				t.Fatalf("expected error but none occurred")
			}

			if len(d.Msgs) == 0 {
				t.Fatalf("no message was printed")
			} else if len(d.Msgs) > 1 {
				t.Fatalf("got %v messages but expected 1", len(d.Msgs))
			} else {
				if reflect.TypeOf(d.Msgs[0]) != reflect.TypeOf(&format.Message{}) {
					t.Errorf("unexpected formatter type: %v", reflect.TypeOf(d.Msgs[0]))
				}
				// TODO content of the message shouldn't be trusted blindly
			}
		})
	}
}

func TestSessionStart(t *testing.T) {
	cases := []struct {
		sessionPre  *unpack.SessionInfo
		user        string
		ok          bool
		sessionPost *unpack.SessionInfo
	}{
		{&unpack.SessionInfo{User: "U"}, "A", false, &unpack.SessionInfo{User: "U"}},
		{&unpack.SessionInfo{User: "U"}, "U", false, &unpack.SessionInfo{User: "U"}},
		{nil, "A", true, &unpack.SessionInfo{User: "A"}},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			files, _ := mock.IO(map[rsrc.Locator][]byte{rsrc.SessionInfo(): []byte("")}, mock.Path)
			s, _ := store.New([][]rsrc.IO{[]rsrc.IO{files}})
			d := mock.NewDisplay()
			cmd := sessionStart{user: c.user}

			err := cmd.Execute(c.sessionPre, s, d)
			if err != nil && c.ok {
				t.Fatalf("unexpected error: %v", err)
			} else if err == nil && !c.ok {
				t.Fatalf("expected error but none occurred")
			}
			if err == nil {
				session, err := unpack.LoadSessionInfo(s)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if !reflect.DeepEqual(session, c.sessionPost) {
					t.Errorf("wrong session was stored: %v != %v", session, c.sessionPost)
				}
			}
		})
	}
}

func TestSessionStop(t *testing.T) {
	cases := []struct {
		sessionPre *unpack.SessionInfo
		ok         bool
	}{
		{&unpack.SessionInfo{User: "U"}, true},
		{nil, false},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			files, _ := mock.IO(map[rsrc.Locator][]byte{rsrc.SessionInfo(): []byte("a")}, mock.Path)
			s, _ := store.New([][]rsrc.IO{[]rsrc.IO{files}})
			d := mock.NewDisplay()
			cmd := sessionStop{}

			err := cmd.Execute(c.sessionPre, s, d)
			if err != nil && c.ok {
				t.Fatalf("unexpected error: %v", err)
			} else if err == nil && !c.ok {
				t.Fatalf("expected error but none occurred")
			}
			if err == nil {
				session, _ := unpack.LoadSessionInfo(s)
				if session != nil {
					t.Errorf("session was not deleted: %v", session)
				}
			}
		})
	}
}
