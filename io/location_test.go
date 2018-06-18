package io

import (
	"fmt"
	"testing"

	"github.com/nilsbu/fastest"
)

func TestFmtURL(t *testing.T) {
	ft := fastest.T{T: t}

	base := "http://ws.audioscrobbler.com/2.0/?format=json&"

	testCases := []struct {
		apiKey APIKey
		rsrc   *Resource
		url    string
	}{
		{"a", NewUserInfo("Týr"), base + "api_key=a&method=user.getInfo&user=T%C3%BDr"},
		{"b", NewArtistInfo("m2"), base + "api_key=b&method=artist.getInfo&artist=m2"},
		{"b", NewUserRecentTracks("X", 3, 86400), base + "api_key=b&method=user.getRecentTracks&user=X&page=3&from=86399&to=172800&limit=200"},
	}

	for i, tc := range testCases {
		s := fmt.Sprintf("#%v: %v", i, tc.rsrc.name)
		ft.Seq(s, func(ft fastest.T) {
			url := fmtURL(tc.rsrc, tc.apiKey)
			ft.Equals(url, tc.url)
		})
	}
}

func TestFmtPath(t *testing.T) {
	ft := fastest.T{T: t}

	testCases := []struct {
		rsrc *Resource
		path string
	}{
		{NewArtistInfo("m2"), ".lastfm/data/artist.getInfo/m2.json"},
		{NewUserInfo("con"), ".lastfm/data/user.getInfo/_con.json"},
		{NewArtistInfo("Týr"), ".lastfm/data/artist.getInfo/T%C3%BDr.json"},
		{NewArtistInfo("A B"), ".lastfm/data/artist.getInfo/A+B.json"},
		{NewUserRecentTracks("X", 3, 86400), ".lastfm/data/user.getRecentTracks/X.86400(3).json"},
		{NewAPIKey(), ".lastfm/util/apikey.json"},
		{NewAllDayPlays("XX"), ".lastfm/user/XX/alldayplays.json"},
		{NewBookmark("zY"), ".lastfm/user/zY/bookmark.json"},
	}

	for i, tc := range testCases {
		s := fmt.Sprintf("#%v: %v", i, tc.rsrc.name)
		ft.Seq(s, func(ft fastest.T) {
			path := fmtPath(tc.rsrc)
			ft.Equals(path, tc.path)
		})
	}
}
