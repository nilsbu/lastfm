package unpack

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/nilsbu/fastest"
)

func TestUserRecentTracks(t *testing.T) {
	ft := fastest.T{T: t}

	testCases := []struct {
		json []byte
		urt  *UserRecentTracks
		err  fastest.Code
	}{
		{
			[]byte("{}"),
			&UserRecentTracks{},
			fastest.OK,
		},
		{
			[]byte(`{"recenttracks":{"track":[],"@attr":{"user":"U","page":"1","perPage":"200","totalPages":"0","total":"0"}}}`),
			&UserRecentTracks{
				RecentTracks: recentTracks{
					Track: make([]track, 0),
					Attr: recentTracksAttr{
						User: "U", Page: 1, PerPage: 200},
				},
			},
			fastest.OK,
		},
		{
			[]byte(`{"recenttracks":{
        "track":[{
          "artist":{"#text":"BTS","mbid":"ac6"},
          "name":"t3",
          "streamable":"0",
          "mbid":"",
          "album":{"#text":"轉","mbid":"2e"},
          "url":"https://example.com",
          "image":[{"#text":"s.png","size":"small"}],
          "date":{"uts":"1526811835","#text":"20 May 2018, 10:23"}}],
        "@attr":{"user":"ÄÖ","page":"1","perPage":"200","totalPages":"1","total":"1"}}}`),
			&UserRecentTracks{
				RecentTracks: recentTracks{
					Track: []track{track{
						Artist: text{"BTS"},
						Name:   "t3",
						Album:  text{"轉"},
						Date:   date{1526811835},
					}},
					Attr: recentTracksAttr{
						User: "ÄÖ", Page: 1, PerPage: 200, TotalPages: 1, Total: 1},
				},
			},
			fastest.OK,
		},
		{
			[]byte(`{"recenttracks":{
        "track":[{
          "artist":{"#text":"BTS","mbid":"ac6"},
          "name":"t3",
          "streamable":"0",
          "mbid":"",
          "album":{"#text":"轉","mbid":"2e"},
          "url":"https://example.com",
          "image":[{"#text":"s.png","size":"small"}],
          "@attr":{"nowplaying":"true"}}],
        "@attr":{"user":"ÄÖ","page":"1","perPage":"200","totalPages":"1","total":"1"}}}`),
			&UserRecentTracks{
				RecentTracks: recentTracks{
					Track: []track{track{
						Artist: text{"BTS"},
						Name:   "t3",
						Album:  text{"轉"},
						Attr:   trackAttr{NowPlaying: true},
					}},
					Attr: recentTracksAttr{
						User: "ÄÖ", Page: 1, PerPage: 200, TotalPages: 1, Total: 1},
				},
			},
			fastest.OK,
		},
	}

	for i, tc := range testCases {
		s := fmt.Sprintf("#%v", i)
		ft.Seq(s, func(ft fastest.T) {
			urt := &UserRecentTracks{}
			err := json.Unmarshal(tc.json, urt)

			ft.Implies(tc.err == fastest.Fail, err != nil)
			ft.Implies(tc.err == fastest.OK, err == nil, err)
			ft.Only(tc.err == fastest.OK)
			ft.DeepEquals(urt, tc.urt)
		})
	}
}

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
