package mock

import (
	"testing"

	"github.com/nilsbu/lastfm/pkg/fail"
	"github.com/nilsbu/lastfm/pkg/rsrc"
)

func TestIO(t *testing.T) {
	cases := []struct {
		resolve   func(loc rsrc.Locator) (string, error)
		files     map[rsrc.Locator][]byte
		loc       rsrc.Locator
		data      []byte
		ctorOK    bool
		writeOK   bool
		writeSev  fail.Severity
		remove    bool
		removeOK  bool
		removeSev fail.Severity
		readOK    bool
	}{
		{ // read what was written
			Path,
			map[rsrc.Locator][]byte{rsrc.APIKey(): nil},
			rsrc.APIKey(),
			[]byte("xxd"),
			true, true, fail.Control,
			false, false, fail.Control,
			true,
		},
		{ // write fails (key not contained in files)
			Path,
			map[rsrc.Locator][]byte{},
			rsrc.APIKey(),
			[]byte("xxd"),
			true, false, fail.Critical,
			true, false, fail.Critical,
			false,
		},
		{ // read from nil is not possible
			Path,
			map[rsrc.Locator][]byte{rsrc.APIKey(): nil},
			rsrc.APIKey(),
			nil,
			true, true, fail.Control,
			false, false, fail.Control,
			false,
		},
		{ // resolve of file path fails
			URL,
			map[rsrc.Locator][]byte{},
			rsrc.APIKey(),
			[]byte(""),
			true, false, fail.Control,
			false, false, fail.Control,
			false,
		},
		{ // unresolvable url in ctor
			URL,
			map[rsrc.Locator][]byte{rsrc.APIKey(): nil},
			rsrc.APIKey(),
			[]byte(""),
			false, false, fail.Control,
			false, false, fail.Control,
			false,
		},
		{ // written and remove
			Path,
			map[rsrc.Locator][]byte{rsrc.APIKey(): nil},
			rsrc.APIKey(),
			[]byte("xxd"),
			true, true, fail.Control,
			true, true, fail.Control,
			false,
		},
		{ // written and remove
			URL,
			map[rsrc.Locator][]byte{},
			rsrc.APIKey(),
			[]byte("xxd"),
			true, false, fail.Control,
			true, false, fail.Control,
			false,
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			io, err := IO(c.files, c.resolve)
			if err != nil {
				if c.ctorOK {
					t.Error("unexprected error in constructor:", err)
				}
				return
			}

			err = io.Write(c.data, c.loc)
			if msg, ok := IsThreatCorrect(err, c.writeOK, c.writeSev); !ok {
				t.Error(msg)
			}

			if c.remove {
				err = io.Remove(c.loc)
				if msg, ok := IsThreatCorrect(err, c.removeOK, c.removeSev); !ok {
					t.Error(msg)
				}
			}

			data, err := io.Read(c.loc)

			if msg, ok := IsThreatCorrect(err, c.readOK, fail.Control); !ok {
				t.Error(msg)
			}
			if err == nil {
				if string(data) != string(c.data) {
					t.Errorf("result does not match:\nresult:   %v\nexpected: %v",
						string(data), string(c.data))
				}
			}
		})
	}
}
