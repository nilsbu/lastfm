package charts

import (
	"math"
	"reflect"
	"testing"

	"github.com/nilsbu/lastfm/pkg/rsrc"
)

func TestChartsSum(t *testing.T) {
	cases := []struct {
		charts Charts
		sums   Charts
	}{
		{
			Charts{
				Headers: Days(rsrc.ParseDay("2000-01-01"), rsrc.ParseDay("2000-01-01")),
				Keys:    []Key{},
				Values:  [][]float64{}},
			Charts{
				Headers: Days(rsrc.ParseDay("2000-01-01"), rsrc.ParseDay("2000-01-01")),
				Keys:    []Key{},
				Values:  [][]float64{}},
		},
		{
			Charts{
				Headers: Days(rsrc.ParseDay("2000-01-01"), rsrc.ParseDay("2000-01-02")),
				Keys:    []Key{simpleKey("X")},
				Values:  [][]float64{{}}},
			Charts{
				Headers: Days(rsrc.ParseDay("2000-01-01"), rsrc.ParseDay("2000-01-02")),
				Keys:    []Key{simpleKey("X")},
				Values:  [][]float64{{}}},
		},
		{
			Charts{
				Headers: Days(rsrc.ParseDay("2000-01-01"), rsrc.ParseDay("2000-01-04")),
				Keys:    []Key{simpleKey("X"), simpleKey("o0o")},
				Values:  [][]float64{{1, 3, 4}, {0, 0, 7}}},
			Charts{
				Headers: Days(rsrc.ParseDay("2000-01-01"), rsrc.ParseDay("2000-01-04")),
				Keys:    []Key{simpleKey("X"), simpleKey("o0o")},
				Values:  [][]float64{{1, 4, 8}, {0, 0, 7}}},
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			sums := c.charts.Sum()

			if !c.sums.Equal(sums) {
				t.Error("charts are wrong")
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
			faded := Charts{
				Headers: Days(rsrc.ParseDay("2000-01-01"), rsrc.ParseDay("2000-01-04")),
				Keys:    []Key{simpleKey("XX")},
				Values:  [][]float64{c.charts},
			}.Fade(c.halflife)
			if 1 != len(faded.Values) {
				t.Fatalf("expected 1 line but got %v", len(faded.Values))
			}
			f := faded.Values[0]
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
			Charts{
				Headers: Days(rsrc.ParseDay("2000-01-01"), rsrc.ParseDay("2000-01-01")),
				Keys:    []Key{},
				Values:  [][]float64{}},
			0,
			Column{},
			false,
		},
		{
			Charts{
				Headers: Days(rsrc.ParseDay("2000-01-01"), rsrc.ParseDay("2000-01-01")),
				Keys:    []Key{simpleKey("X")},
				Values:  [][]float64{{}}},
			0,
			Column{},
			false,
		},
		{
			Charts{
				Headers: Days(rsrc.ParseDay("2000-01-01"), rsrc.ParseDay("2000-01-04")),
				Keys:    []Key{simpleKey("o0o"), simpleKey("lol"), simpleKey("X")},
				Values:  [][]float64{{0, 0, 7}, {1, 2, 3}, {1, 3, 4}}},
			1,
			Column{Score{"X", 3}, Score{"lol", 2}, Score{"o0o", 0}},
			true,
		},
		{
			Charts{
				Headers: Days(rsrc.ParseDay("2000-01-01"), rsrc.ParseDay("2000-01-04")),
				Keys:    []Key{simpleKey("X")},
				Values:  [][]float64{{1, 3, 4}}},
			-1,
			Column{Score{"X", 4}},
			true,
		},
		{
			Charts{
				Headers: Days(rsrc.ParseDay("2000-01-01"), rsrc.ParseDay("2000-01-04")),
				Keys:    []Key{simpleKey("X")},
				Values:  [][]float64{{1, 3, 4}}},
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

func TestChartsFullTitleColumn(t *testing.T) {
	testCases := []struct {
		charts Charts
		i      int
		column Column
		ok     bool
	}{
		{
			Charts{
				Headers: Days(rsrc.ParseDay("2000-01-01"), rsrc.ParseDay("2000-01-04")),
				Keys: []Key{
					NewCustomKey("k0", "a0", "f0"),
					NewCustomKey("k1", "a1", "f1"),
					NewCustomKey("k2", "a2", "f2")},
				Values: [][]float64{{0, 0, 7}, {1, 2, 3}, {1, 3, 4}}},
			1,
			Column{Score{"f2", 3}, Score{"f1", 2}, Score{"f0", 0}},
			true,
		},
	}

	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			column, err := tc.charts.FullTitleColumn(tc.i)
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
				Headers: Days(rsrc.ParseDay("2000-01-01"), rsrc.ParseDay("2000-01-04")),
				Keys:    []Key{simpleKey("o0o"), simpleKey("lol"), simpleKey("X")},
				Values:  [][]float64{{0, 0, 7}, {1, 2, 3}, {1, 3, 4}}},
			map[string]string{"X": "o0o"},
			Charts{
				Headers: Days(rsrc.ParseDay("2000-01-01"), rsrc.ParseDay("2000-01-04")),
				Keys:    []Key{simpleKey("o0o"), simpleKey("lol")},
				Values:  [][]float64{{1, 3, 11}, {1, 2, 3}}},
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			corrected := c.charts.Correct(c.correction)

			if !c.corrected.Equal(corrected) {
				t.Error("charts are wrong")
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
				Headers: Days(rsrc.ParseDay("2000-01-01"), rsrc.ParseDay("2000-01-04")),
				Keys:    []Key{simpleKey("o0o"), simpleKey("lol"), simpleKey("X")},
				Values:  [][]float64{{0, 0, 7}, {1, 2, 3}, {1, 3, 4}}},
			Charts{
				Headers: Days(rsrc.ParseDay("2000-01-01"), rsrc.ParseDay("2000-01-04")),
				Keys:    []Key{simpleKey("o0o"), simpleKey("lol"), simpleKey("X")},
				Values:  [][]float64{{3, 3, 1}, {1, 2, 3}, {1, 1, 2}}},
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
			Charts{
				Headers: Days(rsrc.ParseDay("2000-01-01"), rsrc.ParseDay("2000-01-01")),
				Keys:    []Key{},
				Values:  [][]float64{}},
			[]float64{},
		},
		{
			Charts{
				Headers: Days(rsrc.ParseDay("2000-01-01"), rsrc.ParseDay("2000-01-04")),
				Keys:    []Key{simpleKey("o0o")},
				Values:  [][]float64{{0, 0, 7}}},
			[]float64{0, 0, 7},
		},
		{
			Charts{
				Headers: Days(rsrc.ParseDay("2000-01-01"), rsrc.ParseDay("2000-01-04")),
				Keys:    []Key{simpleKey("o0o"), simpleKey("lol")},
				Values:  [][]float64{{0, 0, 7}, {1, 2, 3}}},
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
			Charts{
				Headers: Days(rsrc.ParseDay("2000-01-01"), rsrc.ParseDay("2000-01-01")),
				Keys:    []Key{},
				Values:  [][]float64{}},
			Column{},
		},
		{
			Charts{
				Headers: Days(rsrc.ParseDay("2000-01-01"), rsrc.ParseDay("2000-01-01")),
				Keys:    []Key{simpleKey("a")},
				Values:  [][]float64{{}}},
			Column{{Name: "a", Score: 0}},
		},
		{
			Charts{
				Headers: Days(rsrc.ParseDay("2000-01-01"), rsrc.ParseDay("2000-01-04")),
				Keys:    []Key{simpleKey("o0o")},
				Values:  [][]float64{{0, 0, 7}}},
			Column{{Name: "o0o", Score: 7}},
		},
		{
			Charts{
				Headers: Days(rsrc.ParseDay("2000-01-01"), rsrc.ParseDay("2000-01-04")),
				Keys:    []Key{simpleKey("o0o"), simpleKey("lol")},
				Values:  [][]float64{{0, 0, 7}, {1, 2, 0}}},
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

type brokenIntervals struct {
	dayIntervals
}

func (h brokenIntervals) Index(day rsrc.Day) int {
	return 0
}

type brokenKey struct {
	simpleKey
}

func (h brokenKey) FullTitle() string {
	return ""
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
				Keys:    []Key{simpleKey("xx")},
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
		{
			"different FullTitle",
			Charts{
				Headers: Days(rsrc.ParseDay("2000-01-01"), rsrc.ParseDay("2000-01-04")),
				Keys:    []Key{simpleKey("xx")},
				Values:  [][]float64{{3, 3, 1}},
			},
			Charts{
				Headers: Days(rsrc.ParseDay("2000-01-01"), rsrc.ParseDay("2000-01-04")),
				Keys:    []Key{brokenKey{simpleKey("xx")}},
				Values:  [][]float64{{3, 3, 1}},
			},
			false,
		},
		{
			"first has no header",
			Charts{
				Headers: nil,
				Keys:    []Key{simpleKey("xx")},
				Values:  [][]float64{{3, 3, 1}},
			},
			Charts{
				Headers: Days(rsrc.ParseDay("2000-01-01"), rsrc.ParseDay("2000-01-04")),
				Keys:    []Key{brokenKey{simpleKey("xx")}},
				Values:  [][]float64{{3, 3, 1}},
			},
			false,
		},
		{
			"second has no header",
			Charts{
				Headers: Days(rsrc.ParseDay("2000-01-01"), rsrc.ParseDay("2000-01-04")),
				Keys:    []Key{simpleKey("xx")},
				Values:  [][]float64{{3, 3, 1}},
			},
			Charts{
				Headers: nil,
				Keys:    []Key{brokenKey{simpleKey("xx")}},
				Values:  [][]float64{{3, 3, 1}},
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
