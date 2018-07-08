package organize

import (
	"reflect"
	"testing"

	"github.com/nilsbu/lastfm/pkg/charts"
	"github.com/nilsbu/lastfm/pkg/rsrc"
	"github.com/nilsbu/lastfm/test/mock"
	"github.com/pkg/errors"
)

func TestMultiError(t *testing.T) {
	err := &multiError{
		"message",
		[]error{errors.New("error 1"), errors.New("error 2")}}

	msg := err.Error()
	str := "message:\n  error 1\n  error 2"
	if msg != str {
		t.Errorf("wrong message:\nhas:\n%v\nwant:\n%v", msg, str)
	}
}

func TestLoadArtistTags(t *testing.T) {
	cases := []struct {
		files   map[rsrc.Locator][]byte
		artists []string
		tags    map[string][]charts.Tag
		ok      bool
	}{
		{
			map[rsrc.Locator][]byte{
				rsrc.ArtistTags("asdf"): nil,
			},
			[]string{"asdf"}, nil, false,
		},
		{
			map[rsrc.Locator][]byte{
				rsrc.ArtistTags("asdf"): []byte(`{"toptags":{"tag":[{"name":"t0","count":100}]}}`),
				rsrc.ArtistTags("basd"): []byte(`{"toptags":{"tag":[{"name":"t0","count":20}]}}`),
				rsrc.TagInfo("t0"):      []byte(`{"tag":{"name":"t0","total":1024,"reach":42}}`),
			},
			[]string{"asdf", "basd"},
			map[string][]charts.Tag{
				"asdf": []charts.Tag{charts.Tag{Name: "t0", Total: 1024, Reach: 42, Weight: 100}},
				"basd": []charts.Tag{charts.Tag{Name: "t0", Total: 1024, Reach: 42, Weight: 20}},
			},
			true,
		},
		{
			map[rsrc.Locator][]byte{
				rsrc.ArtistTags("asdf"): []byte(`{"toptags":{"tag":[{"name":"t0","count":100}]}}`),
				rsrc.TagInfo("t0"):      nil,
			},
			[]string{"asdf"}, nil, false,
		},
		{
			map[rsrc.Locator][]byte{
				rsrc.ArtistTags("asdf"): []byte(`{"toptags":{"tag":[{"name":"UPPER","count":100}]}}`),
				rsrc.TagInfo("UPPER"):   []byte(`{"tag":{"name":"UPPER","total":1024,"reach":42}}`),
			},
			[]string{"asdf"},
			map[string][]charts.Tag{
				"asdf": []charts.Tag{charts.Tag{Name: "upper", Total: 1024, Reach: 42, Weight: 100}},
			},
			true,
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			io, err := mock.IO(c.files, mock.Path)
			if err != nil {
				t.Fatal("setup error")
			}

			tags, err := LoadArtistTags(c.artists, io)
			if err != nil && c.ok {
				t.Error("unexpected error:", err)
			} else if err == nil && !c.ok {
				t.Error("expected error but none occurred")
			}

			if err == nil {
				if !reflect.DeepEqual(tags, c.tags) {
					t.Errorf("wrong data:\nhas:  %v\nwant: %v",
						tags, c.tags)
				}
			}
		})
	}
}
