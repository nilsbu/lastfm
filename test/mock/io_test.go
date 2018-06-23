package mock

import (
	"testing"

	"github.com/nilsbu/lastfm/pkg/rsrc"
)

func TestIO(t *testing.T) {
	cases := []struct {
		resolve func(loc rsrc.Locator) (string, error)
		files   map[rsrc.Locator][]byte
		loc     rsrc.Locator
		data    []byte
		ctorOK  bool
		writeOK bool
		readOK  bool
	}{
		{ // no data
			Path,
			map[rsrc.Locator][]byte{},
			rsrc.APIKey(),
			[]byte(""),
			true, false, false,
		},
		{ // read what was written
			Path,
			map[rsrc.Locator][]byte{rsrc.APIKey(): nil},
			rsrc.APIKey(),
			[]byte("xxd"),
			true, true, true,
		},
		{ // write fails (key not contained in files)
			Path,
			map[rsrc.Locator][]byte{},
			rsrc.APIKey(),
			[]byte("xxd"),
			true, false, false,
		},
		{ // read from nil is not possible
			Path,
			map[rsrc.Locator][]byte{rsrc.APIKey(): nil},
			rsrc.APIKey(),
			nil,
			true, true, false,
		},
		{ // resolve of file path fails
			URL,
			map[rsrc.Locator][]byte{},
			rsrc.APIKey(),
			[]byte(""),
			true, false, false,
		},
		{ // unresolvable url in ctor
			URL,
			map[rsrc.Locator][]byte{rsrc.APIKey(): nil},
			rsrc.APIKey(),
			[]byte(""),
			false, false, false,
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			r, w, err := IO(c.files, c.resolve)
			if err != nil {
				if c.ctorOK {
					t.Error("unexprected error in constructor:", err)
				}
				return
			}

			err = w.Write(c.data, c.loc)
			if err != nil && c.writeOK {
				t.Error("unexpected error during write:", err)
			} else if err == nil && !c.writeOK {
				t.Error("write should have failed but did not")
			}

			data, err := r.Read(c.loc)
			close(r)
			close(w)

			if err != nil && c.readOK {
				t.Error("unexpected error during read:", err)
			} else if err == nil && !c.readOK {
				t.Error("read should have failed but did not")
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
