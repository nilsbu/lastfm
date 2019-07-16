package charts

import (
	"math"
)

type Normalizer interface {
	Normalize(charts Charts) Charts
}

type SimpleNormalizer struct{}

func (SimpleNormalizer) Normalize(charts Charts) Charts {
	return charts.devideBy(charts.Total())
}

func (charts Charts) devideBy(total []float64) Charts {
	return charts.mapLine(func(in []float64, out []float64) {
		for i, x := range in {
			if total[i] > 0 {
				out[i] = x / total[i]
			} else {
				out[i] = 0
			}
		}
	})
}

type GaussianNormalizer struct {
	Sigma      float64
	MirrorBack bool
}

func (n GaussianNormalizer) Normalize(charts Charts) Charts {
	wing := int(2 * n.Sigma)
	kernel := getGaussianKernel(n.Sigma, 2*wing+1)

	total := Charts{
		Headers: charts.Headers,
		Keys:    []Key{simpleKey("total")},
		Values:  [][]float64{charts.Total()},
	}

	// TODO figure out a way to only normalize once

	blurredTotal := total.mapLine(func(in, out []float64) {
		n2 := n
		n2.MirrorBack = true
		n2.normalize(in, out, wing, kernel)
	})

	blurred := charts.mapLine(func(in, out []float64) {
		n.normalize(in, out, wing, kernel)
	})

	return blurred.devideBy(blurredTotal.Values[0])
}

func (n GaussianNormalizer) normalize(
	in, out []float64,
	wing int,
	kernel []float64) {
	for i := range out {
		for j := range kernel {
			jj := i + j - wing
			if jj >= len(in) && !n.MirrorBack {
				continue
			}

			for {
				if jj < 0 {
					jj = -jj
				} else if jj >= len(in) {
					jj = 2*len(in) - jj - 2
				} else {
					break
				}
			}

			out[i] += in[jj] * kernel[j]
		}
	}
}

func getGaussianKernel(sigma float64, width int) []float64 {
	kernel := make([]float64, width)

	var sum float64
	for i := range kernel {
		dx := float64(i - width/2)
		val := 1 / math.Sqrt(2*math.Pi*sigma*sigma) * math.Exp(-0.5*dx*dx/sigma/sigma)
		kernel[i] = val
		sum += val
	}

	for i := range kernel {
		kernel[i] /= sum
	}

	return kernel
}
