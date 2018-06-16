package charts

import (
	"fmt"
	"testing"

	"github.com/nilsbu/fastest"
	"github.com/nilsbu/lastfm/unpack"
)

func TestCompile(t *testing.T) {
	ft := fastest.T{T: t}

	testCases := []struct {
		dps    []unpack.DayPlays
		charts Charts
	}{
		{
			[]unpack.DayPlays{},
			Charts{},
		},
		{
			[]unpack.DayPlays{unpack.DayPlays{}},
			Charts{},
		},
		{
			[]unpack.DayPlays{
				unpack.DayPlays{"ASD": 2},
				unpack.DayPlays{"WASD": 1},
				unpack.DayPlays{"ASD": 13, "WASD": 4},
			},
			Charts{"ASD": []int{2, 0, 13}, "WASD": []int{0, 1, 4}},
		},
	}

	for i, tc := range testCases {
		ft.Seq(fmt.Sprintf("#%v", i), func(ft fastest.T) {
			charts := Compile(tc.dps)

			ft.DeepEquals(charts, tc.charts)
		})
	}
}

func TestSum(t *testing.T) {
	ft := fastest.T{T: t}

	testCases := []struct {
		charts Charts
		sums   Sums
	}{
		{
			Charts{},
			Sums{},
		},
		{
			Charts{"X": []int{}},
			Sums{"X": []int{}},
		},
		{
			Charts{"X": []int{1, 3, 4}, "o0o": []int{0, 0, 7}},
			Sums{"X": []int{1, 4, 8}, "o0o": []int{0, 0, 7}},
		},
	}

	for i, tc := range testCases {
		ft.Seq(fmt.Sprintf("#%v", i), func(ft fastest.T) {
			sums := tc.charts.Sum()

			ft.DeepEquals(sums, tc.sums)
		})
	}
}
