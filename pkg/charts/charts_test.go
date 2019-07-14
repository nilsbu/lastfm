package charts

import (
	"reflect"
	"sort"
	"testing"

	"github.com/nilsbu/lastfm/pkg/rsrc"
)

func TestCompileArtist(t *testing.T) {
	cases := []struct {
		days       []map[string]float64
		registered rsrc.Day
		charts     Charts
	}{
		{
			[]map[string]float64{},
			rsrc.ParseDay("2008-01-01"),
			Charts{
				Headers: Days(rsrc.ParseDay("2008-01-01"), rsrc.ParseDay("2008-01-01")),
				Keys:    []Key{},
				Values:  [][]float64{}},
		},
		{
			[]map[string]float64{{}},
			rsrc.ParseDay("2008-01-01"),
			Charts{
				Headers: Days(rsrc.ParseDay("2008-01-01"), rsrc.ParseDay("2008-01-02")),
				Keys:    []Key{},
				Values:  [][]float64{}},
		},
		{
			[]map[string]float64{
				{"ASD": 2},
				{"WASD": 1},
				{"ASD": 13, "WASD": 4},
			},
			rsrc.ParseDay("2008-01-01"),
			Charts{
				Headers: Days(rsrc.ParseDay("2008-01-01"), rsrc.ParseDay("2008-01-04")),
				Keys:    []Key{simpleKey("ASD"), simpleKey("WASD")},
				Values:  [][]float64{{2, 0, 13}, {0, 1, 4}}},
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			charts := CompileArtists(c.days, c.registered)

			if !c.charts.Equal(charts) {
				t.Error("charts are wrong")
			}
		})
	}
}

