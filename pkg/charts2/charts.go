package charts2

import "github.com/nilsbu/lastfm/pkg/charts"

type Charts struct {
	Headers charts.Interval
	titles  []Title
	Values  map[string][]float64
}
