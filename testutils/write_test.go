package testutils

import (
	"fmt"
	"testing"

	"github.com/nilsbu/fastest"
	"github.com/nilsbu/lastfm/io"
)

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
