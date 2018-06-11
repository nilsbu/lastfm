package io

import (
	"fmt"
	"testing"

	"github.com/nilsbu/fastest"
)

func TestFmtPath(t *testing.T) {
	ft := fastest.T{T: t}

	testCases := []struct {
		rsrc *Resource
		path string
	}{
		{NewArtistInfo("m2"), "data/rawdata/artist.getInfo/m2.json"},
		{NewUserInfo("con"), "data/rawdata/user.getInfo/_con.json"},
		{NewArtistInfo("Týr"), "data/rawdata/artist.getInfo/T%C3%BDr.json"},
		{NewArtistInfo("A B"), "data/rawdata/artist.getInfo/A+B.json"},
		{NewUserRecentTracks("X", 3, 86400), "data/rawdata/user.getRecentTracks/X.86400(3).json"},
	}

	for i, tc := range testCases {
		s := fmt.Sprintf("#%v: %v", i, tc.rsrc.name)
		ft.Seq(s, func(ft fastest.T) {
			path := FmtPath(tc.rsrc)
			ft.Equals(path, tc.path)
		})
	}
}
