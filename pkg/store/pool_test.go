package store

import (
	"testing"

	"github.com/nilsbu/lastfm/pkg/rsrc"
	"github.com/nilsbu/lastfm/test/mock"
)

func TestPool(t *testing.T) {
	cases := []struct {
		files    map[rsrc.Locator][]byte
		loc      rsrc.Locator
		data     []byte
		numIOs   int
		ctorOK   bool
		writeOK  bool
		remove   bool
		removeOK bool
		readOK   bool
	}{
		{
			map[rsrc.Locator][]byte{},
			rsrc.SessionInfo(),
			[]byte("asdf"),
			0,
			false, true,
			false, false,
			false,
		},
		{
			map[rsrc.Locator][]byte{rsrc.SessionInfo(): []byte("asdf")},
			rsrc.SessionInfo(),
			[]byte("asdf"),
			3,
			true, true,
			false, false,
			true,
		},
		{
			map[rsrc.Locator][]byte{},
			rsrc.SessionInfo(),
			[]byte("asdf"),
			1,
			true, false,
			false, false,
			false,
		},
		{
			map[rsrc.Locator][]byte{},
			rsrc.SessionInfo(),
			[]byte("asdf"),
			1,
			true, false,
			true, false,
			false,
		},
		{
			map[rsrc.Locator][]byte{rsrc.SessionInfo(): []byte("asdf")},
			rsrc.SessionInfo(),
			nil,
			3,
			true, true,
			true, true,
			false,
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			io, err := mock.IO(c.files, mock.Path)
			if err != nil {
				t.Fatal("setup failed:", err)
			}

			var ios []rsrc.IO
			for i := 0; i < c.numIOs; i++ {
				ios = append(ios, io)
			}

			p, err := NewPool(ios)
			if err != nil {
				if c.ctorOK {
					t.Error("unexpected error in constructor:", err)
				}
				return
			}

			err = <-p.Write(c.data, c.loc)
			if err != nil && c.writeOK {
				t.Error("unexpected error during write:", err)
			} else if err == nil && !c.writeOK {
				t.Error("expected error during write but none occurred")
			}

			if c.remove {
				err = <-p.Remove(c.loc)
				if err != nil && c.removeOK {
					t.Error("unexpected error during remove:", err)
				} else if err == nil && !c.removeOK {
					t.Error("expected error during remove but none occurred")
				}
			}
			if err != nil {
				return
			}

			readResult := <-p.Read(c.loc)
			data, err := readResult.Data, readResult.Err
			if err != nil && c.readOK {
				t.Error("unexpected error during read:", err)
			} else if err == nil && !c.readOK {
				t.Error("expected error during read but none occurred")
			}

			if string(data) != string(c.data) {
				t.Errorf("wrong result: got '%v', expected '%v'",
					string(data), string(c.data))
			}
		})
	}
}
