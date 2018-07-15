package mock

import (
	"testing"

	"github.com/nilsbu/lastfm/pkg/rsrc"
)

func TestIO(t *testing.T) {
	cases := []struct {
		resolve  func(loc rsrc.Locator) (string, error)
		files    map[rsrc.Locator][]byte
		loc      rsrc.Locator
		data     []byte
		ctorOK   bool
		writeOK  bool
		remove   bool
		removeOK bool
		readOK   bool
	}{
		{ // read what was written
			Path,
			map[rsrc.Locator][]byte{rsrc.APIKey(): nil},
			rsrc.APIKey(),
			[]byte("xxd"),
			true, true,
			false, false,
			true,
		},
		{ // write fails (key not contained in files)
			Path,
			map[rsrc.Locator][]byte{},
			rsrc.APIKey(),
			[]byte("xxd"),
			true, false,
			true, false,
			false,
		},
		{ // read from nil is not possible
			Path,
			map[rsrc.Locator][]byte{rsrc.APIKey(): nil},
			rsrc.APIKey(),
			nil,
			true, true,
			false, false,
			false,
		},
		{ // resolve of file path fails
			URL,
			map[rsrc.Locator][]byte{},
			rsrc.APIKey(),
			[]byte(""),
			true, false,
			false, false,
			false,
		},
		{ // unresolvable url in ctor
			URL,
			map[rsrc.Locator][]byte{rsrc.APIKey(): nil},
			rsrc.APIKey(),
			[]byte(""),
			false, false,
			false, false,
			false,
		},
		{ // written and remove
			Path,
			map[rsrc.Locator][]byte{rsrc.APIKey(): nil},
			rsrc.APIKey(),
			[]byte("xxd"),
			true, true,
			true, true,
			false,
		},
		{ // written and remove
			URL,
			map[rsrc.Locator][]byte{},
			rsrc.APIKey(),
			[]byte("xxd"),
			true, false,
			true, false,
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
			if err != nil && c.writeOK {
				t.Error("unexpected error:", err)
			} else if err == nil && !c.writeOK {
				t.Errorf("expected error but none occurred")
			}

			if c.remove {
				err = io.Remove(c.loc)
				if err != nil && c.removeOK {
					t.Error("unexpected error:", err)
				} else if err == nil && !c.removeOK {
					t.Errorf("expected error but none occurred")
				}
			}

			data, err := io.Read(c.loc)
			if err != nil && c.readOK {
				t.Error("unexpected error:", err)
			} else if err == nil && !c.readOK {
				t.Errorf("expected error but none occurred")
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
