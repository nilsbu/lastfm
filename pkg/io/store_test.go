package io

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/nilsbu/lastfm/pkg/display"
	"github.com/nilsbu/lastfm/pkg/format"
	"github.com/nilsbu/lastfm/pkg/rsrc"
	"github.com/nilsbu/lastfm/test/mock"
)

func TestStoreNew(t *testing.T) {
	cases := []struct {
		numIOs []int
		ok     bool
	}{
		{ // must have at least one layer
			[]int{},
			false,
		},
		{ // layer 0 empty
			[]int{0, 1},
			false,
		},
		{ // ok
			[]int{2, 1},
			true,
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("#%v", i), func(t *testing.T) {
			files := map[rsrc.Locator][]byte{}
			io, _ := mock.IO(files, mock.Path)

			ios := make([][]rsrc.IO, len(c.numIOs))
			for i := range ios {
				x := []rsrc.IO{}
				for j := 0; j < c.numIOs[i]; j++ {
					x = append(x, io)
				}
				ios[i] = x
			}

			s, err := New(ios)
			if err != nil && c.ok {
				t.Fatal("unexpected error:", err)
			} else if err == nil && !c.ok {
				t.Fatal("expected error but non occurred")
			}

			if err == nil && s == nil {
				t.Error("store cannot be nil if no error was returned")
			}
		})
	}
}

