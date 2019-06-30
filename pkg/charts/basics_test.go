package charts

import (
	"math"
	"reflect"
	"testing"
)

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

func TestChartsRank(t *testing.T) {
	cases := []struct {
		charts Charts
		ranks  Charts
	}{
		{
			Charts{
				"o0o": []float64{0, 0, 7},
				"lol": []float64{1, 2, 3},
				"X":   []float64{1, 3, 4}},
			Charts{
				"o0o": []float64{3, 3, 1},
				"lol": []float64{1, 2, 3},
				"X":   []float64{1, 1, 2}},
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			ranks := c.charts.Rank()

			if !reflect.DeepEqual(ranks, c.ranks) {
				t.Errorf("wrong data:\nhas:  %v\nwant: %v",
					ranks, c.ranks)
			}
		})
	}
}

func TestChartsTotal(t *testing.T) {
	cases := []struct {
		charts Charts
		total  []float64
	}{
		{
			Charts{},
			[]float64{},
		},
		{
			Charts{"o0o": []float64{0, 0, 7}},
			[]float64{0, 0, 7},
		},
		{
			Charts{
				"o0o": []float64{0, 0, 7},
				"lol": []float64{1, 2, 3}},
			[]float64{1, 2, 10},
		},
	}

	for _, c := range cases {
		total := c.charts.Total()
		if !reflect.DeepEqual(total, c.total) {
			t.Errorf("wrong data:\nhas:  %v\nwant: %v",
				total, c.total)
		}
	}
}

func TestChartsMax(t *testing.T) {
	cases := []struct {
		charts Charts
		max    Column
	}{
		{
			Charts{},
			Column{},
		},
		{
			Charts{"a": []float64{}},
			Column{{Name: "a", Score: 0}},
		},
		{
			Charts{"o0o": []float64{0, 0, 7}},
			Column{{Name: "o0o", Score: 7}},
		},
		{
			Charts{
				"o0o": []float64{0, 0, 7},
				"lol": []float64{1, 2, 0}},
			Column{
				{Name: "o0o", Score: 7},
				{Name: "lol", Score: 2}},
		},
	}

	for _, c := range cases {
		max := c.charts.Max()
		if !reflect.DeepEqual(max, c.max) {
			t.Errorf("wrong data:\nhas:  %v\nwant: %v",
				max, c.max)
		}
	}
}
