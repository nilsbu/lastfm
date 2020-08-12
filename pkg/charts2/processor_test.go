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
		titles   []string
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
			[]string{"A", "B", "C"}, 4,
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
			[]string{"A", "B", "C"}, 4,
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
			[]string{"A", "B", "C"}, 4,
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
			[]string{"A", "B", "C"}, 4,
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
			[]string{"A", "B", "C"}, 4,
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
			[]string{"A", "B"}, 4,
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
				row := c.lc.Row("A", 0, 4)
				if !reflect.DeepEqual(row, c.rowA04) {
					t.Error("row A 0-4 not equal:", row, "!=", c.rowA04)
				}
			}
			{
				row := c.lc.Row("B", 1, 3)
				if !reflect.DeepEqual(row, c.rowB13) {
					t.Error("row B 1-3 not equal:", row, "!=", c.rowB13)
				}
			}
			{
				col := c.lc.Column([]string{"A", "B"}, 1)
				if !reflect.DeepEqual(col, c.colAB1) {
					t.Error("col A,B 1 not equal:", col, "!=", c.colAB1)
				}
			}
			{
				col := c.lc.Column([]string{"B"}, 3)
				if !reflect.DeepEqual(col, c.colB3) {
					t.Error("col B 3 not equal:", col, "!=", c.colB3)
				}
			}
			{
				data := c.lc.Data([]string{"A", "B"}, 1, 4)
				if !reflect.DeepEqual(data, c.dataAB14) {
					t.Error("data A,B 1-4 not equal:", data, "!=", c.dataAB14)
				}
			}
			{
				titles := c.lc.Titles()
				if !reflect.DeepEqual(titles, c.titles) {
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

func TestEmptyCharts(t *testing.T) {
	c := &Charts{}

	if c.Len() != -1 {
		t.Error("unxecptected len:", c.Len())
	}
}
