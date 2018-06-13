package testutils

import (
	"fmt"
	"testing"

	"github.com/nilsbu/fastest"
	"github.com/nilsbu/lastfm/io"
)

func TestError(t *testing.T) {
	ft := fastest.T{T: t}

	testCases := []struct {
		s string
	}{{""}, {"ß"}}

	for _, tc := range testCases {
		ft.Seq(tc.s, func(ft fastest.T) {
			ft.Equals(strerr(tc.s).Error(), tc.s)
		})
	}
}

// TestReader tests both NewReader() and Read()
func TestReader(t *testing.T) {
	ft := fastest.T{T: t}

	testCases := []struct {
		rsrc *io.Resource
		data []byte
		err  fastest.Code
	}{
		{io.NewUserInfo("AS"), []byte("xx"), fastest.OK},
		// Twice to ensure that requests are repeatable
		{io.NewUserInfo("AS"), []byte("xx"), fastest.OK},
		{io.NewArtistInfo("AS"), nil, fastest.Fail},
		{io.NewUserRecentTracks("D", 1, 86400), []byte("a"), fastest.OK},
		{io.NewUserRecentTracks("D", 2, 0), []byte("c"), fastest.OK},
		{io.NewUserRecentTracks("D", 1, 0), []byte("b"), fastest.OK},
		{io.NewUserRecentTracks("D", 3, 0), nil, fastest.Fail},
	}

	r := NewReader(map[io.Resource][]byte{
		*io.NewUserInfo("AS"):                  []byte("xx"),
		*io.NewUserRecentTracks("D", 1, 86400): []byte("a"),
		*io.NewUserRecentTracks("D", 1, 0):     []byte("b"),
		*io.NewUserRecentTracks("D", 2, 0):     []byte("c"),
	})

	for i, tc := range testCases {
		ft.Seq(fmt.Sprintf("#%d", i), func(ft fastest.T) {
			data, err := r.Read(tc.rsrc)
			ft.Equals(tc.err == fastest.Fail, err != nil)
			ft.Equals(string(tc.data), string(data))
		})
	}
}

// TestWriter tests both NewWriter() and Write()
func TestWriter(t *testing.T) {
	ft := fastest.T{T: t}

	testCases := []struct {
		rsrc *io.Resource
		// in and stored should be []byte but used string for convenience
		in     string
		stored []string
		err    fastest.Code
	}{
		{io.NewUserInfo("AS"), "xx", []string{"xx"}, fastest.OK},
		{io.NewUserInfo("X"), "++", []string{"++"}, fastest.Fail},
		{io.NewUserInfo("AS"), "xy", []string{"xx", "xy"}, fastest.OK},
		{io.NewUserInfo("AS"), "--", []string{"xx", "xy", "--"}, fastest.Fail},
		{io.NewUserInfo("AS"), "++", []string{"xx", "xy", "--", "++"}, fastest.OK},
		{io.NewUserRecentTracks("D", 1, 86400), "", []string{""}, fastest.OK},
	}

	w := NewWriter(map[io.Resource][]bool{
		*io.NewUserInfo("AS"): []bool{true, true, false},
		*io.NewUserInfo("X"):  []bool{false},
	})

	for i, tc := range testCases {
		ft.Seq(fmt.Sprintf("#%d", i), func(ft fastest.T) {
			err := w.Write([]byte(tc.in), tc.rsrc)
			data, ok := w.data[*tc.rsrc]
			ft.Implies(tc.err == fastest.OK, ok)
			ft.Equals(tc.err == fastest.Fail, err != nil)
			ft.Only(ok)

			ft.Equals(len(tc.stored), len(data))
			for j := range tc.stored {
				ft.Equals(tc.stored[j], string(data[j]),
					fmt.Sprintf("Failed at element %v", j))
			}
		})
	}
}
