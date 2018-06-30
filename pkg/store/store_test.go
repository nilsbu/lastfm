package store

import (
	"fmt"
	"testing"

	"github.com/nilsbu/lastfm/pkg/fail"
	"github.com/nilsbu/lastfm/pkg/rsrc"
	"github.com/nilsbu/lastfm/test/mock"
)

func TestStoreNew(t *testing.T) {
	cases := []struct {
		numReaders  []int
		numWriters  []int
		numRemovers []int
		ok          bool
	}{
		{ // must have at least one layer
			[]int{}, []int{}, []int{},
			false,
		},
		{ // uploader missing
			[]int{1, 1}, []int{0, 1}, []int{1, 1},
			false,
		},
		{ // layers are unequal
			[]int{2, 1, 2}, []int{1, 2}, []int{1, 2},
			false,
		},
		{ // ok
			[]int{2, 1}, []int{1, 2}, []int{1, 2},
			true,
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("#%v", i), func(t *testing.T) {
			files := map[rsrc.Locator][]byte{}
			io, _ := mock.IO(files, mock.Path)

			readers := make([][]rsrc.Reader, len(c.numReaders))
			for i := range readers {
				reads := []rsrc.Reader{}
				for j := 0; j < c.numReaders[i]; j++ {
					reads = append(reads, io)
				}
				readers[i] = reads
			}

			writers := make([][]rsrc.Writer, len(c.numWriters))
			for i := range writers {
				writes := []rsrc.Writer{}
				for j := 0; j < c.numWriters[i]; j++ {
					writes = append(writes, io)
				}
				writers[i] = writes
			}

			removers := make([][]rsrc.Remover, len(c.numRemovers))
			for i := range removers {
				rms := []rsrc.Remover{}
				for j := 0; j < c.numRemovers[i]; j++ {
					rms = append(rms, io)
				}
				removers[i] = rms
			}

			s, err := New(readers, writers, removers)
			if str, ok := mock.IsThreatCorrect(err, c.ok, fail.Critical); !ok {
				t.Error(str)
			}
			if err == nil && s == nil {
				t.Error("store cannot be nil if no error was returned")
			}
		})
	}
}

func TestStoreRead(t *testing.T) {
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
		{
			[]map[rsrc.Locator][]byte{
				map[rsrc.Locator][]byte{},
			},
			[]mock.Resolver{mock.Path},
			nil,
			apiKey,
			[][]byte{nil},
			false, fail.Control,
		},
		{
			[]map[rsrc.Locator][]byte{
				map[rsrc.Locator][]byte{userInfo: []byte("xx")},
				map[rsrc.Locator][]byte{userInfo: nil},
			},
			[]mock.Resolver{mock.URL, mock.Path},
			[]byte("xx"),
			userInfo,
			[][]byte{[]byte("xx"), []byte("xx")},
			true, fail.Control,
		},
		{
			[]map[rsrc.Locator][]byte{
				map[rsrc.Locator][]byte{},
				map[rsrc.Locator][]byte{userInfo: []byte("9")},
			},
			[]mock.Resolver{mock.URL, mock.Path},
			[]byte("9"),
			userInfo,
			[][]byte{nil, []byte("9")},
			true, fail.Control,
		},
		{
			[]map[rsrc.Locator][]byte{
				map[rsrc.Locator][]byte{userInfo: nil},
				map[rsrc.Locator][]byte{userInfo: []byte("xx")},
				map[rsrc.Locator][]byte{userInfo: nil},
				map[rsrc.Locator][]byte{userInfo: nil},
			},
			[]mock.Resolver{mock.Path, mock.Path, mock.Path, mock.Path},
			[]byte("xx"),
			userInfo,
			[][]byte{nil, []byte("xx"), []byte("xx"), []byte("xx")},
			true, fail.Control,
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			var readers [][]rsrc.Reader
			var writers [][]rsrc.Writer
			var removers [][]rsrc.Remover
			for i := range c.files {
				io, err := mock.IO(c.files[i], c.locf[i])
				if err != nil {
					t.Fatal("setup error")
				}
				readers = append(readers, []rsrc.Reader{io})
				writers = append(writers, []rsrc.Writer{io})
				removers = append(removers, []rsrc.Remover{io})
			}

			s, err := New(readers, writers, removers)
			if err != nil {
				t.Error("unexpected error in constructor")
			}

			data, err := s.Read(c.loc)
			if str, ok := mock.IsThreatCorrect(err, c.ok, c.sev); !ok {
				t.Error(str)
			}
			if string(data) != string(c.data) {
				t.Errorf("read data is wrong:\nhas:      '%v'\nexpected: '%v'",
					string(data), string(c.data))
			}

			for i, rs := range readers {
				data, err := rs[0].Read(c.loc)

				if err == nil && string(data) != string(c.written[i]) {
					t.Errorf("written data false at level %v:\n"+
						"has:      '%v'\nexpected: '%v'",
						i, string(data), string(c.written[i]))
				}
			}
		})
	}
}

