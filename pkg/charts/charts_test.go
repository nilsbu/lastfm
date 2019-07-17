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
