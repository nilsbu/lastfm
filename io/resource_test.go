package io

import (
	"fmt"
	"testing"

	"github.com/nilsbu/fastest"
)

func TestNewUserInfo(t *testing.T) {
	ft := fastest.T{T: t}

	testCases := []struct {
		name string
	}{
		{"AB"},
		{"A@@"},
	}

	for i, tc := range testCases {
		s := fmt.Sprintf("#%v: %v", i, tc.name)
		ft.Seq(s, func(ft fastest.T) {
			rsrc := NewUserInfo(tc.name)

			ft.Equals(rsrc.main, "user")
			ft.Equals(rsrc.method, "getInfo")
			ft.Equals(rsrc.name, tc.name)
			ft.Equals(rsrc.page, Page(0))
			ft.Equals(len(rsrc.params), 0)
		})
	}
}

func TestNewUserRecentTracks(t *testing.T) {
	ft := fastest.T{T: t}

	testCases := []struct {
		name string
		page Page
		time Midnight
		from string
		to   string
	}{
		{"_X", 1, 86400, "86399", "172800"},
		{"Abra", 2, 8640000, "8639999", "8726400"},
	}

	for i, tc := range testCases {
		s := fmt.Sprintf("#%v: %v", i, tc.name)
		ft.Seq(s, func(ft fastest.T) {
			rsrc := NewUserRecentTracks(tc.name, tc.page, tc.time)

			ft.Equals(rsrc.main, "user")
			ft.Equals(rsrc.method, "getRecentTracks")
			ft.Equals(rsrc.name, tc.name)
			ft.Equals(rsrc.page, tc.page)
			ft.Equals(rsrc.params["from"], tc.from)
			ft.Equals(rsrc.params["to"], tc.to)
		})
	}
}

func TestNewArtistInfo(t *testing.T) {
	ft := fastest.T{T: t}

	testCases := []struct {
		name string
	}{
		{"xy"},
		{"Týr"},
	}

	for i, tc := range testCases {
		s := fmt.Sprintf("#%v: %v", i, tc.name)
		ft.Seq(s, func(ft fastest.T) {
			rsrc := NewArtistInfo(tc.name)

			ft.Equals(rsrc.main, "artist")
			ft.Equals(rsrc.method, "getInfo")
			ft.Equals(rsrc.name, tc.name)
			ft.Equals(rsrc.page, Page(0))
			ft.Equals(len(rsrc.params), 0)
		})
	}
}
