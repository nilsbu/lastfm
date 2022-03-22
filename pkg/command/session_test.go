package command

import (
	"reflect"
	"testing"

	"github.com/nilsbu/lastfm/pkg/format"
	"github.com/nilsbu/lastfm/pkg/io"
	"github.com/nilsbu/lastfm/pkg/pipeline"
	"github.com/nilsbu/lastfm/pkg/rsrc"
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

			pl := pipeline.New(c.session, nil)
			err := cmd.Execute(c.session, nil, pl, d)
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

func TestSession(t *testing.T) {
	cases := []struct {
		name        string
		sessionPre  *unpack.SessionInfo
		cmd         command
		ok          bool
		sessionPost *unpack.SessionInfo
	}{
		{
			"start: session with other user running",
			&unpack.SessionInfo{User: "U"},
			sessionStart{user: "A"},
			false,
			&unpack.SessionInfo{User: "U"},
		},
		{
			"start: session with same user running",
			&unpack.SessionInfo{User: "U"},
			sessionStart{user: "U"},
			false,
			&unpack.SessionInfo{User: "U"},
		},
		{
			"start: successful",
			nil,
			sessionStart{user: "A"},
			true,
			&unpack.SessionInfo{User: "A"},
		},
		{
			"stop: successful",
			&unpack.SessionInfo{User: "U"},
			sessionStop{},
			true,
			nil,
		},
		{
			"stop: no session running",
			nil,
			sessionStop{},
			false,
			nil,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			files, _ := mock.IO(map[rsrc.Locator][]byte{rsrc.SessionInfo(): []byte("")}, mock.Path)
			s, _ := io.NewStore([][]rsrc.IO{{files}})
			d := mock.NewDisplay()

			pl := pipeline.New(c.sessionPre, s)
			err := c.cmd.Execute(c.sessionPre, s, pl, d)
			if err != nil && c.ok {
				t.Fatalf("unexpected error: %v", err)
			} else if err == nil && !c.ok {
				t.Fatalf("expected error but none occurred")
			}
			if err == nil {
				session, err := unpack.LoadSessionInfo(s)
				if c.sessionPost != nil && err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if !reflect.DeepEqual(session, c.sessionPost) {
					t.Errorf("wrong session was stored: %v != %v", session, c.sessionPost)
				}
			}
		})
	}
}
