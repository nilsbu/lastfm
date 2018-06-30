package store

import (
	"testing"

	"github.com/nilsbu/lastfm/pkg/rsrc"
	"github.com/nilsbu/lastfm/test/mock"
)

func TestPool(t *testing.T) {
	cases := []struct {
		files       map[rsrc.Locator][]byte
		loc         rsrc.Locator
		data        []byte
		numReaders  int
		numWriters  int
		numRemovers int
		ctorOK      bool
		writeOK     bool
		remove      bool
		removeOK    bool
		readOK      bool
	}{
		{
			map[rsrc.Locator][]byte{},
			rsrc.SessionID(),
			[]byte("asdf"),
			0, 1, 1,
			false, true,
			false, false,
			true,
		},
		{
			map[rsrc.Locator][]byte{},
			rsrc.SessionID(),
			[]byte("asdf"),
			1, 0, 1,
			false, true,
			false, false,
			true,
		},
		{
			map[rsrc.Locator][]byte{},
			rsrc.SessionID(),
			[]byte("asdf"),
			1, 1, 0,
			false, true,
			false, false,
			true,
		},
		{
			map[rsrc.Locator][]byte{rsrc.SessionID(): []byte("asdf")},
			rsrc.SessionID(),
			[]byte("asdf"),
			3, 3, 1,
			true, true,
			false, false,
			true,
		},
		{
			map[rsrc.Locator][]byte{},
			rsrc.SessionID(),
			[]byte("asdf"),
			1, 1, 1,
			true, false,
			false, false,
			false,
		},
		{
			map[rsrc.Locator][]byte{},
			rsrc.SessionID(),
			[]byte("asdf"),
			1, 1, 1,
			true, false,
			true, false,
			false,
		},
		{
			map[rsrc.Locator][]byte{rsrc.SessionID(): []byte("asdf")},
			rsrc.SessionID(),
			nil,
			3, 3, 1,
			true, true,
			true, true,
			false,
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			io, _ := mock.IO(c.files, mock.Path)

			var readers []rsrc.Reader
			var writers []rsrc.Writer
			var removers []rsrc.Remover

			for i := 0; i < c.numReaders; i++ {
				readers = append(readers, io)
			}
			for i := 0; i < c.numWriters; i++ {
				writers = append(writers, io)
			}
			for i := 0; i < c.numRemovers; i++ {
				removers = append(removers, io)
			}

			p, err := NewPool(readers, writers, removers)
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
