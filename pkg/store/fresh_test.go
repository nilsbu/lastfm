package store

import (
	"testing"

	"github.com/nilsbu/lastfm/pkg/rsrc"
	"github.com/nilsbu/lastfm/test/mock"
)

func TestFresh(t *testing.T) {
	l0, err := mock.IO(map[rsrc.Locator][]byte{
		rsrc.APIKey(): []byte("new"),
	}, mock.Path)
	if err != nil {
		t.Fatal("setup error:", err)
	}

	l1, err := mock.IO(map[rsrc.Locator][]byte{
		rsrc.APIKey(): []byte("old"),
	}, mock.Path)
	if err != nil {
		t.Fatal("setup error:", err)
	}

	store, err := New([][]rsrc.IO{{l0}, {l1}})
	if err != nil {
		t.Fatal("unexpected ctor error:", err)
	}

	fresh := Fresh(store)

	data, err := fresh.Read(rsrc.APIKey())
	if err != nil {
		t.Error("unexpected read error:", err)
	} else if string(data) != "new" {
		t.Errorf("read: wrong data: has '%v', want '%v'", string(data), "new")
	}

	// ensure l1 is overwritten
	data, err = l1.Read(rsrc.APIKey())
	if err != nil {
		t.Fatal("unexpected read error:", err)
	} else if string(data) != "new" {
		t.Errorf("l1: wrong data: has '%v', want '%v'", string(data), "new")
	}

	if err = fresh.Remove(rsrc.APIKey()); err != nil {
		t.Error("expected error in remove but none occurred")
	}

	// both layers must be removed
	if _, err = l0.Read(rsrc.APIKey()); err == nil {
		t.Error("expected error in l0 but none occurred")
	}
	if _, err = l1.Read(rsrc.APIKey()); err == nil {
		t.Error("expected error in l1 but none occurred")
	}

	if err = fresh.Write([]byte("written"), rsrc.APIKey()); err != nil {
		t.Error("expected error in write but none occurred")
	}

	// ensure both layers are written
	data, err = l0.Read(rsrc.APIKey())
	if err != nil {
		t.Fatal("unexpected write error:", err)
	} else if string(data) != "written" {
		t.Errorf("l0: wrong data: has '%v', want '%v'", string(data), "written")
	}
	data, err = l1.Read(rsrc.APIKey())
	if err != nil {
		t.Fatal("unexpected write error:", err)
	} else if string(data) != "written" {
		t.Errorf("l1: wrong data: has '%v', want '%v'", string(data), "written")
	}
}
