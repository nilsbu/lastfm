package unpack

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestArtistInfo(t *testing.T) {
	cases := []struct {
		json []byte
		ai   *ArtistInfo
		ok   bool
	}{
		{
			[]byte(`{"artist":{"name":"ギルガメッシュ","stats":{"listeners":"229"}}`),
			nil, false,
		},
		{
			[]byte(`{"artist":{"name":"ギルガメッシュ","stats":{"listeners":"229","playcount":"999"}}}`),
			&ArtistInfo{
				Artist: artistArtist{Name: "ギルガメッシュ", Stats: artistStats{Listeners: 229, PlayCount: 999}},
			},
			true,
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			artist := &ArtistInfo{}
			err := json.Unmarshal(c.json, artist)

			if err != nil && c.ok {
				t.Error("unexpected error:", err)
			} else if err == nil && !c.ok {
				t.Error("no error occurred but should have")
			}

			if err == nil {
				if !reflect.DeepEqual(c.ai, artist) {
					t.Errorf("false result\nhas:  %v\nwant: %v", artist, c.ai)
				}
			}
		})
	}
}

func TestArtistTags(t *testing.T) {
	cases := []struct {
		json []byte
		at   *ArtistTags
		ok   bool
	}{
		{
			[]byte(`{"toptags":{"tag":[{"count":100,"name":"xyz"},{"count":77,"name":"J-A"}]}}`),
			&ArtistTags{
				TopTags: topTags{[]tag{tag{"xyz", 100}, tag{"J-A", 77}}},
			},
			true,
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			tags := &ArtistTags{}
			err := json.Unmarshal(c.json, tags)

			if err != nil && c.ok {
				t.Error("unexpected error:", err)
			} else if err == nil && !c.ok {
				t.Error("no error occurred but should have")
			}

			if err == nil {
				if !reflect.DeepEqual(c.at, tags) {
					t.Errorf("false result\nhas:  %v\nwant: %v", tags, c.at)
				}
			}
		})
	}
}
