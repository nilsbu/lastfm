package charts

import (
	"math"
	"testing"

	"github.com/nilsbu/lastfm/pkg/rsrc"
)

func SSD(a, b Charts) float64 {
	var ssd float64
	for i := range a.Values {
		aline, bline := a.Values[i], b.Values[i]
		for j := range aline {
			diff := aline[j] - bline[j]
			ssd += diff * diff
		}
	}

	return ssd
}

func TestNormalizer(t *testing.T) {
	weight := 1 / math.Sqrt(2*math.Pi)
	v := []float64{
		weight * 1,
		weight * math.Exp(-.5),
		weight * math.Exp(-2)}

	cases := []struct {
		charts     Charts
		normalizer Normalizer
		normalized Charts
	}{
		{
			Charts{
				Headers: Days(rsrc.ParseDay("2018-01-01"), rsrc.ParseDay("2018-01-01")),
				Keys:    []Key{},
				Values:  [][]float64{}},
			SimpleNormalizer{},
			Charts{
				Headers: Days(rsrc.ParseDay("2018-01-01"), rsrc.ParseDay("2018-01-01")),
				Keys:    []Key{},
				Values:  [][]float64{}},
		},
		{
			Charts{
				Headers: Days(rsrc.ParseDay("2018-01-01"), rsrc.ParseDay("2018-01-04")),
				Keys:    []Key{simpleKey("a"), simpleKey("b")},
				Values:  [][]float64{{3, 4, 7}, {1, 1, 1}}},
			SimpleNormalizer{},
			Charts{
				Headers: Days(rsrc.ParseDay("2018-01-01"), rsrc.ParseDay("2018-01-04")),
				Keys:    []Key{simpleKey("a"), simpleKey("b")},
				Values:  [][]float64{{.75, .8, .875}, {.25, .2, .125}}},
		},
		{
			Charts{
				Headers: Days(rsrc.ParseDay("2018-01-01"), rsrc.ParseDay("2018-01-03")),
				Keys:    []Key{simpleKey("a"), simpleKey("b")},
				Values:  [][]float64{{0, 1}, {0, 0}}},
			SimpleNormalizer{},
			Charts{
				Headers: Days(rsrc.ParseDay("2018-01-01"), rsrc.ParseDay("2018-01-03")),
				Keys:    []Key{simpleKey("a"), simpleKey("b")},
				Values:  [][]float64{{0, 1}, {0, 0}}},
		},
		{
			Charts{
				Headers: Days(rsrc.ParseDay("2018-01-01"), rsrc.ParseDay("2018-01-01")),
				Keys:    []Key{},
				Values:  [][]float64{}},
			GaussianNormalizer{Sigma: 1.0},
			Charts{
				Headers: Days(rsrc.ParseDay("2018-01-01"), rsrc.ParseDay("2018-01-01")),
				Keys:    []Key{},
				Values:  [][]float64{}},
		},
		{
			Charts{
				Headers: Days(rsrc.ParseDay("2018-01-01"), rsrc.ParseDay("2018-01-04")),
				Keys:    []Key{simpleKey("a"), simpleKey("b")},
				Values:  [][]float64{{1, 0, 0}, {1, 1, 1}}},
			GaussianNormalizer{Sigma: 1},
			Charts{
				Headers: Days(rsrc.ParseDay("2018-01-01"), rsrc.ParseDay("2018-01-04")),
				Keys:    []Key{simpleKey("a"), simpleKey("b")},
				Values: [][]float64{{
					(1*v[0] + 0*v[1] + 0*v[2]) / (2*v[0] + 2*v[1] + 2*v[2]),
					(0*v[0] + 1*v[1] + 0*v[2]) / (1*v[0] + 3*v[1] + 1*v[2]),
					(0*v[0] + 0*v[1] + 1*v[2]) / (1*v[0] + 1*v[1] + 2*v[2]),
				}, {
					(1*v[0] + 2*v[1] + 2*v[2]) / (2*v[0] + 2*v[1] + 2*v[2]),
					(1*v[0] + 2*v[1] + 1*v[2]) / (1*v[0] + 3*v[1] + 1*v[2]),
					(1*v[0] + 1*v[1] + 1*v[2]) / (1*v[0] + 1*v[1] + 2*v[2]),
				}}},
		},
		{
			Charts{
				Headers: Days(rsrc.ParseDay("2018-01-01"), rsrc.ParseDay("2018-01-04")),
				Keys:    []Key{simpleKey("a")},
				Values:  [][]float64{{1, 1, 1}}},
			GaussianNormalizer{Sigma: 12, MirrorBack: true},
			Charts{
				Headers: Days(rsrc.ParseDay("2018-01-01"), rsrc.ParseDay("2018-01-04")),
				Keys:    []Key{simpleKey("a")},
				Values:  [][]float64{{1, 1, 1}}},
		},
		{
			Charts{
				Headers: Days(rsrc.ParseDay("2018-01-01"), rsrc.ParseDay("2018-01-04")),
				Keys:    []Key{simpleKey("a")},
				Values:  [][]float64{{1, 1, 1}}},
			SongDurations{},
			Charts{
				Headers: Days(rsrc.ParseDay("2018-01-01"), rsrc.ParseDay("2018-01-04")),
				Keys:    []Key{simpleKey("a")},
				Values:  [][]float64{{1, 1, 1}}},
		},
		{
			Charts{
				Headers: Days(rsrc.ParseDay("2018-01-01"), rsrc.ParseDay("2018-01-04")),
				Keys:    []Key{simpleKey("a")},
				Values:  [][]float64{{1, 1, 1}}},
			SongDurations{"": {"": 60}},
			Charts{
				Headers: Days(rsrc.ParseDay("2018-01-01"), rsrc.ParseDay("2018-01-04")),
				Keys:    []Key{simpleKey("a")},
				Values:  [][]float64{{60, 60, 60}}},
		},
		{
			Charts{
				Headers: Days(rsrc.ParseDay("2018-01-01"), rsrc.ParseDay("2018-01-04")),
				Keys:    []Key{simpleKey("a")},
				Values:  [][]float64{{1, 1, 1}}},
			SongDurations{"": {"": 60}, "a": {"b": 2}},
			Charts{
				Headers: Days(rsrc.ParseDay("2018-01-01"), rsrc.ParseDay("2018-01-04")),
				Keys:    []Key{simpleKey("a")},
				Values:  [][]float64{{60, 60, 60}}},
		},
		{
			Charts{
				Headers: Days(rsrc.ParseDay("2018-01-01"), rsrc.ParseDay("2018-01-04")),
				Keys:    []Key{Song{Artist: "a", Title: "b"}, Song{Artist: "a", Title: "c"}},
				Values:  [][]float64{{1, 1, 1}, {1, 2, 3}}},
			SongDurations{"": {"": 60}, "a": {"b": 2}},
			Charts{
				Headers: Days(rsrc.ParseDay("2018-01-01"), rsrc.ParseDay("2018-01-04")),
				Keys:    []Key{Song{Artist: "a", Title: "b"}, Song{Artist: "a", Title: "c"}},
				Values:  [][]float64{{2, 2, 2}, {60, 120, 180}}},
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			normalized := c.normalizer.Normalize(c.charts)
			ssd := SSD(normalized, c.normalized)
			if ssd > 1e-8 {
				t.Errorf("wrong data (SSD = %v):\nhas:  %v\nwant: %v",
					ssd, normalized, c.normalized)
			}
		})
	}
}
