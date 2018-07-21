package charts

import (
	"math"
)

type Normalizer interface {
	Normalize(charts Charts) Charts
}

type SimpleNormalizer struct{}

func (SimpleNormalizer) Normalize(charts Charts) Charts {
	total := charts.Total()

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
	Sigma       float64
	MirrorFront bool
	MirrorBack  bool
}

func (n GaussianNormalizer) Normalize(charts Charts) Charts {
	wing := int(2 * n.Sigma)
	kernel := getGaussianKernel(n.Sigma, 2*wing+1)

	blurred := charts.mapLine(func(in []float64, out []float64) {
		for i := range out {
			for j := range kernel {
				jj := i + j - wing
				if jj < 0 && !n.MirrorFront {
					continue
				}
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
	})

	return SimpleNormalizer{}.Normalize(blurred)
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
