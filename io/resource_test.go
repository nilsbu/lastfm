package io

import (
	"fmt"
	"testing"

	"github.com/nilsbu/fastest"
)

func TestNewUserInfo(t *testing.T) {
	ft := fastest.T{T: t}

	testCases := []struct {
		name Name
	}{
		{"AB"},
		{"A@@"},
	}

	for i, tc := range testCases {
		s := fmt.Sprintf("#%v: %v", i, tc.name)
		ft.Seq(s, func(ft fastest.T) {
			rsrc := NewUserInfo(tc.name)

			ft.Equals(rsrc.domain, Raw)
			ft.Equals(rsrc.main, "user")
			ft.Equals(rsrc.method, "getInfo")
			ft.Equals(rsrc.name, tc.name)
			ft.Equals(rsrc.page, Page(0))
			ft.Equals(rsrc.time, Midnight(-1))
		})
	}
}

func TestNewUserRecentTracks(t *testing.T) {
	ft := fastest.T{T: t}

	testCases := []struct {
		name Name
		page Page
		time Midnight
	}{
		{"_X", 1, 86400},
		{"Abra", 2, 8640000},
	}

	for i, tc := range testCases {
		s := fmt.Sprintf("#%v: %v", i, tc.name)
		ft.Seq(s, func(ft fastest.T) {
			rsrc := NewUserRecentTracks(tc.name, tc.page, tc.time)

			ft.Equals(rsrc.domain, Raw)
			ft.Equals(rsrc.main, "user")
			ft.Equals(rsrc.method, "getRecentTracks")
			ft.Equals(rsrc.name, tc.name)
			ft.Equals(rsrc.page, tc.page)
			ft.Equals(rsrc.time, tc.time)
		})
	}
}

func TestNewArtistInfo(t *testing.T) {
	ft := fastest.T{T: t}

	testCases := []struct {
		name Name
	}{
		{"xy"},
		{"Týr"},
	}

	for i, tc := range testCases {
		s := fmt.Sprintf("#%v: %v", i, tc.name)
		ft.Seq(s, func(ft fastest.T) {
			rsrc := NewArtistInfo(tc.name)

			ft.Equals(rsrc.domain, Raw)
			ft.Equals(rsrc.main, "artist")
			ft.Equals(rsrc.method, "getInfo")
			ft.Equals(rsrc.name, tc.name)
			ft.Equals(rsrc.page, Page(0))
			ft.Equals(rsrc.time, Midnight(-1))
		})
	}
}
