package charts

import (
	"reflect"
	"testing"

	"github.com/nilsbu/lastfm/pkg/rsrc"
)

func TestChartsSupertags(t *testing.T) {
	cases := []struct {
		charts      Charts
		tags        map[string][]Tag
		supertags   map[string]string
		corrections map[string]string
		tagcharts   Charts
	}{
		{
			Charts{
				Headers: Days(rsrc.ParseDay("2000-01-01"), rsrc.ParseDay("2000-01-01")),
				Keys:    []Key{},
				Values:  [][]float64{}},
			map[string][]Tag{},
			map[string]string{"a": "v", "b": "c", "c": "c", "d": "c", "e": "v"},
			nil,
			Charts{
				Headers: Days(rsrc.ParseDay("2000-01-01"), rsrc.ParseDay("2000-01-01")),
				Keys:    []Key{simpleKey("-"), simpleKey("c"), simpleKey("v")},
				Values:  [][]float64{{}, {}, {}}},
		},
		{
			Charts{
				Headers: Days(rsrc.ParseDay("2000-01-01"), rsrc.ParseDay("2000-01-03")),
				Keys:    []Key{simpleKey("asdf"), simpleKey("bbh")},
				Values:  [][]float64{{7, 1}, {10, 2}}},
			map[string][]Tag{
				"asdf": []Tag{Tag{"b", 23, 1, 100}},
				"bbh":  []Tag{Tag{"d", 500, 21, 100}},
			},
			map[string]string{"a": "v", "b": "c", "c": "c", "d": "c", "e": "v"},
			nil,
			Charts{
				Headers: Days(rsrc.ParseDay("2000-01-01"), rsrc.ParseDay("2000-01-03")),
				Keys:    []Key{simpleKey("-"), simpleKey("c"), simpleKey("v")},
				Values:  [][]float64{{0, 0}, {17, 3}, {0, 0}}},
		},
		{
			Charts{
				Headers: Days(rsrc.ParseDay("2000-01-01"), rsrc.ParseDay("2000-01-03")),
				Keys: []Key{
					simpleKey("asdf"),
					simpleKey("bbh"),
					simpleKey("33w"),
					simpleKey("wer")},
				Values: [][]float64{{7, 1}, {10, 2}, {0, 2}, {7, 9}}},
			map[string][]Tag{
				"asdf": []Tag{
					Tag{"1", 23, 1, 100},
					Tag{"b", 23, 1, 40}},
				"bbh": []Tag{
					Tag{"e", 500, 21, 100},
					Tag{"d", 500, 21, 11}},
				"33w": []Tag{
					Tag{"0", 23, 1, 100}},
				"wer": []Tag{
					Tag{"d", 500, 21, 100}},
			},
			map[string]string{"a": "v", "b": "c", "c": "c", "d": "c", "e": "v"},
			nil,
			Charts{
				Headers: Days(rsrc.ParseDay("2000-01-01"), rsrc.ParseDay("2000-01-03")),
				Keys:    []Key{simpleKey("-"), simpleKey("c"), simpleKey("v")},
				Values:  [][]float64{{0, 2}, {14, 10}, {10, 2}}},
		},
		{ // correction
			Charts{
				Headers: Days(rsrc.ParseDay("2000-01-01"), rsrc.ParseDay("2000-01-03")),
				Keys: []Key{
					simpleKey("asdf")},
				Values: [][]float64{{7, 1}}},
			map[string][]Tag{
				"asdf": []Tag{Tag{"b", 23, 1, 100}},
			},
			map[string]string{"a": "v", "b": "c", "c": "c", "d": "c", "e": "v"},
			map[string]string{"asdf": "v"},
			Charts{
				Headers: Days(rsrc.ParseDay("2000-01-01"), rsrc.ParseDay("2000-01-03")),
				Keys:    []Key{simpleKey("-"), simpleKey("c"), simpleKey("v")},
				Values:  [][]float64{{0, 0}, {0, 0}, {7, 1}}},
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			tagcharts := c.charts.Group(Supertags(c.tags, c.supertags, c.corrections))

			if !c.tagcharts.Equal(tagcharts) {
				t.Error("charts are wrong")
			}
		})
	}
}

