package unpack_test

import (
	"reflect"
	"testing"

	"github.com/nilsbu/lastfm/pkg/rsrc"
	"github.com/nilsbu/lastfm/pkg/unpack"
	"github.com/nilsbu/lastfm/test/mock"
)

func TestLoadAPIKey(t *testing.T) {
	cases := []struct {
		json []byte
		key  string
		ok   bool
	}{
		{[]byte(""), "", false},
		{[]byte(`{`), "", false},
		{[]byte(`{}`), "", false},
		{[]byte(`{"apikey":"0000a2e30000dc430000a2e30000dc43"}`),
			"0000a2e30000dc430000a2e30000dc43", true},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			io, err := mock.IO(
				map[rsrc.Locator][]byte{rsrc.APIKey(): c.json}, mock.Path)
			if err != nil {
				t.Fatal("setup error")
			}

			key, err := unpack.LoadAPIKey(io)
			if err != nil && c.ok {
				t.Error("unexpected error:", err)
			} else if err == nil && !c.ok {
				t.Error("expected error")
			}

			if err == nil {
				if key != c.key {
					t.Errorf("wrong key\nhas:  '%v'\nwant: '%v'", key, c.key)
				}
			}
		})
	}
}

func TestLoadSessionInfo(t *testing.T) {
	cases := []struct {
		json    []byte
		session *unpack.SessionInfo
		ok      bool
	}{
		{[]byte(""), nil, false},
		{[]byte(`{`), nil, false},
		{[]byte(`{}`), nil, false},
		{[]byte(`{"user":"somename"}`), &unpack.SessionInfo{User: "somename", Options: map[string]string{}}, true},
		{
			[]byte(`{"user":"somename","options":[{"name":"k","value":"v"},{"name":"k2","value":"v2"}]}`),
			&unpack.SessionInfo{User: "somename", Options: map[string]string{"k": "v", "k2": "v2"}}, true,
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			io, err := mock.IO(
				map[rsrc.Locator][]byte{rsrc.SessionInfo(): c.json}, mock.Path)
			if err != nil {
				t.Fatal("setup error")
			}

			session, err := unpack.LoadSessionInfo(io)
			if err != nil && c.ok {
				t.Error("unexpected error:", err)
			} else if err == nil && !c.ok {
				t.Error("expected error")
			}

			if err == nil {
				if !reflect.DeepEqual(session, c.session) {
					t.Errorf("wrong data\nhas:  '%v'\nwant: '%v'", session, c.session)
				}
			}
		})
	}
}

func TestWriteSessionInfo(t *testing.T) {
	cases := []struct {
		json    []byte
		session *unpack.SessionInfo
		ok      bool
	}{
		{[]byte(`{"user":"somename","options":[]}`), &unpack.SessionInfo{User: "somename"}, true},
		{
			[]byte(`{"user":"somename","options":[{"name":"k","value":"v"},{"name":"k2","value":"v2"}]}`),
			&unpack.SessionInfo{User: "somename", Options: map[string]string{"k": "v", "k2": "v2"}}, true,
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			io, err := mock.IO(
				map[rsrc.Locator][]byte{rsrc.SessionInfo(): nil}, mock.Path)
			if err != nil {
				t.Fatal("setup error")
			}

			err = unpack.WriteSessionInfo(c.session, io)
			if err != nil && c.ok {
				t.Error("unexpected error:", err)
			} else if err == nil && !c.ok {
				t.Error("expected error")
			}

			if err == nil {
				json, err := io.Read(rsrc.SessionInfo())
				if err != nil {
					t.Fatal("file was not written")
				}

				if string(json) != string(c.json) {
					t.Errorf("wrong data\nhas:  '%v'\nwant: '%v'",
						string(json), string(c.json))
				}
			}
		})
	}
}
