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
		{NewArtistInfo("m2"), ".data/artist.getInfo/m2.json"},
		{NewUserInfo("con"), ".data/user.getInfo/_con.json"},
		{NewArtistInfo("Týr"), ".data/artist.getInfo/T%C3%BDr.json"},
		{NewArtistInfo("A B"), ".data/artist.getInfo/A+B.json"},
		{NewUserRecentTracks("X", 3, 86400), ".data/user.getRecentTracks/X.86400(3).json"},
	}

	for i, tc := range testCases {
		s := fmt.Sprintf("#%v: %v", i, tc.rsrc.name)
		ft.Seq(s, func(ft fastest.T) {
			path := fmtPath(tc.rsrc)
			ft.Equals(path, tc.path)
		})
	}
}
