package charts_test

import (
	"reflect"
	"testing"

	"github.com/nilsbu/lastfm/pkg/charts"
)

func TestDiffCharts(t *testing.T) {
	for _, c := range []struct {
		name    string
		charts  charts.Charts
		prevIdx int
		prev    []struct {
			title charts.Title
			place int
			value float64
		}
	}{
		{
			"empty",
			charts.FromMap(map[string][]float64{}),
			0,
			nil,
		},
		{
			"switch order",
			charts.FromMap(map[string][]float64{
				"A": {1, 5},
				"B": {2, 4},
			}),
			0,
			[]struct {
				title charts.Title
				place int
				value float64
			}{
				{charts.KeyTitle("A"), 1, 1},
				{charts.KeyTitle("B"), 0, 2},
			},
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			ch := charts.NewDiffCharts(c.charts, c.prevIdx)

			if len(ch.Titles()) != len(c.charts.Titles()) {
				t.Errorf("titles: %v != %v", ch.Titles(), c.charts.Titles())
			} else {
				for i, title := range ch.Titles() {
					if title != c.charts.Titles()[i] {
						t.Errorf("titles: %v != %v", ch.Titles(), c.charts.Titles())
					}
				}

				expect, _ := c.charts.Data(ch.Titles(), 0, ch.Len())
				actual, err := ch.Data(ch.Titles(), 0, ch.Len())
				if err != nil {
					t.Errorf("data: %v", err)
				} else if !reflect.DeepEqual(actual, expect) {
					t.Errorf("data: %v != %v", actual, expect)
				}

				for _, title := range ch.Titles() {
					if place, value, err := ch.Previous(title); err != nil {
						t.Errorf("previous: %v", err)
					} else {
						for _, prev := range c.prev {
							if prev.title == title {
								if place != prev.place {
									t.Errorf("previous: place %v != %v", place, prev.place)
								}
								if value != prev.value {
									t.Errorf("previous: value %v != %v", value, prev.value)
								}
							}
						}
					}
				}
			}
		})
	}
}

// mapCharts(map[string][]float64{
// 	"A": {1, 2, 1, 0, 0, 1},
// 	"B": {1, 0, 14, 1, 0, 1},
// 	"C": {2, 2, 1, 1, 0, 0},
// })
