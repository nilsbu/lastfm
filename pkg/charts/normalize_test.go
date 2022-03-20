package charts_test

import (
	"math"
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
	f := 1 / math.Sqrt(2*math.Pi)
	m := []float64{f * math.Exp(0), f * math.Exp(-.5), f * math.Exp(-2)}

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
		{
			"NormalizeGaussian",
			charts.NormalizeGaussian(mapCharts(map[string][]float64{
				"A": {0, 0, 1, 0, 0, 1},
				"B": {1, 0, 0, 2, 0, 0},
			}), 1, 2, true, true),
			mapCharts(map[string][]float64{
				"A": {
					(0*m[0] + 0*m[1] + 1*m[2]) / (1*m[0] + 1*m[1] + 1*m[2]),
					(0*m[0] + 1*m[1] + 0*m[2]) / (0*m[0] + 2*m[1] + 3*m[2]),
					(1*m[0] + 0*m[1] + 0*m[2]) / (1*m[0] + 2*m[1] + 1*m[2]),
					(0*m[0] + 1*m[1] + 1*m[2]) / (2*m[0] + 1*m[1] + 1*m[2]),
					(0*m[0] + 1*m[1] + 2*m[2]) / (0*m[0] + 3*m[1] + 2*m[2]),
					(1*m[0] + 1*m[1] + 0*m[2]) / (1*m[0] + 1*m[1] + 2*m[2]),
				},
				"B": {
					(1*m[0] + 1*m[1] + 0*m[2]) / (1*m[0] + 1*m[1] + 1*m[2]),
					(0*m[0] + 1*m[1] + 3*m[2]) / (0*m[0] + 2*m[1] + 3*m[2]),
					(0*m[0] + 2*m[1] + 1*m[2]) / (1*m[0] + 2*m[1] + 1*m[2]),
					(2*m[0] + 0*m[1] + 0*m[2]) / (2*m[0] + 1*m[1] + 1*m[2]),
					(0*m[0] + 2*m[1] + 0*m[2]) / (0*m[0] + 3*m[1] + 2*m[2]),
					(0*m[0] + 0*m[1] + 2*m[2]) / (1*m[0] + 1*m[1] + 2*m[2]),
					// 1,
				},
			}),
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			checkLazyCharts(t, c.expect, c.actual, 5)
		})
	}
}
