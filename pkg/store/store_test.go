package store

import (
	"fmt"
	"testing"

	"github.com/nilsbu/lastfm/pkg/fail"
	"github.com/nilsbu/lastfm/pkg/io"
	"github.com/nilsbu/lastfm/pkg/rsrc"
	"github.com/nilsbu/lastfm/test/mock"
)

func TestPoolRead(t *testing.T) {
	cases := []struct {
		data       []byte
		numReaders []int
		numWriters []int
		r, d, w    bool
		ctorOK     bool
		ok         bool
	}{
		{ // uploader missing
			[]byte("A"),
			[]int{1, 1}, []int{0, 1},
			true, true, false,
			false, true,
		},
		{ // Read from disk (availability of download doesn't matter)
			[]byte("A"),
			[]int{1, 1}, []int{1, 1},
			true, true, false,
			true, true,
		},
		{
			[]byte("B"),
			[]int{1, 1}, []int{1, 1},
			true, false, false,
			true, true,
		},

		{ // Downloaded and written
			[]byte("C"),
			[]int{1, 1}, []int{1, 1},
			false, true, true,
			true, true,
		},
		{ // Read and download fails
			[]byte("D"),
			[]int{1, 1}, []int{1, 1},
			false, false, false,
			true, false,
		},
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

			readers := make([][]io.Reader, 2)
			writers := make([][]io.Writer, 2)

			reads := []io.Reader{}
			for i := 0; i < c.numReaders[0]; i++ {
				reads = append(reads, d)
			}
			readers[0] = reads

			reads = []io.Reader{}
			for i := 0; i < c.numReaders[1]; i++ {
				reads = append(reads, r)
			}
			readers[1] = reads

			writes := []io.Writer{}
			for i := 0; i < c.numWriters[0]; i++ {
				writes = append(writes, io.FailIO{})
			}
			writers[0] = writes

			writes = []io.Writer{}
			for i := 0; i < c.numWriters[1]; i++ {
				writes = append(writes, w)
			}
			writers[1] = writes

			p, err := New(
				readers,
				writers)
			if err != nil && c.ctorOK {
				t.Error("unexpected error in constructor:", err)
			} else if err == nil && !c.ctorOK {
				t.Error("expected error in constructor but none occurred")
			}
			if err != nil {
				return
			}

			data, err := p.Read(loc)
			if err != nil && c.ok {
				t.Error("unexpected error:", err)
			} else if err == nil && !c.ok {
				t.Error("expected error but none occurred")
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

func TestStoreWrite(t *testing.T) {
	apiKey := rsrc.APIKey()
	userInfo, _ := rsrc.UserInfo("abc")

	cases := []struct {
		files   []map[rsrc.Locator][]byte
		locf    []mock.Resolver
		data    []byte
		loc     rsrc.Locator
		written [][]byte
		ok      bool
		sev     fail.Severity
	}{
		{ // failed write (critical)
			[]map[rsrc.Locator][]byte{
				map[rsrc.Locator][]byte{},
			},
			[]mock.Resolver{mock.Path},
			[]byte("xx"),
			userInfo,
			[][]byte{nil},
			false, fail.Critical,
		},
		{ // not written in layer 0 (no URL for APIKey)
			[]map[rsrc.Locator][]byte{
				map[rsrc.Locator][]byte{},
				map[rsrc.Locator][]byte{apiKey: nil},
			},
			[]mock.Resolver{mock.URL, mock.Path},
			[]byte("xx"),
			apiKey,
			[][]byte{nil, []byte("xx")},
			true, fail.Control,
		},
		{ // written in neither (0 not reached since 1 fails)
			[]map[rsrc.Locator][]byte{
				map[rsrc.Locator][]byte{apiKey: nil},
				map[rsrc.Locator][]byte{},
			},
			[]mock.Resolver{mock.Path, mock.URL},
			[]byte("xx"),
			apiKey,
			[][]byte{nil, nil},
			true, fail.Control,
		},
		{ // written in both layers
			[]map[rsrc.Locator][]byte{
				map[rsrc.Locator][]byte{userInfo: nil},
				map[rsrc.Locator][]byte{userInfo: nil},
			},
			[]mock.Resolver{mock.URL, mock.Path},
			[]byte("xx"),
			userInfo,
			[][]byte{[]byte("xx"), []byte("xx")},
			true, fail.Control,
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			var readers [][]io.Reader
			var writers [][]io.Writer
			for i := range c.files {
				r, w, err := mock.IO(c.files[i], c.locf[i])
				if err != nil {
					t.Fatal("setup error")
				}
				readers = append(readers, []io.Reader{r})
				writers = append(writers, []io.Writer{w})
			}

			s, err := New(readers, writers)
			if err != nil {
				t.Error("unexpected error in constructor")
			}

			err = s.Write(c.data, c.loc)
			if str, ok := mock.IsThreatCorrect(err, c.ok, c.sev); !ok {
				t.Error(str)
			}

			for i, rs := range readers {
				data, err := rs[0].Read(c.loc)

				if err != nil {
					if c.written[i] != nil {
						if string(data) != string(c.written[i]) {
							t.Errorf("written data false at level %v:\n"+
								"has:      '%v'\nexpected: '%v'",
								i, string(data), string(c.written[i]))
						}
					} else {
						if c.written[i] != nil {
							t.Errorf("expected written in layer %v but was not", i)
						}
					}
				}
			}
		})
	}
}
