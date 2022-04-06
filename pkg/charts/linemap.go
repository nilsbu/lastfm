package charts

import (
	"math"

	"github.com/nilsbu/async"
)

type lineMapCharts struct {
	chartsNode
	mapF   lineProc
	foldF  valueProc
	rangeF rangeSpec
}

type valueProc func(i int, line []float64) float64
type lineProc func(line []float64) []float64
type rangeSpec func(size, begin, end int) (b, e int)

func Sum(parent Charts) Charts {
	return &lineMapCharts{
		chartsNode: chartsNode{parent: parent},
		mapF: func(in []float64) []float64 {
			out := make([]float64, len(in))
			acc := 0.0
			for i := range in {
				acc += in[i]
				out[i] = acc
			}
			return out
		},
		foldF: func(i int, line []float64) float64 {
			acc := 0.0
			for j := 0; j <= i; j++ {
				acc += line[j]
			}
			return acc
		},
		rangeF: fromBeginRange,
	}
}

func Fade(parent Charts, hl float64) Charts {
	fac := math.Pow(0.5, 1.0/hl)
	return &lineMapCharts{
		chartsNode: chartsNode{parent: parent},
		mapF: func(in []float64) []float64 {
			out := make([]float64, len(in))
			acc := 0.0
			for i := range in {
				acc *= fac
				acc += in[i]
				out[i] = acc
			}
			return out
		},
		foldF: func(i int, line []float64) float64 {
			acc := 0.0
			for j := 0; j <= i; j++ {
				acc *= fac
				acc += line[j]
			}
			return acc
		},
		rangeF: fromBeginRange,
	}
}

// Max calculates the maximum of the parent charts.
func Max(parent Charts) Charts {
	return &lineMapCharts{
		chartsNode: chartsNode{parent: parent},
		mapF: func(in []float64) []float64 {
			out := make([]float64, len(in))
			acc := -math.MaxFloat64
			for i := range in {
				acc = math.Max(acc, in[i])
				out[i] = acc
			}
			return out
		},
		foldF: func(i int, line []float64) float64 {
			acc := 0.0
			for j := 0; j <= i; j++ {
				acc = math.Max(acc, line[j])
			}
			return acc
		},
		rangeF: fromBeginRange,
	}
}

// Gaussian blurs the data with a Gaussian kernel.
func Gaussian(
	parent Charts,
	sigma float64,
	width int,
	mirrorBegin, mirrorEnd bool) Charts {

	gaussian := make([]float64, width+1)
	fac := 1 / (sigma * math.Sqrt(2*math.Pi))
	for i := 0; i <= width; i++ {
		gaussian[i] = fac * math.Exp(-.5*float64(i*i)/sigma/sigma)
	}

	f := func(i int, line []float64) float64 {
		acc := 0.0
		b := i - width
		if b < 0 {
			if mirrorBegin {
				ee := width - i
				if ee > len(line) {
					ee = len(line)
				}
				for j := 0; j < ee; j++ {
					acc += gaussian[i+j+1] * line[j]
				}
			}

			b = 0
		}
		e := i + width + 1
		if e > len(line) {
			if mirrorEnd {
				bb := 2*len(line) - i - 1 - width
				if bb < 0 {
					bb = 0
				}
				for j := len(line) - 1; j >= bb; j-- {
					acc += gaussian[2*len(line)-(i+j+1)] * line[j]
				}
			}

			e = len(line)
		}

		for j := b; j < e; j++ {
			idx := j - i
			if idx < 0 {
				idx = -idx
			}
			acc += gaussian[idx] * line[j]
		}
		return acc
	}

	return &lineMapCharts{
		chartsNode: chartsNode{parent: parent},
		mapF: func(in []float64) []float64 {
			out := make([]float64, len(in))
			for i := range in {
				out[i] = f(i, in)
			}
			return out
		},
		foldF: f,
		rangeF: func(size, begin, end int) (b, e int) {
			if end+width > size {
				return 0, size
			}
			return 0, end + width
		},
	}
}

// Multiply multiplies the charts by a factor
// TODO test
func Multiply(parent Charts, s float64) Charts {
	return &lineMapCharts{
		chartsNode: chartsNode{parent: parent},
		mapF: func(in []float64) []float64 {
			out := make([]float64, len(in))
			for i := range in {
				out[i] = s * in[i]
			}
			return out
		},
		foldF: func(i int, line []float64) float64 {
			return s*line[len(line)] - 1
		},
		rangeF: fromBeginRange,
	}
}

func (l *lineMapCharts) Data(titles []Title, begin, end int) ([][]float64, error) {
	data := make([][]float64, len(titles))
	rb, re := l.rangeF(l.parent.Len(), begin, end)

	var err error
	if end-begin == 1 && false {
		err = async.Pie(len(titles), func(i int) error {
			if in, err := l.parent.Data([]Title{titles[i]}, rb, re); err != nil {
				return err
			} else {
				data[i] = []float64{l.foldF(end-1, in[0])}
				return nil
			}
		})
	} else {
		err = async.Pie(len(titles), func(i int) error {
			in, err := l.parent.Data([]Title{titles[i]}, rb, re)
			if err != nil {
				return err
			} else {
				out := l.mapF(in[0])
				data[i] = out[begin-rb : end-rb]
				return nil
			}
		})
	}

	if err != nil {
		return nil, err
	} else {
		return data, nil
	}
}
