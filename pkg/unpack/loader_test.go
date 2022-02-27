package unpack

import (
	"testing"

	"github.com/nilsbu/lastfm/pkg/rsrc"
	"github.com/nilsbu/lastfm/test/mock"
)

// caches are mostly tested in combination with tag infos (see there)

func TestCachedLoaderShutdownOnError(t *testing.T) {
	io, err := mock.IO(map[rsrc.Locator][]byte{
		rsrc.TagInfo("error"):   []byte(`{"error":29,"message":"Rate Limit Exceeded"}`),
		rsrc.TagInfo("african"): []byte(`{"tag":{"name":"african","total":55266,"reach":10493}}`),
	}, mock.Path)
	if err != nil {
		t.Fatal("setup error")
	}

	buf := NewCached(io)

	_, err = LoadTagInfo("error", buf)
	if err == nil {
		t.Fatal("expected error but none occurred for tag 'error'")
	}

	_, err = LoadTagInfo("african", buf)
	if err == nil {
		t.Fatal("expected error but none occurred for tag 'african'")
	}
}