func TestChartsSplitBySupertag(t *testing.T) {
	cases := []struct {
		charts      Charts
		tags        map[string][]Tag
		supertags   map[string]string
		corrections map[string]string
		tagcharts   map[string]Charts
	}{
		{
			Charts{
				Headers: Days(rsrc.ParseDay("2000-01-01"), rsrc.ParseDay("2000-01-01")),
				Keys:    []Key{},
				Values:  [][]float64{}},
			map[string][]Tag{},
			map[string]string{"a": "v", "b": "c", "c": "c", "d": "c", "e": "v"},
			nil,
			map[string]Charts{
				"c": Charts{
					Headers: Days(rsrc.ParseDay("2000-01-01"), rsrc.ParseDay("2000-01-01")),
					Keys:    []Key{},
					Values:  [][]float64{}},
				"v": Charts{
					Headers: Days(rsrc.ParseDay("2000-01-01"), rsrc.ParseDay("2000-01-01")),
					Keys:    []Key{},
					Values:  [][]float64{}},
				"-": Charts{
					Headers: Days(rsrc.ParseDay("2000-01-01"), rsrc.ParseDay("2000-01-01")),
					Keys:    []Key{},
					Values:  [][]float64{}},
			},
		},
		{
			Charts{
				Headers: Days(rsrc.ParseDay("2000-01-01"), rsrc.ParseDay("2000-01-03")),
				Keys:    []Key{simpleKey("asdf"), simpleKey("bbh")},
				Values:  [][]float64{{7, 1}, {10, 2}}},
			map[string][]Tag{
				"asdf": []Tag{Tag{"b", 23, 1, 100}},
				"bbh":  []Tag{Tag{"e", 500, 21, 100}},
			},
			map[string]string{"a": "v", "b": "c", "c": "c", "d": "c", "e": "v"},
			nil,
			map[string]Charts{
				"c": Charts{
					Headers: Days(rsrc.ParseDay("2000-01-01"), rsrc.ParseDay("2000-01-03")),
					Keys:    []Key{simpleKey("asdf")},
					Values:  [][]float64{{7, 1}}},
				"v": Charts{
					Headers: Days(rsrc.ParseDay("2000-01-01"), rsrc.ParseDay("2000-01-03")),
					Keys:    []Key{simpleKey("bbh")},
					Values:  [][]float64{{10, 2}}},
				"-": Charts{
					Headers: Days(rsrc.ParseDay("2000-01-01"), rsrc.ParseDay("2000-01-03")),
					Keys:    []Key{},
					Values:  [][]float64{}},
			},
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			tagcharts := c.charts.Split(Supertags(c.tags, c.supertags, c.corrections))

			if len(c.tagcharts) != len(tagcharts) {
				t.Errorf("unexpected length:\nhas:  %v\nwant: %v",
					len(tagcharts), len(c.tagcharts))
			}

			for name, exCharts := range c.tagcharts {
				if acCharts, ok := tagcharts[name]; ok {
					if !exCharts.Equal(acCharts) {
						t.Errorf("charts are wrong @ '%v'", name)
					}
				} else {
					t.Errorf("bucket '%v' missing", name)
				}
			}
		})
	}
}

func TestEmptySupertags(t *testing.T) {
	cha := Charts{
		Headers: Days(rsrc.ParseDay("2000-01-01"), rsrc.ParseDay("2000-01-03")),
		Keys:    []Key{simpleKey("a")},
		Values:  [][]float64{{2, 3}}}
	buckets := cha.Split(Supertags(nil, nil, nil))

	expected := map[string]Charts{"-": Charts{
		Headers: Days(rsrc.ParseDay("2000-01-01"), rsrc.ParseDay("2000-01-03")),
		Keys:    []Key{simpleKey("a")},
		Values:  [][]float64{{2, 3}}}}
	if !reflect.DeepEqual(buckets, expected) {
		t.Errorf("wrong data:\nhas:  %v\nwant: %v", buckets, expected)
	}
}