func TestStoreUpdate(t *testing.T) {
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
		{
			[]map[rsrc.Locator][]byte{
				map[rsrc.Locator][]byte{},
			},
			[]mock.Resolver{mock.Path},
			nil,
			apiKey,
			[][]byte{nil},
			false, fail.Control,
		},
		{
			[]map[rsrc.Locator][]byte{
				map[rsrc.Locator][]byte{userInfo: []byte("xx")},
				map[rsrc.Locator][]byte{userInfo: nil},
			},
			[]mock.Resolver{mock.URL, mock.Path},
			[]byte("xx"),
			userInfo,
			[][]byte{[]byte("xx"), []byte("xx")},
			true, fail.Control,
		},
		{
			[]map[rsrc.Locator][]byte{
				map[rsrc.Locator][]byte{},
				map[rsrc.Locator][]byte{apiKey: []byte("9")},
			},
			[]mock.Resolver{mock.URL, mock.Path},
			[]byte("9"),
			apiKey,
			[][]byte{nil, []byte("9")},
			true, fail.Control,
		},
		{
			[]map[rsrc.Locator][]byte{
				map[rsrc.Locator][]byte{userInfo: nil},
				map[rsrc.Locator][]byte{userInfo: []byte("9")},
			},
			[]mock.Resolver{mock.URL, mock.Path},
			[]byte("9"),
			userInfo,
			[][]byte{nil, []byte("9")},
			true, fail.Control,
		},
		{
			[]map[rsrc.Locator][]byte{
				map[rsrc.Locator][]byte{userInfo: nil},
				map[rsrc.Locator][]byte{userInfo: nil},
			},
			[]mock.Resolver{mock.Path, mock.Path},
			nil,
			userInfo,
			[][]byte{nil, nil},
			false, fail.Control,
		},
		{
			[]map[rsrc.Locator][]byte{
				map[rsrc.Locator][]byte{apiKey: []byte("9")},
				map[rsrc.Locator][]byte{},
			},
			[]mock.Resolver{mock.Path, mock.URL},
			[]byte("9"),
			apiKey,
			[][]byte{[]byte("9"), nil},
			true, fail.Control,
		},
		{
			[]map[rsrc.Locator][]byte{
				map[rsrc.Locator][]byte{userInfo: []byte("9")},
				map[rsrc.Locator][]byte{},
			},
			[]mock.Resolver{mock.Path, mock.Path},
			[]byte("9"),
			userInfo,
			[][]byte{[]byte("9"), nil},
			false, fail.Critical,
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			var readers [][]rsrc.Reader
			var writers [][]rsrc.Writer
			var removers [][]rsrc.Remover
			for i := range c.files {
				io, err := mock.IO(c.files[i], c.locf[i])
				if err != nil {
					t.Fatal("setup error")
				}
				readers = append(readers, []rsrc.Reader{io})
				writers = append(writers, []rsrc.Writer{io})
				removers = append(removers, []rsrc.Remover{io})
			}

			s, err := New(readers, writers, removers)
			if err != nil {
				t.Error("unexpected error in constructor")
			}

			data, err := s.Update(c.loc)
			if str, ok := mock.IsThreatCorrect(err, c.ok, c.sev); !ok {
				t.Error(str)
			}
			if string(data) != string(c.data) {
				t.Errorf("read data is wrong:\nhas:      '%v'\nexpected: '%v'",
					string(data), string(c.data))
			}

			for i, rs := range readers {
				data, err := rs[0].Read(c.loc)

				if err == nil {
					if string(data) != string(c.written[i]) {
						t.Errorf("written data false at level %v:\n"+
							"has:      '%v'\nexpected: '%v'",
							i, string(data), string(c.written[i]))
					}
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
			[][]byte{[]byte("xx"), nil},
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
			var readers [][]rsrc.Reader
			var writers [][]rsrc.Writer
			var removers [][]rsrc.Remover
			for i := range c.files {
				io, err := mock.IO(c.files[i], c.locf[i])
				if err != nil {
					t.Fatal("setup error")
				}
				readers = append(readers, []rsrc.Reader{io})
				writers = append(writers, []rsrc.Writer{io})
				removers = append(removers, []rsrc.Remover{io})
			}

			s, err := New(readers, writers, removers)
			if err != nil {
				t.Error("unexpected error in constructor")
			}

			err = s.Write(c.data, c.loc)
			if str, ok := mock.IsThreatCorrect(err, c.ok, c.sev); !ok {
				t.Error(str)
			}

			for i, rs := range readers {
				data, err := rs[0].Read(c.loc)

				if err == nil && string(data) != string(c.written[i]) {
					t.Errorf("written data false at level %v:\n"+
						"has:      '%v'\nexpected: '%v'",
						i, string(data), string(c.written[i]))
				}
			}
		})
	}
}

