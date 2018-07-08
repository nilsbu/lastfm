package charts

import (
	"fmt"
	"math"
	"testing"

	"github.com/nilsbu/fastest"
	"github.com/nilsbu/lastfm/pkg/unpack"
)

func TestCompile(t *testing.T) {
	ft := fastest.T{T: t}

	testCases := []struct {
		dps    []unpack.PlayCount
		charts Charts
	}{
		{
			[]unpack.PlayCount{},
			Charts{},
		},
		{
			[]unpack.PlayCount{unpack.PlayCount{}},
			Charts{},
		},
		{
			[]unpack.PlayCount{
				unpack.PlayCount{"ASD": 2},
				unpack.PlayCount{"WASD": 1},
				unpack.PlayCount{"ASD": 13, "WASD": 4},
			},
			Charts{"ASD": []float64{2, 0, 13}, "WASD": []float64{0, 1, 4}},
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
		sums   Charts
	}{
		{
			Charts{},
			Charts{},
		},
		{
			Charts{"X": []float64{}},
			Charts{"X": []float64{}},
		},
		{
			Charts{"X": []float64{1, 3, 4}, "o0o": []float64{0, 0, 7}},
			Charts{"X": []float64{1, 4, 8}, "o0o": []float64{0, 0, 7}},
		},
	}

	for i, tc := range testCases {
		ft.Seq(fmt.Sprintf("#%v", i), func(ft fastest.T) {
			sums := tc.charts.Sum()

			ft.DeepEquals(sums, tc.sums)
		})
	}
}

func TestChartsFade(t *testing.T) {
	cases := []struct {
		halflife float64
		charts   []float64
		faded    []float64
	}{
		{
			1.0,
			[]float64{1, 0, 0},
			[]float64{1, 0.5, 0.25},
		},
		{
			2.0,
			[]float64{1, 0, 1},
			[]float64{1, math.Sqrt(0.5), 1.5},
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			faded := Charts{"XX": c.charts}.Fade(c.halflife)
			f := faded["XX"]
			if len(f) != len(c.faded) {
				t.Fatalf("line length false: %v != %v", len(f), len(c.faded))
			}
			for i := 0; i < len(f); i++ {
				if math.Abs(f[i]-c.faded[i]) > 1e-6 {
					t.Errorf("at position %v: %v != %v", i, f[i], c.faded[i])
				}
			}
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
			Charts{"X": []float64{}},
			0,
			Column{},
			false,
		},
		{
			Charts{
				"o0o": []float64{0, 0, 7},
				"lol": []float64{1, 2, 3},
				"X":   []float64{1, 3, 4}},
			1,
			Column{Score{"X", 3}, Score{"lol", 2}, Score{"o0o", 0}},
			true,
		},
		{
			Charts{"X": []float64{1, 3, 4}},
			-1,
			Column{Score{"X", 4}},
			true,
		},
		{
			Charts{"X": []float64{1, 3, 4}},
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
