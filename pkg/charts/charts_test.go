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
				Headers: dayHeaders{rsrc.ParseDay("2008-01-01")},
				Keys:    []Key{},
				Values:  [][]float64{}},
		},
		{
			[]map[string]float64{{}},
			rsrc.ParseDay("2008-01-01"),
			Charts{
				Headers: dayHeaders{rsrc.ParseDay("2008-01-01")},
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
				Headers: dayHeaders{rsrc.ParseDay("2008-01-01")},
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
				Headers: dayHeaders{rsrc.ParseDay("2000-01-01")},
				Keys:    []Key{simpleKey("A")},
				Values:  [][]float64{{}},
			},
			[]map[string]float64{},
		},
		{
			Charts{
				Headers: dayHeaders{rsrc.ParseDay("2000-01-01")},
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
				Headers: dayHeaders{rsrc.ParseDay("2000-01-01")},
				Keys:    []Key{},
				Values:  [][]float64{{}},
			},
			[]string{},
		},
		{
			Charts{
				Headers: dayHeaders{rsrc.ParseDay("2000-01-01")},
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

type brokenHeaders dayHeaders

func (h brokenHeaders) Index(day rsrc.Day) int {
	return 0
}

func (h brokenHeaders) At(index int) rsrc.Day {
	return dayHeaders(h).At(index)
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
				Headers: dayHeaders{rsrc.ParseDay("2000-01-01")},
				Keys:    []Key{},
				Values:  [][]float64{},
			},
			Charts{
				Headers: dayHeaders{rsrc.ParseDay("2000-01-01")},
				Keys:    []Key{},
				Values:  [][]float64{},
			},
			true,
		},
		{
			"equal",
			Charts{
				Headers: dayHeaders{rsrc.ParseDay("2000-01-01")},
				Keys:    []Key{simpleKey("xx"), simpleKey("yy")},
				Values:  [][]float64{{1, 0, 1, 2}, {5, 5, 6, 7}},
			},
			Charts{
				Headers: dayHeaders{rsrc.ParseDay("2000-01-01")},
				Keys:    []Key{simpleKey("yy"), simpleKey("xx")},
				Values:  [][]float64{{5, 5, 6, 7}, {1, 0, 1, 2}},
			},
			true,
		},
		{
			"different date",
			Charts{
				Headers: dayHeaders{rsrc.ParseDay("2000-01-01")},
				Keys:    []Key{simpleKey("xx")},
				Values:  [][]float64{{1}},
			},
			Charts{
				Headers: dayHeaders{rsrc.ParseDay("2000-01-02")},
				Keys:    []Key{simpleKey("xx")},
				Values:  [][]float64{{1}},
			},
			false,
		},
		{
			"different length",
			Charts{
				Headers: dayHeaders{rsrc.ParseDay("2000-01-01")},
				Keys:    []Key{simpleKey("xx")},
				Values:  [][]float64{{1}},
			},
			Charts{
				Headers: dayHeaders{rsrc.ParseDay("2000-01-01")},
				Keys:    []Key{simpleKey("xx")},
				Values:  [][]float64{{1, 2}},
			},
			false,
		},
		{
			"different values",
			Charts{
				Headers: dayHeaders{rsrc.ParseDay("2000-01-01")},
				Keys:    []Key{simpleKey("xx")},
				Values:  [][]float64{{3, 3, 1}},
			},
			Charts{
				Headers: dayHeaders{rsrc.ParseDay("2000-01-01")},
				Keys:    []Key{simpleKey("xx")},
				Values:  [][]float64{{3, 3, 2}},
			},
			false,
		},
		{
			"different key",
			Charts{
				Headers: dayHeaders{rsrc.ParseDay("2000-01-01")},
				Keys:    []Key{simpleKey("xx")},
				Values:  [][]float64{{3, 3, 1}},
			},
			Charts{
				Headers: dayHeaders{rsrc.ParseDay("2000-01-01")},
				Keys:    []Key{simpleKey("xy")},
				Values:  [][]float64{{3, 3, 1}},
			},
			false,
		},
		{
			"different header function, same effective headers",
			Charts{
				Headers: dayHeaders{rsrc.ParseDay("2000-01-01")},
				Keys:    []Key{simpleKey("xx")},
				Values:  [][]float64{{3, 3, 1}},
			},
			Charts{
				Headers: header("d", rsrc.ParseDay("2000-01-01"), rsrc.ParseDay("2000-01-04")),
				Keys:    []Key{simpleKey("xx")},
				Values:  [][]float64{{3, 3, 1}},
			},
			true,
		},
		{
			"broken headers",
			Charts{
				Headers: dayHeaders{rsrc.ParseDay("2000-01-01")},
				Keys:    []Key{simpleKey("xx")},
				Values:  [][]float64{{3, 3, 1}},
			},
			Charts{
				Headers: brokenHeaders(dayHeaders{rsrc.ParseDay("2000-01-01")}),
				Keys:    []Key{simpleKey("xx")},
				Values:  [][]float64{{3, 3, 1}},
			},
			false,
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			eq := c.a.Equal(c.b)

			if c.eq && !eq {
				t.Error("charts not recognized as equal")
			} else if !c.eq && eq {
				t.Error("charts not recognized as unequal")
			}
		})
	}
}
