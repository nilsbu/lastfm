package io

import (
	"fmt"
	"testing"

	"github.com/nilsbu/fastest"
)

func TestGetPages(t *testing.T) {
	ft := fastest.T{T: t}

	const (
		ok int = iota
		fail
	)

	testCases := []struct {
		json  string
		pages int
		err   int
	}{
		{"{\"recenttracks\":{\"track\":[],\"@attr\":{\"user\":\"\",\"page\":\"1\",\"totalPages\":\"2\"}}}", 2, ok},
		{"asd", 0, fail},

		// TODO escape non-ASCII characters
	}

	for i, tc := range testCases {
		s := fmt.Sprintf("#%v", i)
		ft.Seq(s, func(ft fastest.T) {
			pages, err := getPages(tc.json)

			ft.Implies(tc.err == fail, err != nil)
			ft.Only(tc.err == ok)
			ft.Equals(pages, tc.pages)
		})
	}
}
