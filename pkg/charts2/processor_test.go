package charts2

import (
	"reflect"
	"testing"
)

func TestLazyEval(t *testing.T) {
	charts := &Charts{
		Values: map[string][]float64{
			"A": {8, 8, 0, 0},
			"B": {16, 0, 0, 0},
			"C": {1, 1, 2, 1},
		},
		titles: []Title{KeyTitle("A"), KeyTitle("B"), KeyTitle("C")},
	}

	cs := []struct {
		name     string
		lc       LazyCharts
		titles   []Title
		len      int
		rowA04   []float64
		rowB13   []float64
		colAB1   map[string]float64
		colB3    map[string]float64
		dataAB14 map[string][]float64
	}{
		{
			"charts themselves",
			charts,
			[]Title{KeyTitle("A"), KeyTitle("B"), KeyTitle("C")}, 4,
			[]float64{8, 8, 0, 0},
			[]float64{0, 0},
			map[string]float64{"A": 8, "B": 0},
			map[string]float64{"B": 0},
			map[string][]float64{
				"A": {8, 0, 0},
				"B": {0, 0, 0},
			},
		},
		{
			"sum",
			Sum(charts),
			[]Title{KeyTitle("A"), KeyTitle("B"), KeyTitle("C")}, 4,
			[]float64{8, 16, 16, 16},
			[]float64{16, 16},
			map[string]float64{"A": 16, "B": 16},
			map[string]float64{"B": 16},
			map[string][]float64{
				"A": {16, 16, 16},
				"B": {16, 16, 16},
			},
		},
		{
			"sum of sum",
			Sum(Sum(charts)),
			[]Title{KeyTitle("A"), KeyTitle("B"), KeyTitle("C")}, 4,
			[]float64{8, 24, 40, 56},
			[]float64{32, 48},
			map[string]float64{"A": 24, "B": 32},
			map[string]float64{"B": 64},
			map[string][]float64{
				"A": {24, 40, 56},
				"B": {32, 48, 64},
			},
		},
		{
			"fade",
			Fade(charts, 1),
			[]Title{KeyTitle("A"), KeyTitle("B"), KeyTitle("C")}, 4,
			[]float64{8, 12, 6, 3},
			[]float64{8, 4},
			map[string]float64{"A": 12, "B": 8},
			map[string]float64{"B": 2},
			map[string][]float64{
				"A": {12, 6, 3},
				"B": {8, 4, 2},
			},
		},
		{
			"max of fade",
			Max(Fade(charts, 1)),
			[]Title{KeyTitle("A"), KeyTitle("B"), KeyTitle("C")}, 4,
			[]float64{8, 12, 12, 12},
			[]float64{16, 16},
			map[string]float64{"A": 12, "B": 16},
			map[string]float64{"B": 16},
			map[string][]float64{
				"A": {12, 12, 12},
				"B": {16, 16, 16},
			},
		},
		{
			"merge partition",
			&partitionSum{chartsNode: chartsNode{parent: charts},
				partition: map[string]string{
					"A": "A", "B": "B", "C": "B",
				},
				key: func(t Title) string { return t.Key() },
			},
			[]Title{KeyTitle("A"), KeyTitle("B")}, 4,
			[]float64{8, 8, 0, 0},
			[]float64{1, 2},
			map[string]float64{"A": 8, "B": 1},
			map[string]float64{"B": 1},
			map[string][]float64{
				"A": {8, 0, 0},
				"B": {1, 2, 1},
			},
		},
	}

	for _, c := range cs {
		t.Run(c.name, func(t *testing.T) {
			{
				row := c.lc.Row(KeyTitle("A"), 0, 4)
				if !reflect.DeepEqual(row, c.rowA04) {
					t.Error("row A 0-4 not equal:", row, "!=", c.rowA04)
				}
			}
			{
				row := c.lc.Row(KeyTitle("B"), 1, 3)
				if !reflect.DeepEqual(row, c.rowB13) {
					t.Error("row B 1-3 not equal:", row, "!=", c.rowB13)
				}
			}
			{
				col := c.lc.Column([]Title{KeyTitle("A"), KeyTitle("B")}, 1)
				if !reflect.DeepEqual(col, c.colAB1) {
					t.Error("col A,B 1 not equal:", col, "!=", c.colAB1)
				}
			}
			{
				col := c.lc.Column([]Title{KeyTitle("B")}, 3)
				if !reflect.DeepEqual(col, c.colB3) {
					t.Error("col B 3 not equal:", col, "!=", c.colB3)
				}
			}
			{
				data := c.lc.Data([]Title{KeyTitle("A"), KeyTitle("B")}, 1, 4)
				if !reflect.DeepEqual(data, c.dataAB14) {
					t.Error("data A,B 1-4 not equal:", data, "!=", c.dataAB14)
				}
			}
			{
				titles := c.lc.Titles()
				if !areTitlesSame(titles, c.titles) {
					t.Error("not equal:", titles, "!=", c.titles)
				}
			}
			{
				len := c.lc.Len()
				if !reflect.DeepEqual(len, c.len) {
					t.Error("not equal:", len, "!=", c.len)
				}
			}
		})
	}
}

func allEqual(a, b Title) bool {
	return a.String() == b.String() &&
		a.Artist() == b.Artist() &&
		a.Key() == b.Key() &&
		a.Song() == b.Song()
}

func areTitlesSame(a, b []Title) bool {
	if len(a) != len(b) {
		return false
	}
	used := make([]bool, len(a))
	for _, c := range a {
		found := false
		for i, d := range b {
			if allEqual(c, d) {
				if used[i] {
					return false
				}
				used[i] = true
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func TestEmptyCharts(t *testing.T) {
	c := &Charts{}

	if c.Len() != -1 {
		t.Error("unxecptected len:", c.Len())
	}
}
