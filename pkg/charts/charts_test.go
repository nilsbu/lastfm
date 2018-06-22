package charts

import (
	"fmt"
	"testing"

	"github.com/nilsbu/fastest"
	"github.com/nilsbu/lastfm/pkg/unpack"
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

func TestChartsSum(t *testing.T) {
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

func TestChartsColumn(t *testing.T) {
	ft := fastest.T{T: t}

	testCases := []struct {
		charts Charts
		i      int
		column Column
		ok     bool
	}{
		{
			Charts{},
			0,
			Column{},
			false,
		},
		{
			Charts{"X": []int{}},
			0,
			Column{},
			false,
		},
		{
			Charts{
				"o0o": []int{0, 0, 7},
				"lol": []int{1, 2, 3},
				"X":   []int{1, 3, 4}},
			1,
			Column{Score{"X", 3}, Score{"lol", 2}, Score{"o0o", 0}},
			true,
		},
		{
			Charts{"X": []int{1, 3, 4}},
			-1,
			Column{Score{"X", 4}},
			true,
		},
		{
			Charts{"X": []int{1, 3, 4}},
			-4,
			Column{},
			false,
		},
	}

	for i, tc := range testCases {
		ft.Seq(fmt.Sprintf("#%v", i), func(ft fastest.T) {
			column, err := tc.charts.Column(tc.i)

			ft.Implies(err == nil, tc.ok)
			ft.Implies(err != nil, !tc.ok, err)
			ft.DeepEquals(column, tc.column)
		})
	}
}

func TestColumnTop(t *testing.T) {
	ft := fastest.T{T: t}

	testCases := []struct {
		column Column
		n      int
		top    Column
	}{
		{
			Column{},
			0,
			Column{},
		},
		{
			Column{Score{"X", 4}},
			0,
			Column{},
		},
		{
			Column{Score{"X", 3}, Score{"lol", 2}, Score{"o0o", 0}},
			2,
			Column{Score{"X", 3}, Score{"lol", 2}},
		},
		{
			Column{Score{"X", 3}, Score{"lol", 2}, Score{"o0o", 0}},
			4,
			Column{Score{"X", 3}, Score{"lol", 2}, Score{"o0o", 0}},
		},
	}

	for i, tc := range testCases {
		ft.Seq(fmt.Sprintf("#%v", i), func(ft fastest.T) {
			top := tc.column.Top(tc.n)

			ft.DeepEquals(top, tc.top)
		})
	}
}
