package unpack

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/nilsbu/fastest"
)

func TestUserInfo(t *testing.T) {
	ft := fastest.T{T: t}

	testCases := []struct {
		json []byte
		ui   *UserInfo
		err  fastest.Code
	}{
		{
			[]byte(`{"user":{"name":"What","playcount":"1928","registered":{"unixtime":"1144225884"}}}`),
			&UserInfo{
				User: userUser{Name: "What", PlayCount: 1928, Registered: time{1144225884}},
			},
			fastest.OK,
		},
	}

	for i, tc := range testCases {
		s := fmt.Sprintf("#%v", i)
		ft.Seq(s, func(ft fastest.T) {
			ui := &UserInfo{}
			err := json.Unmarshal(tc.json, ui)

			ft.Implies(tc.err == fastest.Fail, err != nil)
			ft.Only(tc.err == fastest.OK)
			ft.DeepEquals(ui, tc.ui)
		})
	}
}

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
	}

	for i, tc := range testCases {
		s := fmt.Sprintf("#%v", i)
		ft.Seq(s, func(ft fastest.T) {
			urt := &UserRecentTracks{}
			err := json.Unmarshal(tc.json, urt)

			ft.Implies(tc.err == fastest.Fail, err != nil)
			ft.Only(tc.err == fastest.OK)
			ft.DeepEquals(urt, tc.urt)
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

func TestBookmark(t *testing.T) {
	// Test is redundant, see organize.TestBookmark
	ft := fastest.T{T: t}

	testCases := []struct {
		json     []byte
		bookmark *Bookmark
	}{
		{
			[]byte(`{"time":"2128-06-11 08:53:20 +0000 UTC","unixtime":"5000000000"}`),
			&Bookmark{5000000000, "2128-06-11 08:53:20 +0000 UTC"},
		},
	}

	for i, tc := range testCases {
		s := fmt.Sprintf("#%v", i)
		ft.Seq(s, func(ft fastest.T) {
			bookmark := &Bookmark{}
			err := json.Unmarshal(tc.json, bookmark)

			ft.Nil(err, err)
			ft.DeepEquals(bookmark, tc.bookmark)
		})
	}
}