func TestStoreRemove(t *testing.T) {
	apiKey := rsrc.APIKey()
	userInfo, _ := rsrc.UserInfo("abc")

	cases := []struct {
		files []map[rsrc.Locator][]byte
		locf  []mock.Resolver
		loc   rsrc.Locator
		exist []bool
		ok    bool
		sev   fail.Severity
	}{
		{ // failed remove (critical)
			[]map[rsrc.Locator][]byte{
				map[rsrc.Locator][]byte{},
			},
			[]mock.Resolver{mock.Path},
			userInfo,
			[]bool{false},
			false, fail.Critical,
		},
		{ // remove both
			[]map[rsrc.Locator][]byte{
				map[rsrc.Locator][]byte{userInfo: []byte("xx")},
				map[rsrc.Locator][]byte{userInfo: []byte("xx")},
			},
			[]mock.Resolver{mock.URL, mock.Path},
			userInfo,
			[]bool{false, false},
			true, fail.Control,
		},
		{ // level 0 not removed since 1 failes
			[]map[rsrc.Locator][]byte{
				map[rsrc.Locator][]byte{apiKey: []byte("xx")},
				map[rsrc.Locator][]byte{},
			},
			[]mock.Resolver{mock.Path, mock.URL},
			apiKey,
			[]bool{true, false},
			true, fail.Control,
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			var readers [][]rsrc.Reader
			var writers [][]rsrc.Writer
			var removers [][]rsrc.Remover
			for i := range c.files {
				io, err := mock.IO(c.files[i], c.locf[i])
				if err != nil {
					t.Fatal("setup error")
				}
				readers = append(readers, []rsrc.Reader{io})
				writers = append(writers, []rsrc.Writer{io})
				removers = append(removers, []rsrc.Remover{io})
			}

			s, err := New(readers, writers, removers)
			if err != nil {
				t.Error("unexpected error in constructor")
			}

			err = s.Remove(c.loc)
			if str, ok := mock.IsThreatCorrect(err, c.ok, c.sev); !ok {
				t.Error(str)
			}

			if err != nil {
				for i, rs := range readers {
					_, err := rs[0].Read(c.loc)

					if err != nil && c.exist[i] {
						t.Errorf("file at level %v does not exists but should", i)
					} else if err == nil && !c.exist[i] {
						t.Errorf("file at level %v exists but should not", i)
					}
				}
			}
		})
	}
}