func TestStoreRead(t *testing.T) {
	cases := []struct {
		files   []map[rsrc.Locator][]byte
		locf    []mock.Resolver
		data    []byte
		loc     rsrc.Locator
		written [][]byte
		ok      bool
	}{
		{
			[]map[rsrc.Locator][]byte{{}},
			[]mock.Resolver{mock.Path},
			nil,
			rsrc.APIKey(),
			[][]byte{nil},
			false,
		},
		{
			[]map[rsrc.Locator][]byte{
				{rsrc.UserInfo("abc"): []byte("xx")},
				{rsrc.UserInfo("abc"): nil},
			},
			[]mock.Resolver{mock.URL, mock.Path},
			[]byte("xx"),
			rsrc.UserInfo("abc"),
			[][]byte{[]byte("xx"), []byte("xx")},
			true,
		},
		{
			[]map[rsrc.Locator][]byte{
				{},
				{rsrc.UserInfo("abc"): []byte("9")},
			},
			[]mock.Resolver{mock.URL, mock.Path},
			[]byte("9"),
			rsrc.UserInfo("abc"),
			[][]byte{nil, []byte("9")},
			true,
		},
		{
			[]map[rsrc.Locator][]byte{
				{rsrc.UserInfo("abc"): nil},
				{rsrc.UserInfo("abc"): []byte("xx")},
				{rsrc.UserInfo("abc"): nil},
				{rsrc.UserInfo("abc"): nil},
			},
			[]mock.Resolver{mock.Path, mock.Path, mock.Path, mock.Path},
			[]byte("xx"),
			rsrc.UserInfo("abc"),
			[][]byte{nil, []byte("xx"), []byte("xx"), []byte("xx")},
			true,
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			var ios [][]rsrc.IO
			for i := range c.files {
				io, err := mock.IO(c.files[i], c.locf[i])
				if err != nil {
					t.Fatal("setup error")
				}
				ios = append(ios, []rsrc.IO{io})
			}

			s, err := New(ios)
			if err != nil {
				t.Fatal("unexpected error in constructor:", err)
			}

			data, err := s.Read(c.loc)
			if err != nil && c.ok {
				t.Fatal("unexpected error:", err)
			} else if err == nil && !c.ok {
				t.Fatal("expected error but non occurred")
			}

			if string(data) != string(c.data) {
				t.Errorf("read data is wrong:\nhas:      '%v'\nexpected: '%v'",
					string(data), string(c.data))
			}

			for i, io := range ios {
				data, err := io[0].Read(c.loc)

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
	cases := []struct {
		files   []map[rsrc.Locator][]byte
		locf    []mock.Resolver
		data    []byte
		loc     rsrc.Locator
		written [][]byte
		ok      bool
	}{
		{
			[]map[rsrc.Locator][]byte{{}},
			[]mock.Resolver{mock.Path},
			nil,
			rsrc.APIKey(),
			[][]byte{nil},
			false,
		},
		{
			[]map[rsrc.Locator][]byte{
				{rsrc.UserInfo("abc"): []byte("xx")},
				{rsrc.UserInfo("abc"): nil},
			},
			[]mock.Resolver{mock.URL, mock.Path},
			[]byte("xx"),
			rsrc.UserInfo("abc"),
			[][]byte{[]byte("xx"), []byte("xx")},
			true,
		},
		{
			[]map[rsrc.Locator][]byte{
				{},
				{rsrc.APIKey(): []byte("9")},
			},
			[]mock.Resolver{mock.URL, mock.Path},
			[]byte("9"),
			rsrc.APIKey(),
			[][]byte{nil, []byte("9")},
			true,
		},
		{
			[]map[rsrc.Locator][]byte{
				{rsrc.UserInfo("abc"): nil},
				{rsrc.UserInfo("abc"): []byte("9")},
			},
			[]mock.Resolver{mock.URL, mock.Path},
			[]byte("9"),
			rsrc.UserInfo("abc"),
			[][]byte{nil, []byte("9")},
			true,
		},
		{
			[]map[rsrc.Locator][]byte{
				{rsrc.UserInfo("abc"): nil},
				{rsrc.UserInfo("abc"): nil},
			},
			[]mock.Resolver{mock.Path, mock.Path},
			nil,
			rsrc.UserInfo("abc"),
			[][]byte{nil, nil},
			false,
		},
		{
			[]map[rsrc.Locator][]byte{
				{rsrc.APIKey(): []byte("9")},
				{},
			},
			[]mock.Resolver{mock.Path, mock.URL},
			[]byte("9"),
			rsrc.APIKey(),
			[][]byte{[]byte("9"), nil},
			true,
		},
		{
			[]map[rsrc.Locator][]byte{
				{rsrc.UserInfo("abc"): []byte("9")},
				{},
			},
			[]mock.Resolver{mock.Path, mock.Path},
			[]byte("9"),
			rsrc.UserInfo("abc"),
			[][]byte{[]byte("9"), nil},
			true,
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			var ios [][]rsrc.IO
			for i := range c.files {
				io, err := mock.IO(c.files[i], c.locf[i])
				if err != nil {
					t.Fatal("setup error")
				}
				ios = append(ios, []rsrc.IO{io})
			}

			s, err := New(ios)
			if err != nil {
				t.Fatal("unexpected error in constructor:", err)
			}

			data, err := s.Update(c.loc)
			if err != nil && c.ok {
				t.Fatal("unexpected error:", err)
			} else if err == nil && !c.ok {
				t.Fatal("expected error but non occurred")
			}

			if string(data) != string(c.data) {
				t.Errorf("read data is wrong:\nhas:      '%v'\nexpected: '%v'",
					string(data), string(c.data))
			}

			for i, io := range ios {
				data, err := io[0].Read(c.loc)

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
	cases := []struct {
		files   []map[rsrc.Locator][]byte
		locf    []mock.Resolver
		data    []byte
		loc     rsrc.Locator
		written [][]byte
	}{
		{ // failed write (critical)
			[]map[rsrc.Locator][]byte{{}},
			[]mock.Resolver{mock.Path},
			[]byte("xx"),
			rsrc.UserInfo("abc"),
			[][]byte{nil},
		},
		{ // not written in layer 0 (no URL for APIKey)
			[]map[rsrc.Locator][]byte{
				{},
				{rsrc.APIKey(): nil},
			},
			[]mock.Resolver{mock.URL, mock.Path},
			[]byte("xx"),
			rsrc.APIKey(),
			[][]byte{nil, []byte("xx")},
		},
		{ // written in neither (0 not reached since 1 fails)
			[]map[rsrc.Locator][]byte{
				{rsrc.APIKey(): nil},
				{},
			},
			[]mock.Resolver{mock.Path, mock.URL},
			[]byte("xx"),
			rsrc.APIKey(),
			[][]byte{[]byte("xx"), nil},
		},
		{ // written in both layers
			[]map[rsrc.Locator][]byte{
				{rsrc.UserInfo("abc"): nil},
				{rsrc.UserInfo("abc"): nil},
			},
			[]mock.Resolver{mock.URL, mock.Path},
			[]byte("xx"),
			rsrc.UserInfo("abc"),
			[][]byte{[]byte("xx"), []byte("xx")},
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			var ios [][]rsrc.IO
			for i := range c.files {
				io, err := mock.IO(c.files[i], c.locf[i])
				if err != nil {
					t.Fatal("setup error")
				}
				ios = append(ios, []rsrc.IO{io})
			}

			s, err := New(ios)
			if err != nil {
				t.Fatal("unexpected error in constructor:", err)
			}

			err = s.Write(c.data, c.loc)
			if err != nil {
				t.Fatal("unexpected error:", err)
			}

			for i, io := range ios {
				data, err := io[0].Read(c.loc)

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
	cases := []struct {
		files []map[rsrc.Locator][]byte
		locf  []mock.Resolver
		loc   rsrc.Locator
		exist []bool
	}{
		{ // failed remove (critical)
			[]map[rsrc.Locator][]byte{{}},
			[]mock.Resolver{mock.Path},
			rsrc.UserInfo("abc"),
			[]bool{false},
		},
		{ // remove both
			[]map[rsrc.Locator][]byte{
				{rsrc.UserInfo("abc"): []byte("xx")},
				{rsrc.UserInfo("abc"): []byte("xx")},
			},
			[]mock.Resolver{mock.URL, mock.Path},
			rsrc.UserInfo("abc"),
			[]bool{false, false},
		},
		{ // level 0 not removed since 1 failes
			[]map[rsrc.Locator][]byte{
				{rsrc.APIKey(): []byte("xx")},
				{},
			},
			[]mock.Resolver{mock.Path, mock.URL},
			rsrc.APIKey(),
			[]bool{true, false},
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			var ios [][]rsrc.IO
			for i := range c.files {
				io, err := mock.IO(c.files[i], c.locf[i])
				if err != nil {
					t.Fatal("setup error")
				}
				ios = append(ios, []rsrc.IO{io})
			}

			s, err := New(ios)
			if err != nil {
				t.Fatal("unexpected error in constructor:", err)
			}

			err = s.Remove(c.loc)
			if err != nil {
				t.Fatal("unexpected error:", err)
			}

			if err != nil {
				for i, io := range ios {
					_, err := io[0].Read(c.loc)

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

func TestStoreObserver(t *testing.T) {
	cases := []struct {
		files []map[rsrc.Locator][]byte
		locf  []mock.Resolver
		data  []byte
		loc   rsrc.Locator
		msgs  [][]format.Formatter
	}{
		{ // failed write (critical)
			[]map[rsrc.Locator][]byte{{}},
			[]mock.Resolver{mock.Path},
			[]byte("xx"),
			rsrc.UserInfo("abc"),
			[][]format.Formatter{{
				&format.Message{Msg: "r: 0/0, w: 0/1, rm: 0/0"},
				&format.Message{Msg: "r: 0/0, w: 1/1, rm: 0/0"},
			}},
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			var ios [][]rsrc.IO
			var ds []*mock.Display
			var fChans []chan format.Formatter
			var fChansOut []chan<- format.Formatter
			var quits []chan bool

			for i := range c.files {
				io, err := mock.IO(c.files[i], c.locf[i])
				if err != nil {
					t.Fatal("setup error")
				}
				ios = append(ios, []rsrc.IO{io})

				fChan := make(chan format.Formatter)
				quit := make(chan bool)
				d := mock.NewDisplay()
				go func(
					d display.Display,
					fChan chan format.Formatter,
					quit chan bool,
				) {
					for msg := range fChan {
						d.Display(msg)
					}
					quit <- true
				}(d, fChan, quit)

				ds = append(ds, d)
				fChans = append(fChans, fChan)
				fChansOut = append(fChansOut, fChan)
				quits = append(quits, quit)
			}
			s, err := NewObserved(ios, fChansOut)
			if err != nil {
				t.Error("unexpected error in constructor")
			}

			err = s.Write(c.data, c.loc)
			if err != nil {
				t.Error("unexpected error:", err)
			}

			for i := range ios {
				close(fChans[i])
				<-quits[i]

				if len(c.msgs[i]) != len(ds[i].Msgs) {
					t.Fatalf("expected %v messages but got %v", len(c.msgs[i]), len(ds[i].Msgs))
				}
				for j, expect := range c.msgs[i] {
					if !reflect.DeepEqual(expect, ds[i].Msgs[j]) {
						t.Errorf("expect %v, got %v", expect, ds[i].Msgs[j])
					}
				}
			}
		})
	}
}

func TestStoreWrongObserverCount(t *testing.T) {
	var ios [][]rsrc.IO
	io, err := mock.IO(map[rsrc.Locator][]byte{}, mock.Path)
	if err != nil {
		t.Fatal("setup error")
	}
	ios = append(ios, []rsrc.IO{io})

	_, err = NewObserved(ios, dumpChans(len(ios)+2))
	if err == nil {
		t.Fatal("expected error in constructor")
	}
}
