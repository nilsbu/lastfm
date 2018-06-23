package store

import (
	"fmt"
	"testing"

	"github.com/nilsbu/lastfm/pkg/io"
	"github.com/nilsbu/lastfm/pkg/rsrc"
	"github.com/nilsbu/lastfm/test/mock"
)

func TestPoolRead(t *testing.T) {
	cases := []struct {
		data    []byte
		r, d, w bool
		ok      bool
	}{
		// Read from disk (availability of download doesn't matter)
		{[]byte("A"), true, true, false, true},
		{[]byte("B"), true, false, false, true},
		// Downloaded and written
		{[]byte("C"), false, true, true, true},
		// Read and download fails
		{[]byte("D"), false, false, false, false},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("#%v", i), func(t *testing.T) {
			loc, _ := rsrc.UserInfo("sss")

			var files, web map[rsrc.Locator][]byte
			if c.r {
				files = map[rsrc.Locator][]byte{loc: c.data}
			} else {
				files = map[rsrc.Locator][]byte{loc: nil}
			}
			if c.d {
				web = map[rsrc.Locator][]byte{loc: c.data}
			} else {
				web = map[rsrc.Locator][]byte{}
			}

			r, w, _ := mock.IO(files, mock.Path)
			d, _, _ := mock.IO(web, mock.URL)

			p := New(
				[]io.Reader{d},
				[]io.Reader{r},
				[]io.Writer{w})

			data, err := p.Read(loc)
			if err != nil && c.ok {
				t.Error("unexpected error:", err)
			} else if err == nil && !c.ok {
				t.Error("expected and error but none occurred")
			}

			if err == nil {
				if string(data) != string(c.data) {
					t.Errorf("read data is wrong\nread:     %v\nexpected: %v",
						string(data), string(c.data))
				}

				written, err := r.Read(loc)
				if err != nil {
					t.Error("unexpected error while reading witten data:", err)
				}
				if string(written) != string(c.data) {
					t.Errorf("read data is wrong\nread:     %v\nexpected: %v",
						string(written), string(c.data))
				}
			}
		})
	}
}
