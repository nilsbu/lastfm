package unpack

import (
	"fmt"
	"testing"

	"github.com/nilsbu/fastest"
)

func TestUnmarshalUserRecentTracks(t *testing.T) {
	ft := fastest.T{T: t}

	// TODO use fastest
	const (
		ok int = iota
		fail
	)

	testCases := []struct {
		json []byte
		urt  *UserRecentTracks
		err  int
	}{
		{
			[]byte("{}"),
			&UserRecentTracks{},
			ok,
		},
		{
			[]byte(`{"recenttracks":{"track":[],"@attr":{"user":"U","page":"1","perPage":"200","totalPages":"0","total":"0"}}}`),
			&UserRecentTracks{
				RecentTracks: recentTracks{
					Track: make([]Track, 0),
					Attr: recentTracksAttr{
						User: "U", Page: 1, PerPage: 200},
				},
			},
			ok,
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
					Track: []Track{Track{
						Artist: Text{"BTS"},
						Name:   "t3",
						Album:  Text{"轉"},
						Date:   Date{1526811835},
					}},
					Attr: recentTracksAttr{
						User: "ÄÖ", Page: 1, PerPage: 200, TotalPages: 1, Total: 1},
				},
			},
			ok,
		},
	}

	for i, tc := range testCases {
		s := fmt.Sprintf("#%v", i)
		ft.Seq(s, func(ft fastest.T) {
			urt, err := UnmarshalUserRecentTracks(tc.json)

			ft.Implies(tc.err == fail, err != nil)
			ft.Only(tc.err == ok)
			ft.DeepEquals(urt, tc.urt)
		})
	}
}

func TestUnmarshalAPIKey(t *testing.T) {
	ft := fastest.T{T: t}

	const (
		ok int = iota
		fail
	)

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
			key, err := UnmarshalAPIKey(tc.json)

			ft.Equals(tc.err == fastest.Fail, err != nil)
			ft.DeepEquals(key, tc.key)
		})
	}
}
