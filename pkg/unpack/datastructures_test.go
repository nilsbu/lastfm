package unpack

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/nilsbu/fastest"
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

func TestAPIKey(t *testing.T) {
	ft := fastest.T{T: t}

	testCases := []struct {
		json []byte
		key  *APIKey
		err  fastest.Code
	}{
		{
			[]byte(`{"apikey":"asdf97"}`),
			&APIKey{"asdf97"},
			fastest.OK,
		},
	}

	for i, tc := range testCases {
		s := fmt.Sprintf("#%v", i)
		ft.Seq(s, func(ft fastest.T) {
			key := &APIKey{}
			err := json.Unmarshal(tc.json, key)

			ft.Equals(tc.err == fastest.Fail, err != nil)
			ft.DeepEquals(key, tc.key)
		})
	}
}

func TestSessionID(t *testing.T) {
	ft := fastest.T{T: t}

	testCases := []struct {
		json []byte
		sid  *SessionID
		err  fastest.Code
	}{
		{
			[]byte(`{"user":"somename"}`),
			&SessionID{"somename"},
			fastest.OK,
		},
	}

	for i, tc := range testCases {
		s := fmt.Sprintf("#%v", i)
		ft.Seq(s, func(ft fastest.T) {
			sid := &SessionID{}
			err := json.Unmarshal(tc.json, sid)

			ft.Equals(tc.err == fastest.Fail, err != nil)
			ft.DeepEquals(sid, tc.sid)
		})
	}
}
