package charts2

import (
	"math"
	"reflect"
	"testing"
)

func TestLazyCharts(t *testing.T) {
	root := &charts{
		values: map[string][]float64{
			"A": {8, 8, 0, 0},
			"B": {16, 0, 0, 0},
			"C": {1, 1, 2, 1},
		},
		titles: []Title{KeyTitle("A"), KeyTitle("B"), KeyTitle("C")},
	}

	songs := [][]Song{
		{
			{Artist: "A", Title: "a", Duration: 3},
			{Artist: "A", Title: "b", Duration: 4},
			{Artist: "A", Title: "a", Duration: 3},
			{Artist: "B", Title: "x", Duration: 4.5},
		},
		{
			{Artist: "A", Title: "a", Duration: 3},
			{Artist: "A", Title: "b", Duration: 4},
		},
		{
			{Artist: "B", Title: "x", Duration: 4.5},
			{Artist: "B", Title: "y", Duration: 1},
		},
		{},
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
			root,
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
			Sum(root),
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
			Sum(Sum(root)),
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
			Fade(root, 1),
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
			Max(Fade(root, 1)),
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
			Group(
				root,
				KeyPartition([][2]Title{
					{KeyTitle("A"), KeyTitle("A")},
					{KeyTitle("B"), KeyTitle("B")},
					{KeyTitle("C"), KeyTitle("B")},
				}),
			),
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
		{
			"artist charts",
			Artists(songs),
			[]Title{ArtistTitle("A"), ArtistTitle("B")}, 4,
			[]float64{3, 2, 0, 0},
			[]float64{0, 2},
			map[string]float64{"A": 2, "B": 0},
			map[string]float64{"B": 0},
			map[string][]float64{
				"A": {2, 0, 0},
				"B": {0, 2, 0},
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
				if !eqColWithKeyTitle(c.colAB1, col) {
					t.Error("col A,B 1 not equal:", c.colAB1, "!=", col)
				}
			}
			{
				col := c.lc.Column([]Title{KeyTitle("B")}, 3)
				if !eqColWithKeyTitle(c.colB3, col) {
					t.Error("col B 3 not equal:", c.colB3, "!=", col)
				}
			}
			{
				data := c.lc.Data([]Title{KeyTitle("A"), KeyTitle("B")}, 1, 4)
				if !eqDataWithKeyTitle(c.dataAB14, data) {
					t.Error("data A,B 1-4 not equal:", c.dataAB14, "!=", data)
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

func eqColWithKeyTitle(expect map[string]float64, actual TitleValueMap) bool {
	if len(expect) != len(actual) {
		return false
	}

	for t, v := range expect {
		if tv, ok := actual[t]; ok {
			if tv.Value != v {
				return false
			}
			if !allEqual(KeyTitle(t), tv.Title) {
				return false
			}
		} else {
			return false
		}
	}

	return true
}

func eqDataWithKeyTitle(expect map[string][]float64, actual TitleLineMap) bool {
	if len(expect) != len(actual) {
		return false
	}

	for t, v := range expect {
		if tv, ok := actual[t]; ok {
			if !reflect.DeepEqual(v, tv.Line) {
				return false
			}
			if !allEqual(KeyTitle(t), tv.Title) {
				return false
			}
		} else {
			return false
		}
	}

	return true
}

func allEqual(a, b Title) bool {
	return a.String() == b.String() &&
		a.Key() == b.Key() &&
		a.Artist() == b.Artist() &&
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
	c := &charts{}

	if c.Len() != -1 {
		t.Error("unxecptected len:", c.Len())
	}
}

func TestGaussian(t *testing.T) {
	root := &charts{
		values: map[string][]float64{
			"A": {0, 0, 0, 1, 0},
			"B": {0, 1, 0, 1, 0},
		},
		titles: []Title{KeyTitle("A"), KeyTitle("B")},
	}

	f := 1 / math.Sqrt(2*math.Pi)
	m := []float64{math.Exp(0), math.Exp(-.5), math.Exp(-2)}

	cs := []struct {
		name   string
		lc     LazyCharts
		expect LazyCharts
	}{
		{
			"mirror none",
			Gaussian(root, 1, 2, false, false),
			&charts{
				values: map[string][]float64{
					"A": {0, f * m[2], f * m[1], f * m[0], f * m[1]},
					"B": {f * m[1], f * (m[2] + m[0]), 2 * f * m[1], f * (m[2] + m[0]), f * m[1]},
				},
				titles: []Title{KeyTitle("A"), KeyTitle("B")},
			},
		},
		{
			"mirror begin",
			Gaussian(root, 1, 2, true, false),
			&charts{
				values: map[string][]float64{
					"A": {0, f * m[2], f * m[1], f * m[0], f * m[1]},
					"B": {f * (m[1] + m[2]), f * (m[2] + m[0]), 2 * f * m[1], f * (m[2] + m[0]), f * m[1]},
				},
				titles: []Title{KeyTitle("A"), KeyTitle("B")},
			},
		},
		{
			"mirror both",
			Gaussian(root, 1, 2, true, true),
			&charts{
				values: map[string][]float64{
					"A": {0, f * m[2], f * m[1], f * m[0], f * (m[1] + m[2])},
					"B": {f * (m[1] + m[2]), f * (m[2] + m[0]), 2 * f * m[1], f * (m[2] + m[0]), f * (m[1] + m[2])},
				},
				titles: []Title{KeyTitle("A"), KeyTitle("B")},
			},
		},
	}

	for _, c := range cs {
		t.Run(c.name, func(t *testing.T) {
			// Rows
			for _, title := range c.expect.Titles() {
				x := c.expect.Row(title, 0, c.expect.Len())
				a := c.lc.Row(title, 0, c.expect.Len())
				eq := true
				for i := range x {
					if math.Abs(x[i]-a[i]) > 1e-6 {
						eq = false
						break
					}
				}
				if !eq {
					t.Errorf("row(%v): expect=%v, actual=%v", title, x, a)
				}
			}

			// Columns
			for i := 0; i < c.expect.Len(); i++ {
				x := c.expect.Column(c.expect.Titles(), i)
				a := c.lc.Column(c.expect.Titles(), i)

				for k := range x {
					if math.Abs(x[k].Value-a[k].Value) > 1e-6 {
						t.Errorf("col(%v, %v): expect=%v, actual=%v",
							k, i, x[k].Value, a[k].Value)
					}
				}
			}

			// TODO data
		})
	}
}