func TestCompileSongs(t *testing.T) {
	cases := []struct {
		days       [][]Song
		registered rsrc.Day
		charts     Charts
	}{
		{
			[][]Song{},
			rsrc.ParseDay("2008-01-01"),
			Charts{
				Headers: Days(rsrc.ParseDay("2008-01-01"), rsrc.ParseDay("2008-01-01")),
				Keys:    []Key{},
				Values:  [][]float64{}},
		},
		{
			[][]Song{{}},
			rsrc.ParseDay("2008-01-01"),
			Charts{
				Headers: Days(rsrc.ParseDay("2008-01-01"), rsrc.ParseDay("2008-01-02")),
				Keys:    []Key{},
				Values:  [][]float64{}},
		},
		{
			[][]Song{
				{
					{Artist: "A", Title: "s", Album: "x"},
					{Artist: "A", Title: "s", Album: "x"},
					{Artist: "A", Title: "t", Album: "x"},
					{Artist: "B", Title: "s", Album: "x"},
				},
				{
					{Artist: "A", Title: "t", Album: "x"},
					{Artist: "C", Title: "w", Album: "x"},
				},
			},
			rsrc.ParseDay("2008-01-01"),
			Charts{
				Headers: Days(rsrc.ParseDay("2008-01-01"), rsrc.ParseDay("2008-01-03")),
				Keys: []Key{
					Song{Artist: "A", Title: "s", Album: "x"},
					Song{Artist: "A", Title: "t", Album: "x"},
					Song{Artist: "B", Title: "s", Album: "x"},
					Song{Artist: "C", Title: "w", Album: "x"},
				},
				Values: [][]float64{
					{2, 0}, {1, 1}, {1, 0}, {0, 1}}},
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			charts := CompileSongs(c.days, c.registered)

			if err := c.charts.AssertEqual(charts); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestChartsUnravelDays(t *testing.T) {
	cases := []struct {
		charts Charts
		days   []map[string]float64
	}{
		{
			Charts{},
			[]map[string]float64{},
		},
		{
			Charts{
				Headers: Days(rsrc.ParseDay("2000-01-01"), rsrc.ParseDay("2000-01-01")),
				Keys:    []Key{simpleKey("A")},
				Values:  [][]float64{{}},
			},
			[]map[string]float64{},
		},
		{
			Charts{
				Headers: Days(rsrc.ParseDay("2000-01-01"), rsrc.ParseDay("2000-01-04")),
				Keys:    []Key{simpleKey("ASD"), simpleKey("WASD")},
				Values:  [][]float64{{2, 0, 13}, {0, 1, 4}},
			},
			[]map[string]float64{
				{"ASD": 2},
				{"WASD": 1},
				{"ASD": 13, "WASD": 4},
			},
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			days := c.charts.UnravelDays()

			if !reflect.DeepEqual(days, c.days) {
				t.Errorf("wrong data:\nhas:  %v\nwant: %v", days, c.days)
			}
		})
	}
}

func TestChartsGetKeys(t *testing.T) {
	cases := []struct {
		charts Charts
		keys   []string
	}{
		{
			Charts{
				Headers: Days(rsrc.ParseDay("2000-01-01"), rsrc.ParseDay("2000-01-01")),
				Keys:    []Key{},
				Values:  [][]float64{{}},
			},
			[]string{},
		},
		{
			Charts{
				Headers: Days(rsrc.ParseDay("2000-01-01"), rsrc.ParseDay("2000-01-03")),
				Keys:    []Key{simpleKey("xx"), simpleKey("yy")},
				Values:  [][]float64{{32, 45}, {32, 45}}},
			[]string{"xx", "yy"},
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			keys := c.charts.GetKeys()

			sort.Strings(keys)
			sort.Strings(c.keys)
			if !reflect.DeepEqual(keys, c.keys) {
				t.Errorf("wrong data (sorted):\nhas:  %v\nwant: %v",
					keys, c.keys)
			}
		})
	}
}

type brokenIntervals struct {
	dayIntervals
}

func (h brokenIntervals) Index(day rsrc.Day) int {
	return 0
}

func TestChartsEqual(t *testing.T) {
	cases := []struct {
		name string
		a    Charts
		b    Charts
		eq   bool
	}{
		{
			"empty",
			Charts{
				Headers: Days(rsrc.ParseDay("2000-01-01"), rsrc.ParseDay("2000-01-01")),
				Keys:    []Key{},
				Values:  [][]float64{},
			},
			Charts{
				Headers: Days(rsrc.ParseDay("2000-01-01"), rsrc.ParseDay("2000-01-01")),
				Keys:    []Key{},
				Values:  [][]float64{},
			},
			true,
		},
		{
			"equal",
			Charts{
				Headers: Days(rsrc.ParseDay("2000-01-01"), rsrc.ParseDay("2000-01-05")),
				Keys:    []Key{simpleKey("xx"), simpleKey("yy")},
				Values:  [][]float64{{1, 0, 1, 2}, {5, 5, 6, 7}},
			},
			Charts{
				Headers: Days(rsrc.ParseDay("2000-01-01"), rsrc.ParseDay("2000-01-05")),
				Keys:    []Key{simpleKey("yy"), simpleKey("xx")},
				Values:  [][]float64{{5, 5, 6, 7}, {1, 0, 1, 2}},
			},
			true,
		},
		{
			"different date",
			Charts{
				Headers: Days(rsrc.ParseDay("2000-01-01"), rsrc.ParseDay("2000-01-02")),
				Keys:    []Key{simpleKey("xx")},
				Values:  [][]float64{{1}},
			},
			Charts{
				Headers: Days(rsrc.ParseDay("2000-01-02"), rsrc.ParseDay("2000-01-02")),
				Keys:    []Key{simpleKey("xx")},
				Values:  [][]float64{{1}},
			},
			false,
		},
		{
			"different length",
			Charts{
				Headers: Days(rsrc.ParseDay("2000-01-01"), rsrc.ParseDay("2000-01-02")),
				Keys:    []Key{simpleKey("xx")},
				Values:  [][]float64{{1}},
			},
			Charts{
				Headers: Days(rsrc.ParseDay("2000-01-01"), rsrc.ParseDay("2000-01-03")),
				Keys:    []Key{simpleKey("xx")},
				Values:  [][]float64{{1, 2}},
			},
			false,
		},
		{
			"different values",
			Charts{
				Headers: Days(rsrc.ParseDay("2000-01-01"), rsrc.ParseDay("2000-01-04")),
				Keys:    []Key{simpleKey("xx")},
				Values:  [][]float64{{3, 3, 1}},
			},
			Charts{
				Headers: Days(rsrc.ParseDay("2000-01-01"), rsrc.ParseDay("2000-01-04")),
				Keys:    []Key{simpleKey("xx")},
				Values:  [][]float64{{3, 3, 2}},
			},
			false,
		},
		{
			"different key",
			Charts{
				Headers: Days(rsrc.ParseDay("2000-01-01"), rsrc.ParseDay("2000-01-04")),
				Keys:    []Key{simpleKey("xx")},
				Values:  [][]float64{{3, 3, 1}},
			},
			Charts{
				Headers: Days(rsrc.ParseDay("2000-01-01"), rsrc.ParseDay("2000-01-04")),
				Keys:    []Key{simpleKey("xy")},
				Values:  [][]float64{{3, 3, 1}},
			},
			false,
		},
		{
			"different artist",
			Charts{
				Headers: Days(rsrc.ParseDay("2000-01-01"), rsrc.ParseDay("2000-01-04")),
				Keys:    []Key{simpleKey("xx")},
				Values:  [][]float64{{3, 3, 1}},
			},
			Charts{
				Headers: Days(rsrc.ParseDay("2000-01-01"), rsrc.ParseDay("2000-01-04")),
				Keys:    []Key{tagKey("xx")},
				Values:  [][]float64{{3, 3, 1}},
			},
			false,
		},
		// TODO test FullTitle
		{
			"different begin",
			Charts{
				Headers: Days(rsrc.ParseDay("2000-01-01"), rsrc.ParseDay("2000-01-04")),
				Keys:    []Key{simpleKey("xx")},
				Values:  [][]float64{{3, 3, 1}},
			},
			Charts{
				Headers: Days(rsrc.ParseDay("2001-01-01"), rsrc.ParseDay("2001-01-04")),
				Keys:    []Key{simpleKey("xx")},
				Values:  [][]float64{{3, 3, 1}},
			},
			false,
		},
		{
			"different end",
			Charts{
				Headers: Months(rsrc.ParseDay("2000-01-01"), rsrc.ParseDay("2000-04-01"), 1),
				Keys:    []Key{simpleKey("xx")},
				Values:  [][]float64{{3, 3, 1}},
			},
			Charts{
				Headers: Days(rsrc.ParseDay("2000-01-01"), rsrc.ParseDay("2000-01-04")),
				Keys:    []Key{simpleKey("xy")},
				Values:  [][]float64{{3, 3, 1}},
			},
			false,
		},
		{
			"broken index in headers",
			Charts{
				Headers: Days(rsrc.ParseDay("2000-01-01"), rsrc.ParseDay("2000-01-04")),
				Keys:    []Key{simpleKey("xx")},
				Values:  [][]float64{{3, 3, 1}},
			},
			Charts{
				Headers: brokenIntervals{dayIntervals{intervalsBase{
					begin: rsrc.ParseDay("2000-01-01"),
					n:     3,
					step:  1,
				}}},
				Keys:   []Key{simpleKey("xx")},
				Values: [][]float64{{3, 3, 1}},
			},
			false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			eq := c.a.Equal(c.b)

			if c.eq && !eq {
				t.Error("charts not recognized as equal (a first)")
			} else if !c.eq && eq {
				t.Error("charts not recognized as unequal (a first)")
			}

			eq = c.b.Equal(c.a)

			if c.eq && !eq {
				t.Error("charts not recognized as equal (b first)")
			} else if !c.eq && eq {
				t.Error("charts not recognized as unequal (b first)")
			}

			err := c.a.AssertEqual(c.b)

			if err == nil && !c.eq {
				t.Error("expected error but non occurred (a first)")
			} else if err != nil && c.eq {
				t.Errorf("unexpected error (a first): %v", err)
			}

			err = c.b.AssertEqual(c.a)

			if err == nil && !c.eq {
				t.Error("expected error but non occurred (b first)")
			} else if err != nil && c.eq {
				t.Errorf("unexpected error (b first): %v", err)
			}
		})
	}
}
