package charts

import (
	"math"
	"reflect"
	"sort"
	"testing"
)

func TestCompile(t *testing.T) {
	cases := []struct {
		days   []Charts
		charts Charts
	}{
		{
			[]Charts{},
			Charts{},
		},
		{
			[]Charts{Charts{}},
			Charts{},
		},
		{
			[]Charts{
				Charts{"ASD": []float64{2}},
				Charts{"WASD": []float64{1}},
				Charts{"ASD": []float64{13}, "WASD": []float64{4}},
			},
			Charts{"ASD": []float64{2, 0, 13}, "WASD": []float64{0, 1, 4}},
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			charts := Compile(c.days)

			if !reflect.DeepEqual(charts, c.charts) {
				t.Errorf("wrong data:\nhas:  %v\nwant: %v",
					charts, c.charts)
			}
		})
	}
}

func TestChartsSum(t *testing.T) {
	cases := []struct {
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

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			sums := c.charts.Sum()

			if !reflect.DeepEqual(sums, c.sums) {
				t.Errorf("wrong data:\nhas:  %v\nwant: %v",
					sums, c.sums)
			}
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

	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			column, err := tc.charts.Column(tc.i)
			if err != nil && tc.ok {
				t.Error("unexpected error:", err)
			} else if err == nil && !tc.ok {
				t.Error("expected error but none occurred")
			}

			if err == nil {
				if !reflect.DeepEqual(column, tc.column) {
					t.Errorf("wrong data:\nhas:  %v\nwant: %v",
						column, tc.column)
				}
			}
		})
	}
}

func TestColumnTop(t *testing.T) {
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

	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			top := tc.column.Top(tc.n)

			if !reflect.DeepEqual(top, tc.top) {
				t.Errorf("wrong data:\nhas:  %v\nwant: %v", top, tc.top)
			}
		})
	}
}

func TestColumnSum(t *testing.T) {
	cases := []struct {
		col Column
		sum float64
	}{
		{Column{}, 0},
		{Column{{"a", 10}, {"b", 2.5}}, 12.5},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			sum := c.col.Sum()
			if sum != c.sum {
				t.Errorf("got %v, expected %v", sum, c.sum)
			}
		})
	}
}

func TestChartsKeys(t *testing.T) {
	cases := []struct {
		charts Charts
		keys   []string
	}{
		{
			Charts{},
			[]string{},
		},
		{
			Charts{
				"xx": []float64{32, 45},
				"yy": []float64{32, 45},
			},
			[]string{"xx", "yy"},
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			keys := c.charts.Keys()

			sort.Strings(keys)
			sort.Strings(c.keys)
			if !reflect.DeepEqual(keys, c.keys) {
				t.Errorf("wrong data (sorted):\nhas:  %v\nwant: %v",
					keys, c.keys)
			}
		})
	}
}

func TestChartsCorrect(t *testing.T) {
	cases := []struct {
		charts     Charts
		correction map[string]string
		corrected  Charts
	}{
		{
			Charts{
				"o0o": []float64{0, 0, 7},
				"lol": []float64{1, 2, 3},
				"X":   []float64{1, 3, 4}},
			map[string]string{"X": "o0o"},
			Charts{
				"o0o": []float64{1, 3, 11},
				"lol": []float64{1, 2, 3},
			},
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			corrected := c.charts.Correct(c.correction)

			if !reflect.DeepEqual(corrected, c.corrected) {
				t.Errorf("wrong data:\nhas:  %v\nwant: %v",
					corrected, c.corrected)
			}
		})
	}
}
