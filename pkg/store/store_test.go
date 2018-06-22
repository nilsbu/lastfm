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
			rs, _ := rsrc.UserInfo("sss")
			path, _ := rs.Path()
			url, _ := rs.URL(mock.APIKey)

			var files, web map[string][]byte
			if c.r {
				files = map[string][]byte{path: c.data}
			} else {
				files = map[string][]byte{path: nil}
			}
			if c.d {
				web = map[string][]byte{url: c.data}
			} else {
				web = map[string][]byte{}
			}

			r, w := mock.FileIO(files)
			d := mock.Downloader(web)

			p := New(
				[]io.Reader{d},
				[]io.Reader{r},
				[]io.Writer{w})

			data, err := p.Read(rs)
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
				if string(files[path]) != string(c.data) {
					t.Errorf("written data is wrong\nread:     %v\nexpected: %v",
						string(files[path]), string(c.data))
				}
			}
		})
	}
}
