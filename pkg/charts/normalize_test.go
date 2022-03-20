package charts_test

import (
	"sort"
	"testing"

	"github.com/nilsbu/lastfm/pkg/charts"
)

// TODO create helper file
func mapCharts(data map[string][]float64) charts.LazyCharts {
	keys := []string{}
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	titles := make([]charts.Title, len(keys))
	for i, k := range keys {
		titles[i] = charts.KeyTitle(k)
	}

	return charts.FromMap(data)
}

func TestNormalizer(t *testing.T) {
	for _, c := range []struct {
		name           string
		actual, expect charts.LazyCharts
	}{
		{
			"NormalizeColumn",
			charts.NormalizeColumn(mapCharts(map[string][]float64{
				"A": {1, 2, 1, 0, 0, 1},
				"B": {1, 0, 14, 1, 0, 1},
				"C": {2, 2, 1, 1, 0, 0},
			})),
			mapCharts(map[string][]float64{
				"A": {.25, .5, .0625, 0, 0, .5},
				"B": {.25, 0, .875, .5, 0, .5},
				"C": {.5, .5, .0625, .5, 0, 0},
			}),
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			checkLazyCharts(t, c.expect, c.actual, 5)
		})
	}
}
