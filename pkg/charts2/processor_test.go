package charts2

import (
	"math"
	"math/rand"
	"reflect"
	"testing"

	"github.com/nilsbu/lastfm/pkg/rsrc"
)

// TODO merge with LazyCharts
func TestLazyChartsPartial(t *testing.T) {
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
	return a != nil && b != nil &&
		a.String() == b.String() &&
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

func checkTitle(t *testing.T, x, a Title) {
	if x.String() != a.String() {
		t.Fatalf("String(): expect=%v, actual=%v",
			x.String(), a.String())
	}
	if x.Key() != a.Key() {
		t.Fatalf("Key(): expect=%v, actual=%v",
			x.Key(), a.Key())
	}
	if x.Artist() != a.Artist() {
		t.Fatalf("Artist(): expect=%v, actual=%v",
			x.Artist(), a.Artist())
	}
	if x.Song() != a.Song() {
		t.Fatalf("Song(): expect=%v, actual=%v",
			x.Song(), a.Song())
	}
}

func checkMeta(t *testing.T, expect, actual LazyCharts) {

	if expect.Len() != actual.Len() {
		t.Fatalf("len differs: expect=%v, actual %v", expect.Len(), actual.Len())
	}

	tx := expect.Titles()
	ta := actual.Titles()
	if len(tx) != len(ta) {
		t.Fatalf("number of titles differs: expect=%v, actual %v",
			len(tx), len(ta))
	}
	for i := range tx {
		checkTitle(t, tx[i], ta[i])
	}
}

func ranges(size, nRand int) [][2]int {
	ranges := [][2]int{{0, size}}
	for i := 0; i < nRand; i++ {
		b := rand.Int() % (size - 1)
		s := rand.Int() % (size - b)
		ranges = append(ranges, [2]int{b, b + s})
	}
	return ranges
}

func checkRows(t *testing.T, expect, actual LazyCharts, ranges [][2]int) {

	for _, be := range ranges {
		for _, title := range expect.Titles() {
			x := expect.Row(title, be[0], be[1])
			a := actual.Row(title, be[0], be[1])

			if len(a) != be[1]-be[0] {
				t.Fatalf("row length: expect=%v-%v=%v, actual=%v",
					be[1], be[0], be[1]-be[0], len(a))
			}
			if len(x) != be[1]-be[0] {
				t.Fatalf("row length: expect=%v-%v=%v, ground truth=%v",
					be[1], be[0], be[1]-be[0], len(x))
			}

			eq := true
			for i := range x {
				if math.Abs(x[i]-a[i]) > 1e-6 {
					eq = false
					break
				}
			}
			if !eq {
				t.Errorf("row(%v, [%v-%v]): expect=%v, actual=%v",
					title, be[0], be[1], x, a)
			}
		}
	}
}

func sets(titles []Title, nRand int) [][]Title {
	sets := [][]Title{titles}

	for i := 0; i < nRand; i++ {
		set := []Title{}
		set = append(set, titles...)

		rand.Shuffle(len(set), func(i, j int) {
			set[i], set[j] = set[j], set[i]
		})

		n := rand.Int() % len(set)
		sets = append(sets, set[:n])
	}

	return sets
}

func checkCols(t *testing.T, expect, actual LazyCharts, sets [][]Title) {

	for _, set := range sets {
		for i := 0; i < expect.Len(); i++ {
			x := expect.Column(set, i)
			a := actual.Column(set, i)

			if len(a) != len(set) {
				t.Fatalf("col length: expect=%v, actual=%v",
					len(set), len(a))
			}
			if len(x) != len(set) {
				t.Fatalf("col length: expect=%v, ground truth=%v",
					len(set), len(x))
			}

			for k := range x {
				checkTitle(t, x[k].Title, a[k].Title)
				if math.Abs(x[k].Value-a[k].Value) > 1e-6 {
					t.Errorf("col(%v, %v): expect=%v, actual=%v",
						k, i, x[k].Value, a[k].Value)
				}
			}
		}
	}
}

func checkData(t *testing.T, expect, actual LazyCharts,
	ranges [][2]int, sets [][]Title) {

	for i := range sets {
		set := sets[i]
		b, e := ranges[i][0], ranges[i][1]

		x := expect.Data(set, b, e)
		a := actual.Data(set, b, e)

		for k := range x {
			rowX := x[k]
			rowA := a[k]
			checkTitle(t, rowX.Title, rowA.Title)

			if len(rowA.Line) != e-b {
				t.Fatalf("data row length: expect=%v-%v=%v, actual=%v",
					e, b, e-b, len(rowA.Line))
			}
			if len(rowA.Line) != e-b {
				t.Fatalf("data row length: expect=%v-%v=%v, ground truth=%v",
					e, b, e-b, len(rowX.Line))
			}

			eq := true
			for i := range rowX.Line {
				if math.Abs(rowX.Line[i]-rowA.Line[i]) > 1e-6 {
					eq = false
					break
				}
			}
			if !eq {
				t.Errorf("row(%v, [%v-%v]): expect=%v, actual=%v",
					k, b, e, rowX.Line, rowA.Line)
			}
		}
	}
}

func checkLazyCharts(t *testing.T, expect, actual LazyCharts, nRand int) {
	checkMeta(t, expect, actual)
	checkRows(t, expect, actual, ranges(expect.Len(), nRand))
	checkCols(t, expect, actual, sets(expect.Titles(), nRand))
	checkData(t, expect, actual,
		ranges(expect.Len(), nRand), sets(expect.Titles(), nRand))
}

func TestLazyCharts(t *testing.T) {
	root := &charts{
		values: map[string][]float64{
			"A": {0, 0, 0, 1, 0},
			"B": {0, 1, 0, 1, 0},
		},
		titles: []Title{KeyTitle("A"), KeyTitle("B")},
	}

	f := 1 / math.Sqrt(2*math.Pi)
	m := []float64{f * math.Exp(0), f * math.Exp(-.5), f * math.Exp(-2)}

	cs := []struct {
		name   string
		actual LazyCharts
		expect LazyCharts
	}{
		{
			"Gaussian mirror none",
			Gaussian(root, 1, 2, false, false),
			&charts{
				values: map[string][]float64{
					"A": {0, m[2], m[1], m[0], m[1]},
					"B": {m[1], m[2] + m[0], 2 * m[1], m[2] + m[0], m[1]},
				},
				titles: []Title{KeyTitle("A"), KeyTitle("B")},
			},
		},
		{
			"Gaussian mirror begin",
			Gaussian(root, 1, 2, true, false),
			&charts{
				values: map[string][]float64{
					"A": {0, m[2], m[1], m[0], m[1]},
					"B": {m[1] + m[2], m[2] + m[0], 2 * m[1], m[2] + m[0], m[1]},
				},
				titles: []Title{KeyTitle("A"), KeyTitle("B")},
			},
		},
		{
			"Gaussian mirror both",
			Gaussian(root, 1, 2, true, true),
			&charts{
				values: map[string][]float64{
					"A": {0, m[2], m[1], m[0], m[1] + m[2]},
					"B": {m[1] + m[2], m[2] + m[0], 2 * m[1], m[2] + m[0], m[1] + m[2]},
				},
				titles: []Title{KeyTitle("A"), KeyTitle("B")},
			},
		},
		{
			"ColumnSum",
			ColumnSum(&charts{
				values: map[string][]float64{
					"A": {0, 1, 2, 3, 4},
					"B": {2, 2, 2, -1, -1},
					"":  {0, 0, 0, 0, 0},
				},
				titles: []Title{KeyTitle("A"), KeyTitle("B"), KeyTitle("")},
			}),
			&charts{
				values: map[string][]float64{
					"total": {2, 3, 4, 2, 3},
				},
				titles: []Title{StringTitle("total")},
			},
		},
		{
			"cached Gaussian",
			Gaussian(root, 1, 2, false, false),
			Cache(Gaussian(root, 1, 2, false, false)),
		},
		{
			"Interval",
			Interval(&charts{
				values: map[string][]float64{
					"A": {0, 1, 2, 3, 4},
				},
				titles: []Title{KeyTitle("A")},
			},
				Range{
					Begin:      rsrc.ParseDay("2020-01-01"),
					End:        rsrc.ParseDay("2020-01-03"),
					Registered: rsrc.ParseDay("2019-12-31")},
			),
			&charts{
				values: map[string][]float64{
					"A": {1, 2},
				},
				titles: []Title{KeyTitle("A")},
			},
		},
	}

	for _, c := range cs {
		t.Run(c.name, func(t *testing.T) {
			checkLazyCharts(t, c.expect, c.actual, 5)
		})
	}
}

func TestCacheIncremental(t *testing.T) {
	expect := &charts{
		values: map[string][]float64{
			"A": {1, 2, 3, 4, 5},
			"B": {11, 13, 12, 15, 14},
		},
		titles: []Title{KeyTitle("A"), KeyTitle("B")},
	}

	// Row
	{
		actual := Cache(expect)
		ranges := [][2]int{
			{2, 2}, {1, 2}, {3, 4}, {0, 5}, {0, 5},
		}

		checkRows(t, expect, actual, ranges)
	}

	// Column
	{
		actual := Cache(expect)
		sets := [][]Title{
			{}, {KeyTitle("A")}, {KeyTitle("B")}, {KeyTitle("A"), KeyTitle("B")},
		}

		checkCols(t, expect, actual, sets)
	}

	// Data
	{
		actual := Cache(expect)
		ranges := [][2]int{
			{1, 2}, {3, 4}, {0, 5}, {0, 5},
		}
		sets := [][]Title{
			{}, {KeyTitle("A")}, {KeyTitle("B")}, {KeyTitle("A"), KeyTitle("B")},
		}

		checkData(t, expect, actual, ranges, sets)
	}
}

func TestTop(t *testing.T) {
	for _, c := range []struct {
		name   string
		charts LazyCharts
		n      int
		titles []Title
	}{
		// {
		// 	"empty",
		// 	FromMap(map[string][]float64{}),
		// 	2,
		// 	[]Title{},
		// },
		// {
		// 	"2 out of 3",
		// 	FromMap(map[string][]float64{
		// 		"A": {1, 3},
		// 		"B": {1, 1},
		// 		"C": {0, 2},
		// 	}),
		// 	2,
		// 	[]Title{KeyTitle("A"), KeyTitle("C")},
		// },
		// {
		// 	"n > len",
		// 	FromMap(map[string][]float64{
		// 		"A": {1, 3},
		// 		"B": {1, 1},
		// 		"C": {0, 2},
		// 	}),
		// 	4,
		// 	[]Title{KeyTitle("A"), KeyTitle("C"), KeyTitle("B")},
		// },
		{
			"many",
			FromMap(map[string][]float64{
				"A": {1, 3},
				"B": {1, 1},
				"C": {0, 6},
				"D": {0, 0},
				"E": {0, 8},
				"F": {0, 11},
				"G": {0, 2},
				"H": {0, -1},
				"I": {0, 99},
			}),
			3,
			[]Title{KeyTitle("I"), KeyTitle("F"), KeyTitle("E")},
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			titles := Top(c.charts, c.n)
			if !areTitlesSame(c.titles, titles) {
				t.Errorf("expect: %v\nactual: %v", c.titles, titles)
			}
		})
	}

}
