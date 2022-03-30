package charts_test

import (
	"math"
	"math/rand"
	"reflect"
	"testing"

	"github.com/nilsbu/lastfm/pkg/charts"
	"github.com/nilsbu/lastfm/pkg/info"
	"github.com/nilsbu/lastfm/pkg/rsrc"
)

// TODO merge with LazyCharts
func TestLazyChartsPartial(t *testing.T) {
	root := charts.FromMap(map[string][]float64{
		"A": {8, 8, 0, 0},
		"B": {16, 0, 0, 0},
		"C": {1, 1, 2, 1},
	})

	songs := [][]info.Song{
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
		lc       charts.Charts
		titles   []charts.Title
		len      int
		rowA04   []float64
		rowB13   []float64
		colAB1   []float64
		colB3    []float64
		dataAB14 [][]float64
	}{
		{
			"charts themselves",
			root,
			[]charts.Title{charts.KeyTitle("A"), charts.KeyTitle("B"), charts.KeyTitle("C")}, 4,
			[]float64{8, 8, 0, 0},
			[]float64{0, 0},
			[]float64{8, 0},
			[]float64{0},
			[][]float64{
				{8, 0, 0},
				{0, 0, 0},
			},
		},
		{
			"sum",
			charts.Sum(root),
			[]charts.Title{charts.KeyTitle("A"), charts.KeyTitle("B"), charts.KeyTitle("C")}, 4,
			[]float64{8, 16, 16, 16},
			[]float64{16, 16},
			[]float64{16, 16},
			[]float64{16},
			[][]float64{
				{16, 16, 16},
				{16, 16, 16},
			},
		},
		{
			"sum of sum",
			charts.Sum(charts.Sum(root)),
			[]charts.Title{charts.KeyTitle("A"), charts.KeyTitle("B"), charts.KeyTitle("C")}, 4,
			[]float64{8, 24, 40, 56},
			[]float64{32, 48},
			[]float64{24, 32},
			[]float64{64},
			[][]float64{
				{24, 40, 56},
				{32, 48, 64},
			},
		},
		{
			"fade",
			charts.Fade(root, 1),
			[]charts.Title{charts.KeyTitle("A"), charts.KeyTitle("B"), charts.KeyTitle("C")}, 4,
			[]float64{8, 12, 6, 3},
			[]float64{8, 4},
			[]float64{12, 8},
			[]float64{2},
			[][]float64{
				{12, 6, 3},
				{8, 4, 2},
			},
		},
		{
			"max of fade",
			charts.Max(charts.Fade(root, 1)),
			[]charts.Title{charts.KeyTitle("A"), charts.KeyTitle("B"), charts.KeyTitle("C")}, 4,
			[]float64{8, 12, 12, 12},
			[]float64{16, 16},
			[]float64{12, 16},
			[]float64{16},
			[][]float64{
				{12, 12, 12},
				{16, 16, 16},
			},
		},
		{
			"merge partition",
			charts.Group(
				root,
				charts.KeyPartition([][2]charts.Title{
					{charts.KeyTitle("A"), charts.KeyTitle("A")},
					{charts.KeyTitle("B"), charts.KeyTitle("B")},
					{charts.KeyTitle("C"), charts.KeyTitle("B")},
				}),
			),
			[]charts.Title{charts.KeyTitle("A"), charts.KeyTitle("B")}, 4,
			[]float64{8, 8, 0, 0},
			[]float64{1, 2},
			[]float64{8, 1},
			[]float64{1},
			[][]float64{
				{8, 0, 0},
				{1, 2, 1},
			},
		},
		{
			"artist charts",
			charts.Artists(songs),
			[]charts.Title{charts.ArtistTitle("A"), charts.ArtistTitle("B")}, 4,
			[]float64{3, 2, 0, 0},
			[]float64{0, 2},
			[]float64{2, 0},
			[]float64{0},
			[][]float64{
				{2, 0, 0},
				{0, 2, 0},
			},
		},
	}

	for _, c := range cs {
		t.Run(c.name, func(t *testing.T) {
			{
				row, _ := c.lc.Data([]charts.Title{charts.KeyTitle("A")}, 0, 4)
				if !reflect.DeepEqual(row[0], c.rowA04) {
					t.Error("row A 0-4 not equal:", row[0], "!=", c.rowA04)
				}
			}
			{
				row, _ := c.lc.Data([]charts.Title{charts.KeyTitle("B")}, 1, 3)
				if !reflect.DeepEqual(row[0], c.rowB13) {
					t.Error("row B 1-3 not equal:", row[0], "!=", c.rowB13)
				}
			}
			{
				col_, _ := c.lc.Data([]charts.Title{charts.KeyTitle("A"), charts.KeyTitle("B")}, 1, 2)
				col := make([]float64, len(col_))
				for i, c := range col_ {
					col[i] = c[0]
				}
				if !reflect.DeepEqual(c.colAB1, col) {
					t.Error("col A,B 1 not equal:", c.colAB1, "!=", col)
				}
			}
			{
				col_, _ := c.lc.Data([]charts.Title{charts.KeyTitle("B")}, 3, 4)
				col := make([]float64, len(col_))
				for i, c := range col_ {
					col[i] = c[0]
				}
				if !reflect.DeepEqual(c.colB3, col) {
					t.Error("col B 3 not equal:", c.colB3, "!=", col)
				}
			}
			{
				data, _ := c.lc.Data([]charts.Title{charts.KeyTitle("A"), charts.KeyTitle("B")}, 1, 4)
				if !reflect.DeepEqual(c.dataAB14, data) {
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

func allEqual(a, b charts.Title) bool {
	return a != nil && b != nil &&
		a.String() == b.String() &&
		a.Key() == b.Key() &&
		a.Artist() == b.Artist()
}

func areTitlesSame(a, b []charts.Title) bool {
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
	c := charts.FromMap(map[string][]float64{})

	if c.Len() != -1 {
		t.Error("unxecptected len:", c.Len())
	}
}

func checkTitle(t *testing.T, x, a charts.Title) {
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
}

func checkMeta(t *testing.T, expect, actual charts.Charts) {

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

func checkRows(t *testing.T, expect, actual charts.Charts, ranges [][2]int) {
	// Since Rows() doesn't exist anymore, this is just another Data() test
	for _, be := range ranges {
		for _, title := range expect.Titles() {
			xs, _ := expect.Data([]charts.Title{title}, be[0], be[1])
			as, _ := actual.Data([]charts.Title{title}, be[0], be[1])
			x, a := xs[0], as[0]

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

func sets(titles []charts.Title, nRand int) [][]charts.Title {
	sets := [][]charts.Title{titles}

	for i := 0; i < nRand; i++ {
		set := []charts.Title{}
		set = append(set, titles...)

		rand.Shuffle(len(set), func(i, j int) {
			set[i], set[j] = set[j], set[i]
		})

		n := rand.Int() % len(set)
		sets = append(sets, set[:n])
	}

	return sets
}

func checkCols(t *testing.T, expect, actual charts.Charts, sets [][]charts.Title) {

	for _, set := range sets {
		for i := 0; i < expect.Len(); i++ {
			x, _ := expect.Data(set, i, i+1)
			a, _ := actual.Data(set, i, i+1)

			if len(a) != len(set) {
				t.Fatalf("col length: expect=%v, actual=%v",
					len(set), len(a))
			}
			if len(x) != len(set) {
				t.Fatalf("col length: expect=%v, ground truth=%v",
					len(set), len(x))
			}

			for k := range x {
				if math.Abs(x[k][0]-a[k][0]) > 1e-6 {
					t.Errorf("col(%v, %v): expect=%v, actual=%v",
						k, i, x[k], a[k])
				}
			}
		}
	}
}

func checkData(t *testing.T, expect, actual charts.Charts,
	ranges [][2]int, sets [][]charts.Title) {

	for i := range sets {
		set := sets[i]
		b, e := ranges[i][0], ranges[i][1]

		x, _ := expect.Data(set, b, e)
		a, _ := actual.Data(set, b, e)

		for k := range x {
			rowX := x[k]
			rowA := a[k]
			if len(rowA) != e-b {
				t.Fatalf("data row length: expect=%v-%v=%v, actual=%v",
					e, b, e-b, len(rowA))
			}
			if len(rowA) != e-b {
				t.Fatalf("data row length: expect=%v-%v=%v, ground truth=%v",
					e, b, e-b, len(rowX))
			}

			eq := true
			for i := range rowX {
				if math.Abs(rowX[i]-rowA[i]) > 1e-6 {
					eq = false
					break
				}
			}
			if !eq {
				t.Errorf("row(%v, [%v-%v]): expect=%v, actual=%v",
					k, b, e, rowX, rowA)
			}
		}
	}
}

func checkLazyCharts(t *testing.T, expect, actual charts.Charts, nRand int) {
	checkMeta(t, expect, actual)
	checkRows(t, expect, actual, ranges(expect.Len(), nRand))
	checkCols(t, expect, actual, sets(expect.Titles(), nRand))
	checkData(t, expect, actual,
		ranges(expect.Len(), nRand), sets(expect.Titles(), nRand))
}

func TestLazyCharts(t *testing.T) {
	root := charts.FromMap(map[string][]float64{
		"A": {0, 0, 0, 1, 0},
		"B": {0, 1, 0, 1, 0},
	})

	f := 1 / math.Sqrt(2*math.Pi)
	m := []float64{f * math.Exp(0), f * math.Exp(-.5), f * math.Exp(-2)}

	cs := []struct {
		name   string
		actual charts.Charts
		expect charts.Charts
	}{
		{
			"Gaussian mirror none",
			charts.Gaussian(root, 1, 2, false, false),
			charts.FromMap(map[string][]float64{
				"A": {0, m[2], m[1], m[0], m[1]},
				"B": {m[1], m[2] + m[0], 2 * m[1], m[2] + m[0], m[1]},
			}),
		},
		{
			"Gaussian mirror begin",
			charts.Gaussian(root, 1, 2, true, false),
			charts.FromMap(map[string][]float64{
				"A": {0, m[2], m[1], m[0], m[1]},
				"B": {m[1] + m[2], m[2] + m[0], 2 * m[1], m[2] + m[0], m[1]},
			}),
		},
		{
			"Gaussian mirror both",
			charts.Gaussian(root, 1, 2, true, true),
			charts.FromMap(map[string][]float64{
				"A": {0, m[2], m[1], m[0], m[1] + m[2]},
				"B": {m[1] + m[2], m[2] + m[0], 2 * m[1], m[2] + m[0], m[1] + m[2]},
			}),
		},
		{
			"ColumnSum",
			charts.ColumnSum(charts.FromMap(map[string][]float64{
				"A": {0, 1, 2, 3, 4},
				"B": {2, 2, 2, -1, -1},
				"":  {0, 0, 0, 0, 0},
			})),
			charts.FromMap(map[string][]float64{
				"total": {2, 3, 4, 2, 3},
			}),
		},
		{
			"cached Gaussian",
			charts.Gaussian(root, 1, 2, false, false),
			charts.Cache(charts.Gaussian(root, 1, 2, false, false)),
		},
		{
			"Interval",
			charts.Interval(charts.FromMap(map[string][]float64{
				"A": {0, 1, 2, 3, 4},
			}),
				charts.Range{
					Begin:      rsrc.ParseDay("2020-01-01"),
					End:        rsrc.ParseDay("2020-01-03"),
					Registered: rsrc.ParseDay("2019-12-31")},
			),
			charts.FromMap(map[string][]float64{
				"A": {1, 2},
			}),
		},
	}

	for _, c := range cs {
		t.Run(c.name, func(t *testing.T) {
			checkLazyCharts(t, c.expect, c.actual, 5)
		})
	}
}

func TestCacheIncremental(t *testing.T) {
	expect := charts.FromMap(map[string][]float64{
		"A": {1, 2, 3, 4, 5},
		"B": {11, 13, 12, 15, 14},
	})

	// Row
	{
		actual := charts.Cache(expect)
		ranges := [][2]int{
			{2, 2}, {1, 2}, {3, 4}, {0, 5}, {0, 5},
		}

		checkRows(t, expect, actual, ranges)
	}

	// Column
	{
		actual := charts.Cache(expect)
		sets := [][]charts.Title{
			{}, {charts.KeyTitle("A")}, {charts.KeyTitle("B")}, {charts.KeyTitle("A"), charts.KeyTitle("B")},
		}

		checkCols(t, expect, actual, sets)
	}

	// Data
	{
		actual := charts.Cache(expect)
		ranges := [][2]int{
			{1, 2}, {3, 4}, {0, 5}, {0, 5},
		}
		sets := [][]charts.Title{
			{}, {charts.KeyTitle("A")}, {charts.KeyTitle("B")}, {charts.KeyTitle("A"), charts.KeyTitle("B")},
		}

		checkData(t, expect, actual, ranges, sets)
	}
}

func TestTop(t *testing.T) {
	for _, c := range []struct {
		name   string
		charts charts.Charts
		n      int
		titles []charts.Title
	}{
		// TODO reenter these tests
		{
			"empty",
			charts.FromMap(map[string][]float64{}),
			2,
			[]charts.Title{},
		},
		{
			"2 out of 3",
			charts.FromMap(map[string][]float64{
				"A": {1, 3},
				"B": {1, 1},
				"C": {0, 2},
			}),
			2,
			[]charts.Title{charts.KeyTitle("A"), charts.KeyTitle("C")},
		},
		{
			"n > len",
			charts.FromMap(map[string][]float64{
				"A": {1, 3},
				"B": {1, 1},
				"C": {0, 2},
			}),
			4,
			[]charts.Title{charts.KeyTitle("A"), charts.KeyTitle("C"), charts.KeyTitle("B")},
		},
		{
			"drop zero",
			charts.FromMap(map[string][]float64{
				"A": {2},
				"B": {3},
				"C": {0},
			}),
			4,
			[]charts.Title{charts.KeyTitle("B"), charts.KeyTitle("A")},
		},
		{
			"many",
			charts.FromMap(map[string][]float64{
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
			[]charts.Title{charts.KeyTitle("I"), charts.KeyTitle("F"), charts.KeyTitle("E")},
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			titles, _ := charts.Top(c.charts, c.n)
			if !areTitlesSame(c.titles, titles) {
				t.Errorf("expect: %v\nactual: %v", c.titles, titles)
			}
		})
	}

}
