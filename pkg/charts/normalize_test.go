package charts

import (
	"math"
	"testing"
)

func SSD(a, b Charts) float64 {
	var ssd float64
	for key := range a {
		aline, bline := a[key], b[key]
		for i := range aline {
			diff := aline[i] - bline[i]
			ssd += diff * diff
		}
	}

	return ssd
}

func TestNormalizer(t *testing.T) {
	v := []float64{1, math.Exp(-.5), math.Exp(-2)}

	cases := []struct {
		charts     Charts
		normalizer Normalizer
		normalized Charts
	}{
		{
			Charts{},
			SimpleNormalizer{},
			Charts{},
		},
		{
			Charts{"a": []float64{3, 4, 7}, "b": []float64{1, 1, 1}},
			SimpleNormalizer{},
			Charts{"a": []float64{.75, .8, .875}, "b": []float64{.25, .2, .125}},
		},
		{
			Charts{"a": []float64{0, 1}, "b": []float64{0, 0}},
			SimpleNormalizer{},
			Charts{"a": []float64{0, 1}, "b": []float64{0, 0}},
		},
		{
			Charts{},
			GaussianNormalizer{Sigma: 1},
			Charts{},
		},
		{
			Charts{"a": []float64{1, 0, 0},
				"b": []float64{1, 1, 1}},
			GaussianNormalizer{Sigma: 1},
			Charts{
				"a": []float64{
					(1*v[0] + 0*v[1] + 0*v[2]) / (2*v[0] + 1*v[1] + 1*v[2]),
					(0*v[0] + 1*v[1] + 0*v[2]) / (1*v[0] + 3*v[1] + 0*v[2]),
					(0*v[0] + 0*v[1] + 1*v[2]) / (1*v[0] + 1*v[1] + 2*v[2]),
				},
				"b": []float64{
					(1*v[0] + 1*v[1] + 1*v[2]) / (2*v[0] + 1*v[1] + 1*v[2]),
					(1*v[0] + 2*v[1] + 0*v[2]) / (1*v[0] + 3*v[1] + 0*v[2]),
					(1*v[0] + 1*v[1] + 1*v[2]) / (1*v[0] + 1*v[1] + 2*v[2]),
				}},
		},
		{
			Charts{"a": []float64{1, 1, 1}},
			GaussianNormalizer{Sigma: 12, MirrorFront: true, MirrorBack: true},
			Charts{"a": []float64{1, 1, 1}},
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			normalized := c.normalizer.Normalize(c.charts)
			if SSD(normalized, c.normalized) > 1e-8 {
				t.Errorf("wrong data:\nhas:  %v\nwant: %v", normalized, c.normalized)
			}
		})
	}
}
