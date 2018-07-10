package charts

import (
	"reflect"
	"testing"
)

func TestChartsSupertags(t *testing.T) {
	cases := []struct {
		charts    Charts
		tags      map[string][]Tag
		supertags map[string]string
		tagcharts Charts
	}{
		{
			Charts{},
			map[string][]Tag{},
			map[string]string{"a": "v", "b": "c", "c": "c", "d": "c", "e": "v"},
			Charts{
				"c": []float64{},
				"v": []float64{},
				"":  []float64{},
			},
		},
		{
			Charts{
				"asdf": []float64{7, 1},
				"bbh":  []float64{10, 2},
			},
			map[string][]Tag{
				"asdf": []Tag{Tag{"b", 23, 1, 100}},
				"bbh":  []Tag{Tag{"d", 500, 21, 100}},
			},
			map[string]string{"a": "v", "b": "c", "c": "c", "d": "c", "e": "v"},
			Charts{
				"c": []float64{17, 3},
				"v": []float64{0, 0},
				"":  []float64{0, 0},
			},
		},
		{
			Charts{
				"asdf": []float64{7, 1},
				"bbh":  []float64{10, 2},
				"33w":  []float64{0, 2},
				"wer":  []float64{7, 9},
			},
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
			Charts{
				"c": []float64{14, 10},
				"v": []float64{10, 2},
				"":  []float64{0, 2},
			},
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			tagcharts := c.charts.Supertags(c.tags, c.supertags)

			if !reflect.DeepEqual(tagcharts, c.tagcharts) {
				t.Errorf("wrong data:\nhas:  %v\nwant: %v",
					tagcharts, c.tagcharts)
			}
		})
	}
}

func TestChartsSplitBySupertag(t *testing.T) {
	cases := []struct {
		charts    Charts
		tags      map[string][]Tag
		supertags map[string]string
		tagcharts map[string]Charts
	}{
		{
			Charts{},
			map[string][]Tag{},
			map[string]string{"a": "v", "b": "c", "c": "c", "d": "c", "e": "v"},
			map[string]Charts{
				"c": Charts{},
				"v": Charts{},
				"":  Charts{},
			},
		},
		{
			Charts{
				"asdf": []float64{7, 1},
				"bbh":  []float64{10, 2},
			},
			map[string][]Tag{
				"asdf": []Tag{Tag{"b", 23, 1, 100}},
				"bbh":  []Tag{Tag{"e", 500, 21, 100}},
			},
			map[string]string{"a": "v", "b": "c", "c": "c", "d": "c", "e": "v"},
			map[string]Charts{
				"c": Charts{"asdf": []float64{7, 1}},
				"v": Charts{"bbh": []float64{10, 2}},
				"":  Charts{},
			},
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			tagcharts := c.charts.SplitBySupertag(c.tags, c.supertags)

			if !reflect.DeepEqual(tagcharts, c.tagcharts) {
				t.Errorf("wrong data:\nhas:  %v\nwant: %v",
					tagcharts, c.tagcharts)
			}
		})
	}
}
